package git

import (
	"os"
	"path/filepath"
	"testing"

	gogit "github.com/go-git/go-git/v5"
)

func TestAddStagesDeletedFile(t *testing.T) {
	repoPath := t.TempDir()
	repo, err := gogit.PlainInit(repoPath, false)
	if err != nil {
		t.Fatal(err)
	}
	worktree, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(repoPath, "obsolete.patch")
	if err := os.WriteFile(path, []byte("obsolete\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := worktree.Add("obsolete.patch"); err != nil {
		t.Fatal(err)
	}
	if _, err := worktree.Commit("add patch", &gogit.CommitOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}

	if err := Add(worktree, []string{"obsolete.patch"}); err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	status, err := worktree.Status()
	if err != nil {
		t.Fatal(err)
	}
	if status.File("obsolete.patch").Staging != gogit.Deleted {
		t.Fatalf("staging status = %v, want deleted", status.File("obsolete.patch").Staging)
	}
}
