package agentpatchfixer

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const systemPrompt = "Repair the semantic intent of a failed Git mail patch by editing the current upstream checkout. " +
	"Preserve existing business logic except for the original patch's intended change. Read the original patch first, " +
	"then inspect only the relevant current file ranges. If the target already satisfies or supersedes the patch intent, " +
	"make no edits and explain the exact evidence. Otherwise prefer targeted replacements over rewriting unchanged content. " +
	"Edit only explicitly allowed files. The harness preserves the original author, date, subject, and commit message; " +
	"do not recreate commit metadata or emit patch text. Once the intent is satisfied, review git diff once, confirm no " +
	"unrelated drift, and finish immediately."

type modelFactory func(context.Context, string, string) (model, error)

type model interface {
	Converse(context.Context, modelRequest) (modelResponse, error)
}

type modelRequest struct {
	System    string
	Messages  []modelMessage
	Tools     []toolDefinition
	MaxTokens int32
}

type modelResponse struct {
	Message      modelMessage
	StopReason   string
	InputTokens  int
	OutputTokens int
}

type modelMessage struct {
	Role    string
	Content []contentBlock
}

type contentBlock struct {
	Text       string      `json:"text,omitempty"`
	ToolUse    *toolUse    `json:"tool_use,omitempty"`
	ToolResult *toolResult `json:"tool_result,omitempty"`
	raw        any
}

type toolUse struct {
	ID    string         `json:"id"`
	Name  string         `json:"name"`
	Input map[string]any `json:"input"`
}

type toolResult struct {
	ToolUseID string `json:"tool_use_id"`
	Content   string `json:"content"`
	IsError   bool   `json:"is_error,omitempty"`
}

type toolDefinition struct {
	Name        string
	Description string
	InputSchema map[string]any
}

type agentOutcome struct {
	StopReason string         `json:"stop_reason"`
	Summary    string         `json:"summary"`
	Turns      int            `json:"turns"`
	Tokens     int            `json:"total_tokens"`
	LastReply  []contentBlock `json:"-"`
}

func execute(ctx context.Context, request Request, runDir string, factory modelFactory) (*Result, error) {
	patchContent, err := os.ReadFile(request.FailedPatchPath)
	if err != nil {
		return nil, fmt.Errorf("reading failed patch: %v", err)
	}
	patchText := string(patchContent)
	metadata, err := parsePatchMetadata(ctx, patchText)
	if err != nil {
		return nil, err
	}
	allowedFiles, err := patchFiles(patchText)
	if err != nil {
		return nil, err
	}
	for _, editableFile := range request.EditableFiles {
		allowedFiles[filepath.ToSlash(filepath.Clean(editableFile))] = true
	}
	allowedRoots := parentDirectories(allowedFiles)

	workspaceParent, err := os.MkdirTemp("", "patch-fixer-workspace-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(workspaceParent)

	workspacePath, err := prepareWorkspace(ctx, request.UpstreamRepoPath, workspaceParent)
	if err != nil {
		return nil, fmt.Errorf("preparing patch workspace: %v", err)
	}
	patchWorkspace, err := newWorkspace(workspacePath, allowedFiles, allowedRoots, patchText)
	if err != nil {
		return nil, err
	}
	seed, err := seedWorkspace(ctx, workspacePath, request.FailedPatchPath)
	if err != nil {
		return nil, err
	}

	region := strings.TrimSpace(os.Getenv("AWS_REGION"))
	if region == "" {
		region = strings.TrimSpace(os.Getenv("AWS_DEFAULT_REGION"))
	}
	if region == "" {
		region = "us-west-2"
	}
	client, err := factory(ctx, modelID(), region)
	if err != nil {
		return nil, err
	}
	prompt := buildPrompt(request, metadata, allowedFiles, allowedRoots, seed)

	if diagnosticsFull() {
		if err := writeJSON(filepath.Join(runDir, "model-request.json"), map[string]any{
			"schema_version": requestSchemaVersion,
			"model_id":       modelID(),
			"system_prompt":  systemPrompt,
			"prompt":         prompt,
			"limits": map[string]int{
				"turns":        request.MaxTurns,
				"total_tokens": request.MaxTotalTokens,
			},
			"reject_files": seed.RejectFiles,
			"tools":        toolNames(),
		}); err != nil {
			return nil, err
		}
	}

	outcome, err := runAgent(ctx, client, patchWorkspace, prompt, request.MaxTurns, request.MaxTotalTokens, runDir)
	if err != nil {
		return nil, err
	}
	removeRejectFiles(workspacePath, seed.RejectFiles)

	result := &Result{
		SchemaVersion: requestSchemaVersion,
		StopReason:    outcome.StopReason,
		Summary:       boundedTail(outcome.Summary, 2000),
		Turns:         outcome.Turns,
		TotalTokens:   outcome.Tokens,
	}
	if outcome.StopReason != "end_turn" && outcome.StopReason != "stop_sequence" {
		result.Status = "agent_incomplete"
		diff, diffErr := patchWorkspace.showDiff(ctx)
		if diffErr == nil {
			_ = os.WriteFile(filepath.Join(runDir, "incomplete.diff"), []byte(diff), 0o600)
		}
	} else {
		candidatePath := filepath.Join(runDir, "candidate.patch")
		editedFiles, err := generateCandidate(ctx, workspacePath, metadata, candidatePath, allowedFiles, allowedRoots)
		if err != nil {
			return nil, err
		}
		result.EditedFiles = editedFiles
		if len(editedFiles) == 0 {
			result.Status = "no_changes"
		} else {
			result.Status = "candidate_generated"
			result.CandidatePatchPath = filepath.Base(candidatePath)
		}
	}

	if diagnosticsFull() {
		if err := writeJSON(filepath.Join(runDir, "model-response.json"), map[string]any{
			"schema_version":       result.SchemaVersion,
			"status":               result.Status,
			"stop_reason":          result.StopReason,
			"summary":              result.Summary,
			"candidate_patch_path": result.CandidatePatchPath,
			"edited_files":         result.EditedFiles,
			"turns":                result.Turns,
			"total_tokens":         result.TotalTokens,
			"message":              outcome.LastReply,
		}); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func runAgent(
	ctx context.Context,
	client model,
	workspace *workspace,
	prompt string,
	maxTurns int,
	maxTotalTokens int,
	runDir string,
) (agentOutcome, error) {
	messages := []modelMessage{{
		Role:    "user",
		Content: []contentBlock{{Text: prompt}},
	}}
	var outcome agentOutcome
	for turn := 1; turn <= maxTurns; turn++ {
		response, err := client.Converse(ctx, modelRequest{
			System:    systemPrompt,
			Messages:  messages,
			Tools:     patchTools(),
			MaxTokens: 16384,
		})
		if err != nil {
			return agentOutcome{}, fmt.Errorf("calling Bedrock Converse: %v", err)
		}
		outcome.Turns = turn
		outcome.Tokens += response.InputTokens + response.OutputTokens
		outcome.StopReason = response.StopReason
		outcome.Summary = responseText(response.Message)
		outcome.LastReply = response.Message.Content
		messages = append(messages, response.Message)

		if diagnosticsFull() {
			if err := appendTrace(filepath.Join(runDir, "trace.jsonl"), map[string]any{
				"turn":          turn,
				"stop_reason":   response.StopReason,
				"input_tokens":  response.InputTokens,
				"output_tokens": response.OutputTokens,
				"text":          responseText(response.Message),
				"tools":         requestedToolNames(response.Message),
			}); err != nil {
				return agentOutcome{}, err
			}
		}
		if outcome.Tokens > maxTotalTokens {
			outcome.StopReason = "limit_tokens"
			return outcome, nil
		}

		switch response.StopReason {
		case "end_turn", "stop_sequence":
			return outcome, nil
		case "tool_use":
			toolResults := executeToolUses(ctx, workspace, response.Message)
			if len(toolResults) == 0 {
				outcome.StopReason = "malformed_tool_use"
				return outcome, nil
			}
			messages = append(messages, modelMessage{Role: "user", Content: toolResults})
		default:
			return outcome, nil
		}
	}
	outcome.StopReason = "limit_turns"
	return outcome, nil
}

func executeToolUses(ctx context.Context, workspace *workspace, message modelMessage) []contentBlock {
	var results []contentBlock
	for _, block := range message.Content {
		if block.ToolUse == nil {
			continue
		}
		output, err := executeTool(ctx, workspace, *block.ToolUse)
		result := &toolResult{ToolUseID: block.ToolUse.ID, Content: output}
		if err != nil {
			result.Content = err.Error()
			result.IsError = true
		}
		results = append(results, contentBlock{ToolResult: result})
	}
	return results
}

func executeTool(ctx context.Context, workspace *workspace, request toolUse) (string, error) {
	switch request.Name {
	case "list_files":
		return workspace.listFiles(stringInput(request.Input, "pattern", "**/*"))
	case "read_file":
		return workspace.readFile(
			requiredStringInput(request.Input, "path"),
			intInput(request.Input, "start_line", 1),
			intInput(request.Input, "end_line", 200),
		)
	case "search_repository":
		return workspace.searchRepository(
			requiredStringInput(request.Input, "query"),
			stringInput(request.Input, "path", "."),
		)
	case "read_original_patch":
		return workspace.readOriginalPatch(
			intInput(request.Input, "start_line", 1),
			intInput(request.Input, "end_line", 200),
		), nil
	case "replace_text":
		return workspace.replaceText(
			requiredStringInput(request.Input, "path"),
			requiredStringInput(request.Input, "old_text"),
			requiredStringInput(request.Input, "new_text"),
		)
	case "replace_lines":
		return workspace.replaceLines(
			requiredStringInput(request.Input, "path"),
			intInput(request.Input, "start_line", 0),
			intInput(request.Input, "end_line", 0),
			requiredStringInput(request.Input, "new_text"),
		)
	case "write_file":
		return workspace.writeFile(
			requiredStringInput(request.Input, "path"),
			requiredStringInput(request.Input, "content"),
		)
	case "delete_file":
		return workspace.deleteFile(requiredStringInput(request.Input, "path"))
	case "show_diff":
		return workspace.showDiff(ctx)
	default:
		return "", fmt.Errorf("unknown tool %q", request.Name)
	}
}

func buildPrompt(
	request Request,
	metadata patchMetadata,
	allowedFiles map[string]bool,
	allowedRoots []string,
	seed patchSeed,
) string {
	files := make([]string, 0, len(allowedFiles))
	for file := range allowedFiles {
		files = append(files, file)
	}
	sort.Strings(files)
	roots := strings.Join(allowedRoots, ", ")
	if roots == "" {
		roots = "(none)"
	}
	rejects := strings.Join(seed.RejectFiles, ", ")
	if rejects == "" {
		rejects = "(none)"
	}
	return fmt.Sprintf(
		"Goal: repair this failed patch without changing unrelated business behavior.\n"+
			"Project: %s\nOld revision: %s\nNew revision: %s\nFailed patch: %s\n"+
			"Allowed files: %s\nAllowed directories for renamed files: %s\n"+
			"Original author: %s <%s>\nOriginal subject: %s\nPatch application error:\n%s\n\n"+
			"Git already applied every clean hunk from the original patch to this isolated checkout.\n"+
			"Reject files that still need resolution: %s\n"+
			"Read the reject files and relevant current file ranges first. Preserve the pre-applied changes "+
			"and resolve only the rejected intent. For large rejected sections, use replace_lines in small "+
			"ranges instead of embedding entire old and new files in replace_text. Leave the checkout unchanged "+
			"if the target already satisfies or supersedes the intent. Review the final diff once and finish.",
		request.ProjectName,
		request.CurrentRevision,
		request.TargetRevision,
		filepath.Base(request.FailedPatchPath),
		strings.Join(files, ", "),
		roots,
		metadata.AuthorName,
		metadata.AuthorEmail,
		metadata.Subject,
		boundedTail(request.FailureOutput, 12000),
		rejects,
	)
}

func parentDirectories(files map[string]bool) []string {
	roots := make(map[string]bool)
	for file := range files {
		parent := filepath.ToSlash(filepath.Dir(file))
		if parent != "." {
			roots[parent] = true
		}
	}
	result := make([]string, 0, len(roots))
	for root := range roots {
		result = append(result, root)
	}
	sort.Strings(result)
	return result
}

func patchTools() []toolDefinition {
	stringProperty := func(description string) map[string]any {
		return map[string]any{"type": "string", "description": description}
	}
	integerProperty := func(description string, minimum int) map[string]any {
		return map[string]any{"type": "integer", "description": description, "minimum": minimum}
	}
	object := func(properties map[string]any, required ...string) map[string]any {
		schema := map[string]any{
			"type":                 "object",
			"properties":           properties,
			"additionalProperties": false,
		}
		if len(required) > 0 {
			schema["required"] = required
		}
		return schema
	}
	return []toolDefinition{
		{Name: "list_files", Description: "List repository files matching a glob pattern.", InputSchema: object(map[string]any{"pattern": stringProperty("Glob pattern")})},
		{Name: "read_file", Description: "Read a bounded line range from a repository file.", InputSchema: object(map[string]any{
			"path": stringProperty("Repository-relative path"), "start_line": integerProperty("First line", 1), "end_line": integerProperty("Last line", 1),
		}, "path")},
		{Name: "search_repository", Description: "Search repository text using a regular expression.", InputSchema: object(map[string]any{
			"query": stringProperty("Regular expression"), "path": stringProperty("Repository-relative file or directory"),
		}, "query")},
		{Name: "read_original_patch", Description: "Read a bounded line range from the original failed mail patch.", InputSchema: object(map[string]any{
			"start_line": integerProperty("First line", 1), "end_line": integerProperty("Last line", 1),
		})},
		{Name: "replace_text", Description: "Replace one exact text occurrence in an allowed file.", InputSchema: object(map[string]any{
			"path": stringProperty("Repository-relative path"), "old_text": stringProperty("Exact existing text"), "new_text": stringProperty("Replacement text"),
		}, "path", "old_text", "new_text")},
		{Name: "replace_lines", Description: "Replace an inclusive line range in an allowed file.", InputSchema: object(map[string]any{
			"path": stringProperty("Repository-relative path"), "start_line": integerProperty("First line", 1), "end_line": integerProperty("Last line", 1), "new_text": stringProperty("Replacement text"),
		}, "path", "start_line", "end_line", "new_text")},
		{Name: "write_file", Description: "Write complete UTF-8 content to an allowed file.", InputSchema: object(map[string]any{
			"path": stringProperty("Repository-relative path"), "content": stringProperty("Complete file content"),
		}, "path", "content")},
		{Name: "delete_file", Description: "Delete an allowed file.", InputSchema: object(map[string]any{
			"path": stringProperty("Repository-relative path"),
		}, "path")},
		{Name: "show_diff", Description: "Show current uncommitted changes.", InputSchema: object(map[string]any{})},
	}
}

func toolNames() []string {
	tools := patchTools()
	names := make([]string, 0, len(tools))
	for _, tool := range tools {
		names = append(names, tool.Name)
	}
	return names
}

func requestedToolNames(message modelMessage) []string {
	var names []string
	for _, block := range message.Content {
		if block.ToolUse != nil {
			names = append(names, block.ToolUse.Name)
		}
	}
	return names
}

func responseText(message modelMessage) string {
	var text []string
	for _, block := range message.Content {
		if block.Text != "" {
			text = append(text, block.Text)
		}
	}
	return strings.Join(text, "\n")
}

func requiredStringInput(input map[string]any, key string) string {
	value, _ := input[key].(string)
	return value
}

func stringInput(input map[string]any, key, fallback string) string {
	if value := requiredStringInput(input, key); value != "" {
		return value
	}
	return fallback
}

func intInput(input map[string]any, key string, fallback int) int {
	switch value := input[key].(type) {
	case int:
		return value
	case int32:
		return int(value)
	case int64:
		return int(value)
	case float64:
		return int(value)
	case json.Number:
		result, err := value.Int64()
		if err == nil {
			return int(result)
		}
	}
	return fallback
}

func appendTrace(path string, event map[string]any) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	if err := json.NewEncoder(writer).Encode(map[string]any{"event": event}); err != nil {
		return err
	}
	return writer.Flush()
}
