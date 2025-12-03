package fixpatches

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/types"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

// ExtractPatchContext builds context for LLM from .rej files for a specific patch.
func ExtractPatchContext(rejFiles []string, patchFile string, projectPath string, attempt int, patchResult *types.PatchApplicationResult) (*types.PatchContext, error) {
	if len(rejFiles) == 0 {
		return nil, fmt.Errorf("no rejection files provided")
	}

	if patchFile == "" {
		return nil, fmt.Errorf("patch file not specified")
	}

	ctx := &types.PatchContext{
		FailedHunks:       make([]types.FailedHunk, 0),
		CurrentFileState:  make(map[string]string),
		PreviousAttempts:  make([]string, 0),
		ApplicationResult: patchResult,
		AllFileContexts:   make(map[string]string),
	}

	// Extract project name from path (e.g., "projects/aquasecurity/trivy" -> "aquasecurity/trivy")
	pathParts := strings.Split(projectPath, string(filepath.Separator))
	if len(pathParts) >= 2 {
		ctx.ProjectName = filepath.Join(pathParts[len(pathParts)-2], pathParts[len(pathParts)-1])
	}

	logger.Info("Extracting patch context", "patch_file", filepath.Base(patchFile), "rej_files", len(rejFiles))

	// Read and parse the specified patch file
	patchContent, err := os.ReadFile(patchFile)
	if err != nil {
		return nil, fmt.Errorf("reading patch file %s: %v", patchFile, err)
	}

	ctx.OriginalPatch = string(patchContent)

	// NEW: Parse patch to find all files and extract context for each
	patchFiles, err := parsePatchFiles(string(patchContent))
	if err != nil {
		logger.Info("Warning: could not parse patch files", "error", err)
	} else {
		logger.Info("Parsed patch files", "count", len(patchFiles))

		// OPTIMIZATION: Only extract context for files that actually failed
		// Clean files don't need full context - they'll be copied from original patch
		failedFileMap := make(map[string]bool)
		for _, rejFile := range rejFiles {
			// Extract base filename from .rej file path
			baseName := filepath.Base(strings.TrimSuffix(rejFile, ".rej"))
			failedFileMap[baseName] = true
		}

		// Filter to only failed files
		failedFiles := make([]PatchFile, 0)
		for _, file := range patchFiles {
			baseName := filepath.Base(file.Path)
			if failedFileMap[baseName] {
				failedFiles = append(failedFiles, file)
			}
		}

		logger.Info("Categorized patch files",
			"total", len(patchFiles),
			"failed", len(failedFiles),
			"clean", len(patchFiles)-len(failedFiles))

		// Use PRISTINE content from patchResult if available
		// This is critical: pristine content was captured BEFORE git apply modified the files
		if patchResult != nil && len(patchResult.PristineContent) > 0 {
			logger.Info("Using pristine content from before patch application", "files", len(patchResult.PristineContent))
			// Extract context ONLY for failed files to save tokens
			allFileContexts := extractContextFromPristine(failedFiles, patchResult.PristineContent)
			ctx.AllFileContexts = allFileContexts
			logger.Info("Extracted context from pristine files (failed only)",
				"count", len(allFileContexts),
				"token_savings_estimate", (len(patchFiles)-len(failedFiles))*700)
		} else {
			// Fallback: read from current files (may be modified)
			logger.Info("Warning: no pristine content available, reading from current files")
			// Determine the repository path from the .rej files
			var repoPath string
			if len(rejFiles) > 0 {
				rejFileDir := filepath.Dir(rejFiles[0])
				repoPath = rejFileDir
				logger.Info("Determined repository path from .rej file", "repo_path", repoPath)
			} else {
				repoPath = projectPath
				logger.Info("No .rej files to determine repo path, using project path", "repo_path", repoPath)
			}

			// Extract context ONLY for failed files
			allFileContexts := extractContextForAllFiles(failedFiles, repoPath)
			ctx.AllFileContexts = allFileContexts
			logger.Info("Extracted context for failed files only", "count", len(allFileContexts))
		}
	}

	// Extract patch metadata from headers
	if err := extractPatchMetadata(ctx, string(patchContent)); err != nil {
		logger.Info("Warning: failed to extract patch metadata", "error", err)
		// Continue anyway - metadata is helpful but not critical
	}

	// Parse each .rej file
	for i, rejFile := range rejFiles {
		logger.Info("Parsing rejection file", "file", rejFile, "index", i+1)

		hunks, err := parseRejectionFile(rejFile, projectPath)
		if err != nil {
			logger.Info("Warning: failed to parse rejection file", "file", rejFile, "error", err)
			continue
		}

		ctx.FailedHunks = append(ctx.FailedHunks, hunks...)
	}

	if len(ctx.FailedHunks) == 0 {
		return nil, fmt.Errorf("no failed hunks extracted from rejection files")
	}

	// Extract context from current files for each failed hunk
	for i := range ctx.FailedHunks {
		hunk := &ctx.FailedHunks[i]
		context, err := extractFileContext(hunk.FilePath, hunk.LineNumber, projectPath)
		if err != nil {
			logger.Info("Warning: failed to extract file context", "file", hunk.FilePath, "error", err)
			hunk.Context = fmt.Sprintf("(Unable to read file context: %v)", err)
		} else {
			hunk.Context = context
		}

		// Extract expected vs actual comparison
		if err := extractExpectedVsActual(hunk); err != nil {
			logger.Info("Warning: failed to extract expected vs actual comparison", "file", hunk.FilePath, "error", err)
			// Continue anyway - this is supplementary information
		}
	}

	// Estimate token count
	ctx.TokenCount = estimateTokenCount(ctx)

	logger.Info("Context extraction complete",
		"hunks", len(ctx.FailedHunks),
		"estimated_tokens", ctx.TokenCount)

	return ctx, nil
}

// extractPatchMetadata extracts From, Date, and Subject from patch headers.
func extractPatchMetadata(ctx *types.PatchContext, patchContent string) error {
	scanner := bufio.NewScanner(strings.NewReader(patchContent))

	// Regex patterns for patch headers
	fromPattern := regexp.MustCompile(`^From:\s*(.+)$`)
	datePattern := regexp.MustCompile(`^Date:\s*(.+)$`)
	subjectPattern := regexp.MustCompile(`^Subject:\s*(.+)$`)

	var inSubject bool
	for scanner.Scan() {
		line := scanner.Text()

		// Stop at the first diff marker
		if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "diff --git") {
			break
		}

		if match := fromPattern.FindStringSubmatch(line); match != nil {
			ctx.PatchAuthor = match[1]
			inSubject = false
		} else if match := datePattern.FindStringSubmatch(line); match != nil {
			ctx.PatchDate = match[1]
			inSubject = false
		} else if match := subjectPattern.FindStringSubmatch(line); match != nil {
			ctx.PatchSubject = match[1]
			inSubject = true
		} else if inSubject && strings.HasPrefix(line, " ") {
			// Continuation of Subject line (starts with space)
			ctx.PatchSubject += "\n" + line
		} else if inSubject && line == "" {
			// Empty line marks end of Subject
			inSubject = false
		}
	}

	// Extract intent from complete subject (remove [PATCH] prefix if present)
	if ctx.PatchSubject != "" {
		subject := strings.TrimSpace(ctx.PatchSubject)
		subject = strings.TrimPrefix(subject, "[PATCH]")
		subject = strings.TrimSpace(subject)
		ctx.PatchIntent = subject
	}

	if ctx.PatchAuthor == "" || ctx.PatchDate == "" || ctx.PatchSubject == "" {
		return fmt.Errorf("incomplete patch metadata: author=%q, date=%q, subject=%q",
			ctx.PatchAuthor, ctx.PatchDate, ctx.PatchSubject)
	}

	return nil
}

// parseRejectionFile parses a .rej file and extracts failed hunks.
func parseRejectionFile(rejFile string, projectPath string) ([]types.FailedHunk, error) {
	content, err := os.ReadFile(rejFile)
	if err != nil {
		return nil, fmt.Errorf("reading rejection file: %v", err)
	}

	// Determine the original file path from .rej file
	// .rej files are named like "go.mod.rej" and are in the same directory as the original file
	filePath := strings.TrimSuffix(rejFile, ".rej")

	hunks := make([]types.FailedHunk, 0)
	lines := strings.Split(string(content), "\n")

	var currentHunk *types.FailedHunk
	hunkIndex := 0

	for i, line := range lines {
		// Detect hunk header: @@ -line,count +line,count @@
		if strings.HasPrefix(line, "@@") {
			// Save previous hunk if exists
			if currentHunk != nil {
				hunks = append(hunks, *currentHunk)
			}

			// Start new hunk
			hunkIndex++
			currentHunk = &types.FailedHunk{
				FilePath:      filePath,
				HunkIndex:     hunkIndex,
				OriginalLines: make([]string, 0),
			}

			// Extract line number from hunk header
			// Format: @@ -old_start,old_count +new_start,new_count @@
			lineNumPattern := regexp.MustCompile(`@@\s+-(\d+)`)
			if match := lineNumPattern.FindStringSubmatch(line); match != nil {
				if lineNum, err := fmt.Sscanf(match[1], "%d", &currentHunk.LineNumber); err == nil {
					_ = lineNum
				}
			}

			currentHunk.OriginalLines = append(currentHunk.OriginalLines, line)
		} else if currentHunk != nil {
			// Add line to current hunk
			currentHunk.OriginalLines = append(currentHunk.OriginalLines, line)
		} else if i == 0 {
			// First line might be a header comment, skip it
			continue
		}
	}

	// Save last hunk
	if currentHunk != nil {
		hunks = append(hunks, *currentHunk)
	}

	return hunks, nil
}

// extractFileContext reads ±50 lines around the conflict point in the current file.
func extractFileContext(filePath string, lineNumber int, projectPath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("reading file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	totalLines := len(lines)

	// Calculate context window (±50 lines)
	contextWindow := 50
	var startLine, endLine int

	// Handle case where lineNumber is beyond the file length
	// This can happen when the file has changed significantly since the patch was created
	if lineNumber > totalLines {
		// Provide context from the end of the file instead
		startLine = totalLines - contextWindow
		if startLine < 0 {
			startLine = 0
		}
		endLine = totalLines
	} else {
		// Normal case: extract context around the line number
		startLine = lineNumber - contextWindow
		if startLine < 0 {
			startLine = 0
		}

		endLine = lineNumber + contextWindow
		if endLine > totalLines {
			endLine = totalLines
		}
	}

	// Extract context lines
	extractedLines := lines[startLine:endLine]
	context := strings.Join(extractedLines, "\n")

	// Add line number markers for clarity
	var result string
	if lineNumber > totalLines {
		result = fmt.Sprintf("Note: Patch expects line %d, but file only has %d lines. Showing end of file (lines %d-%d) of %s:\n%s",
			lineNumber, totalLines, startLine+1, endLine, filepath.Base(filePath), context)
	} else {
		result = fmt.Sprintf("Lines %d-%d of %s:\n%s", startLine+1, endLine, filepath.Base(filePath), context)
	}

	return result, nil
}

// estimateTokenCount provides a rough estimate of token usage.
// Rule of thumb: ~3 characters per token for code (more conservative than 4 for prose).
func estimateTokenCount(ctx *types.PatchContext) int {
	totalChars := 0

	// Count original patch
	totalChars += len(ctx.OriginalPatch)

	// Count failed hunks
	for _, hunk := range ctx.FailedHunks {
		for _, line := range hunk.OriginalLines {
			totalChars += len(line)
		}
		totalChars += len(hunk.Context)
	}

	// Count all file contexts
	for _, context := range ctx.AllFileContexts {
		totalChars += len(context)
	}

	// Count build error if present
	totalChars += len(ctx.BuildError)

	// Count previous attempts
	for _, attempt := range ctx.PreviousAttempts {
		totalChars += len(attempt)
	}

	// Estimate tokens (3 chars per token for code - more conservative)
	tokens := totalChars / 3

	// Add overhead for prompt structure and template (~2000 tokens)
	tokens += 2000

	return tokens
}

// PruneContext reduces context size if it exceeds token limits.
// This is critical for staying within Claude's context window.
func PruneContext(ctx *types.PatchContext, maxTokens int) error {
	logger.Info("Pruning context", "current_tokens", ctx.TokenCount, "max_tokens", maxTokens)

	if ctx.TokenCount <= maxTokens {
		logger.Info("Context within limits, no pruning needed")
		return nil
	}

	// Strategy: Prioritize most relevant context
	// 1. Keep patch metadata (small, critical)
	// 2. Keep failed hunks (essential)
	// 3. Reduce file context (can be trimmed)
	// 4. Reduce original patch (keep only relevant sections)
	// 5. Reduce previous attempts (keep only last 2)

	// Step 1: Reduce file context from ±50 lines to ±25 lines
	if ctx.TokenCount > maxTokens {
		logger.Info("Reducing file context window")
		for i := range ctx.FailedHunks {
			hunk := &ctx.FailedHunks[i]
			// Re-extract with smaller window
			context, err := extractFileContextWithWindow(hunk.FilePath, hunk.LineNumber, 25)
			if err != nil {
				logger.Info("Warning: failed to re-extract context", "file", hunk.FilePath, "error", err)
			} else {
				hunk.Context = context
			}
		}
		ctx.TokenCount = estimateTokenCount(ctx)
		logger.Info("After reducing file context", "tokens", ctx.TokenCount)
	}

	// Step 2: Further reduce file context to ±10 lines if still too large
	if ctx.TokenCount > maxTokens {
		logger.Info("Further reducing file context window")
		for i := range ctx.FailedHunks {
			hunk := &ctx.FailedHunks[i]
			context, err := extractFileContextWithWindow(hunk.FilePath, hunk.LineNumber, 10)
			if err != nil {
				logger.Info("Warning: failed to re-extract context", "file", hunk.FilePath, "error", err)
			} else {
				hunk.Context = context
			}
		}
		ctx.TokenCount = estimateTokenCount(ctx)
		logger.Info("After further reducing file context", "tokens", ctx.TokenCount)
	}

	// Step 2.5: Reduce AllFileContexts window size (these can be huge for go.mod/go.sum)
	if ctx.TokenCount > maxTokens && len(ctx.AllFileContexts) > 0 {
		logger.Info("Reducing AllFileContexts window size")
		// For each file context, reduce the window to ±5 lines instead of ±10
		for filename, context := range ctx.AllFileContexts {
			lines := strings.Split(context, "\n")
			if len(lines) > 20 {
				// Keep only first 10 lines (header + some context)
				ctx.AllFileContexts[filename] = strings.Join(lines[:10], "\n") + "\n...(truncated for brevity)..."
				logger.Info("Truncated file context", "file", filename, "original_lines", len(lines), "kept_lines", 10)
			}
		}
		ctx.TokenCount = estimateTokenCount(ctx)
		logger.Info("After reducing AllFileContexts", "tokens", ctx.TokenCount)
	}

	// Step 3: Aggressively reduce AllFileContexts - clear them entirely if needed
	if ctx.TokenCount > maxTokens && len(ctx.AllFileContexts) > 0 {
		logger.Info("Clearing AllFileContexts to save tokens")
		ctx.AllFileContexts = make(map[string]string)
		ctx.TokenCount = estimateTokenCount(ctx)
		logger.Info("After clearing AllFileContexts", "tokens", ctx.TokenCount)
	}

	// Step 4: Trim original patch to only include relevant sections
	if ctx.TokenCount > maxTokens {
		logger.Info("Trimming original patch")
		// Keep only the diff sections, remove commit message details
		lines := strings.Split(ctx.OriginalPatch, "\n")
		var trimmedLines []string
		inDiff := false
		for _, line := range lines {
			if strings.HasPrefix(line, "diff --git") || strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "@@") {
				inDiff = true
			}
			if inDiff {
				trimmedLines = append(trimmedLines, line)
			}
		}
		if len(trimmedLines) > 0 {
			ctx.OriginalPatch = strings.Join(trimmedLines, "\n")
		}
		ctx.TokenCount = estimateTokenCount(ctx)
		logger.Info("After trimming original patch", "tokens", ctx.TokenCount)
	}

	// Step 5: Keep only last 2 previous attempts
	if ctx.TokenCount > maxTokens && len(ctx.PreviousAttempts) > 2 {
		logger.Info("Reducing previous attempts history")
		ctx.PreviousAttempts = ctx.PreviousAttempts[len(ctx.PreviousAttempts)-2:]
		ctx.TokenCount = estimateTokenCount(ctx)
		logger.Info("After reducing previous attempts", "tokens", ctx.TokenCount)
	}

	// Step 6: Truncate build error to last 200 lines if still too large
	if ctx.TokenCount > maxTokens && len(ctx.BuildError) > 0 {
		logger.Info("Truncating build error")
		lines := strings.Split(ctx.BuildError, "\n")
		if len(lines) > 200 {
			ctx.BuildError = strings.Join(lines[len(lines)-200:], "\n")
		}
		ctx.TokenCount = estimateTokenCount(ctx)
		logger.Info("After truncating build error", "tokens", ctx.TokenCount)
	}

	// Step 7: Last resort - truncate original patch to first 50% if still too large
	if ctx.TokenCount > maxTokens {
		logger.Info("Last resort: truncating original patch")
		lines := strings.Split(ctx.OriginalPatch, "\n")
		if len(lines) > 100 {
			keepLines := len(lines) / 2
			ctx.OriginalPatch = strings.Join(lines[:keepLines], "\n") + "\n...(truncated for token limit)..."
			logger.Info("Truncated original patch", "original_lines", len(lines), "kept_lines", keepLines)
		}
		ctx.TokenCount = estimateTokenCount(ctx)
		logger.Info("After truncating original patch", "tokens", ctx.TokenCount)
	}

	// Final check
	if ctx.TokenCount > maxTokens {
		return fmt.Errorf("unable to prune context below %d tokens (current: %d)", maxTokens, ctx.TokenCount)
	}

	logger.Info("Context pruning complete", "final_tokens", ctx.TokenCount)
	return nil
}

// extractFileContextWithWindow reads a specific number of lines around the conflict point.
func extractFileContextWithWindow(filePath string, lineNumber int, windowSize int) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("reading file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	totalLines := len(lines)

	// Calculate context window
	var startLine, endLine int

	// Handle case where lineNumber is beyond the file length
	if lineNumber > totalLines {
		// Provide context from the end of the file instead
		startLine = totalLines - windowSize
		if startLine < 0 {
			startLine = 0
		}
		endLine = totalLines
	} else {
		// Normal case: extract context around the line number
		startLine = lineNumber - windowSize
		if startLine < 0 {
			startLine = 0
		}

		endLine = lineNumber + windowSize
		if endLine > totalLines {
			endLine = totalLines
		}
	}

	// Extract context lines
	extractedLines := lines[startLine:endLine]
	context := strings.Join(extractedLines, "\n")

	// Add line number markers for clarity
	var result string
	if lineNumber > totalLines {
		result = fmt.Sprintf("Note: Patch expects line %d, but file only has %d lines. Showing end of file (lines %d-%d) of %s:\n%s",
			lineNumber, totalLines, startLine+1, endLine, filepath.Base(filePath), context)
	} else {
		result = fmt.Sprintf("Lines %d-%d of %s:\n%s", startLine+1, endLine, filepath.Base(filePath), context)
	}

	return result, nil
}

// extractExpectedVsActual extracts what the patch expects vs what's actually in the file.
// This helps the LLM understand why the patch failed (whitespace, line shifts, content changes).
func extractExpectedVsActual(hunk *types.FailedHunk) error {
	// Parse the .rej file content to extract expected context lines
	// Context lines in unified diff format start with a space
	// Lines to be removed start with '-'
	// Lines to be added start with '+'

	expectedLines := make([]string, 0)
	for _, line := range hunk.OriginalLines {
		// Skip hunk header
		if strings.HasPrefix(line, "@@") {
			continue
		}
		// Context lines (what the patch expects to find) start with space or are removal lines
		if len(line) > 0 && (line[0] == ' ' || line[0] == '-') {
			// Remove the diff marker to get the actual line content
			if len(line) > 1 {
				expectedLines = append(expectedLines, line[1:])
			} else {
				expectedLines = append(expectedLines, "")
			}
		}
	}

	hunk.ExpectedContext = expectedLines

	// Read the actual file content at the target location
	content, err := os.ReadFile(hunk.FilePath)
	if err != nil {
		return fmt.Errorf("reading file: %v", err)
	}

	lines := strings.Split(string(content), "\n")

	// Extract actual lines around the target location
	// Use a window that matches the expected context size
	contextSize := len(expectedLines)
	if contextSize == 0 {
		contextSize = 10 // Default fallback
	}

	startLine := hunk.LineNumber - 1 // Convert to 0-based index
	if startLine < 0 {
		startLine = 0
	}

	endLine := startLine + contextSize
	if endLine > len(lines) {
		endLine = len(lines)
	}

	// Also check a few lines before in case line numbers shifted
	searchStart := startLine - 5
	if searchStart < 0 {
		searchStart = 0
	}

	actualLines := make([]string, 0)
	if startLine < len(lines) {
		actualLines = lines[startLine:endLine]
	}

	hunk.ActualContext = actualLines

	// Identify specific differences
	differences := make([]string, 0)

	// Compare line by line
	maxLen := len(expectedLines)
	if len(actualLines) > maxLen {
		maxLen = len(actualLines)
	}

	for i := 0; i < maxLen; i++ {
		var expected, actual string
		if i < len(expectedLines) {
			expected = expectedLines[i]
		}
		if i < len(actualLines) {
			actual = actualLines[i]
		}

		if expected != actual {
			// Check for whitespace-only differences
			if strings.TrimSpace(expected) == strings.TrimSpace(actual) {
				if expected == "" && actual != "" {
					differences = append(differences, fmt.Sprintf("Line %d: Patch expects blank line, but file has: %q", i+1, actual))
				} else if expected != "" && actual == "" {
					differences = append(differences, fmt.Sprintf("Line %d: Patch expects %q, but file has blank line", i+1, expected))
				} else {
					differences = append(differences, fmt.Sprintf("Line %d: Whitespace difference (content matches)", i+1))
				}
			} else {
				// Content difference
				if expected == "" {
					differences = append(differences, fmt.Sprintf("Line %d: Patch expects blank line, file has: %q", i+1, actual))
				} else if actual == "" {
					differences = append(differences, fmt.Sprintf("Line %d: File has blank line, patch expects: %q", i+1, expected))
				} else {
					differences = append(differences, fmt.Sprintf("Line %d: Content differs\n  Expected: %q\n  Actual:   %q", i+1, expected, actual))
				}
			}
		}
	}

	// Check for line count mismatch
	if len(expectedLines) != len(actualLines) {
		differences = append(differences, fmt.Sprintf("Line count mismatch: patch expects %d lines, file has %d lines", len(expectedLines), len(actualLines)))
	}

	hunk.Differences = differences

	logger.Info("Extracted expected vs actual comparison",
		"file", filepath.Base(hunk.FilePath),
		"expected_lines", len(expectedLines),
		"actual_lines", len(actualLines),
		"differences", len(differences))

	return nil
}

// PatchFile represents a file being modified in a patch.
type PatchFile struct {
	Path       string
	LineRanges []LineRange // Lines being changed
}

// LineRange represents a range of lines in a file.
type LineRange struct {
	Start int
	End   int
}

// parsePatchFiles extracts all files and their changed line ranges from a patch.
func parsePatchFiles(patchContent string) ([]PatchFile, error) {
	files := make([]PatchFile, 0)

	scanner := bufio.NewScanner(strings.NewReader(patchContent))
	var currentFile *PatchFile

	for scanner.Scan() {
		line := scanner.Text()

		// New file: diff --git a/file b/file
		if strings.HasPrefix(line, "diff --git") {
			if currentFile != nil {
				files = append(files, *currentFile)
			}
			// Extract filename from "diff --git a/file b/file"
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				filename := strings.TrimPrefix(parts[3], "b/")
				currentFile = &PatchFile{
					Path:       filename,
					LineRanges: make([]LineRange, 0),
				}
			}
		}

		// Hunk header: @@ -10,5 +12,7 @@
		if strings.HasPrefix(line, "@@") && currentFile != nil {
			// Parse to get new line range (+12,7 means start at 12, 7 lines)
			re := regexp.MustCompile(`\+(\d+),(\d+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 3 {
				start := 0
				count := 0
				fmt.Sscanf(matches[1], "%d", &start)
				fmt.Sscanf(matches[2], "%d", &count)
				currentFile.LineRanges = append(currentFile.LineRanges, LineRange{
					Start: start,
					End:   start + count,
				})
			} else {
				// Handle single line case: @@ -10 +12 @@
				re := regexp.MustCompile(`\+(\d+)`)
				matches := re.FindStringSubmatch(line)
				if len(matches) >= 2 {
					start := 0
					fmt.Sscanf(matches[1], "%d", &start)
					currentFile.LineRanges = append(currentFile.LineRanges, LineRange{
						Start: start,
						End:   start + 1,
					})
				}
			}
		}
	}

	// Don't forget last file
	if currentFile != nil {
		files = append(files, *currentFile)
	}

	return files, nil
}

// extractContextForAllFiles extracts current file content for all files in patch.
func extractContextForAllFiles(patchFiles []PatchFile, projectPath string) map[string]string {
	contexts := make(map[string]string)

	for _, patchFile := range patchFiles {
		filePath := filepath.Join(projectPath, patchFile.Path)

		// Read file
		content, err := os.ReadFile(filePath)
		if err != nil {
			logger.Info("Warning: could not read file", "file", patchFile.Path, "error", err)
			continue
		}

		lines := strings.Split(string(content), "\n")

		// Extract context around each changed line range
		var contextLines []string
		for _, lineRange := range patchFile.LineRanges {
			// Get ±10 lines around the change
			start := lineRange.Start - 10
			if start < 0 {
				start = 0
			}
			end := lineRange.End + 10
			if end > len(lines) {
				end = len(lines)
			}

			contextLines = append(contextLines, fmt.Sprintf("Lines %d-%d:", start+1, end))
			for i := start; i < end; i++ {
				contextLines = append(contextLines, lines[i])
			}
			contextLines = append(contextLines, "") // Blank line between ranges
		}

		contexts[patchFile.Path] = strings.Join(contextLines, "\n")
	}

	return contexts
}

// extractContextFromPristine extracts context from pristine file content (before git apply).
// This is the CORRECT way to extract context because it shows the LLM the original state,
// not the state after git apply has already modified some files.
func extractContextFromPristine(patchFiles []PatchFile, pristineContent map[string]string) map[string]string {
	contexts := make(map[string]string)

	for _, patchFile := range patchFiles {
		// Get pristine content for this file
		content, ok := pristineContent[patchFile.Path]
		if !ok {
			logger.Info("Warning: no pristine content for file", "file", patchFile.Path)
			continue
		}

		lines := strings.Split(content, "\n")

		// Extract context around each changed line range
		var contextLines []string
		for _, lineRange := range patchFile.LineRanges {
			// Get ±10 lines around the change
			start := lineRange.Start - 10
			if start < 0 {
				start = 0
			}
			end := lineRange.End + 10
			if end > len(lines) {
				end = len(lines)
			}

			contextLines = append(contextLines, fmt.Sprintf("Lines %d-%d:", start+1, end))
			for i := start; i < end; i++ {
				contextLines = append(contextLines, lines[i])
			}
			contextLines = append(contextLines, "") // Blank line between ranges
		}

		contexts[patchFile.Path] = strings.Join(contextLines, "\n")
	}

	return contexts
}
