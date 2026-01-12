package fixpatches

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/types"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

// BedrockResponse represents the response from Bedrock API.
type BedrockResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// convertToInferenceProfile converts a model ID to an inference profile ID if needed.
// Claude Sonnet 4.5 and newer models require using inference profiles instead of direct model IDs.
// Inference profiles provide cross-region routing and better availability.
func convertToInferenceProfile(modelID string, region string) string {
	// Map of model IDs that require inference profiles
	// Format: model-id -> inference-profile-id
	// Note: Inference profile IDs keep the full date-based version, just add "us." or "global." prefix
	inferenceProfileMap := map[string]string{
		"anthropic.claude-sonnet-4-5-20250929-v1:0": "us.anthropic.claude-sonnet-4-5-20250929-v1:0",
		"anthropic.claude-3-7-sonnet-20250219-v1:0": "us.anthropic.claude-3-7-sonnet-20250219-v1:0", // 1M tokens/min default!
		"anthropic.claude-3-5-sonnet-20241022-v2:0": "us.anthropic.claude-3-5-sonnet-20241022-v2:0",
		"anthropic.claude-sonnet-4-20250514-v1:0":   "us.anthropic.claude-sonnet-4-20250514-v1:0",
		"anthropic.claude-opus-4-20250514-v1:0":     "us.anthropic.claude-opus-4-20250514-v1:0",
		"anthropic.claude-opus-4-1-20250805-v1:0":   "us.anthropic.claude-opus-4-1-20250805-v1:0",
		"anthropic.claude-3-5-haiku-20241022-v1:0":  "us.anthropic.claude-3-5-haiku-20241022-v1:0",
	}

	// Check if this model needs an inference profile
	if profileID, needsProfile := inferenceProfileMap[modelID]; needsProfile {
		return profileID
	}

	// For older models (Claude 3.0, 3.5 v1) that work with direct model IDs, return as-is
	return modelID
}

// Global client to reuse across calls (avoids recreating client on every retry)
var globalBedrockClient *bedrockruntime.Client
var globalModelOrProfile string
var lastRequestTime time.Time
var requestMutex sync.Mutex

// initBedrockClient initializes the Bedrock client once and reuses it.
func initBedrockClient(model string) (*bedrockruntime.Client, string, error) {
	// Force us-west-2 region for Bedrock API endpoint - matches CodeBuild region.
	// Note: Cross-region inference (CRIS) may route requests to other regions,
	// but IAM policy allows all regions for foundation models.
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-west-2"),
		config.WithRetryMaxAttempts(3),
		config.WithRetryMode(aws.RetryModeAdaptive),
	)
	if err != nil {
		return nil, "", fmt.Errorf("loading AWS config: %v", err)
	}

	modelOrProfile := convertToInferenceProfile(model, cfg.Region)

	// Reuse client if model hasn't changed
	if globalBedrockClient != nil && globalModelOrProfile == modelOrProfile {
		return globalBedrockClient, globalModelOrProfile, nil
	}

	// Model changed or first initialization
	logger.Info("Initializing Bedrock client", "model", model, "profile", modelOrProfile, "region", cfg.Region)

	// Create new client
	globalBedrockClient = bedrockruntime.NewFromConfig(cfg)
	globalModelOrProfile = modelOrProfile

	return globalBedrockClient, globalModelOrProfile, nil
}

// waitForRateLimit ensures we don't exceed Bedrock's rate limits.
// Bedrock has a 4 requests/min limit for cross-region inference profiles.
// This means we need at least 15 seconds between requests.
func waitForRateLimit() {
	requestMutex.Lock()
	defer requestMutex.Unlock()

	// Calculate time since last request
	timeSinceLastRequest := time.Since(lastRequestTime)

	// Bedrock limit for Claude 3.7 Sonnet: 250 requests/min
	// We use 15s between requests for safety (4 req/min) to avoid any rate limit issues
	minTimeBetweenRequests := 15 * time.Second

	if timeSinceLastRequest < minTimeBetweenRequests {
		waitTime := minTimeBetweenRequests - timeSinceLastRequest
		logger.Info("Rate limiting: waiting to respect Bedrock limits",
			"wait_seconds", waitTime.Seconds(),
			"time_since_last_request", timeSinceLastRequest.Seconds())
		time.Sleep(waitTime)
	}

	// Update last request time
	lastRequestTime = time.Now()
}

// CallBedrockForPatchFix invokes Bedrock with patch context to generate a fix.
func CallBedrockForPatchFix(ctx *types.PatchContext, model string, attempt int) (*types.PatchFix, error) {
	logger.Info("Calling Bedrock API", "model", model, "attempt", attempt)

	// Initialize or reuse existing client
	client, modelOrProfile, err := initBedrockClient(model)
	if err != nil {
		return nil, err
	}

	logger.Info("Initialized Bedrock client", "model", model, "profile", modelOrProfile, "region", "us-west-2")

	// Build the prompt
	prompt := BuildPrompt(ctx, attempt)

	// Estimate input tokens (use conservative 3 chars/token for code)
	estimatedInputTokens := len(prompt) / 3
	logger.Info("Prompt built", "length", len(prompt), "estimated_tokens", estimatedInputTokens)

	// Safety check: ensure we're not exceeding Claude's 200K input limit
	if estimatedInputTokens > 200000 {
		return nil, fmt.Errorf("prompt too large: %d estimated tokens (limit: 200K) - this should have been caught by context pruning", estimatedInputTokens)
	}

	// Write prompt to debug file/S3 for inspection
	promptDebugFile := fmt.Sprintf("/tmp/llm-prompt-attempt-%d.txt", attempt)
	if err := writeDebugFile(promptDebugFile, []byte(prompt), "prompt", attempt); err != nil {
		logger.Info("Warning: failed to write prompt debug file", "error", err)
	}

	// Construct Bedrock request for Claude
	systemPrompt := `You are an expert at resolving Git patch conflicts. Your task is to fix failed patch hunks by analyzing the original intent and the current code state.

Rules:
1. Preserve the original patch intent exactly
2. Preserve the original patch metadata (From, Date, Subject) exactly
3. Only modify the diff content to resolve the conflict
4. Maintain code style and formatting
5. Output ONLY the corrected patch in unified diff format with complete headers
6. Do not add explanations or commentary`

	// Calculate max_tokens based on patch size
	// Use patch size as proxy: larger patches need more output tokens
	// Conservative estimate: patch size in chars / 3 * 2 (for output expansion)
	patchSize := len(ctx.OriginalPatch)
	maxTokens := (patchSize / 3) * 2

	// Clamp to reasonable bounds
	// With extended output feature enabled, we can use up to 128K tokens
	if maxTokens < 8192 {
		maxTokens = 8192 // Minimum for any patch
	}
	if maxTokens > 100000 {
		maxTokens = 100000 // Stay well under 128K limit for safety
	}

	logger.Info("Calculated max_tokens for patch",
		"patch_size_bytes", patchSize,
		"max_tokens", maxTokens)

	requestBody := map[string]any{
		"anthropic_version": "bedrock-2023-05-31",
		"max_tokens":        maxTokens, // Dynamic based on patch size
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"system": systemPrompt,
		// Enable extended output feature for Claude models
		// This allows up to 128K output tokens instead of the default 8K limit
		"anthropic_beta": []string{"output-128k-2025-02-19"},
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshaling request body: %v", err)
	}

	// Invoke model with retry logic and exponential backoff
	// Bedrock rate limits for Claude Sonnet 4.5 (cross-region inference profile):
	// - Requests per minute: 4 (L-4A6BFAB1)
	// - Tokens per minute: 4,000 (L-F4DDD3EB)
	// - Max tokens per day: 144M (L-381AD9EE)
	//
	// With 4 requests/min, we need at least 15 seconds between requests (60s / 4 = 15s)
	// To be safe and account for clock skew, we use 20s as the minimum wait time
	var response *bedrockruntime.InvokeModelOutput
	maxRetries := 5 // Give multiple chances with proper backoff

	for i := 0; i < maxRetries; i++ {
		// Log the attempt
		if i > 0 {
			logger.Info("Retrying Bedrock API call", "attempt", i+1, "max_retries", maxRetries)
		}

		// CRITICAL: Wait for rate limit before making request
		// This ensures we never exceed 4 requests/min
		waitForRateLimit()

		response, err = client.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
			ModelId:     aws.String(modelOrProfile),
			ContentType: aws.String("application/json"),
			Body:        requestBodyBytes,
		})

		if err == nil {
			logger.Info("Bedrock API call succeeded", "attempt", i+1)
			break
		}

		// Log the error
		logger.Info("Bedrock API call failed", "attempt", i+1, "max_retries", maxRetries, "error", err.Error())

		if i < maxRetries-1 {
			// Exponential backoff starting at 20s to respect 4 requests/min limit
			// Wait times: 20s, 40s, 80s, 160s
			// This ensures we stay well under the 4 requests/min limit (15s minimum)
			waitTime := time.Duration(20*(1<<uint(i))) * time.Second
			logger.Info("Waiting before retry to respect rate limits",
				"wait_seconds", waitTime.Seconds(),
				"rate_limit", "4 requests/min for Claude Sonnet 4.5")
			time.Sleep(waitTime)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("invoking Bedrock after %d retries: %v", maxRetries, err)
	}

	// Parse response
	var result BedrockResponse
	if err := json.Unmarshal(response.Body, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling Bedrock response: %v", err)
	}

	if len(result.Content) == 0 {
		return nil, fmt.Errorf("empty response from Bedrock")
	}

	responseText := result.Content[0].Text

	// CRITICAL: Log actual token usage from Bedrock API
	// This shows the REAL token count, not our estimate
	logger.Info("✅ Bedrock API Success - Actual Token Usage",
		"input_tokens_actual", result.Usage.InputTokens,
		"output_tokens_actual", result.Usage.OutputTokens,
		"total_tokens", result.Usage.InputTokens+result.Usage.OutputTokens,
		"response_length_bytes", len(responseText),
		"model_limit", "1M input tokens (if quota enabled)")

	// Compare with our estimate
	tokenEstimateAccuracy := float64(result.Usage.InputTokens) / float64(estimatedInputTokens) * 100
	logger.Info("Token Estimation Accuracy",
		"estimated_input_tokens", estimatedInputTokens,
		"actual_input_tokens", result.Usage.InputTokens,
		"accuracy_percent", fmt.Sprintf("%.1f%%", tokenEstimateAccuracy))

	// Write response to debug file/S3 for inspection
	responseDebugFile := fmt.Sprintf("/tmp/llm-response-attempt-%d.txt", attempt)
	if err := writeDebugFile(responseDebugFile, []byte(responseText), "response", attempt); err != nil {
		logger.Info("Warning: failed to write response debug file", "error", err)
	}

	// Check if response was truncated
	if result.Usage.OutputTokens >= maxTokens {
		logger.Info("Response truncated: hit max_tokens limit",
			"output_tokens", result.Usage.OutputTokens,
			"max_tokens", maxTokens)
		return nil, fmt.Errorf("LLM response truncated at %d tokens (limit: %d) - patch output too large, consider reducing input context",
			result.Usage.OutputTokens, maxTokens)
	}

	// Extract patch from response
	patch := extractPatchFromResponse(responseText)
	if patch == "" {
		return nil, fmt.Errorf("no patch found in Bedrock response")
	}

	// Validate patch format and metadata
	if err := validatePatchFormat(patch, ctx); err != nil {
		return nil, fmt.Errorf("invalid patch format: %v", err)
	}

	logger.Info("Bedrock API call complete",
		"input_tokens", result.Usage.InputTokens,
		"output_tokens", result.Usage.OutputTokens)

	return &types.PatchFix{
		Patch:      patch,
		TokensUsed: result.Usage.InputTokens + result.Usage.OutputTokens,
		Cost:       0, // Cost tracking removed
	}, nil
}

// extractPatchFromResponse extracts the patch content from LLM response.
// The LLM might wrap the patch in markdown code blocks or add explanations.
func extractPatchFromResponse(response string) string {
	// Look for patch content between ```diff or ``` markers
	if strings.Contains(response, "```") {
		// Extract content between code blocks
		parts := strings.Split(response, "```")
		for i, part := range parts {
			// Skip the first part (before first ```)
			if i == 0 {
				continue
			}
			// Remove language identifier if present (e.g., "diff\n")
			part = strings.TrimPrefix(part, "diff\n")
			part = strings.TrimPrefix(part, "diff ")
			part = strings.TrimSpace(part)

			// Check if this looks like a patch (starts with From or diff)
			if strings.HasPrefix(part, "From ") || strings.HasPrefix(part, "diff --git") {
				return part
			}
		}
	}

	// If no code blocks, look for patch markers
	lines := strings.Split(response, "\n")
	var patchLines []string
	inPatch := false

	for _, line := range lines {
		// Start of patch
		if strings.HasPrefix(line, "From ") || strings.HasPrefix(line, "diff --git") {
			inPatch = true
		}

		if inPatch {
			patchLines = append(patchLines, line)
		}
	}

	if len(patchLines) > 0 {
		return strings.Join(patchLines, "\n")
	}

	// Fallback: return the whole response if it looks like a patch
	if strings.Contains(response, "diff --git") || strings.Contains(response, "From ") {
		return strings.TrimSpace(response)
	}

	return ""
}

// BuildPrompt constructs the LLM prompt from context.
func BuildPrompt(ctx *types.PatchContext, attempt int) string {
	var prompt strings.Builder

	// Project information
	prompt.WriteString(fmt.Sprintf("## Project: %s\n\n", ctx.ProjectName))

	// Original patch metadata
	prompt.WriteString("## Original Patch Metadata\n")
	if ctx.PatchAuthor != "" {
		prompt.WriteString(fmt.Sprintf("From: %s\n", ctx.PatchAuthor))
	}
	if ctx.PatchDate != "" {
		prompt.WriteString(fmt.Sprintf("Date: %s\n", ctx.PatchDate))
	}
	if ctx.PatchSubject != "" {
		prompt.WriteString(fmt.Sprintf("Subject: %s\n", ctx.PatchSubject))
	}
	prompt.WriteString("\n")

	// Original patch intent
	if ctx.PatchIntent != "" {
		prompt.WriteString("## Original Patch Intent\n")
		prompt.WriteString(fmt.Sprintf("%s\n\n", ctx.PatchIntent))
	}

	// Failed hunks
	for i, hunk := range ctx.FailedHunks {
		prompt.WriteString(fmt.Sprintf("## Failed Hunk #%d in %s\n\n", hunk.HunkIndex, hunk.FilePath))

		// What the patch tried to do
		prompt.WriteString("### What the patch tried to do:\n")
		prompt.WriteString("```diff\n")
		for _, line := range hunk.OriginalLines {
			prompt.WriteString(line + "\n")
		}
		prompt.WriteString("```\n\n")

		// Current state of the file (context around conflict)
		prompt.WriteString(fmt.Sprintf("### Current file state (around line %d):\n", hunk.LineNumber))
		prompt.WriteString("```\n")
		prompt.WriteString(hunk.Context)
		prompt.WriteString("\n```\n\n")

		// Add separator between hunks
		if i < len(ctx.FailedHunks)-1 {
			prompt.WriteString("---\n\n")
		}
	}

	// Show context for files that applied with offset or cleanly
	// (Failed files already have detailed context in the hunks section above)
	if len(ctx.AllFileContexts) > 0 {
		hasNonFailedFiles := false

		// Check if there are any non-failed files to show
		for filename := range ctx.AllFileContexts {
			hasFailed := false
			for _, hunk := range ctx.FailedHunks {
				if strings.Contains(hunk.FilePath, filename) {
					hasFailed = true
					break
				}
			}
			if !hasFailed {
				hasNonFailedFiles = true
				break
			}
		}

		if hasNonFailedFiles {
			prompt.WriteString("## Files That Applied Successfully\n\n")
			prompt.WriteString("These files applied but may need line number adjustments:\n\n")

			for filename, context := range ctx.AllFileContexts {
				// Check if this file has a .rej (failed)
				hasFailed := false
				for _, hunk := range ctx.FailedHunks {
					if strings.Contains(hunk.FilePath, filename) {
						hasFailed = true
						break
					}
				}

				// SKIP failed files - we already showed their context in hunks above
				if hasFailed {
					continue
				}

				prompt.WriteString(fmt.Sprintf("### %s\n\n", filename))

				// Check if this file has offset
				hasOffset := false
				offsetAmount := 0
				if ctx.ApplicationResult != nil {
					if offset, ok := ctx.ApplicationResult.OffsetFiles[filename]; ok {
						hasOffset = true
						offsetAmount = offset
					}
				}

				// Show status
				if hasOffset {
					prompt.WriteString(fmt.Sprintf("Applied with offset: +%d lines\n", offsetAmount))
					prompt.WriteString("Include this file in your fixed patch with updated line numbers.\n\n")
				} else {
					prompt.WriteString("Applied cleanly - include in fixed patch as-is.\n\n")
				}

				// Show pristine content (BEFORE git apply modified it)
				prompt.WriteString("```\n")
				prompt.WriteString(context)
				prompt.WriteString("\n```\n\n")
			}
		}
	}

	// Original patch for reference (with warnings about what needs fixing)
	// NOTE: We put this BEFORE the error so the error is closer to the task
	prompt.WriteString("## Original Patch (For Reference)\n\n")

	// Identify which files failed and which succeeded with offset
	failedFiles := make(map[string]bool)
	for _, hunk := range ctx.FailedHunks {
		failedFiles[filepath.Base(hunk.FilePath)] = true
	}

	// Show status of each file
	if ctx.ApplicationResult != nil {
		prompt.WriteString("**Patch Application Status:**\n")

		// List failed files
		if len(failedFiles) > 0 {
			failedList := make([]string, 0, len(failedFiles))
			for file := range failedFiles {
				failedList = append(failedList, file)
			}
			prompt.WriteString(fmt.Sprintf("- ❌ FAILED (needs fixing): %s\n", strings.Join(failedList, ", ")))
		}

		// List offset files
		if len(ctx.ApplicationResult.OffsetFiles) > 0 {
			for file, offset := range ctx.ApplicationResult.OffsetFiles {
				prompt.WriteString(fmt.Sprintf("- ⚠️  APPLIED WITH OFFSET (needs line number update): %s (offset: %d lines)\n", file, offset))
			}
		}
		prompt.WriteString("\n")
	}

	// For attempt 1: include full original patch
	// For attempt 2+: include only failed file portions to save tokens
	if attempt == 1 {
		prompt.WriteString("**Full Original Patch:**\n")
		prompt.WriteString("```diff\n")
		prompt.WriteString(ctx.OriginalPatch)
		prompt.WriteString("\n```\n\n")
	} else {
		// Extract only the failed files from the original patch
		prompt.WriteString("**Original Patch (Failed Files Only):**\n")
		prompt.WriteString("```diff\n")

		// Get list of failed files
		failedFileNames := make(map[string]bool)
		for _, hunk := range ctx.FailedHunks {
			fileName := filepath.Base(hunk.FilePath)
			failedFileNames[fileName] = true
		}

		// Extract diffs for failed files only
		if len(failedFileNames) > 0 {
			failedDiffs := extractFileDiffsFromPatch(ctx.OriginalPatch, failedFileNames)
			if failedDiffs != "" {
				prompt.WriteString(failedDiffs)
			} else {
				// Fallback: include full patch if extraction fails
				prompt.WriteString(ctx.OriginalPatch)
			}
		}

		prompt.WriteString("\n```\n\n")

		// Count total files vs failed files
		totalFiles := strings.Count(ctx.OriginalPatch, "diff --git")
		prompt.WriteString(fmt.Sprintf("ℹ️  Note: Showing only %d failed file(s). %d other files applied successfully.\n\n",
			len(failedFileNames), totalFiles-len(failedFileNames)))
	}

	prompt.WriteString("⚠️  **Important**: This patch was created against an OLD version of the code.\n")
	prompt.WriteString("Some files may have changed (version bumps, line shifts, etc.).\n")
	prompt.WriteString("Use the 'Expected vs Actual' sections above to see what changed.\n\n")

	// CRITICAL: Put error information HERE, right before the task
	// This ensures the LLM sees the error immediately before generating the fix
	if attempt > 1 && ctx.BuildError != "" {
		prompt.WriteString("---\n\n")
		prompt.WriteString(fmt.Sprintf("# 🚨 CRITICAL: Attempt #%d Failed With This Error\n\n", attempt-1))

		prompt.WriteString("**Your previous patch failed to apply with this error:**\n")
		prompt.WriteString("```\n")
		// Limit build error to last 500 lines
		errorLines := strings.Split(ctx.BuildError, "\n")
		if len(errorLines) > 500 {
			prompt.WriteString("...(truncated)...\n")
			errorLines = errorLines[len(errorLines)-500:]
		}
		prompt.WriteString(strings.Join(errorLines, "\n"))
		prompt.WriteString("\n```\n\n")

		prompt.WriteString("**🎯 Your primary goal:**\n")
		prompt.WriteString("Fix the SPECIFIC error shown above. The error message tells you:\n")
		prompt.WriteString("- Which line number failed (e.g., 'corrupt patch at line 276')\n")
		prompt.WriteString("- What type of error occurred (corrupt patch, missing header, etc.)\n")
		prompt.WriteString("- This is the MOST IMPORTANT context for your fix\n\n")

		prompt.WriteString("**Common causes of these errors:**\n")
		prompt.WriteString("- 'corrupt patch at line X': Patch format is malformed (missing newlines, truncated content)\n")
		prompt.WriteString("- 'patch fragment without header': Missing 'diff --git' or file headers\n")
		prompt.WriteString("- 'does not apply': Line numbers or content don't match current file\n\n")

		prompt.WriteString("---\n\n")
	}

	// Reflection for later attempts
	if attempt >= 3 {
		prompt.WriteString(fmt.Sprintf("## 🤔 Reflection Required (Attempt #%d)\n\n", attempt))
		prompt.WriteString(fmt.Sprintf("This is your %s attempt. Before providing the fix, first explain:\n", ordinal(attempt)))
		prompt.WriteString("1. What SPECIFIC error occurred (see the error above)\n")
		prompt.WriteString("2. Why that error happened (patch format issue? line number mismatch?)\n")
		prompt.WriteString("3. What SPECIFIC changes you'll make to fix it\n\n")
		prompt.WriteString("Then provide the corrected patch.\n\n")
	}

	// Task instructions
	prompt.WriteString("## Task\n")
	prompt.WriteString("Generate a corrected patch that:\n")
	prompt.WriteString("1. Preserves the exact metadata (From, Date, Subject) from the original patch\n")
	prompt.WriteString("2. Includes ALL files from the original patch (both failed and offset files)\n")
	prompt.WriteString("3. For FAILED files: Fix them using the 'Expected vs Actual' context above\n")
	prompt.WriteString("4. For OFFSET files: Update line numbers to match current file state\n")
	prompt.WriteString("5. Uses RELATIVE file paths NOT absolute paths\n")
	prompt.WriteString("6. Will compile successfully\n\n")

	prompt.WriteString("## How to Generate the Fix\n\n")

	prompt.WriteString("**Step 1: Understand the Intent**\n")
	prompt.WriteString("Look at 'What the patch tried to do' to understand the semantic change being made.\n\n")

	prompt.WriteString("**Step 2: Use Current File State**\n")
	prompt.WriteString("The 'Expected vs Actual' sections show you:\n")
	prompt.WriteString("- What the original patch expected (OLD version)\n")
	prompt.WriteString("- What's actually in the file NOW (NEW version)\n")
	prompt.WriteString("- The specific differences between them\n\n")
	prompt.WriteString("You MUST use the ACTUAL CURRENT content as your starting point, not the expected content.\n\n")

	prompt.WriteString("**Step 3: Find the Semantic Location**\n")
	prompt.WriteString("Find where in the CURRENT file the change should be applied:\n")
	prompt.WriteString("- Use the 'Current file content' section to see the broader context\n")
	prompt.WriteString("- Match based on semantic meaning (package names, function names, etc.)\n")
	prompt.WriteString("- Don't rely on line numbers from the original patch - they may have shifted\n\n")

	prompt.WriteString("**Step 4: Generate the Patch**\n")
	prompt.WriteString("Create a patch that:\n")
	prompt.WriteString("- Uses context lines from the CURRENT file (complete, not truncated)\n")
	prompt.WriteString("- Uses CURRENT line numbers\n")
	prompt.WriteString("- Makes the SAME semantic change as the original patch intended\n")
	prompt.WriteString("- Preserves exact formatting and whitespace from the current file\n")
	prompt.WriteString("- **CRITICAL**: Match the EXACT indentation (tabs/spaces) from the context lines shown above\n")
	prompt.WriteString("- **CRITICAL**: DO NOT copy context lines from the original patch - they are from an OLD version\n")
	prompt.WriteString("- **CRITICAL**: Use ONLY the context lines from the 'Current file state' section above\n\n")

	prompt.WriteString("Output format (unified diff with complete headers):\n")
	prompt.WriteString("```\n")
	prompt.WriteString("From <commit-hash> Mon Sep 17 00:00:00 2001\n")
	if ctx.PatchAuthor != "" {
		prompt.WriteString(fmt.Sprintf("From: %s\n", ctx.PatchAuthor))
	}
	if ctx.PatchDate != "" {
		prompt.WriteString(fmt.Sprintf("Date: %s\n", ctx.PatchDate))
	}
	if ctx.PatchSubject != "" {
		prompt.WriteString(fmt.Sprintf("Subject: %s\n", ctx.PatchSubject))
	}
	prompt.WriteString("\n---\n")
	prompt.WriteString(" file1.ext | X +/-\n")
	prompt.WriteString(" file2.ext | Y +/-\n")
	prompt.WriteString(" N files changed, X insertions(+), Y deletions(-)\n\n")
	prompt.WriteString("diff --git a/file1.ext b/file1.ext\n")
	prompt.WriteString("...\n")
	prompt.WriteString("```\n")

	return prompt.String()
}

// validatePatchFormat validates that the patch has required metadata and format.
func validatePatchFormat(patch string, ctx *types.PatchContext) error {
	// Check for required patch headers
	if !strings.Contains(patch, "From ") && !strings.Contains(patch, "diff --git") {
		return fmt.Errorf("patch missing required headers (From or diff --git)")
	}

	// Validate patch metadata is preserved (if original had it)
	if ctx.PatchAuthor != "" {
		if !strings.Contains(patch, ctx.PatchAuthor) {
			logger.Info("Warning: patch author not preserved in LLM output",
				"expected", ctx.PatchAuthor)
			// Don't fail - this is a warning, not a hard error
		}
	}

	if ctx.PatchDate != "" {
		if !strings.Contains(patch, ctx.PatchDate) {
			logger.Info("Warning: patch date not preserved in LLM output",
				"expected", ctx.PatchDate)
		}
	}

	if ctx.PatchSubject != "" {
		// Check if subject is preserved (might be slightly reformatted)
		subjectCore := strings.TrimPrefix(ctx.PatchSubject, "[PATCH]")
		subjectCore = strings.TrimSpace(subjectCore)
		if !strings.Contains(patch, subjectCore) {
			logger.Info("Warning: patch subject not preserved in LLM output",
				"expected", subjectCore)
		}
	}

	// Check for diff content
	if !strings.Contains(patch, "@@") {
		return fmt.Errorf("patch missing diff hunks (no @@ markers found)")
	}

	// Check for basic diff structure
	hasMinus := strings.Contains(patch, "---")
	hasPlus := strings.Contains(patch, "+++")
	if !hasMinus || !hasPlus {
		return fmt.Errorf("patch missing file markers (--- or +++)")
	}

	logger.Info("Patch format validation passed")
	return nil
}

// ordinal returns the ordinal string for a number (1st, 2nd, 3rd, etc.)
func ordinal(n int) string {
	suffix := "th"
	switch n % 10 {
	case 1:
		if n%100 != 11 {
			suffix = "st"
		}
	case 2:
		if n%100 != 12 {
			suffix = "nd"
		}
	case 3:
		if n%100 != 13 {
			suffix = "rd"
		}
	}
	return fmt.Sprintf("%d%s", n, suffix)
}

// extractFileDiffsFromPatch extracts only the diffs for specified files from a patch.
// This is used to reduce token usage in retry attempts by only including failed files.
func extractFileDiffsFromPatch(patch string, fileNames map[string]bool) string {
	if len(fileNames) == 0 {
		return ""
	}

	var result strings.Builder
	lines := strings.Split(patch, "\n")

	inTargetFile := false
	currentFileName := ""
	var currentFileDiff strings.Builder

	for i, line := range lines {
		// Check for new file diff
		if strings.HasPrefix(line, "diff --git") {
			// Save previous file if it was a target
			if inTargetFile && currentFileDiff.Len() > 0 {
				result.WriteString(currentFileDiff.String())
				result.WriteString("\n")
			}

			// Reset for new file
			currentFileDiff.Reset()
			inTargetFile = false

			// Extract filename from "diff --git a/path/to/file.go b/path/to/file.go"
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				// Get the b/ path (destination)
				filePath := strings.TrimPrefix(parts[3], "b/")
				currentFileName = filepath.Base(filePath)

				// Check if this is a file we want
				if fileNames[currentFileName] {
					inTargetFile = true
					currentFileDiff.WriteString(line)
					currentFileDiff.WriteString("\n")
				}
			}
		} else if inTargetFile {
			// Include all lines for target files
			currentFileDiff.WriteString(line)
			if i < len(lines)-1 {
				currentFileDiff.WriteString("\n")
			}
		}
	}

	// Don't forget the last file
	if inTargetFile && currentFileDiff.Len() > 0 {
		result.WriteString(currentFileDiff.String())
	}

	return result.String()
}

// writeDebugFile writes debug content to local file and optionally uploads to S3.
// If DEBUG_BUCKET env var is set, uploads to S3 with path: {project}/{pr}/{type}-attempt-{attempt}.txt
// Always writes to local /tmp for immediate debugging.
func writeDebugFile(localPath string, content []byte, fileType string, attempt int) error {
	// Always write to local file first
	if err := os.WriteFile(localPath, content, 0644); err != nil {
		return fmt.Errorf("writing local file: %v", err)
	}
	logger.Info("Wrote debug file locally", "file", localPath)

	// Check if S3 upload is configured
	debugBucket := os.Getenv("DEBUG_BUCKET")
	if debugBucket == "" {
		return nil // S3 upload not configured, that's fine
	}

	// Get project and PR from environment (set by Lambda handler)
	project := os.Getenv("PROJECT_NAME")
	pr := os.Getenv("PR_NUMBER")
	if project == "" || pr == "" {
		logger.Info("Skipping S3 upload: PROJECT_NAME or PR_NUMBER not set")
		return nil
	}

	// Upload to S3
	ctx := context.Background()

	// Allow specifying bucket region via env var to avoid redirect issues
	// If not set, uses default region from AWS config
	bucketRegion := os.Getenv("DEBUG_BUCKET_REGION")
	var cfg aws.Config
	var err error

	if bucketRegion != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(bucketRegion))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}

	if err != nil {
		return fmt.Errorf("loading AWS config: %v", err)
	}

	// Create S3 client with path-style addressing to avoid redirect issues
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	s3Key := fmt.Sprintf("%s/%s/%s-attempt-%d.txt", project, pr, fileType, attempt)

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(debugBucket),
		Key:    aws.String(s3Key),
		Body:   bytes.NewReader(content),
	})

	if err != nil {
		return fmt.Errorf("uploading to S3: %v", err)
	}

	logger.Info("Uploaded debug file to S3", "bucket", debugBucket, "key", s3Key, "region", bucketRegion)
	return nil
}
