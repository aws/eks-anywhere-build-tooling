package files

import (
	"bytes"
	"io/fs"
	"text/template"

	"github.com/pkg/errors"
)

func ExecuteTemplate(content string, data interface{}) ([]byte, error) {
	temp := template.New("tmpl")
	temp, err := temp.Parse(content)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing template")
	}

	var buf bytes.Buffer
	err = temp.Execute(&buf, data)
	if err != nil {
		return nil, errors.Wrap(err, "Error substituting values for template")
	}
	return buf.Bytes(), nil
}

func WriteTemplate(path, content string, data interface{}, permission fs.FileMode) error {
	b, err := ExecuteTemplate(content, data)
	if err != nil {
		return errors.Wrap(err, "Error executing template")
	}
	return Write(path, b, permission)
}
