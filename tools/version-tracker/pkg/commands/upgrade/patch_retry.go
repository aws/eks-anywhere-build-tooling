package upgrade

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/commands/agentpatchfixer"
)

func prepareSourceValidationRetry(
	ctx context.Context,
	repoPath, patchesDir, baseRevision string,
	originalPatches map[string][]byte,
	changedPaths []string,
	sourceErr *repairSourceValidationError,
) (applied, total int, returnErr error) {
	if baseRevision == "" {
		return 0, 0, fmt.Errorf("cannot retry source validation without the upstream base revision")
	}
	var retryPath string
	for _, path := range slices.Backward(changedPaths) {
		content, found := originalPatches[filepath.Base(path)]
		if !found {
			continue
		}
		files, err := agentpatchfixer.PatchFiles(ctx, string(content))
		if err != nil {
			return 0, 0, fmt.Errorf("reading original repair intent: %w", err)
		}
		if slices.Contains(files, filepath.ToSlash(sourceErr.Path)) {
			retryPath = path
			break
		}
	}
	if retryPath == "" {
		return 0, 0, fmt.Errorf("no repaired patch owns the invalid source file %q", sourceErr.Path)
	}
	if err := abortPatchApplication(repoPath); err != nil {
		return 0, 0, err
	}
	if _, err := repairGitOutput(ctx, repoPath, "reset", "--hard", baseRevision); err != nil {
		return 0, 0, err
	}
	if err := os.MkdirAll(patchesDir, 0o755); err != nil {
		return 0, 0, err
	}
	if err := writeFileAtomically(retryPath, originalPatches[filepath.Base(retryPath)], 0o644); err != nil {
		return 0, 0, err
	}
	paths, err := patchFilesInDirectory(patchesDir)
	if err != nil {
		return 0, 0, err
	}
	index := slices.Index(paths, retryPath)
	if index < 0 {
		return 0, 0, fmt.Errorf("restored repair patch is absent from its series")
	}
	if index > 0 {
		args := append([]string{"am", "--committer-date-is-author-date"}, paths[:index]...)
		if _, err := repairGitOutput(ctx, repoPath, args...); err != nil {
			return 0, 0, fmt.Errorf("replaying patches before source retry: %w", err)
		}
	}
	return index, len(paths), nil
}

// These targets depend only on checkout stamps, not patch file timestamps.
// Invalidate them so every candidate is replayed from its upstream revision.
func invalidatePatchCheckoutMarkers(repoPath string) error {
	entries, err := os.ReadDir(repoPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading upstream checkout markers: %w", err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if name != "eks-anywhere-patched" && !strings.HasPrefix(name, "eks-anywhere-checkout-") {
			continue
		}
		if entry.IsDir() {
			return fmt.Errorf("upstream checkout marker %q is a directory", name)
		}
		if err := os.Remove(filepath.Join(repoPath, name)); err != nil {
			return fmt.Errorf("invalidating upstream checkout marker: %w", err)
		}
	}
	return nil
}
