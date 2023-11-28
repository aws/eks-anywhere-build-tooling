package builder

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	k8syaml "sigs.k8s.io/yaml"
	"strings"

	releasev1 "github.com/aws/eks-anywhere/release/api/v1alpha1"
	eksDreleasev1 "github.com/aws/eks-distro-build-tooling/release/api/v1alpha1"
)

func (b *BuildOptions) DownloadArtifacts() error {
	// Default arch if not set
	if b.Arch == "" {
		b.Arch = amd64
	}
	// Clone build tooling from the latest release branch
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Error retrieving current working directory: %v", err)
	}

	buildToolingRepoPath := getBuildToolingPath(cwd)
	bundles, _, err := b.prepBuildToolingRepo(buildToolingRepoPath)
	if err != nil {
		return err
	}
	supportedReleaseBranches := GetSupportedReleaseBranches()
	if !SliceContains(supportedReleaseBranches, b.ReleaseChannel) {
		cleanup(buildToolingRepoPath)
		log.Fatalf("release-channel should be one of %v", supportedReleaseBranches)
	}

	kubeVersion := strings.ReplaceAll(b.ReleaseChannel, "-", ".")
	log.Printf("Setting Kubernetes Version: %s", kubeVersion)

	// Download eks-a artifacts
	eksaArtifactsPath := filepath.Join(cwd, eksAnywhereArtifactsDirName)
	eksdArtifactsPath := filepath.Join(cwd, eksDistroArtifactsDirName)
	if b.Force {
		cleanup(eksaArtifactsPath)
		cleanup(eksdArtifactsPath)
	}
	versionBundle, err := b.downloadEKSAArtifacts(eksaArtifactsPath, kubeVersion, bundles)
	if err != nil {
		return fmt.Errorf("failed to download eks-a artifacts: %v", err)
	}

	// Download eks-d artifacts
	eksDReleaseManifestUrl := versionBundle.EksD.EksDReleaseUrl
	if err := b.downloadEKSDArtifacts(eksdArtifactsPath, eksDReleaseManifestUrl); err != nil {
		return fmt.Errorf("failed to download eks-d artifacts: %v", err)
	}
	cleanup(buildToolingRepoPath)
	log.Println("All artifacts were successfully downloaded")
	log.Printf("Please find EKS-A artifacts under %s and EKS-D artifacts under %s directories\n", eksAnywhereArtifactsDirName, eksDistroArtifactsDirName)
	return nil
}

func (b *BuildOptions) downloadEKSAArtifacts(outputPath, kubeVersion string, bundles *releasev1.Bundles) (*releasev1.VersionsBundle, error) {
	var versionBundle releasev1.VersionsBundle
	for _, bundle := range bundles.Spec.VersionsBundles {
		if bundle.KubeVersion == kubeVersion {
			versionBundle = bundle
			break
		}
	}

	eksDBundle := versionBundle.EksD
	artifactsUrl := map[string]string{
		"containerd": eksDBundle.Containerd.URI,
		"crictl":     eksDBundle.Crictl.URI,
		"etcdadm":    eksDBundle.Etcdadm.URI,
	}
	for artifact, url := range artifactsUrl {
		log.Printf("Downloading EKS-A Artifact: %s\n", artifact)
		if err := downloadArtifact(outputPath, url); err != nil {
			return nil, err
		}
	}
	return &versionBundle, nil
}

func (b *BuildOptions) downloadEKSDArtifacts(outputPath, manifestUrl string) error {
	releases := &eksDreleasev1.Release{}
	manifestData, err := readFileFromURL(manifestUrl)
	if err != nil {
		return fmt.Errorf("failed to read manifest file from url: %v", err)
	}
	if err = k8syaml.Unmarshal(manifestData, releases); err != nil {
		return err
	}

	var fullKubeVersion, eksDBaseUrl string
	for _, component := range releases.Status.Components {
		if component.Name == "kubernetes" {
			fullKubeVersion = component.GitTag
			log.Printf("Full Kube Version: %s\n", fullKubeVersion)

			for _, asset := range component.Assets {
				if asset.Name == fmt.Sprintf("bin/linux/%s/kube-apiserver.tar", b.Arch) {
					apiServerUrl := asset.Archive.URI
					eksDBaseUrl = strings.ReplaceAll(apiServerUrl, "/kube-apiserver.tar", "")
				}
			}
		}

		if component.Name == "etcd" || component.Name == "cni-plugins" {
			for _, asset := range component.Assets {
				if asset.Name == fmt.Sprintf("%s-linux-%s-%s.tar.gz", component.Name, b.Arch, component.GitTag) {
					log.Printf("Downloading EKS-D Artifact: %s", component.Name)
					if err = downloadArtifact(outputPath, asset.Archive.URI); err != nil {
						return err
					}
				}
			}
		}
	}
	log.Printf("EKS-D Kubernetes Base Url: %s\n", eksDBaseUrl)

	eksDKubeArtifacts := []string{
		"kube-apiserver.tar",
		"kube-scheduler.tar",
		"kube-controller-manager.tar",
		"kube-proxy.tar",
		"pause.tar",
		"coredns.tar",
		"etcd.tar",
		"kubeadm",
		"kubelet",
		"kubectl",
	}
	for _, artifact := range eksDKubeArtifacts {
		artifactUrl, err := url.JoinPath(eksDBaseUrl, artifact)
		if err != nil {
			return fmt.Errorf("failed to construct url for artifact %s: %v", artifact, err)
		}
		log.Printf("Downloading EKS-D Kubernetes Artifact: %s\n", artifact)
		if err := downloadArtifact(outputPath, artifactUrl); err != nil {
			return err
		}
	}
	return nil
}

func getArtifactFilePathFromUrl(artifactUrl string) (string, error) {
	u, err := url.Parse(artifactUrl)
	if err != nil {
		return "", err
	}
	// Remove the leading /
	urlPath := strings.TrimPrefix(u.Path, "/")
	return urlPath, nil
}

func downloadArtifact(outputPath, artifactUrl string) error {
	artifactsPath, err := getArtifactFilePathFromUrl(artifactUrl)
	if err != nil {
		return fmt.Errorf("failed to get artifacts path from url: %v", err)
	}
	if err := downloadFile(filepath.Join(outputPath, artifactsPath), artifactUrl); err != nil {
		return err
	}
	return nil
}
