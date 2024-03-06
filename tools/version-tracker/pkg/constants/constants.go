package constants

import (
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/types"
)

// Constants used across the version-tracker source code.
const (
	BaseRepoOwnerEnvvar                     = "BASE_REPO_OWNER"
	HeadRepoOwnerEnvvar                     = "HEAD_REPO_OWNER"
	GitHubTokenEnvvar                       = "GITHUB_TOKEN"
	CommitAuthorNameEnvvar                  = "COMMIT_AUTHOR_NAME"
	CommitAuthorEmailEnvvar                 = "COMMIT_AUTHOR_EMAIL"
	DefaultCommitAuthorName                 = "EKS Distro PR Bot"
	DefaultCommitAuthorEmail                = "aws-model-rocket-bots+eksdistroprbot@amazon.com"
	BuildToolingRepoName                    = "eks-anywhere-build-tooling"
	DefaultBaseRepoOwner                    = "aws"
	BuildToolingRepoURL                     = "https://github.com/%s/eks-anywhere-build-tooling"
	ReadmeFile                              = "README.md"
	ReadmeUpdateScriptFile                  = "build/lib/readme_check.sh"
	LicenseBoilerplateFile                  = "hack/boilerplate.yq.txt"
	EKSDistroLatestReleasesFile             = "EKSD_LATEST_RELEASES"
	EKSDistroProdReleaseNumberFileFormat    = "release/%s/production/RELEASE"
	KubernetesGitTagFileFormat              = "projects/kubernetes/kubernetes/%s/GIT_TAG"
	SkippedProjectsFile                     = "SKIPPED_PROJECTS"
	UpstreamProjectsTrackerFile             = "UPSTREAM_PROJECTS.yaml"
	SupportedReleaseBranchesFile            = "release/SUPPORTED_RELEASE_BRANCHES"
	GitTagFile                              = "GIT_TAG"
	GoVersionFile                           = "GOLANG_VERSION"
	ChecksumsFile                           = "CHECKSUMS"
	AttributionsFilePattern                 = "*ATTRIBUTION.txt"
	PatchesDirectory                        = "patches"
	BottlerocketReleasesFile                = "BOTTLEROCKET_RELEASES"
	BottlerocketContainerMetadataFileFormat = "BOTTLEROCKET_%s_CONTAINER_METADATA"
	BottlerocketHostContainersTOMLFile      = "sources/models/shared-defaults/public-host-containers.toml"
	CiliumImageRepository                   = "public.ecr.aws/isovalent/cilium"
	GithubPerPage                           = 100
	datetimeFormat                          = "%Y-%m-%dT%H:%M:%SZ"
	MainBranchName                          = "main"
	BaseRepoHeadRevision                    = "refs/remotes/origin/main"
	EKSDistroUpgradePullRequestBody         = `This PR bumps EKS Distro releases to the latest available release versions.

/hold
/area dependencies

By submitting this pull request, I confirm that you can use, modify, copy, and redistribute this contribution, under the terms of your choice.`
	DefaultUpgradePullRequestBody = `This PR bumps %[1]s/%[2]s to the latest Git revision.

[Compare changes](https://github.com/%[1]s/%[2]s/compare/%[3]s...%[4]s)
[Release notes](https://github.com/%[1]s/%[2]s/releases/%[4]s)

/hold
/area dependencies

By submitting this pull request, I confirm that you can use, modify, copy, and redistribute this contribution, under the terms of your choice.`
	BottlerocketUpgradePullRequestBody = `This PR bumps Bottlerocket releases to the latest Git revision.

[Compare changes](https://github.com/bottlerocket-os/bottlerocket/compare/%[1]s...%[2]s)
[Release notes](https://github.com/bottlerocket-os/bottlerocket/releases/%[2]s)

/hold
/area dependencies

By submitting this pull request, I confirm that you can use, modify, copy, and redistribute this contribution, under the terms of your choice.`

	CombinedImageBuilderBottlerocketUpgradePullRequestBody = `This PR bumps kubernetes-sigs/image-builder and Bottlerocket releases to the latest Git revision.

[Compare changes for image-builder](https://github.com/kubernetes-sigs/image-builder/compare/%[1]s...%[2]s)
[Release notes for image-builder](https://github.com/kubernetes-sigs/image-builder/releases/%[2]s)

[Compare changes for Bottlerocket](https://github.com/bottlerocket-os/bottlerocket/compare/%[3]s...%[4]s)
[Release notes for Bottlerocket](https://github.com/bottlerocket-os/bottlerocket/releases/%[4]s)

/hold
/area dependencies

By submitting this pull request, I confirm that you can use, modify, copy, and redistribute this contribution, under the terms of your choice.`
	PatchesCommentBody = `# This pull request is incomplete!
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
		"aquasecurity/harbor-scanner-trivy": {
			AssetName:                "harbor-scanner-trivy_%s_Linux_x86_64.tar.gz",
			BinaryName:               "scanner-trivy",
			Extract:                  true,
			TrimLeadingVersionPrefix: true,
		},
		"aquasecurity/trivy": {
			AssetName:                "trivy_%s_Linux-64bit.tar.gz",
			BinaryName:               "trivy",
			Extract:                  true,
			TrimLeadingVersionPrefix: true,
		},
		"aws/rolesanywhere-credential-helper": {
			OverrideAssetURL:         "https://rolesanywhere.amazonaws.com/releases/%s/X86_64/Linux/aws_signing_helper",
			AssetName:                "aws_signing_helper",
			BinaryName:               "aws_signing_helper",
			Extract:                  false,
			TrimLeadingVersionPrefix: true,
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
		"distribution/distribution": {
			AssetName:                "registry_%s_linux_amd64.tar.gz",
			BinaryName:               "registry",
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
		"prometheus/prometheus": {
			AssetName:                "prometheus-%s.linux-amd64.tar.gz",
			BinaryName:               "prometheus-%s.linux-amd64/prometheus",
			Extract:                  true,
			TrimLeadingVersionPrefix: true,
		},
		"prometheus/node_exporter": {
			AssetName:                "node_exporter-%s.linux-amd64.tar.gz",
			BinaryName:               "node_exporter-%s.linux-amd64/node_exporter",
			Extract:                  true,
			TrimLeadingVersionPrefix: true,
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
		"aws/etcdadm-bootstrap-provider": {
			SourceOfTruthFile:     "go.mod",
			GoVersionSearchString: `go (1\.\d\d)`,
		},
		"aws/etcdadm-controller": {
			SourceOfTruthFile:     "go.mod",
			GoVersionSearchString: `go (1\.\d\d)`,
		},
		"brancz/kube-rbac-proxy": {
			SourceOfTruthFile:     ".github/workflows/build.yml",
			GoVersionSearchString: `go-version: '(1\.\d\d)\.\d+'`,
		},
		"emissary-ingress/emissary": {
			SourceOfTruthFile:     "go.mod",
			GoVersionSearchString: `go (1\.\d\d)`,
		},
		"goharbor/harbor": {
			SourceOfTruthFile:     "Makefile",
			GoVersionSearchString: `GOBUILDIMAGE=golang:(1\.\d\d)`,
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
		"metallb/metallb": {
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

	ProjectsWithUnconventionalUpgradeFlows = []string{
		"cilium/cilium",
		"kubernetes-sigs/image-builder",
	}

	BottlerocketImageFormats = []string{"ami", "ova", "raw"}

	BottlerocketHostContainers = []string{"admin", "control"}

	CiliumImageDirectories = []string{"cilium", "operator-generic", "cilium-chart"}
)
