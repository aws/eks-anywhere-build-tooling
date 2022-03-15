package utils

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"os/exec"

	"github.com/pkg/errors"
)

const (
	ApiclientBinary        = "apiclient"
	BootstrapContainerName = "kubeadm-bootstrap"
)

func DisableBootstrapContainer() error {
	cmd := exec.Command(ApiclientBinary, "set", "host-containers."+BootstrapContainerName+".enabled=false")
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "Error disabling bootstrap container")
	}
	return nil
}

// Use Gzip to decompress
func GUnzipBytes(data []byte) ([]byte, error) {
	// Write gzipped data to the client
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer gr.Close()
	uncompressedData, err := ioutil.ReadAll(gr)
	if err != nil {
		return nil, err
	}
	return uncompressedData, nil
}

// Use Gzip to compress
func GzipBytes(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(data); err != nil {
		return []byte{}, errors.Wrap(err, "failed to gzip bytes")
	}

	if err := gz.Close(); err != nil {
		return []byte{}, errors.Wrap(err, "failed to gzip bytes")
	}

	return buf.Bytes(), nil
}
