package upgrade

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizePatchNumbersPreservesMetadataAndOrdersRemainingSeries(t *testing.T) {
	directory := t.TempDir()
	for name, subject := range map[string]string{
		"0001-first.patch": "[PATCH] First change",
		"0003-last.patch":  "[PATCH 3/3] Last change\n continued title",
	} {
		content := "From deadbeef Mon Sep 17 00:00:00 2001\n" +
			"From: Test <test@example.com>\nDate: Thu, 1 Jan 2026 00:00:00 +0000\n" +
			"Subject: " + subject + "\n\nMessage body.\nSubject: [PATCH 7/8] Body example\n\n" +
			"Signed-off-by: Test <test@example.com>\n"
		writeRepairTestFile(t, filepath.Join(directory, name), []byte(content), 0o640)
	}
	changed, err := normalizePatchNumbers(directory)
	if err != nil || len(changed) != 2 {
		t.Fatalf("normalizePatchNumbers() = %v, error %v", changed, err)
	}
	for name, prefix := range map[string]string{
		"0001-first.patch": "[PATCH 1/2] First change",
		"0003-last.patch":  "[PATCH 2/2] Last change\n continued title",
	} {
		content, err := os.ReadFile(filepath.Join(directory, name))
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range []string{
			"From: Test <test@example.com>", "Date: Thu, 1 Jan 2026 00:00:00 +0000",
			"Subject: " + prefix, "\n\nMessage body.\nSubject: [PATCH 7/8] Body example\n",
			"Signed-off-by: Test <test@example.com>",
		} {
			if !strings.Contains(string(content), want) {
				t.Fatalf("numbering changed metadata/body: %q, missing %q", content, want)
			}
		}
		info, err := os.Stat(filepath.Join(directory, name))
		if err != nil || info.Mode().Perm() != 0o640 {
			t.Fatalf("patch mode changed: %v, error %v", info, err)
		}
	}
	if changed, err := normalizePatchNumbers(directory); err != nil || len(changed) != 0 {
		t.Fatalf("second normalization changed files: %v, error %v", changed, err)
	}
}

func TestNormalizePatchNumbersUsesSinglePatchConvention(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "0004-only.patch")
	if err := os.WriteFile(path, []byte("From: Test <test@example.com>\nSubject: [PATCH 4/4] Change\n\nBody\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := normalizePatchNumbers(directory); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(content), "Subject: [PATCH] Change\n") {
		t.Fatalf("single patch header = %q, error %v", content, err)
	}
}

func TestNormalizePatchNumbersPreservesUnsupportedHeaders(t *testing.T) {
	if paths, err := normalizePatchNumbers(filepath.Join(t.TempDir(), "removed")); err != nil || len(paths) != 0 {
		t.Fatalf("removed series = %v, error %v", paths, err)
	}
	directory := t.TempDir()
	for _, content := range []string{
		"From: Test <test@example.com>\n\nBody\n",
		"From: Test <test@example.com>\nSubject: [PATCH v2 1/3] Change\n\nBody\n",
		"Subject: [PATCH] No mail boundary",
	} {
		path := filepath.Join(directory, "0001-unusual.patch")
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		changed := finishPatchNumbering(directory, []string{path})
		if len(changed) != 1 {
			t.Fatalf("cosmetic format blocked accepted repair: %v", changed)
		}
		stored, err := os.ReadFile(path)
		if err != nil || string(stored) != content {
			t.Fatalf("unsupported header altered: %q, error %v", stored, err)
		}
	}
}

func TestFinishPatchNumberingDoesNotRejectValidatedRepairOnFormattingError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(path, []byte("dummy"), 0o644); err != nil {
		t.Fatal(err)
	}
	changed := finishPatchNumbering(path, []string{"accepted.patch"})
	if len(changed) != 1 || changed[0] != "accepted.patch" {
		t.Fatalf("numbering rejected accepted repair: %v", changed)
	}
}
