package agentpatchfixer

import (
	"context"
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

type agentInvocation struct {
	SchemaVersion  int      `json:"schema_version"`
	WorkspacePath  string   `json:"workspace_path"`
	AllowedFiles   []string `json:"allowed_files"`
	OriginalPatch  string   `json:"original_patch"`
	SystemPrompt   string   `json:"system_prompt"`
	Prompt         string   `json:"prompt"`
	ModelID        string   `json:"model_id"`
	Region         string   `json:"region"`
	MaxTurns       int      `json:"max_turns"`
	MaxTotalTokens int      `json:"max_total_tokens"`
	Diagnostics    bool     `json:"diagnostics"`
}

type agentOutcome struct {
	SchemaVersion int    `json:"schema_version"`
	Status        string `json:"status"`
	StopReason    string `json:"stop_reason"`
	Summary       string `json:"summary"`
	Turns         int    `json:"turns"`
	Tokens        int    `json:"total_tokens"`
	Message       any    `json:"message,omitempty"`
}

func execute(ctx context.Context, request Request, runDir string) (*Result, error) {
	patchContent, err := os.ReadFile(request.FailedPatchPath)
	if err != nil {
		return nil, fmt.Errorf("reading failed patch: %v", err)
	}
	patchText := string(patchContent)
	metadata, err := parsePatchMetadata(ctx, patchText)
	if err != nil {
		return nil, err
	}
	allowedFiles, err := patchFiles(ctx, patchText)
	if err != nil {
		return nil, err
	}
	for _, editableFile := range request.EditableFiles {
		allowedFiles[filepath.ToSlash(filepath.Clean(editableFile))] = true
	}

	workspaceParent, err := os.MkdirTemp("", "patch-fixer-workspace-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(workspaceParent)

	workspacePath, err := prepareWorkspace(ctx, request.UpstreamRepoPath, workspaceParent)
	if err != nil {
		return nil, fmt.Errorf("preparing patch workspace: %v", err)
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
	prompt := buildPrompt(request, metadata, allowedFiles, seed)
	invocation := agentInvocation{
		SchemaVersion:  requestSchemaVersion,
		WorkspacePath:  workspacePath,
		AllowedFiles:   sortedFiles(allowedFiles),
		OriginalPatch:  patchText,
		SystemPrompt:   systemPrompt,
		Prompt:         prompt,
		ModelID:        modelID(),
		Region:         region,
		MaxTurns:       request.MaxTurns,
		MaxTotalTokens: request.MaxTotalTokens,
		Diagnostics:    diagnosticsFull(),
	}

	if diagnosticsFull() {
		if err := writeJSON(filepath.Join(runDir, "model-request.json"), map[string]any{
			"schema_version": requestSchemaVersion,
			"model_id":       invocation.ModelID,
			"system_prompt":  systemPrompt,
			"prompt":         prompt,
			"limits": map[string]int{
				"turns":        request.MaxTurns,
				"total_tokens": request.MaxTotalTokens,
			},
			"reject_files":  seed.RejectFiles,
			"missing_files": seed.MissingFiles,
			"tools":         toolNames(),
		}); err != nil {
			return nil, err
		}
	}

	outcome, err := runAgentProcess(ctx, invocation, runDir)
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
		diff, diffErr := showWorkspaceDiff(ctx, workspacePath)
		if diffErr == nil {
			_ = os.WriteFile(filepath.Join(runDir, "incomplete.diff"), []byte(diff), 0o600)
		}
	} else {
		candidatePath := filepath.Join(runDir, "candidate.patch")
		editedFiles, err := generateCandidate(ctx, workspacePath, metadata, candidatePath, allowedFiles)
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
			"message":              outcome.Message,
		}); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func buildPrompt(
	request Request,
	metadata patchMetadata,
	allowedFiles map[string]bool,
	seed patchSeed,
) string {
	rejects := strings.Join(seed.RejectFiles, ", ")
	if rejects == "" {
		rejects = "(none)"
	}
	missingFiles := strings.Join(seed.MissingFiles, ", ")
	if missingFiles == "" {
		missingFiles = "(none)"
	}
	return fmt.Sprintf(
		"Goal: repair this failed patch without changing unrelated business behavior.\n"+
			"Project: %s\nOld revision: %s\nNew revision: %s\nFailed patch: %s\n"+
			"Allowed files: %s\n"+
			"Original author: %s <%s>\nOriginal subject: %s\nPatch application error:\n%s\n\n"+
			"Git already applied every clean hunk from the original patch to this isolated checkout.\n"+
			"Reject files that still need resolution: %s\n"+
			"Patch targets missing from the current checkout: %s\n"+
			"Read the reject files and relevant current file ranges first. For missing targets, search for the "+
			"current location or determine whether the original intent is obsolete. Preserve the pre-applied changes "+
			"and resolve only the rejected intent. For large rejected sections, use replace_lines in small "+
			"ranges instead of embedding entire old and new files in replace_text. Leave the checkout unchanged "+
			"if the target already satisfies or supersedes the intent. Review the final diff once and finish.",
		request.ProjectName,
		request.CurrentRevision,
		request.TargetRevision,
		filepath.Base(request.FailedPatchPath),
		strings.Join(sortedFiles(allowedFiles), ", "),
		metadata.AuthorName,
		metadata.AuthorEmail,
		metadata.Subject,
		boundedTail(request.FailureOutput, 12000),
		rejects,
		missingFiles,
	)
}

func sortedFiles(files map[string]bool) []string {
	result := make([]string, 0, len(files))
	for file := range files {
		result = append(result, file)
	}
	sort.Strings(result)
	return result
}

func toolNames() []string {
	return []string{
		"list_files",
		"read_file",
		"search_repository",
		"read_original_patch",
		"replace_text",
		"replace_lines",
		"write_file",
		"delete_file",
		"show_diff",
	}
}
