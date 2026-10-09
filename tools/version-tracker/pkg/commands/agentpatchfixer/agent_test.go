package agentpatchfixer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecuteKeepsConciseRepairSummary(t *testing.T) {
	tests := []struct {
		name    string
		edit    bool
		status  string
		summary string
	}{
		{
			name:    "repaired patch",
			edit:    true,
			status:  "candidate_generated",
			summary: "Change: Preserve useful behavior on the new release.\nReview: Inspect retry behavior.",
		},
		{
			name:    "obsolete patch",
			status:  "no_changes",
			summary: "Change: Removed obsolete patch; its intent is already upstream.\nReview: Original intent retained.",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := initializeRepository(t, "current\n")
			patchPath := filepath.Join(t.TempDir(), "failed.patch")
			if err := os.WriteFile(patchPath, []byte(testPatch), 0o644); err != nil {
				t.Fatal(err)
			}
			runner := fakeAgentRunner(t, "end_turn", test.edit)
			runnerPath := strings.TrimPrefix(runner, "python3 ")
			script, err := os.ReadFile(runnerPath)
			if err != nil {
				t.Fatal(err)
			}
			scriptText := strings.Replace(string(script), `"summary": "done"`, fmt.Sprintf(`"summary": %q`, test.summary), 1)
			if err := os.WriteFile(runnerPath, []byte(scriptText), 0o755); err != nil {
				t.Fatal(err)
			}
			t.Setenv(AgentCommandEnv, runner)
			runDir := t.TempDir()
			result, err := execute(context.Background(), Request{
				SchemaVersion:    requestSchemaVersion,
				ProjectName:      "example/project",
				UpstreamRepoPath: source,
				FailedPatchPath:  patchPath,
				MaxTurns:         4,
				MaxTotalTokens:   1000,
			}, runDir)
			if err != nil {
				t.Fatal(err)
			}
			if result.SchemaVersion != requestSchemaVersion || result.Status != test.status || result.Summary != test.summary {
				t.Fatalf("result = %#v, want %s with unchanged summary", result, test.status)
			}
			invocationBytes, err := os.ReadFile(filepath.Join(runDir, "agent-request.json"))
			if err != nil {
				t.Fatal(err)
			}
			var invocation agentInvocation
			if err := json.Unmarshal(invocationBytes, &invocation); err != nil {
				t.Fatal(err)
			}
			if invocation.SystemPrompt != systemPrompt {
				t.Fatalf("invocation prompt = %q, want concise review instructions", invocation.SystemPrompt)
			}
		})
	}
}
