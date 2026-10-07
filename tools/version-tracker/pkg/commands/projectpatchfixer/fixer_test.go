package projectpatchfixer

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRunReturnsUnhandledForGenericProject(t *testing.T) {
	result, err := Run(context.Background(), Request{ProjectName: "example/project"})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Handled {
		t.Fatal("Run() handled = true, want false")
	}
}

func TestRunSkipsProjectsExcludedFromGenericAgentRepair(t *testing.T) {
	for _, projectName := range []string{
		"goharbor/harbor",
		"kubernetes-sigs/cluster-api",
	} {
		t.Run(projectName, func(t *testing.T) {
			result, err := Run(context.Background(), Request{ProjectName: projectName})
			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}
			if !result.Handled || !result.SkipAgent || result.Reason == "" {
				t.Fatalf("Run() result = %#v, want handled agent skip with reason", result)
			}
		})
	}
}

func TestRunRoutesAutoscalerWithoutAgent(t *testing.T) {
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	patchesDir := filepath.Join(tempDir, "existing")
	for _, directory := range []string{sourceDir, patchesDir} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	original := regenerateAutoscaler
	regenerateAutoscaler = func(_ context.Context, _ Request, outputDir string) ([]string, error) {
		patch := filepath.Join(outputDir, "0001-generated.patch")
		if err := os.WriteFile(patch, []byte("patch\n"), 0o600); err != nil {
			return nil, err
		}
		return []string{patch}, nil
	}
	t.Cleanup(func() {
		regenerateAutoscaler = original
	})
	t.Setenv("PATCH_FIXER_RESULT_DIR", filepath.Join(tempDir, "results"))
	result, err := Run(context.Background(), Request{
		ProjectName:    "kubernetes/autoscaler",
		UpstreamRepo:   sourceDir,
		PatchesDir:     patchesDir,
		TargetRevision: "cluster-autoscaler-1.36.0",
		ReleaseBranch:  "1-36",
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !result.Handled || result.PatchCount != 1 {
		t.Fatalf("Run() result = %#v", result)
	}
}

func TestRunCancelsAutoscalerCommandWithContext(t *testing.T) {
	tempDir := t.TempDir()
	original := regenerateAutoscaler
	regenerateAutoscaler = func(ctx context.Context, _ Request, _ string) ([]string, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}
	t.Cleanup(func() {
		regenerateAutoscaler = original
	})
	t.Setenv("PATCH_FIXER_RESULT_DIR", filepath.Join(tempDir, "results"))
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := Run(ctx, Request{
		ProjectName:    "kubernetes/autoscaler",
		UpstreamRepo:   filepath.Join(tempDir, "source"),
		PatchesDir:     filepath.Join(tempDir, "patches"),
		TargetRevision: "example-revision",
		ReleaseBranch:  "example-branch",
	})
	if err == nil {
		t.Fatal("Run() error = nil, want context cancellation")
	}
	if time.Since(start) > time.Second {
		t.Fatalf("cancellation took too long: %s", time.Since(start))
	}
}
