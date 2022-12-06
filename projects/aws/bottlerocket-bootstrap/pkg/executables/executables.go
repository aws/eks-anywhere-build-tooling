//go:generate ../../hack/tools/bin/mockgen -destination ./mocks/executables_mock.go -package mocks . Executable

package executables

import (
	"os/exec"

	"github.com/pkg/errors"
)

type Executable interface {
	Execute(args ...string) (out []byte, err error)
}

type executable struct {
	binary string
}

func NewExecutable(binary string) Executable {
	return &executable{
		binary: binary,
	}
}

func (e *executable) Execute(args ...string) ([]byte, error) {
	cmd := exec.Command(e.binary, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrapf(err, "Error running command: %v, output: %s\n", cmd, string(out))
	}
	return out, nil
}
