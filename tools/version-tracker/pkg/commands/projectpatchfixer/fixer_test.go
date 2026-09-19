package projectpatchfixer

import (
	"context"
	"os"
	"path/filepath"
	"testing"
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
	commandPath := filepath.Join(tempDir, "regenerate.sh")
	command := `#!/bin/sh
while [ "$#" -gt 0 ]; do
  case "$1" in
    --output) output="$2"; shift 2 ;;
    --manifest) manifest="$2"; shift 2 ;;
    *) shift ;;
  esac
done
patch="$output/0001-generated.patch"
printf 'patch\n' > "$patch"
printf '{"patches":["%s"]}\n' "$patch" > "$manifest"
`
	if err := os.WriteFile(commandPath, []byte(command), 0o755); err != nil {
		t.Fatal(err)
	}

	t.Setenv(autoscalerCommandEnv, commandPath)
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
