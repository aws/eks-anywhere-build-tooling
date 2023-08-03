package tar

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ExtractTarGz extracts the contents of the given tarball.
func ExtractTarGz(tarballDownloadPath string, gzipStream io.Reader) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)
	var header *tar.Header
	for header, err = tarReader.Next(); err == nil; header, err = tarReader.Next() {
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(filepath.Join(tarballDownloadPath, header.Name), 0o755); err != nil {
				return fmt.Errorf("creating directory from archive: %v", err)
			}
		case tar.TypeReg:
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
		default:
			return fmt.Errorf("unknown type in tar header: %b in %s", header.Typeflag, header.Name)
		}
	}
	if err != io.EOF {
		return fmt.Errorf("advancing to next entry in archive: %v", err)
	}
	return nil
}
