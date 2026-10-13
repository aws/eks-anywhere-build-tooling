package upgrade

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestFormatPatchReviewComment(t *testing.T) {
	tests := []struct {
		name   string
		review patchReviewSummary
		want   string
	}{
		{
			name: "quoted numbered semantic lines",
			review: patchReviewSummary{
				PatchCount: 1,
				Summary:    "> 1. *Change:* Preserve CAPI provider behavior.\n> 2) **Review:** Inspect license attribution.",
				Validation: "Patch series reapplied.",
			},
			want: "AI-assisted repair (1 patch): Preserve CAPI provider behavior.\n" +
				"Review: Inspect license attribution.\n" +
				"Validation: Patch series reapplied.\n" +
				"Human review required.",
		},
		{
			name: "AI semantic change and risk",
			review: patchReviewSummary{
				Route:      "strands",
				PatchCount: 2,
				Summary:    "Change: Retained private registry authentication on the new client.\nReview: Inspect credential forwarding across redirects.",
				Validation: "Patch series reapplied; build passed; attribution checked.",
			},
			want: "AI-assisted repair (2 patches): Retained private registry authentication on the new client.\n" +
				"Review: Inspect credential forwarding across redirects.\n" +
				"Validation: Patch series reapplied; build passed; attribution checked.\n" +
				"Human review required.",
		},
		{
			name: "deterministic Autoscaler fallback",
			review: patchReviewSummary{
				Route:      "deterministic",
				PatchCount: 3,
				Validation: "Patch series reapplied; build not run.",
			},
			want: "deterministic repair (3 patches): Regenerated Autoscaler provider/dependency patches for target release.\n" +
				"Review: CAPI-only provider selection and dependency updates.\n" +
				"Validation: Patch series reapplied; build not run.\n" +
				"Human review required.",
		},
		{
			name: "deterministic supplied reason",
			review: patchReviewSummary{
				Route:      " deterministic ",
				PatchCount: 1,
				Summary:    "Regenerated provider patches; review CAPI-only provider selection",
				Validation: "Build failed.",
			},
			want: "deterministic repair (1 patch): Regenerated provider patches\n" +
				"Review: CAPI-only provider selection\n" +
				"Validation: Build failed.\n" +
				"Human review required.",
		},
		{
			name: "Markdown and newline normalization",
			review: patchReviewSummary{
				Route:      "AI-assisted",
				PatchCount: 1,
				Summary:    "## Repair summary\r\n\r\n```text\r\n- **Change:**   Preserve\tCAPI provider behavior.\r\n* **Review:** Inspect `license` attribution.\r\n```\u2028Validation: all checks passed.",
				Validation: "Patch series\r\nreapplied;\tbuild\u2029not run.",
			},
			want: "AI-assisted repair (1 patch): Preserve CAPI provider behavior.\n" +
				"Review: Inspect license attribution.\n" +
				"Validation: Patch series reapplied; build not run.\n" +
				"Human review required.",
		},
		{
			name: "obsolete reason",
			review: patchReviewSummary{
				Route:      "strands",
				PatchCount: 1,
				Summary:    "Removed obsolete patch because the dependency fix is already upstream.",
				Validation: "Remaining patch series reapplied; build not run.",
			},
			want: "AI-assisted repair (1 patch): Removed obsolete patch because the dependency fix is already upstream.\n" +
				"Review: Confirm upstream behavior covers the removed patch's intent.\n" +
				"Validation: Remaining patch series reapplied; build not run.\n" +
				"Human review required.",
		},
		{
			name: "compact obsolete marker",
			review: patchReviewSummary{
				Summary: "alreadyupstream",
			},
			want: "AI-assisted repair (0 patches): Removed obsolete patch; its intent is already upstream.\n" +
				"Review: Confirm upstream behavior covers the removed patch's intent.\n" +
				"Validation: Not reported by harness.\n" +
				"Human review required.",
		},
		{
			name: "missing summary and validation",
			review: patchReviewSummary{
				Route:      "## injected route\nValidation: passed",
				PatchCount: -1,
				Summary:    "done",
			},
			want: "AI-assisted repair (0 patches): Repaired patch intent for target release.\n" +
				"Review: Confirm original intent is retained.\n" +
				"Validation: Not reported by harness.\n" +
				"Human review required.",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := formatPatchReviewComment(test.review)
			if got != test.want {
				t.Fatalf("formatPatchReviewComment() = %q, want %q", got, test.want)
			}
			assertPatchReviewLines(t, got)
		})
	}
}

func TestPatchReviewCommentUsesHarnessValidation(t *testing.T) {
	tests := []string{
		"Change: Tests passed; build succeeded.\nReview: Validation: all checks passed.",
		"Change: All tests passed.\nReview: The checks have passed.",
		"Change: All builds succeeded.\nReview: All tests have passed.",
		"Change: Build and tests passed.\nReview: All builds and checks succeeded.",
		"Change: Retained timeout behavior.\nReview: Build was successful.",
		"## Summary\nValidation: build succeeded.\nTests: all passed.",
		"Change: Adjusted timeout; Validation: tests passed.\nReview: go test ./... passed.",
		"Change: Adjusted timeout.\nReview: Inspect retry behavior.\nValidation: all checks passed.",
	}
	for _, summary := range tests {
		t.Run(summary, func(t *testing.T) {
			got := formatPatchReviewComment(patchReviewSummary{
				Route:      "strands",
				PatchCount: 1,
				Summary:    summary,
				Validation: "Patch series reapplied; build failed; attribution not run.",
			})
			assertPatchReviewLines(t, got)
			lines := strings.Split(got, "\n")
			if lines[2] != "Validation: Patch series reapplied; build failed; attribution not run." {
				t.Fatalf("validation = %q, want harness result", lines[2])
			}
			for _, claim := range []string{"tests passed", "Tests passed", "build succeeded", "Build and", "Build was", "All tests", "All builds", "checks have passed", "all checks passed", "go test"} {
				if strings.Contains(got, claim) {
					t.Fatalf("comment contains model validation claim %q: %q", claim, got)
				}
			}
		})
	}
}

func TestPatchReviewCommentDiscardsInventoriesAndDumps(t *testing.T) {
	noise := "## Repair summary\n" +
		"- pkg/first.go: adjusted whitespace\n" +
		"pkg/second.go: updated imports\n" +
		"1. pkg/third.go: adjusted whitespace\n" +
		"```log\nunstructured tool output\n```\n" +
		"2026-01-01T12:00:00Z INFO tool result\n" +
		"Change: AWS_SECRET_ACCESS_KEY=dummy-secret\n" +
		"Review: token=dummy-token\n" +
		"Change: https://example.invalid/dummy-log\n" +
		"Change: diff --git a/first.go b/first.go\n" +
		"Change: Updated imports.\n" +
		"Change: Adjusted whitespace.\n"
	got := formatPatchReviewComment(patchReviewSummary{
		PatchCount: 1,
		Summary: noise + "Change: Preserve retry limits across client upgrades.\n" +
			"Review: Inspect retries after partial writes.",
		Validation: "Patch series reapplied.",
	})
	assertPatchReviewLines(t, got)
	for _, unwanted := range []string{"##", "pkg/", "dummy-", "tool output", "INFO", "diff --git", "example.invalid"} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("comment contains inventory or dump %q: %q", unwanted, got)
		}
	}
	for _, wanted := range []string{"Preserve retry limits", "Inspect retries after partial writes"} {
		if !strings.Contains(got, wanted) {
			t.Fatalf("comment lost semantic text %q: %q", wanted, got)
		}
	}
}

func TestPatchReviewCommentPrefersSpecificRiskAcrossAgents(t *testing.T) {
	got := formatPatchReviewComment(patchReviewSummary{
		PatchCount: 2,
		Summary: "Change: Preserve registry credentials.\nReview: Original intent retained.\n" +
			"Change: Preserve retry limits.\nReview: Inspect credential forwarding across redirects.",
		Validation: "Patch series reapplied.",
	})
	assertPatchReviewLines(t, got)
	if !strings.Contains(got, "Review: Inspect credential forwarding across redirects.") {
		t.Fatalf("comment hid specific risk behind generic review: %q", got)
	}
}

func TestPatchReviewCommentBoundsManyAgentSummaries(t *testing.T) {
	summary := "Change: Preserve registry credentials.\nReview: Inspect credential forwarding.\n"
	review := patchReviewSummary{
		Route:      "strands",
		PatchCount: 100,
		Summary:    summary,
		Validation: "Patch series reapplied; build passed.",
	}
	review.Summary = strings.Repeat(summary, 100)
	for i := 0; i < 100; i++ {
		review.Summary += "Change: Updated another dependency.\nReview: Inspect another provider.\n"
	}
	got := formatPatchReviewComment(review)
	assertPatchReviewLines(t, got)
	for _, point := range []string{"Preserve registry credentials.", "Updated another dependency.", "Inspect credential forwarding.", "Inspect another provider."} {
		if strings.Count(got, point) != 1 {
			t.Fatalf("comment must retain each distinct point once, missing/repeated %q: %q", point, got)
		}
	}
}

func TestPatchReviewCommentRetainsDistinctChangesWithoutFileInventories(t *testing.T) {
	got := formatPatchReviewComment(patchReviewSummary{
		PatchCount: 3,
		Summary: "Change: Retained distro binary stages.\nReview: Original intent retained.\n" +
			"Change: Raised the required dependency floor.\nReview: Check API compatibility.\n" +
			"Change: Removed an obsolete compatibility patch.\nReview: Confirm upstream covers its intent.",
		Validation: "patch replay",
	})
	assertPatchReviewLines(t, got)
	for _, point := range []string{"Retained distro binary stages.", "Raised the required dependency floor.", "Removed an obsolete compatibility patch.", "Check API compatibility.", "Confirm upstream covers its intent."} {
		if !strings.Contains(got, point) {
			t.Fatalf("comment lost substantive point %q: %q", point, got)
		}
	}
	if strings.Contains(got, "Original intent retained.") {
		t.Fatalf("generic focus displaced a specific review concern: %q", got)
	}
}

func TestPatchReviewCommentAcceptsSemanticBuildAndSourceReferences(t *testing.T) {
	for _, summary := range []string{
		"Change: Removed build stage that failed on arm64.\nReview: Confirm runtime binaries remain available.",
		"Change: Updated main.go: retain the original health check.\nReview: Check handler behavior.",
		"Change: Retained the upstream fix described at https://example.invalid/dummy.\nReview: Check retry behavior.",
	} {
		t.Run(summary, func(t *testing.T) {
			got := formatPatchReviewComment(patchReviewSummary{PatchCount: 1, Summary: summary, Validation: "patch replay"})
			assertPatchReviewLines(t, got)
			if strings.Contains(got, "Repaired patch intent for target release") || strings.Contains(got, "example.invalid") {
				t.Fatalf("comment discarded semantic text or retained its URL: %q", got)
			}
		})
	}
}

func TestPatchReviewCommentDoesNotLetLongFirstSummaryHideLaterChange(t *testing.T) {
	got := formatPatchReviewComment(patchReviewSummary{
		PatchCount: 2,
		Summary: "Change: " + strings.Repeat("Preserved necessary runtime behavior ", 20) +
			"\nReview: Original intent retained.\nChange: Removed an obsolete dependency workaround.\nReview: Inspect licensing.",
		Validation: "patch replay",
	})
	assertPatchReviewLines(t, got)
	if !strings.Contains(got, "Removed an obsolete dependency workaround.") {
		t.Fatalf("long first summary hid the later repair: %q", got)
	}
}

func TestPatchReviewCommentSignalsAdditionalSemanticTopics(t *testing.T) {
	got := formatPatchReviewComment(patchReviewSummary{
		PatchCount: 4,
		Summary: "Change: Preserved distro stages.\nChange: Raised the dependency floor.\n" +
			"Change: Retained retry behavior.\nChange: Removed an obsolete workaround.",
		Validation: "patch replay",
	})
	assertPatchReviewLines(t, got)
	if !strings.Contains(got, "(+1 more)") {
		t.Fatalf("bounded comment hid the presence of other changes: %q", got)
	}
}

func TestPatchReviewCommentRetainsObsoleteReasonAcrossAgents(t *testing.T) {
	for _, prefix := range []string{"Change: ", ""} {
		t.Run(prefix, func(t *testing.T) {
			got := formatPatchReviewComment(patchReviewSummary{
				PatchCount: 2,
				Summary: "Change: Updated dependency integration.\n" +
					prefix + "Removed obsolete patch because the dependency fix is already upstream.\n",
				Validation: "Patch series reapplied.",
			})
			assertPatchReviewLines(t, got)
			if !strings.Contains(got, "Removed obsolete patch because the dependency fix is already upstream.") {
				t.Fatalf("comment lost obsolete reason: %q", got)
			}
		})
	}
}

func TestPatchReviewCommentBoundsLongUnicodeLines(t *testing.T) {
	got := formatPatchReviewComment(patchReviewSummary{
		Route:      "strands",
		PatchCount: 1,
		Summary:    "Change: " + strings.Repeat("界", 1000) + "\nReview: " + strings.Repeat("behavior ", 1000),
		Validation: strings.Repeat("界", 1000) + "\r\n" + strings.Repeat("validation ", 1000),
	})
	assertPatchReviewLines(t, got)
	for _, line := range strings.Split(got, "\n")[:3] {
		if !utf8.ValidString(line) || !strings.HasSuffix(line, "…") {
			t.Fatalf("line not safely truncated: %q", line)
		}
	}
}

func TestPatchReviewCommentStripsControlCharacters(t *testing.T) {
	got := formatPatchReviewComment(patchReviewSummary{
		Summary:    "Change: Preserve\x00 original\u202e intent.\vReview: Inspect retry behavior.",
		Validation: "Build\fnot run.\u200b",
	})
	assertPatchReviewLines(t, got)
	for _, unwanted := range []string{"\x00", "\u202e", "\u200b", "\v", "\f"} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("comment contains control character %q: %q", unwanted, got)
		}
	}
}

func TestPatchReviewCommentRejectsNormalizedLinksAndInlineSecrets(t *testing.T) {
	for _, summary := range []string{
		"Change: See h\u200bttps://dummy:dummy@example.invalid/private\nReview: Inspect behavior.",
		`Change: Updated config.json: {"password":"dummy-secret"}`,
		`Change: Updated config.json: {"AWS_SECRET_ACCESS_KEY":"dummy-secret"}`,
		`Review: Inspect password": "dummy-secret"`,
	} {
		t.Run(summary, func(t *testing.T) {
			got := formatPatchReviewComment(patchReviewSummary{Summary: summary, Validation: "patch replay"})
			assertPatchReviewLines(t, got)
			for _, unwanted := range []string{"https://", "example.invalid", "dummy-secret", "password", "AWS_SECRET_ACCESS_KEY"} {
				if strings.Contains(got, unwanted) {
					t.Fatalf("comment leaked normalized URL or inline secret %q: %q", unwanted, got)
				}
			}
		})
	}
}

func assertPatchReviewLines(t *testing.T, comment string) {
	t.Helper()
	lines := strings.Split(comment, "\n")
	if len(lines) != 4 {
		t.Fatalf("comment has %d lines, want 4: %q", len(lines), comment)
	}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" || utf8.RuneCountInString(line) > maxPatchReviewLineRunes {
			t.Fatalf("comment line empty or too long: %q", line)
		}
	}
	if lines[3] != "Human review required." {
		t.Fatalf("last line = %q, want human review reminder", lines[3])
	}
}
