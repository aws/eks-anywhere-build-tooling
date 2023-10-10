package builder

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	_, err = prepBuildToolingRepo(buildToolingRepoPath, "", b.Force)
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
	releaseManifestUrl := getEksAReleasesManifestURL()
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
