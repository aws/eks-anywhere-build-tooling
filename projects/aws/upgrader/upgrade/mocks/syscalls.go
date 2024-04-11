//go:generate ../../hack/tools/bin/mockgen -destination ./syscalls_mock.go -package mocks . SysCalls
package mocks

import (
	context "context"
	os "os"
)

// SysCalls provides a list of system calls for a host.
// Mock implementations can aid in testing without making an underlying os call.
type SysCalls interface {
	WriteFile(string, []byte, os.FileMode) error
	ReadFile(string) ([]byte, error)
	OpenFile(string, int, os.FileMode) (*os.File, error)
	Stat(string) (os.FileInfo, error)
	Executable() (string, error)
	ExecCommand(context.Context, string, ...string) ([]byte, error)
	MkdirAll(string, os.FileMode) error
}
