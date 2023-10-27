package builder

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	releasev1 "github.com/aws/eks-anywhere/release/api/v1alpha1"
	"github.com/ghodss/yaml"
	k8syaml "sigs.k8s.io/yaml"
)

type EksDReleases struct {
	Releases []EksDRelease `yaml:"releases"`
}

type EksDRelease struct {
	Branch      string `yaml:"branch"`
	Number      string `yaml:"number"`
	KubeVersion string `yaml:"kubeVersion"`
}

func (bo *BuildOptions) prepBuildToolingRepo(buildToolingRepoPath string) (string, error) {
	// Clone build tooling repo
	if bo.Force {
		// Clean up build tooling repo in cwd
		cleanup(buildToolingRepoPath)
	}

	gitCommitFromBundle, detectedEksaVersion, err := bo.getGitCommitFromBundle()
	if err != nil {
		return "", fmt.Errorf("Error getting git commit from bundle: %v", err)
	}
	if codebuild != "true" {
		err = cloneRepo(bo.getBuildToolingRepoUrl(), buildToolingRepoPath)
		if err != nil {
			return "", fmt.Errorf("Error cloning build tooling repo: %v", err)
		}
		log.Println("Cloned eks-anywhere-build-tooling repo")

		err = checkoutRepo(buildToolingRepoPath, gitCommitFromBundle)
		if err != nil {
			return "", fmt.Errorf("Error checking out build tooling repo at commit %s: %v", gitCommitFromBundle, err)
		}
		log.Printf("Checked out eks-anywhere-build-tooling repo at commit %s\n", gitCommitFromBundle)
	} else {
		buildToolingRepoPath = os.Getenv(codebuildSourceDirectoryEnvVar)
		log.Println("Using repo checked out from code commit")
	}
	return detectedEksaVersion, nil
}

func (bo *BuildOptions) getBuildToolingRepoUrl() string {
	if bo.AirGapped {
		switch bo.Hypervisor {
		case VSphere:
			return bo.VsphereConfig.EksABuildToolingRepoUrl
		case Baremetal:
			return bo.BaremetalConfig.EksABuildToolingRepoUrl
		case Nutanix:
			return bo.NutanixConfig.EksABuildToolingRepoUrl
		case CloudStack:
			return bo.CloudstackConfig.EksABuildToolingRepoUrl
		}
	}
	return buildToolingRepoUrl
}

func cloneRepo(cloneUrl, destination string) error {
	log.Printf("Cloning eks-anywhere-build-tooling from %s...", cloneUrl)
	cloneRepoCommandSequence := fmt.Sprintf("git clone %s %s", cloneUrl, destination)
	cmd := exec.Command("bash", "-c", cloneRepoCommandSequence)
	return execCommandWithStreamOutput(cmd)
}

func checkoutRepo(gitRoot, commit string) error {
	log.Printf("Checking out commit %s for build...\n", commit)
	checkoutRepoCommandSequence := fmt.Sprintf("git -C %s checkout %s", gitRoot, commit)
	cmd := exec.Command("bash", "-c", checkoutRepoCommandSequence)
	return execCommandWithStreamOutput(cmd)
}

func execCommandWithStreamOutput(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("Executing command: %v\n", cmd)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	return nil
}

func executeMakeBuildCommand(buildCommand string, envVars ...string) error {
	cmd := exec.Command("bash", "-c", buildCommand)
	cmd.Env = os.Environ()
	for _, envVar := range envVars {
		cmd.Env = append(cmd.Env, envVar)
	}
	return execCommandWithStreamOutput(cmd)
}

func cleanup(buildToolingDir string) {
	if codebuild == "true" {
		return
	}

	log.Print("Cleaning up cache build files")
	err := os.RemoveAll(buildToolingDir)
	if err != nil {
		log.Fatalf("Error cleaning up build tooling dir: %v", err)
	}
}

func GetSupportedReleaseBranches() []string {
	buildToolingPath, err := getRepoRoot()
	if err != nil {
		log.Fatalf(err.Error())
	}

	supportedBranchesFile := filepath.Join(buildToolingPath, supportedReleaseBranchesFileName)
	supportedBranchesFileData, err := os.ReadFile(supportedBranchesFile)
	supportReleaseBranches := strings.Split(string(supportedBranchesFileData), "\n")

	return supportReleaseBranches
}

func getEksDReleaseBranchesWithNumber() (map[string]string, error) {
	buildToolingPath, err := getRepoRoot()
	if err != nil {
		log.Fatalf(err.Error())
	}

	eksDReleaseBranchesFile := filepath.Join(buildToolingPath, eksDLatestReleasesFileName)
	eksDReleaseBranchesFileData, err := os.ReadFile(eksDReleaseBranchesFile)

	eksDReleaseBranchesWithNumber := make(map[string]string)
	var eksDReleaseBranchesDataWithNumber EksDReleases
	err = yaml.Unmarshal(eksDReleaseBranchesFileData, &eksDReleaseBranchesDataWithNumber)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling EKSD_LATEST_RELEASES file: %v", err)
	}

	for _, eksdRelease := range eksDReleaseBranchesDataWithNumber.Releases {
		eksDReleaseBranchesWithNumber[eksdRelease.Branch] = eksdRelease.Number
	}

	return eksDReleaseBranchesWithNumber, nil
}

func downloadFile(filepath, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func getBuildToolingPath(cwd string) string {
	buildToolingRepoPath := filepath.Join(cwd, "eks-anywhere-build-tooling")
	if codebuild == "true" {
		buildToolingRepoPath = os.Getenv(codebuildSourceDirectoryEnvVar)
	}
	return buildToolingRepoPath
}

func getRepoRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error retrieving current working directory: %v", err)
	}
	buildToolingPath := getBuildToolingPath(cwd)
	cmd := exec.Command("git", "-C", buildToolingPath, "rev-parse", "--show-toplevel")
	commandOut, err := execCommand(cmd)
	if err != nil {
		return "", err
	}
	return commandOut, nil
}

func getManifestRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error retrieving current working directory: %v", err)
	}
	return filepath.Join(cwd, manifestsDirName), nil
}

func getAirGapCmdEnvVars(cloneUrl, eksAVersion, eksDReleaseBranch string) ([]string, error) {
	manifestRoot, err := getManifestRoot()
	if err != nil {
		return nil, err
	}
	manifestRoot = fmt.Sprintf("file://%s", manifestRoot)
	// EKS-A Manifest file path
	eksABundlesFilePath := filepath.Join(manifestRoot, fmt.Sprintf(eksAnywhereBundlesFileNameFormat, eksAVersion))
	cmdEnvVars := []string{fmt.Sprintf("%s=%s", eksABundlesURLEnvVar, eksABundlesFilePath)}

	// EKS-D Manifest file path
	eksDManifestFilePath := filepath.Join(manifestRoot, fmt.Sprintf(eksDistroManifestFileNameFormat, eksDReleaseBranch))
	cmdEnvVars = append(cmdEnvVars, fmt.Sprintf("%s=%s", eksDManifestURLEnvVar, eksDManifestFilePath))

	// Upstream clone url
	cmdEnvVars = append(cmdEnvVars, fmt.Sprintf("%s=%s", cloneUrlEnvVar, cloneUrl))
	return cmdEnvVars, nil
}

func SliceContains(s []string, str string) bool {
	for _, elem := range s {
		if elem == str {
			return true
		}
	}
	return false
}

func execCommand(cmd *exec.Cmd) (string, error) {
	log.Printf("Executing command: %v\n", cmd)
	commandOutput, err := cmd.CombinedOutput()
	commandOutputStr := strings.TrimSpace(string(commandOutput))

	if err != nil {
		return commandOutputStr, fmt.Errorf("failed to run command: %v", err)
	}
	return commandOutputStr, nil
}

func (bo *BuildOptions) getGitCommitFromBundle() (string, string, error) {
	manifestDirPath, err := getManifestRoot()
	if err != nil {
		return "", "", err
	}

	eksAReleasesManifestURL, err := getEksAReleasesManifestURL(bo.AirGapped)
	if err != nil {
		return "", "", err
	}
	var releasesManifestContents []byte
	if bo.ManifestTarball != "" {
		eksAManifestFile := filepath.Join(manifestDirPath, eksAnywhereManifestFileName)
		log.Printf("Reading EKS-A Manifest file: %s", eksAManifestFile)
		releasesManifestContents, err = os.ReadFile(eksAManifestFile)
		if err != nil {
			return "", "", err
		}
	} else {
		log.Printf("Reading EKS-A Manifest file: %s", eksAReleasesManifestURL)
		releasesManifestContents, err = readFileFromURL(eksAReleasesManifestURL)
		if err != nil {
			return "", "", err
		}
	}

	releases := &releasev1.Release{}
	if err = k8syaml.Unmarshal(releasesManifestContents, releases); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal release manifest from [%s]: %v", eksAReleasesManifestURL, err)
	}

	var eksAReleaseVersion, bundleManifestUrl string
	var foundRelease bool
	if os.Getenv(eksaUseDevReleaseEnvVar) == "true" {
		eksAReleaseVersion = devEksaReleaseVersion
		log.Printf("EKSA_USE_DEV_RELEASE set to true, using EKS-A dev release version: %s", eksAReleaseVersion)
	} else if bo.EKSAReleaseVersion != "" {
		eksAReleaseVersion = bo.EKSAReleaseVersion
		log.Printf("EKS-A release version provided: %s", eksAReleaseVersion)
	} else if eksaVersion != "" {
		eksAReleaseVersion = eksaVersion
		log.Printf("No EKS-A release version provided, defaulting to EKS-A version configured at build time: %s", eksAReleaseVersion)
	} else {
		eksAReleaseVersion = releases.Spec.LatestVersion
		log.Printf("No EKS-A release version provided, defaulting to latest EKS-A version: %s", eksAReleaseVersion)
	}

	for _, r := range releases.Spec.Releases {
		if r.Version == eksAReleaseVersion {
			foundRelease = true
			bundleManifestUrl = r.BundleManifestUrl
			break
		}
	}

	if !foundRelease {
		// if release was not found, this is probably a dev release version which we need
		// to use a prefix match to find since the version in the release.yaml
		// will be v0.0.0-dev+build.7423
		for _, r := range releases.Spec.Releases {
			if strings.Contains(r.Version, eksAReleaseVersion) {
				foundRelease = true
				bundleManifestUrl = r.BundleManifestUrl
				break
			}
		}
	}
	if !foundRelease {
		return "", "", fmt.Errorf("version %s is not a valid EKS-A release", eksAReleaseVersion)
	}
	var bundleManifestContents []byte
	if bo.ManifestTarball == "" {
		log.Printf("Fetching git commit from bundle manifest: %s", bundleManifestUrl)
		bundleManifestContents, err = readFileFromURL(bundleManifestUrl)
		if err != nil {
			return "", "", err
		}
	} else {
		bundleManifestFile := filepath.Join(manifestDirPath, fmt.Sprintf(eksAnywhereBundlesFileNameFormat, eksAReleaseVersion))
		log.Printf("Fetching git commit from bundle manifest: %s", bundleManifestFile)
		bundleManifestContents, err = os.ReadFile(bundleManifestFile)
		if err != nil {
			return "", "", err
		}
	}

	bundles := &releasev1.Bundles{}
	if err = k8syaml.Unmarshal(bundleManifestContents, bundles); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal bundles manifest from [%s]: %v", bundleManifestUrl, err)
	}

	return bundles.Spec.VersionsBundles[0].EksD.GitCommit, eksAReleaseVersion, nil
}

func readFileFromURL(url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating http GET request for downloading file: %v", err)
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSHandshakeTimeout = 60 * time.Second
	client := &http.Client{
		Transport: transport,
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed reading file from URL [%s]: %v", url, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading contents of URL body [%s]: %v", url, err)
	}

	return data, nil
}

func getEksAReleasesManifestURL(airgapped bool) (string, error) {
	if os.Getenv(eksaUseDevReleaseEnvVar) != "true" {
		if eksaReleaseManifest != "" {
			return eksaReleaseManifest, nil
		}
		if airgapped {
			manifestRoot, err := getManifestRoot()
			if err != nil {
				return "", nil
			}
			return fmt.Sprintf("file://%s/%s", manifestRoot, eksAnywhereManifestFileName), nil
		}
		return prodEksaReleaseManifestURL, nil
	}

	// using a dev release, allow branch_name env var to
	// override manifest url
	branchName, ok := os.LookupEnv(branchNameEnvVar)
	if !ok {
		branchName = mainBranch
	}

	if branchName != mainBranch {
		return fmt.Sprintf(devBranchEksaReleaseManifestURL, branchName), nil
	}

	return devEksaReleaseManifestURL, nil
}

// setRhsmProxy takes the proxy config, parses it and sets the appropriate config on rhsm config
func setRhsmProxy(proxy *ProxyConfig, rhsm *RhsmConfig) error {
	if proxy.HttpProxy != "" {
		host, port, err := parseUrl(proxy.HttpProxy)
		if err != nil {
			return err
		}
		rhsm.ProxyHostname = host
		rhsm.ProxyPort = port
	}

	return nil
}

// parseUrl takes a http endpoint and returns hostname, ports and error
func parseUrl(endpoint string) (string, string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", "", err
	}

	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return "", "", err
	}
	return host, port, nil
}

// createTarball takes in a filename and creates a tarball of the contents in the provided path
func createTarball(fileName, path string) error {
	// Create a new tarball.
	tarball, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer tarball.Close()

	// Create a tar writer.
	tw := tar.NewWriter(tarball)
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Create a new tar header.
		header, err := tar.FileInfoHeader(info, path)
		if err != nil {
			return err
		}

		// Write the header to the tar writer.
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// Copy the file contents to the tar writer.
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(tw, file)
		if err != nil {
			return err
		}
		return nil
	})

	// Close the tar writer.
	err = tw.Close()
	if err != nil {
		panic(err)
	}

	return nil
}

func replaceStringInFile(filePath, oldString, newString string) error {
	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	replacedString := strings.ReplaceAll(string(fileContents), oldString, newString)
	if err = os.Remove(filePath); err != nil {
		return err
	}
	err = os.WriteFile(filePath, []byte(replacedString), 0o755)
	if err != nil {
		return err
	}
	return nil
}
