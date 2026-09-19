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
	"strconv"
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
	RejectFiles []string
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

func patchFiles(patchText string) (map[string]bool, error) {
	lines := strings.Split(patchText, "\n")
	separator := -1
	for index, line := range lines {
		if line == "---" {
			separator = index
		}
	}
	if separator < 0 {
		return nil, fmt.Errorf("mail patch is missing the format-patch separator")
	}
	diffText := strings.Join(lines[separator+1:], "\n")
	for _, unsupported := range []string{"new file mode 120000", "old mode 120000", "GIT binary patch"} {
		if strings.Contains(diffText, unsupported) {
			return nil, fmt.Errorf("symlink and binary patches are not supported by the agent fixer")
		}
	}
	if strings.Contains(diffText, "rename from ") || strings.Contains(diffText, "rename to ") {
		return nil, fmt.Errorf("rename patches are not supported by the agent fixer")
	}

	files := make(map[string]bool)
	for _, line := range strings.Split(diffText, "\n") {
		if !strings.HasPrefix(line, "diff --git ") {
			continue
		}
		parts := splitGitHeader(line)
		if len(parts) >= 4 {
			files[strings.TrimPrefix(parts[3], "b/")] = true
		}
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("patch contains no diff files")
	}
	return files, nil
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
	_, err := runCommandAllowFailure(ctx, workspacePath, nil, nil, "git", "apply", "--reject", "--whitespace=nowarn", patchPath)
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
	if err != nil && len(rejects) == 0 {
		return patchSeed{}, fmt.Errorf("unable to seed workspace from failed patch: %v", err)
	}
	return patchSeed{RejectFiles: rejects}, nil
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
	if _, err := runCommand(ctx, workspacePath, nil, nil, "git", "diff", "--check"); err != nil {
		return nil, err
	}
	status, err := runCommand(ctx, workspacePath, nil, nil, "git", "status", "--porcelain")
	if err != nil {
		return nil, err
	}
	var editedFiles []string
	for _, line := range strings.Split(strings.TrimSuffix(status, "\n"), "\n") {
		if len(line) > 3 {
			editedFiles = append(editedFiles, filepath.ToSlash(line[3:]))
		}
	}
	if len(editedFiles) == 0 {
		return nil, nil
	}
	sort.Strings(editedFiles)
	for _, editedFile := range editedFiles {
		if !allowedFiles[editedFile] {
			return nil, fmt.Errorf("candidate changed file outside allowlist: %s", editedFile)
		}
	}

	if _, err := runCommand(ctx, workspacePath, nil, nil, "git", "add", "--all"); err != nil {
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

func splitGitHeader(line string) []string {
	var parts []string
	for len(line) > 0 {
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		if line[0] != '"' {
			if index := strings.IndexByte(line, ' '); index >= 0 {
				parts = append(parts, line[:index])
				line = line[index+1:]
			} else {
				parts = append(parts, line)
				break
			}
			continue
		}
		closed := false
		for index := 1; index < len(line); index++ {
			if line[index] == '"' && line[index-1] != '\\' {
				value, err := strconv.Unquote(line[:index+1])
				if err == nil {
					parts = append(parts, value)
				}
				line = line[index+1:]
				closed = true
				break
			}
		}
		if !closed {
			break
		}
	}
	return parts
}
