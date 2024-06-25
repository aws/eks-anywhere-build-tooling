package tar

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractFileFromTarball extracts the specified file from the given tarball.
func ExtractFileFromTarball(tarballDownloadPath string, gzipStream io.Reader, targetFile string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)
	var header *tar.Header
	for header, err = tarReader.Next(); err == nil; header, err = tarReader.Next() {
		if header.Name == targetFile {
			if strings.Contains(header.Name, "/") {
				err = os.MkdirAll(filepath.Join(tarballDownloadPath, filepath.Dir(header.Name)), 0o755)
				if err != nil {
					return fmt.Errorf("creating parent directory for archive contents: %v", err)
				}
			}
			outFile, err := os.Create(filepath.Join(tarballDownloadPath, header.Name))
			if err != nil {
				return fmt.Errorf("creating file from archive: %v", err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("copying file to output destination: %v", err)
			}
			if err := outFile.Close(); err != nil {
				return fmt.Errorf("closing output destination file descriptor: %v", err)
			}
		}
	}
	if err != io.EOF {
		return fmt.Errorf("advancing to next entry in archive: %v", err)
	}
	return nil
}
