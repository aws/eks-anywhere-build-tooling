package utils

import (
	"encoding/base64"
	"fmt"
	"io/fs"
	"io/ioutil"
	// "os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
)

func EnableStaticPods(path string) ([]*v1.Pod, error) {
	var podDefinitions []*v1.Pod
	fmt.Println("Enabling static pods on host")

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return podDefinitions, errors.Wrap(err, "error reading from manifest path")
	}

	for _, f := range files {
		if !isPodFile(f) {
			continue
		}

		baseFileName := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		// Read the manifest and base64 encode it
		fileData, err := ioutil.ReadFile(filepath.Join(path, f.Name()))
		if err != nil {
			return podDefinitions, errors.Wrapf(err, "Error reading file %s", f.Name())
		}
		b64Manifest := base64.StdEncoding.EncodeToString(fileData)
		fmt.Println("Enabling static pod " + baseFileName)
		fmt.Println("-------------------------------")
		fmt.Printf("Manifest string: \n%s\n", string(fileData))
		fmt.Println("-------------------------------")
		fmt.Printf("Encoded string: \n%s\n", b64Manifest)
		// cmd := exec.Command("bash", "-c", "apiclient set \"kubernetes.static-pods."+baseFileName+".manifest\"=\""+b64Manifest+"\" \"kubernetes.static-pods."+baseFileName+".enabled\"=true")
		// out, err := cmd.CombinedOutput()
		// if err != nil {
		// 	fmt.Printf("Apiclient command for static pod failed, output: %s\n", string(out))
		// 	return podDefinitions, errors.Wrapf(err, "error running apiclient command for static pod: %v", err)
		// }

		// Parse manifest from file and add to array
		podDef, err := UnmarshalPodDefinition(fileData)
		if err != nil {
			return podDefinitions, errors.Wrap(err, "Error getting pod def from manifest")
		}
		podDefinitions = append(podDefinitions, podDef)
	}
	return podDefinitions, nil
}

var podFileExtensions = map[string]struct{}{".yaml": {}, ".manifest": {}, "": {}}

func isPodFile(f fs.FileInfo) bool {
	_, ok := podFileExtensions[filepath.Ext(f.Name())]
	return ok
}
