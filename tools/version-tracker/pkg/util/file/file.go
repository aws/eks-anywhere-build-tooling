package file

import (
	"io"
	"net/http"
	"os"
)

// Download downloads a file from the given URL to the destination filepath.
func Download(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
