package builder

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	releasev1 "github.com/aws/eks-anywhere/release/api/v1alpha1"
	k8syaml "sigs.k8s.io/yaml"
)

// DownloadManifests connects to the internet, clones a fresh copy of eks-a-build-tooling repo
// and pulls EKS-A and EKS-D artifacts for all supported releases and creates an archive.
func (b *BuildOptions) DownloadManifests() error {
	// Clone build tooling from the latest release branch
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Error retrieving current working directory: %v", err)
	}

	buildToolingRepoPath := getBuildToolingPath(cwd)
	_, err = b.prepBuildToolingRepo(buildToolingRepoPath)
	if err != nil {
		return err
	}

	// Download eks-d manifests
	manifestsPath := filepath.Join(cwd, "manifests")
	if err = downloadEKSDManifests(manifestsPath); err != nil {
		return err
	}

	// Download eks-a manifests
	if err = downloadEKSAManifests(manifestsPath); err != nil {
		return err
	}

	// Create tarball of the downloaded manifests
	log.Println("Creating tarball of downloaded manifests")
	if err = createTarball(manifestsTarballName, manifestsPath); err != nil {
		return err
	}
	log.Printf("Manifest tarball %s was successfully created", manifestsTarballName)

	// Clean up manifests directory
	if err = os.RemoveAll(manifestsPath); err != nil {
		return err
	}
	cleanup(buildToolingRepoPath)
	return nil
}

func downloadEKSDManifests(outputPath string) error {
	// Get all eksDReleases with their release numbers
	eksDReleaseBranchesWithNumber, err := getEksDReleaseBranchesWithNumber()
	if err != nil {
		return err
	}

	// Create output directory path
	if err = os.MkdirAll(outputPath, 0755); err != nil {
		return err
	}

	for branch, number := range eksDReleaseBranchesWithNumber {
		manifestUrl := fmt.Sprintf("https://%s/kubernetes-%s/kubernetes-%s-eks-%s.yaml", eksDistroProdDomain, branch, branch, number)
		manifestFileName := fmt.Sprintf(eksDistroManifestFileNameFormat, branch)
		manifestFilePath := filepath.Join(outputPath, manifestFileName)
		log.Printf("Downloading eks-d manifest for release branch %s, release number %s", branch, number)
		if err := downloadFile(manifestFilePath, manifestUrl); err != nil {
			return err
		}
	}

	return nil
}

func downloadEKSAManifests(outputPath string) error {
	// Download Release manifest
	releaseManifestUrl, err := getEksAReleasesManifestURL(false)
	if err != nil {
		return err
	}

	// Create output directory path
	if err = os.MkdirAll(outputPath, 0755); err != nil {
		return err
	}

	releaseManifestPath := filepath.Join(outputPath, eksAnywhereManifestFileName)
	log.Printf("Downloading eks-a release manifest")
	if err := downloadFile(releaseManifestPath, releaseManifestUrl); err != nil {
		return err
	}

	// Download all the bundles from release manifest
	releaseManifestData, err := os.ReadFile(releaseManifestPath)
	if err != nil {
		return err
	}

	releases := &releasev1.Release{}
	if err = k8syaml.Unmarshal(releaseManifestData, releases); err != nil {
		return err
	}
	for _, r := range releases.Spec.Releases {
		bundleName := fmt.Sprintf(eksAnywhereBundlesFileNameFormat, r.Version)
		bundleFilePath := filepath.Join(outputPath, bundleName)
		log.Printf("Downloading eks-a bundles manifest for eks-a version: %s", r.Version)
		if err = downloadFile(bundleFilePath, r.BundleManifestUrl); err != nil {
			return err
		}
	}
	return nil
}

// extractTarball extracts the provided tarball onto the path
func extractTarball(tarball, path string) error {
	tarFile, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	// Create a new directory to extract the tarball into.
	if err = os.RemoveAll(path); err != nil {
		return err
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}

	// Create a tar reader for the tarball.
	tarReader := tar.NewReader(tarFile)

	// Iterate over the tarball's entries.
	for {
		// Get the next header from the tarball.
		header, err := tarReader.Next()
		if err == io.EOF {
			// We've reached the end of the tarball.
			break
		}
		if err != nil {
			return err
		}

		// Create a new file for the tar entry.
		outFile, err := os.Create(filepath.Join(path, header.Name))
		if err != nil {
			return err
		}
		defer outFile.Close()

		// Copy the tar entry's contents to the new file.
		_, err = io.Copy(outFile, tarReader)
		if err != nil {
			return err
		}

		outFile.Chmod(os.FileMode(header.Mode))
	}

	return nil
}

func extractAndPrepManifestTarball(tarballFile, privateEksDServerDomain, privateEksAServerDomain string) error {
	log.Println("Manifest tarball provided, extracting to directory")
	manifestDir, err := getManifestRoot()
	if err != nil {
		return fmt.Errorf("Error retrieving manifest root")
	}
	if err = extractTarball(tarballFile, manifestDir); err != nil {
		return fmt.Errorf("Error extracting tarball: %v", err)
	}

	// Replacing eks-a bundles manifest hostname
	// These endpoints are used for artifacts like containerd, crictl & etcdadm
	// Find all eks-a bundles manifest and replace hostname
	files, err := os.ReadDir(manifestDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() {
			// Find EKS-D Manifests and replace prod distro domain with private server's
			absFilePath := filepath.Join(manifestDir, file.Name())
			eksDManifestFilePattern := strings.ReplaceAll(eksDistroManifestFileNameFormat, "%s", "*")
			eksDMatch, err := filepath.Match(eksDManifestFilePattern, file.Name())
			if err != nil {
				return err
			}

			if eksDMatch {
				log.Printf("Replacing EKS-D domain name in file: %s\n", file.Name())
				if err = replaceStringInFile(absFilePath, fmt.Sprintf("https://%s", eksDistroProdDomain), privateEksDServerDomain); err != nil {
					return err
				}
			}

			// Find EKS-A Bundles Manifests and replace prod anywhere domain with private server's
			eksABundlesManifestFilePattern := strings.ReplaceAll(eksAnywhereBundlesFileNameFormat, "%s", "*")
			eksAMatch, err := filepath.Match(eksABundlesManifestFilePattern, file.Name())
			if err != nil {
				return err
			}

			fmt.Printf("File: %s, Match: %t, Pattern: %s\n", file.Name(), eksAMatch, eksABundlesManifestFilePattern)
			if eksAMatch {
				log.Printf("Replacing EKS-A domain name in file: %s\n", file.Name())
				if err = replaceStringInFile(absFilePath, fmt.Sprintf("https://%s", eksAnywhereAssetsProdDomain), privateEksAServerDomain); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
