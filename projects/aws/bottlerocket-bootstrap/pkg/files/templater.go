package files

import (
	"bytes"
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
