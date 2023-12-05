package constants

import (
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/types"
)

// Constants used across the version-tracker source code.
const (
	BaseRepoOwnerEnvvar         = "BASE_REPO_OWNER"
	HeadRepoOwnerEnvvar         = "HEAD_REPO_OWNER"
	GitHubTokenEnvvar           = "GITHUB_TOKEN"
	CommitAuthorNameEnvvar      = "COMMIT_AUTHOR_NAME"
	CommitAuthorEmailEnvvar     = "COMMIT_AUTHOR_EMAIL"
	DefaultCommitAuthorName     = "EKS Distro PR Bot"
	DefaultCommitAuthorEmail    = "aws-model-rocket-bots+eksdistroprbot@amazon.com"
	BuildToolingRepoName        = "eks-anywhere-build-tooling"
	BuildToolingRepoURL         = "https://github.com/aws/eks-anywhere-build-tooling"
	ReadmeFile                  = "README.md"
	ReadmeUpdateScriptFile      = "build/lib/readme_check.sh"
	LicenseBoilerplateFile      = "hack/boilerplate.yq.txt"
	SkippedProjectsFile         = "SKIPPED_PROJECTS"
	UpstreamProjectsTrackerFile = "UPSTREAM_PROJECTS.yaml"
	GitTagFile                  = "GIT_TAG"
	GoVersionFile               = "GOLANG_VERSION"
	ChecksumsFile               = "CHECKSUMS"
	AttributionsFilePattern     = "*ATTRIBUTION.txt"
	PatchesDirectory            = "patches"
	GithubPerPage               = 100
	datetimeFormat              = "%Y-%m-%dT%H:%M:%SZ"
	MainBranchName              = "main"
	BaseRepoHeadRevision        = "refs/remotes/origin/main"
	PullRequestBody             = `This PR bumps %[1]s/%[2]s to the latest Git revision, along with other updates such as Go version, checksums and attribution files.
	
[Compare changes](https://github.com/%[1]s/%[2]s/compare/%s...%s)
	
By submitting this pull request, I confirm that you can use, modify, copy, and redistribute this contribution, under the terms of your choice.`
	PullRequestHoldLabel           = "do-not-merge/hold"
	PullRequestWorkInProgressLabel = "do-not-merge/work-in-progress"
	PatchesCommentBody             = `# This pull request is incomplete!
The project being upgraded in this pull request needs changes to patches that cannot be handled automatically. A developer will need to regenerate the patches locally and update the pull request. In addition to patches, the checksums and attribution file(s) corresponding to the project will need to be updated.`
)

var (
	// ProjectReleaseAssets is the mapping of project name to release tarball configurations.
	ProjectReleaseAssets = map[string]types.ReleaseTarball{
		"apache/cloudstack-cloudmonkey": {
			AssetName:  "cmk.linux.x86-64",
			BinaryName: "cmk.linux.x86-64",
			Extract:    false,
		},
		"cert-manager/cert-manager": {
			AssetName:  "cmctl-linux-amd64.tar.gz",
			BinaryName: "cmctl",
			Extract:    true,
		},
		"containerd/containerd": {
			AssetName:                "containerd-%s-linux-amd64.tar.gz",
			BinaryName:               "bin/containerd",
			Extract:                  true,
			TrimLeadingVersionPrefix: true,
		},
		"fluxcd/flux2": {
			AssetName:                "flux_%s_linux_amd64.tar.gz",
			BinaryName:               "flux",
			Extract:                  true,
			TrimLeadingVersionPrefix: true,
		},
		"helm/helm": {
			OverrideAssetURL: "https://get.helm.sh/helm-%s-linux-amd64.tar.gz",
			AssetName:        "helm-%s-linux-amd64.tar.gz",
			BinaryName:       "linux-amd64/helm",
			Extract:          true,
		},
		"kubernetes-sigs/cluster-api": {
			AssetName:  "clusterctl-linux-amd64",
			BinaryName: "clusterctl-linux-amd64",
			Extract:    false,
		},
		"kubernetes-sigs/cri-tools": {
			AssetName:  "crictl-%s-linux-amd64.tar.gz",
			BinaryName: "crictl",
			Extract:    true,
		},
		"kubernetes-sigs/kind": {
			AssetName:  "kind-linux-amd64",
			BinaryName: "kind-linux-amd64",
			Extract:    false,
		},
		"opencontainers/runc": {
			AssetName:  "runc.amd64",
			BinaryName: "runc.amd64",
			Extract:    false,
		},
		"rancher/local-path-provisioner": {
			AssetName:  "local-path-provisioner-amd64",
			BinaryName: "local-path-provisioner-amd64",
			Extract:    false,
		},
		"replicatedhq/troubleshoot": {
			AssetName:  "support-bundle_linux_amd64.tar.gz",
			BinaryName: "support-bundle",
			Extract:    true,
		},
		"vmware/govmomi": {
			AssetName:  "govc_Linux_x86_64.tar.gz",
			BinaryName: "govc",
			Extract:    true,
		},
	}

	// ProjectGoVersionSourceOfTruth is the mapping of project name to Go version source of truth files configuration.
	ProjectGoVersionSourceOfTruth = map[string]types.GoVersionSourceOfTruth{
		"brancz/kube-rbac-proxy": {
			SourceOfTruthFile:     ".github/workflows/build.yml",
			GoVersionSearchString: `go-version: '(1\.\d\d)\.\d+'`,
		},
		"fluxcd/helm-controller": {
			SourceOfTruthFile:     "Dockerfile",
			GoVersionSearchString: `ARG GO_VERSION=(1\.\d\d)`,
		},
		"fluxcd/kustomize-controller": {
			SourceOfTruthFile:     "Dockerfile",
			GoVersionSearchString: `ARG GO_VERSION=(1\.\d\d)`,
		},
		"fluxcd/notification-controller": {
			SourceOfTruthFile:     "Dockerfile",
			GoVersionSearchString: `ARG GO_VERSION=(1\.\d\d)`,
		},
		"fluxcd/source-controller": {
			SourceOfTruthFile:     "Dockerfile",
			GoVersionSearchString: `ARG GO_VERSION=(1\.\d\d)`,
		},
		"kube-vip/kube-vip": {
			SourceOfTruthFile:     "Dockerfile",
			GoVersionSearchString: `golang:(1\.\d\d)`,
		},
		"kubernetes-sigs/cluster-api-provider-cloudstack": {
			SourceOfTruthFile:     "go.mod",
			GoVersionSearchString: `go (1\.\d\d)`,
		},
		"kubernetes-sigs/cluster-api-provider-vsphere": {
			SourceOfTruthFile:     "Makefile",
			GoVersionSearchString: `GO_VERSION \?= (1\.\d\d)\.\d+`,
		},
		"kubernetes-sigs/etcdadm": {
			SourceOfTruthFile:     "Makefile",
			GoVersionSearchString: `GO_IMAGE \?= golang:(1\.\d\d)`,
		},
		"nutanix-cloud-native/cluster-api-provider-nutanix": {
			SourceOfTruthFile:     "go.mod",
			GoVersionSearchString: `go (1\.\d\d)`,
		},
		"tinkerbell/boots": {
			SourceOfTruthFile:     "go.mod",
			GoVersionSearchString: `go (1\.\d\d)`,
		},
		"tinkerbell/cluster-api-provider-tinkerbell": {
			SourceOfTruthFile:     "Makefile",
			GoVersionSearchString: `GOLANG_VERSION := (1\.\d\d)`,
		},
		"tinkerbell/hegel": {
			SourceOfTruthFile:     ".github/workflows/ci.yaml",
			GoVersionSearchString: `GO_VERSION: "(1\.\d\d)"`,
		},
		"tinkerbell/tink": {
			SourceOfTruthFile:     ".github/workflows/ci.yaml",
			GoVersionSearchString: `GO_VERSION: "(1\.\d\d)"`,
		},
	}
)
