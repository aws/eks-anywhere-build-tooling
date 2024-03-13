//go:generate ../hack/tools/bin/mockgen -destination ./mocks/syscalls_mock.go -package mocks . SysCalls

package upgrade

import (
	"context"
	"os"
	"os/exec"
)

// SysCalls serves different os level calls on a node.
type SysCalls struct {
	WriteFile   func(string, []byte, os.FileMode) error
	ReadFile    func(string) ([]byte, error)
	OpenFile    func(string, int, os.FileMode) (*os.File, error)
	Stat        func(string) (os.FileInfo, error)
	Executable  func() (string, error)
	ExecCommand func(context.Context, string, ...string) ([]byte, error)
	MkdirAll    func(string, os.FileMode) error
}

func ExecCommand(ctx context.Context, name string, arg ...string) ([]byte, error) {
	return exec.CommandContext(ctx, name, arg...).CombinedOutput()
}

func NewSysCalls() SysCalls {
	return SysCalls{
		WriteFile:   os.WriteFile,
		ReadFile:    os.ReadFile,
		OpenFile:    os.OpenFile,
		Stat:        os.Stat,
		Executable:  os.Executable,
		ExecCommand: ExecCommand,
		MkdirAll:    os.MkdirAll,
	}
}
