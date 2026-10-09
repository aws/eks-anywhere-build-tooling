package upgrade

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"slices"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

var patchSubjectPrefix = regexp.MustCompile(`(?m)^Subject:[ \t]+\[PATCH(?:[ \t]+[0-9]+/[0-9]+)?\][ \t]*`)

// normalizePatchNumbers updates only mail subject prefixes, keeping filenames
// and semantic commit metadata intact after repairs or obsolete-patch removal.
func normalizePatchNumbers(directory string) ([]string, error) {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	paths, err := patchFilesInDirectory(directory)
	if err != nil {
		return nil, err
	}
	var changed []string
	for index, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return changed, err
		}
		end := bytes.Index(content, []byte("\n\n"))
		if end < 0 {
			logger.Info("Skipping cosmetic patch numbering", "patch", path, "reason", "unsupported mail header boundary")
			continue
		}
		header := content[:end]
		if !patchSubjectPrefix.Match(header) {
			logger.Info("Skipping cosmetic patch numbering", "patch", path, "reason", "unsupported subject prefix")
			continue
		}
		prefix := "Subject: [PATCH] "
		if len(paths) > 1 {
			prefix = fmt.Sprintf("Subject: [PATCH %d/%d] ", index+1, len(paths))
		}
		updated := append(patchSubjectPrefix.ReplaceAll(header, []byte(prefix)), content[end:]...)
		if bytes.Equal(updated, content) {
			continue
		}
		info, err := os.Stat(path)
		if err != nil {
			return changed, err
		}
		if err := writeFileAtomically(path, updated, info.Mode().Perm()); err != nil {
			return changed, err
		}
		changed = append(changed, path)
	}
	return changed, nil
}

func finishPatchNumbering(directory string, changed []string) []string {
	numbered, err := normalizePatchNumbers(directory)
	if err != nil {
		logger.Info("Cosmetic patch numbering incomplete", "error", err)
	}
	paths := append(slices.Clone(changed), numbered...)
	slices.Sort(paths)
	return slices.Compact(paths)
}
