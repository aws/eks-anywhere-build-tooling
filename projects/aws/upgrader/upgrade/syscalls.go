//go:generate ../hack/tools/bin/mockgen -destination ./mocks/syscalls_mock.go -package mocks . SysCalls

package upgrade

import (
	"context"
	"io/fs"
	"os"
	"os/exec"
)

type SysCalls interface {
	WriteFile(string, []byte, os.FileMode) error
	ReadFile(string) ([]byte, error)
	OpenFile(string, int, os.FileMode) (*os.File, error)
	Stat(string) (os.FileInfo, error)
	Executable() (string, error)
	ExecCommand(context.Context, string, ...string) ([]byte, error)
}

type sysCalls struct {
	SysCalls
}

func (s *sysCalls) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (s *sysCalls) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (s *sysCalls) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (s *sysCalls) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

func (s *sysCalls) Executable() (string, error) {
	return os.Executable()
}

func (s *sysCalls) ExecCommand(ctx context.Context, name string, arg ...string) ([]byte, error) {
	return exec.CommandContext(ctx, name, arg...).CombinedOutput()
}

func NewSysCalls() SysCalls {
	return &sysCalls{}
}
