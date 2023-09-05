package builder

import (
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
	"sigs.k8s.io/yaml"
)

func cloneRepo(cloneUrl, destination string) error {
	log.Println("Cloning eks-anywhere-build-tooling...")
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

	supportedBranchesFile := filepath.Join(buildToolingPath, "release/SUPPORTED_RELEASE_BRANCHES")
	supportedBranchesFileData, err := os.ReadFile(supportedBranchesFile)
	supportReleaseBranches := strings.Split(string(supportedBranchesFileData), "\n")

	return supportReleaseBranches
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
	eksAReleasesManifestURL := getEksAReleasesManifestURL()
	releasesManifestContents, err := readFileFromURL(eksAReleasesManifestURL)
	if err != nil {
		return "", "", err
	}

	releases := &releasev1.Release{}
	if err = yaml.Unmarshal(releasesManifestContents, releases); err != nil {
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
	log.Printf("Fetching git commit from bundle manifest: %s", bundleManifestUrl)

	bundleManifestContents, err := readFileFromURL(bundleManifestUrl)
	if err != nil {
		return "", "", err
	}

	bundles := &releasev1.Bundles{}
	if err = yaml.Unmarshal(bundleManifestContents, bundles); err != nil {
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

func getEksAReleasesManifestURL() string {
	if os.Getenv(eksaUseDevReleaseEnvVar) != "true" {
		if eksaReleaseManifest != "" {
			return eksaReleaseManifest
		}

		return prodEksaReleaseManifestURL
	}

	// using a dev release, allow branch_name env var to
	// override manifest url
	branchName, ok := os.LookupEnv(branchNameEnvVar)
	if !ok {
		branchName = mainBranch
	}

	if branchName != mainBranch {
		return fmt.Sprintf(devBranchEksaReleaseManifestURL, branchName)
	}

	return devEksaReleaseManifestURL
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
