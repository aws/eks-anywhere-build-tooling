package writer

import "os"

type Writer struct {
	filename string
	File     *os.File
}

func New(filename string) (*Writer, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return &Writer{
		File:     file,
		filename: filename}, nil
}

func (w *Writer) WriteLine(content string) error {
	_, err := w.File.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}