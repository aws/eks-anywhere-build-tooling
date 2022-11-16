package files

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func Write(path string, content []byte, permission fs.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o640); err != nil {
		return errors.Wrap(err, "Error creating directory")
	}
	if err := ioutil.WriteFile(path, content, permission); err != nil {
		return errors.Wrapf(err, "Error writing file: %s", path)
	}
	return nil
}
