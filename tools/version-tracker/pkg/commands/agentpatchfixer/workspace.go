package agentpatchfixer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const maxDiagnosticDiffBytes = 50000

type patchMetadata struct {
	AuthorName  string
	AuthorEmail string
	AuthorDate  string
	Subject     string
	Body        string
}

type patchSeed struct {
	RejectFiles  []string
	MissingFiles []string
}

func parsePatchMetadata(ctx context.Context, patchText string) (patchMetadata, error) {
	tempDir, err := os.MkdirTemp("", "patch-fixer-mailinfo-")
	if err != nil {
		return patchMetadata{}, err
	}
	defer os.RemoveAll(tempDir)

	messagePath := filepath.Join(tempDir, "message")
	patchPath := filepath.Join(tempDir, "patch")
	output, err := runCommand(ctx, tempDir, strings.NewReader(patchText), nil, "git", "mailinfo", messagePath, patchPath)
	if err != nil {
		return patchMetadata{}, fmt.Errorf("unable to parse patch metadata: %v", err)
	}
	fields := make(map[string]string)
	for _, line := range strings.Split(output, "\n") {
		key, value, found := strings.Cut(line, ": ")
		if found {
			fields[key] = value
		}
	}
	body, err := os.ReadFile(messagePath)
	if err != nil {
		return patchMetadata{}, err
	}
	metadata := patchMetadata{
		AuthorName:  fields["Author"],
		AuthorEmail: fields["Email"],
		AuthorDate:  fields["Date"],
		Subject:     fields["Subject"],
		Body:        strings.TrimSpace(string(body)),
	}
	if metadata.AuthorName == "" || metadata.AuthorEmail == "" || metadata.AuthorDate == "" || metadata.Subject == "" {
		return patchMetadata{}, fmt.Errorf("patch metadata is incomplete")
	}
	return metadata, nil
}

func patchFiles(ctx context.Context, patchText string) (map[string]bool, error) {
	diffText, numstat, err := parsePatchDiff(ctx, patchText)
	if err != nil {
		return nil, err
	}
	for _, unsupported := range []string{"new file mode 120000", "old mode 120000", "GIT binary patch"} {
		if strings.Contains(diffText, unsupported) {
			return nil, fmt.Errorf("symlink and binary patches are not supported by the agent fixer")
		}
	}
	if strings.Contains(diffText, "rename from ") || strings.Contains(diffText, "rename to ") {
		return nil, fmt.Errorf("rename patches are not supported by the agent fixer")
	}

	files := make(map[string]bool)
	for _, record := range strings.Split(strings.TrimSuffix(numstat, "\x00"), "\x00") {
		if record == "" {
			continue
		}
		fields := strings.SplitN(record, "\t", 3)
		if len(fields) != 3 || fields[2] == "" {
			return nil, fmt.Errorf("invalid git apply --numstat record %q", record)
		}
		files[filepath.ToSlash(fields[2])] = true
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("patch contains no diff files")
	}
	return files, nil
}

func parsePatchDiff(ctx context.Context, patchText string) (string, string, error) {
	lines := strings.Split(patchText, "\n")
	var lastErr error
	for index, line := range lines {
		if !strings.HasPrefix(line, "diff --git ") {
			continue
		}
		diffText := strings.Join(lines[index:], "\n")
		numstat, err := runCommand(
			ctx,
			os.TempDir(),
			strings.NewReader(diffText),
			nil,
			"git",
			"apply",
			"--numstat",
			"-z",
		)
		if err == nil {
			return diffText, numstat, nil
		}
		lastErr = err
	}
	if lastErr != nil {
		return "", "", fmt.Errorf("parsing patch file list: %v", lastErr)
	}
	return "", "", fmt.Errorf("patch contains no diff")
}

func prepareWorkspace(ctx context.Context, sourceRepo, parentDir string) (string, error) {
	target := filepath.Join(parentDir, "workspace")
	if _, err := runCommand(ctx, parentDir, nil, nil, "git", "clone", "--shared", "--no-checkout", sourceRepo, target); err != nil {
		return "", err
	}
	for _, args := range [][]string{
		{"checkout", "--detach", "HEAD"},
		{"reset", "--hard", "HEAD"},
		{"clean", "-fdx"},
	} {
		if _, err := runCommand(ctx, target, nil, nil, "git", args...); err != nil {
			return "", err
		}
	}
	return target, nil
}

func seedWorkspace(ctx context.Context, workspacePath, patchPath string) (patchSeed, error) {
	output, err := runCommandAllowFailure(ctx, workspacePath, nil, nil, "git", "apply", "--reject", "--whitespace=nowarn", patchPath)
	var rejects []string
	walkErr := filepath.WalkDir(workspacePath, func(current string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".rej") {
			relative, relErr := filepath.Rel(workspacePath, current)
			if relErr != nil {
				return relErr
			}
			rejects = append(rejects, filepath.ToSlash(relative))
		}
		return nil
	})
	if walkErr != nil {
		return patchSeed{}, walkErr
	}
	sort.Strings(rejects)
	missingFiles := missingFilesFromApplyOutput(output)
	if err != nil && len(rejects) == 0 {
		hasChanges, statusErr := workspaceHasChanges(ctx, workspacePath)
		if statusErr != nil {
			return patchSeed{}, statusErr
		}
		if !hasChanges && len(missingFiles) == 0 {
			return patchSeed{}, fmt.Errorf("unable to seed workspace from failed patch: %v", err)
		}
	}
	return patchSeed{RejectFiles: rejects, MissingFiles: missingFiles}, nil
}

func removeRejectFiles(workspacePath string, rejectFiles []string) {
	for _, relativePath := range rejectFiles {
		_ = os.Remove(filepath.Join(workspacePath, filepath.FromSlash(relativePath)))
	}
}

func generateCandidate(
	ctx context.Context,
	workspacePath string,
	metadata patchMetadata,
	outputPath string,
	allowedFiles map[string]bool,
) ([]string, error) {
	editedFiles, err := changedFiles(ctx, workspacePath)
	if err != nil {
		return nil, err
	}
	if len(editedFiles) == 0 {
		return nil, nil
	}
	for _, editedFile := range editedFiles {
		if !allowedFiles[editedFile] {
			return nil, fmt.Errorf("candidate changed file outside allowlist: %s", editedFile)
		}
	}

	if _, err := runCommand(ctx, workspacePath, nil, nil, "git", "add", "--all"); err != nil {
		return nil, err
	}
	if _, err := runCommand(ctx, workspacePath, nil, nil, "git", "diff", "--cached", "--check"); err != nil {
		return nil, err
	}
	message := metadata.Subject
	if metadata.Body != "" {
		message += "\n\n" + metadata.Body
	}
	messagePath := filepath.Join(workspacePath, ".git", "patch-fixer-message")
	if err := os.WriteFile(messagePath, []byte(message+"\n"), 0o600); err != nil {
		return nil, err
	}
	env := map[string]string{
		"GIT_AUTHOR_NAME":     metadata.AuthorName,
		"GIT_AUTHOR_EMAIL":    metadata.AuthorEmail,
		"GIT_AUTHOR_DATE":     metadata.AuthorDate,
		"GIT_COMMITTER_NAME":  "EKS Distro PR Bot",
		"GIT_COMMITTER_EMAIL": "aws-model-rocket-bots+eksdistroprbot@amazon.com",
	}
	if _, err := runCommand(ctx, workspacePath, nil, env, "git", "commit", "--file", messagePath); err != nil {
		return nil, err
	}
	patch, err := runCommand(ctx, workspacePath, nil, nil, "git", "format-patch", "-1", "--stdout", "--no-signature")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(outputPath, []byte(patch), 0o600); err != nil {
		return nil, err
	}
	return editedFiles, nil
}

func showWorkspaceDiff(ctx context.Context, workspacePath string) (string, error) {
	output, err := runCommand(ctx, workspacePath, nil, nil, "git", "diff", "--")
	if len(output) > maxDiagnosticDiffBytes {
		output = output[:maxDiagnosticDiffBytes]
	}
	return output, err
}

func runCommand(ctx context.Context, dir string, stdin io.Reader, extraEnv map[string]string, name string, args ...string) (string, error) {
	output, err := runCommandAllowFailure(ctx, dir, stdin, extraEnv, name, args...)
	if err != nil {
		return "", err
	}
	return output, nil
}

func runCommandAllowFailure(ctx context.Context, dir string, stdin io.Reader, extraEnv map[string]string, name string, args ...string) (string, error) {
	command := exec.CommandContext(ctx, name, args...)
	command.Dir = dir
	command.Stdin = stdin
	command.Env = os.Environ()
	for key, value := range extraEnv {
		command.Env = append(command.Env, key+"="+value)
	}
	var output bytes.Buffer
	command.Stdout = &output
	command.Stderr = &output
	err := command.Run()
	if err != nil {
		return output.String(), fmt.Errorf("command failed (%s %s): %s", name, strings.Join(args, " "), strings.TrimSpace(output.String()))
	}
	return output.String(), nil
}

func pathEscapesRoot(path string) bool {
	cleaned := filepath.Clean(path)
	return cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator))
}

func workspaceHasChanges(ctx context.Context, workspacePath string) (bool, error) {
	files, err := changedFiles(ctx, workspacePath)
	if err != nil {
		return false, err
	}
	return len(files) > 0, nil
}

func changedFiles(ctx context.Context, workspacePath string) ([]string, error) {
	status, err := runCommand(
		ctx,
		workspacePath,
		nil,
		nil,
		"git",
		"status",
		"--porcelain=v1",
		"-z",
		"--untracked-files=all",
	)
	if err != nil {
		return nil, err
	}

	files := make(map[string]bool)
	records := strings.Split(status, "\x00")
	for index := 0; index < len(records); index++ {
		record := records[index]
		if record == "" {
			continue
		}
		if len(record) < 4 || record[2] != ' ' {
			return nil, fmt.Errorf("invalid git status record %q", record)
		}
		statusCode := record[:2]
		files[filepath.ToSlash(record[3:])] = true
		if strings.ContainsAny(statusCode, "RC") {
			index++
			if index >= len(records) || records[index] == "" {
				return nil, fmt.Errorf("git status rename record is missing its source path")
			}
			files[filepath.ToSlash(records[index])] = true
		}
	}

	result := make([]string, 0, len(files))
	for file := range files {
		result = append(result, file)
	}
	sort.Strings(result)
	return result, nil
}

func missingFilesFromApplyOutput(output string) []string {
	missing := make(map[string]bool)
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "error: ") {
			continue
		}
		detail := strings.TrimPrefix(line, "error: ")
		for _, suffix := range []string{": No such file or directory", ": does not exist in index"} {
			if path, found := strings.CutSuffix(detail, suffix); found && path != "" {
				missing[filepath.ToSlash(path)] = true
			}
		}
	}
	result := make([]string, 0, len(missing))
	for file := range missing {
		result = append(result, file)
	}
	sort.Strings(result)
	return result
}
