package projectpatchfixer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const autoscalerCleanupTimeout = 30 * time.Second

type autoscalerPatchMetadata struct {
	authorName  string
	authorEmail string
	authorDate  string
	subject     string
	body        string
}

var regenerateAutoscaler = regenerateAutoscalerPatches

func fixAutoscaler(ctx context.Context, request Request) (Result, error) {
	runDir, err := newRunDirectory(request.ProjectName)
	if err != nil {
		return Result{Handled: true}, err
	}
	outputDir, err := filepath.Abs(filepath.Join(runDir, "patches"))
	if err != nil {
		return Result{Handled: true, RunDirectory: runDir}, fmt.Errorf("resolving autoscaler patch output directory: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return Result{Handled: true}, fmt.Errorf("creating autoscaler patch output directory: %v", err)
	}

	patches, err := regenerateAutoscaler(ctx, request, outputDir)
	if err != nil {
		return Result{Handled: true, RunDirectory: runDir}, fmt.Errorf("regenerating autoscaler patches: %v", err)
	}
	if len(patches) == 0 {
		return Result{Handled: true, RunDirectory: runDir}, fmt.Errorf("autoscaler regenerator produced no patches")
	}
	for _, patch := range patches {
		if filepath.Clean(filepath.Dir(patch)) != filepath.Clean(outputDir) {
			return Result{Handled: true, RunDirectory: runDir}, fmt.Errorf("autoscaler patch escaped output directory: %s", patch)
		}
		if _, err := os.Stat(patch); err != nil {
			return Result{Handled: true, RunDirectory: runDir}, fmt.Errorf("autoscaler patch %q is not accessible: %v", patch, err)
		}
	}

	manifest, err := json.MarshalIndent(map[string][]string{"patches": patches}, "", "  ")
	if err != nil {
		return Result{Handled: true, RunDirectory: runDir}, fmt.Errorf("marshalling autoscaler patch manifest: %v", err)
	}
	manifest = append(manifest, '\n')
	if err := os.WriteFile(filepath.Join(runDir, "manifest.json"), manifest, 0o600); err != nil {
		return Result{Handled: true, RunDirectory: runDir}, fmt.Errorf("writing autoscaler patch manifest: %v", err)
	}

	return Result{
		Handled:      true,
		PatchesDir:   outputDir,
		PatchCount:   len(patches),
		RunDirectory: runDir,
	}, nil
}

func regenerateAutoscalerPatches(ctx context.Context, request Request, outputDir string) (patches []string, returnErr error) {
	source, err := filepath.Abs(request.UpstreamRepo)
	if err != nil {
		return nil, fmt.Errorf("resolving autoscaler source: %v", err)
	}
	patchesDir, err := filepath.Abs(request.PatchesDir)
	if err != nil {
		return nil, fmt.Errorf("resolving autoscaler patches: %v", err)
	}
	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		return nil, fmt.Errorf("resolving autoscaler output: %v", err)
	}

	existing, err := autoscalerPatchFiles(patchesDir)
	if err != nil {
		return nil, err
	}
	if len(existing) < 2 {
		return nil, fmt.Errorf("autoscaler requires at least provider and go.mod patches")
	}
	if err := restoreAutoscalerSource(ctx, source, request.TargetRevision); err != nil {
		return nil, err
	}
	defer func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), autoscalerCleanupTimeout)
		defer cancel()
		if err := restoreAutoscalerSource(cleanupCtx, source, request.TargetRevision); err != nil {
			returnErr = errors.Join(returnErr, fmt.Errorf("restoring autoscaler source: %v", err))
		}
	}()

	if _, err := runAutoscalerCommand(ctx, source, nil, nil, "git", "config", "user.name", "Prow Bot"); err != nil {
		return nil, err
	}
	if _, err := runAutoscalerCommand(ctx, source, nil, nil, "git", "config", "user.email", "prow@amazonaws.com"); err != nil {
		return nil, err
	}

	transformed, err := transformAutoscalerRouter(source)
	if err != nil {
		return nil, err
	}
	if !transformed {
		transformed, err = transformAutoscalerBuilder(source)
		if err != nil {
			return nil, err
		}
	}
	if !transformed {
		return nil, fmt.Errorf("unsupported autoscaler source: no router_all.go or builder_all.go")
	}

	firstPatch := existing[0]
	firstOutput := filepath.Join(outputDir, filepath.Base(firstPatch))
	metadata, err := parseAutoscalerPatchMetadata(ctx, firstPatch)
	if err != nil {
		return nil, err
	}
	written, err := commitAndFormatAutoscalerPatch(ctx, source, metadata, firstOutput, nil)
	if err != nil {
		return nil, err
	}
	if !written {
		return nil, fmt.Errorf("autoscaler provider transform produced no changes")
	}
	patches = append(patches, firstOutput)

	for _, staticPatch := range existing[1 : len(existing)-1] {
		outputPath := filepath.Join(outputDir, filepath.Base(staticPatch))
		if strings.Contains(strings.ToLower(filepath.Base(staticPatch)), "gce") {
			changed, err := removeAutoscalerGCEDependencies(source)
			if err != nil {
				return nil, err
			}
			if !changed {
				continue
			}
			metadata, err := parseAutoscalerPatchMetadata(ctx, staticPatch)
			if err != nil {
				return nil, err
			}
			written, err := commitAndFormatAutoscalerPatch(ctx, source, metadata, outputPath, nil)
			if err != nil {
				return nil, err
			}
			if written {
				patches = append(patches, outputPath)
			}
			continue
		}

		if _, err := runAutoscalerCommand(ctx, source, nil, nil, "git", "am", "--committer-date-is-author-date", staticPatch); err != nil {
			return nil, err
		}
		if err := copyAutoscalerPatch(staticPatch, outputPath); err != nil {
			return nil, err
		}
		patches = append(patches, outputPath)
	}

	dependencyPatch := existing[len(existing)-1]
	if err := removeAutoscalerCloudProviderDirectories(source); err != nil {
		return nil, err
	}
	clusterAutoscaler := filepath.Join(source, "cluster-autoscaler")
	goCommand := autoscalerGoBinary(request.ReleaseBranch, patchesDir)
	if _, err := runAutoscalerCommand(ctx, clusterAutoscaler, nil, nil, goCommand, "mod", "tidy"); err != nil {
		return nil, err
	}
	goModPaths := []string{"cluster-autoscaler/go.mod", "cluster-autoscaler/go.sum"}
	statusArgs := append([]string{"status", "--porcelain", "--"}, goModPaths...)
	status, err := runAutoscalerCommand(ctx, source, nil, nil, "git", statusArgs...)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(status) != "" {
		outputPath := filepath.Join(outputDir, filepath.Base(dependencyPatch))
		metadata, err := parseAutoscalerPatchMetadata(ctx, dependencyPatch)
		if err != nil {
			return nil, err
		}
		written, err := commitAndFormatAutoscalerPatch(ctx, source, metadata, outputPath, goModPaths)
		if err != nil {
			return nil, err
		}
		if written {
			patches = append(patches, outputPath)
		}
	}

	if err := restoreAutoscalerSource(ctx, source, request.TargetRevision); err != nil {
		return nil, err
	}
	for _, patch := range patches {
		if _, err := runAutoscalerCommand(ctx, source, nil, nil, "git", "am", "--committer-date-is-author-date", patch); err != nil {
			return nil, fmt.Errorf("validating regenerated autoscaler patch %q: %v", filepath.Base(patch), err)
		}
	}
	return patches, nil
}

func transformAutoscalerRouter(source string) (bool, error) {
	router := filepath.Join(source, "cluster-autoscaler", "cloudprovider", "router")
	routerAll := filepath.Join(router, "router_all.go")
	if _, err := os.Stat(routerAll); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if err := pruneAutoscalerProviderFiles(router, "router_", map[string]bool{
		"router_all.go":        true,
		"router_clusterapi.go": true,
	}); err != nil {
		return false, err
	}

	content, err := os.ReadFile(routerAll)
	if err != nil {
		return false, err
	}
	lines := strings.Split(strings.TrimSuffix(string(content), "\n"), "\n")
	transformed := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, `_ "`) &&
			strings.Contains(trimmed, "/cloudprovider/") &&
			!strings.Contains(trimmed, `/clusterapi"`) {
			continue
		}
		if strings.Contains(line, "builder.SetDefaultCloudProvider(") {
			line = leadingWhitespace(line) + "builder.SetDefaultCloudProvider(cloudprovider.ClusterAPIProviderName)"
		}
		transformed = append(transformed, line)
	}
	return true, os.WriteFile(routerAll, []byte(strings.Join(transformed, "\n")+"\n"), 0o644)
}

func transformAutoscalerBuilder(source string) (bool, error) {
	builder := filepath.Join(source, "cluster-autoscaler", "cloudprovider", "builder")
	builderAll := filepath.Join(builder, "builder_all.go")
	if _, err := os.Stat(builderAll); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if err := pruneAutoscalerProviderFiles(builder, "builder_", map[string]bool{
		"builder_all.go":        true,
		"builder_clusterapi.go": true,
	}); err != nil {
		return false, err
	}

	content, err := os.ReadFile(builderAll)
	if err != nil {
		return false, err
	}
	lines := strings.Split(strings.TrimSuffix(string(content), "\n"), "\n")
	transformed := make([]string, 0, len(lines))
	inImport := false
	inProviders := false
	inSwitch := false
	keepCase := true
	switchDepth := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "import (" {
			inImport = true
			transformed = append(transformed, line)
			continue
		}
		if inImport && trimmed == ")" {
			inImport = false
			transformed = append(transformed, line)
			continue
		}
		if inImport {
			if strings.Contains(trimmed, "/cloudprovider/") && !strings.Contains(trimmed, `/clusterapi"`) {
				continue
			}
			transformed = append(transformed, line)
			continue
		}
		if strings.Contains(line, "var AvailableCloudProviders") {
			inProviders = true
			transformed = append(transformed, line)
			continue
		}
		if inProviders {
			if trimmed == "}" {
				inProviders = false
				transformed = append(transformed, line)
			} else if strings.Contains(line, "ClusterAPIProviderName") || !strings.Contains(line, "ProviderName") {
				transformed = append(transformed, line)
			}
			continue
		}
		if strings.HasPrefix(trimmed, "const DefaultCloudProvider =") {
			transformed = append(transformed, leadingWhitespace(line)+"const DefaultCloudProvider = cloudprovider.ClusterAPIProviderName")
			continue
		}
		if strings.Contains(line, "switch opts.CloudProviderName") {
			inSwitch = true
			switchDepth = strings.Count(line, "{") - strings.Count(line, "}")
			keepCase = true
			transformed = append(transformed, line)
			continue
		}
		if inSwitch {
			if strings.HasPrefix(trimmed, "case cloudprovider.") {
				keepCase = strings.Contains(line, "ClusterAPIProviderName")
			}
			switchDepth += strings.Count(line, "{") - strings.Count(line, "}")
			if switchDepth <= 0 {
				inSwitch = false
				transformed = append(transformed, line)
				continue
			}
			if keepCase {
				transformed = append(transformed, line)
			}
			continue
		}
		transformed = append(transformed, line)
	}
	return true, os.WriteFile(builderAll, []byte(strings.Join(transformed, "\n")+"\n"), 0o644)
}

func pruneAutoscalerProviderFiles(directory, prefix string, keep map[string]bool) error {
	matches, err := filepath.Glob(filepath.Join(directory, prefix+"*.go"))
	if err != nil {
		return err
	}
	for _, path := range matches {
		if !keep[filepath.Base(path)] {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
	}
	return nil
}

func removeAutoscalerGCEDependencies(source string) (bool, error) {
	changed := false
	root := filepath.Join(source, "cluster-autoscaler")
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		original := string(content)
		if !strings.Contains(original, "LocalSSDDiskSizeProvider") && !strings.Contains(original, "localssdsize") {
			return nil
		}
		updated := strings.ReplaceAll(original, "opts.GCEOptions.LocalSSDDiskSizeProvider", "nil")
		lines := strings.Split(strings.TrimSuffix(updated, "\n"), "\n")
		filtered := lines[:0]
		for _, line := range lines {
			if !strings.Contains(line, "localssdsize") && !strings.Contains(line, "LocalSSDDiskSizeProvider") {
				filtered = append(filtered, line)
			}
		}
		updated = strings.Join(filtered, "\n") + "\n"
		if updated != original {
			if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
				return err
			}
			changed = true
		}
		return nil
	})
	return changed, err
}

func removeAutoscalerCloudProviderDirectories(source string) error {
	cloudProvider := filepath.Join(source, "cluster-autoscaler", "cloudprovider")
	entries, err := os.ReadDir(cloudProvider)
	if err != nil {
		return err
	}
	keep := map[string]bool{
		"builder": true, "router": true, "clusterapi": true, "mocks": true, "test": true,
	}
	for _, entry := range entries {
		if entry.IsDir() && !keep[entry.Name()] {
			if err := os.RemoveAll(filepath.Join(cloudProvider, entry.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}

func parseAutoscalerPatchMetadata(ctx context.Context, path string) (autoscalerPatchMetadata, error) {
	patch, err := os.ReadFile(path)
	if err != nil {
		return autoscalerPatchMetadata{}, err
	}
	tempDir, err := os.MkdirTemp("", "autoscaler-patch-mailinfo-")
	if err != nil {
		return autoscalerPatchMetadata{}, err
	}
	defer os.RemoveAll(tempDir)

	messagePath := filepath.Join(tempDir, "message")
	diffPath := filepath.Join(tempDir, "patch")
	output, err := runAutoscalerCommand(ctx, tempDir, bytes.NewReader(patch), nil, "git", "mailinfo", messagePath, diffPath)
	if err != nil {
		return autoscalerPatchMetadata{}, fmt.Errorf("parsing patch metadata in %q: %v", path, err)
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
		return autoscalerPatchMetadata{}, err
	}
	metadata := autoscalerPatchMetadata{
		authorName: fields["Author"], authorEmail: fields["Email"], authorDate: fields["Date"],
		subject: fields["Subject"], body: strings.TrimSpace(string(body)),
	}
	if metadata.authorName == "" || metadata.authorEmail == "" || metadata.authorDate == "" || metadata.subject == "" {
		return autoscalerPatchMetadata{}, fmt.Errorf("incomplete patch metadata in %q", path)
	}
	return metadata, nil
}

func commitAndFormatAutoscalerPatch(
	ctx context.Context,
	repo string,
	metadata autoscalerPatchMetadata,
	outputPath string,
	paths []string,
) (bool, error) {
	addArgs := []string{"add", "--all"}
	statusArgs := []string{"status", "--porcelain"}
	if paths != nil {
		addArgs = append([]string{"add", "--"}, paths...)
		statusArgs = append(statusArgs, "--")
		statusArgs = append(statusArgs, paths...)
	}
	if _, err := runAutoscalerCommand(ctx, repo, nil, nil, "git", addArgs...); err != nil {
		return false, err
	}
	status, err := runAutoscalerCommand(ctx, repo, nil, nil, "git", statusArgs...)
	if err != nil {
		return false, err
	}
	if strings.TrimSpace(status) == "" {
		return false, nil
	}

	message := metadata.subject
	if metadata.body != "" {
		message += "\n\n" + metadata.body
	}
	messagePath := filepath.Join(repo, ".git", "patch-fixer-message")
	if err := os.WriteFile(messagePath, []byte(message+"\n"), 0o600); err != nil {
		return false, err
	}
	extraEnv := map[string]string{
		"GIT_AUTHOR_NAME": metadata.authorName, "GIT_AUTHOR_EMAIL": metadata.authorEmail,
		"GIT_AUTHOR_DATE": metadata.authorDate, "GIT_COMMITTER_NAME": "EKS Distro PR Bot",
		"GIT_COMMITTER_EMAIL": "aws-model-rocket-bots+eksdistroprbot@amazon.com",
	}
	if _, err := runAutoscalerCommand(ctx, repo, nil, extraEnv, "git", "commit", "--file", messagePath); err != nil {
		return false, err
	}
	patch, err := runAutoscalerCommand(ctx, repo, nil, nil, "git", "format-patch", "-1", "--stdout", "--no-signature")
	if err != nil {
		return false, err
	}
	return true, os.WriteFile(outputPath, []byte(patch), 0o600)
}

func restoreAutoscalerSource(ctx context.Context, source, revision string) error {
	_, _ = runAutoscalerCommand(ctx, source, nil, nil, "git", "am", "--abort")
	if _, err := runAutoscalerCommand(ctx, source, nil, nil, "git", "checkout", "--force", revision); err != nil {
		return err
	}
	_, err := runAutoscalerCommand(ctx, source, nil, nil, "git", "clean", "-fdx")
	return err
}

func autoscalerPatchFiles(directory string) ([]string, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	var patches []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".patch") {
			patches = append(patches, filepath.Join(directory, entry.Name()))
		}
	}
	sort.Strings(patches)
	return patches, nil
}

func autoscalerGoBinary(releaseBranch, patchesDir string) string {
	projectRoot := filepath.Dir(filepath.Dir(patchesDir))
	version, err := os.ReadFile(filepath.Join(projectRoot, releaseBranch, "GOLANG_VERSION"))
	if err == nil {
		versionedGo := filepath.Join("/go", "go"+strings.TrimSpace(string(version)), "bin", "go")
		if _, err := os.Stat(versionedGo); err == nil {
			return versionedGo
		}
	}
	return "go"
}

func copyAutoscalerPatch(source, destination string) error {
	content, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	info, err := os.Stat(source)
	if err != nil {
		return err
	}
	return os.WriteFile(destination, content, info.Mode().Perm())
}

func leadingWhitespace(value string) string {
	return value[:len(value)-len(strings.TrimLeft(value, " \t"))]
}

func runAutoscalerCommand(
	ctx context.Context,
	directory string,
	stdin io.Reader,
	extraEnv map[string]string,
	name string,
	args ...string,
) (string, error) {
	command := exec.CommandContext(ctx, name, args...)
	command.Dir = directory
	command.Stdin = stdin
	command.Env = autoscalerCommandEnvironment(extraEnv)
	var output bytes.Buffer
	command.Stdout = &output
	command.Stderr = &output
	if err := command.Run(); err != nil {
		return output.String(), fmt.Errorf(
			"command failed (%s %s): %v\n%s",
			name,
			strings.Join(args, " "),
			err,
			strings.TrimSpace(output.String()),
		)
	}
	return output.String(), nil
}

func autoscalerCommandEnvironment(extra map[string]string) []string {
	if len(extra) == 0 {
		return os.Environ()
	}
	environment := make([]string, 0, len(os.Environ())+len(extra))
	for _, entry := range os.Environ() {
		name, _, found := strings.Cut(entry, "=")
		if found {
			if _, replaced := extra[name]; replaced {
				continue
			}
		}
		environment = append(environment, entry)
	}
	for name, value := range extra {
		environment = append(environment, name+"="+value)
	}
	return environment
}
