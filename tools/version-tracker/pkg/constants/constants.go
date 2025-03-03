package constants

import (
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/types"
)

// Constants used across the version-tracker source code.
const (
	BranchNameEnvVar                          = "BRANCH_NAME"
	BaseRepoOwnerEnvvar                       = "BASE_REPO_OWNER"
	HeadRepoOwnerEnvvar                       = "HEAD_REPO_OWNER"
	GitHubTokenEnvvar                         = "GITHUB_TOKEN"
	CommitAuthorNameEnvvar                    = "COMMIT_AUTHOR_NAME"
	CommitAuthorEmailEnvvar                   = "COMMIT_AUTHOR_EMAIL"
	ReleaseBranchEnvvar                       = "RELEASE_BRANCH"
	DefaultCommitAuthorName                   = "EKS Distro PR Bot"
	DefaultCommitAuthorEmail                  = "aws-model-rocket-bots+eksdistroprbot@amazon.com"
	BuildToolingRepoName                      = "eks-anywhere-build-tooling"
	EKSDistroBuildToolingRepoName             = "eks-distro-build-tooling"
	AWSOrgName                                = "aws"
	BottlerocketOrgName                       = "bottlerocket-os"
	BottlerocketRepoName                      = "bottlerocket"
	BuildToolingRepoURL                       = "https://github.com/%s/eks-anywhere-build-tooling"
	ReadmeFile                                = "README.md"
	ReadmeUpdateScriptFile                    = "build/lib/readme_check.sh"
	LicenseBoilerplateFile                    = "hack/boilerplate.yq.txt"
	BottlerocketTargetsJSONURLFormat          = "https://updates.bottlerocket.aws/2020-07-07/%s-k8s-%s/x86_64/%s.targets.json"
	BottlerocketTimestampJSONURLFormat        = "https://updates.bottlerocket.aws/2020-07-07/%s-k8s-%s/x86_64/timestamp.json"
	BottlerocketAMIImageTargetFormat          = "bottlerocket-%s-k8s-%s-x86_64-%s.img.lz4"
	BottlerocketOVAImageTargetFormat          = "bottlerocket-%s-k8s-%s-x86_64-%s.ova"
	BottlerocketRawImageTargetFormat          = "bottlerocket-%s-k8s-%s-x86_64-%s.img.lz4"
	EKSDistroLatestReleasesFile               = "EKSD_LATEST_RELEASES"
	EKSDistroReleaseChannelsFileURLFormat     = "https://distro.eks.amazonaws.com/releasechannels/%s.yaml"
	EKSDistroReleaseManifestURLFormat         = "https://distro.eks.amazonaws.com/kubernetes-%[1]s/kubernetes-%[1]s-eks-%d.yaml"
	SkippedProjectsFile                       = "SKIPPED_PROJECTS"
	UpstreamProjectsTrackerFile               = "UPSTREAM_PROJECTS.yaml"
	SupportedReleaseBranchesFile              = "release/SUPPORTED_RELEASE_BRANCHES"
	GitTagFile                                = "GIT_TAG"
	GoVersionFile                             = "GOLANG_VERSION"
	ChecksumsFile                             = "CHECKSUMS"
	AttributionsFilePattern                   = "*ATTRIBUTION.txt"
	EKSDistroBaseTagFilesPattern              = "EKS_DISTRO*TAG_FILE"
	EKSDistroBaseUpdatedPackagesFileFormat    = "eks-distro-base-updates/%s/update_packages-%s"
	BuildDirectory                            = "build"
	ManifestsDirectory                        = "manifests"
	PatchesDirectory                          = "patches"
	FailedPatchApplyMarker                    = "patch does not apply"
	DoesNotExistInIndexMarker                 = "does not exist in index"
	SemverRegex                               = `v?(?P<major>0|[1-9]\d*)(\.|_)(?P<minor>0|[1-9]\d*)((\.|_)(?P<patch>0|[1-9]\d*))?(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?`
	FailedPatchApplyRegex                     = "Patch failed at .*"
	FailedPatchFilesRegex                     = "error: (.*): patch does not apply"
	DoesNotExistInIndexFilesRegex             = "error: (.*): does not exist in index"
	GitDescribeRegex                          = `v?\d+\.\d+\.\d+(-([0-9]+)-g.*)?`
	BottlerocketReleasesFile                  = "BOTTLEROCKET_RELEASES"
	BottlerocketContainerMetadataFileFormat   = "BOTTLEROCKET_%s_CONTAINER_METADATA"
	BottlerocketHostContainersTOMLFile        = "sources/shared-defaults/public-host-containers.toml"
	BottlerocketHostContainerSourceImageRegex = "schnauzer-v2 render --template '(.*)'"
	CertManagerManifestYAMLFile               = "cert-manager.yaml"
	CiliumImageRepository                     = "public.ecr.aws/isovalent/cilium"
	EnvoyImageRepository                      = "public.ecr.aws/appmesh/aws-appmesh-envoy"
	EKSDistroBaseTagsYAMLFile                 = "EKS_DISTRO_TAG_FILE.yaml"
	AL2023Suffix                              = "-al2023"
	TagFileSuffix                             = "_TAG_FILE"
	KindNodeImageBuildArgsScriptFile          = "node-image-build-args.sh"
	GithubPerPage                             = 100
	datetimeFormat                            = "%Y-%m-%dT%H:%M:%SZ"
	MainBranchName                            = "main"
	BaseRepoHeadRevisionPattern               = "refs/remotes/origin/%s"
	EKSDistroUpgradePullRequestBody           = `This PR bumps EKS Distro releases to the latest available release versions.

/hold
/area dependencies

By submitting this pull request, I confirm that you can use, modify, copy, and redistribute this contribution, under the terms of your choice.`
	EKSDistroBuildToolingUpgradePullRequestBody = `This PR updates the base image tag in tag file(s) with the tag of the newly-built EKS Distro base image and/or its minimal variants.

%s

/hold
/area dependencies

By submitting this pull request, I confirm that you can use, modify, copy, and redistribute this contribution, under the terms of your choice.`
	DefaultUpgradePullRequestBody = `This PR bumps %[1]s/%[2]s to the latest Git revision.

[Compare changes](https://github.com/%[1]s/%[2]s/compare/%[3]s...%[4]s)
[Release notes](https://github.com/%[1]s/%[2]s/releases/%[4]s)

%s

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
	FailedPatchesCommentBody = `## Failed patch details
**Only %d/%d patches were applied!**
%s
The following files in the above patch did not apply successfully:
%s`
	FailedUpgradeCommentBody = `# This pull request is incomplete and requires manual intervention from a team member!
The following steps in the upgrade flow failed:
%s`
)

var (
	// ProjectReleaseAssets is the mapping of project name to release tarball configurations.
	ProjectReleaseAssets = map[string]types.ReleaseTarball{
		"apache/cloudstack-cloudmonkey": {
			AssetName:  "cmk.linux.x86-64",
			BinaryName: "cmk.linux.x86-64",
			Extract:    false,
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
		"goharbor/harbor-scanner-trivy": {
			AssetName:                "harbor-scanner-trivy_%s_Linux_x86_64.tar.gz",
			BinaryName:               "scanner-trivy",
			Extract:                  true,
			TrimLeadingVersionPrefix: true,
		},
		"helm/helm": {
			OverrideAssetURL: "https://get.helm.sh/helm-%s-linux-amd64.tar.gz",
			AssetName:        "helm-%s-linux-amd64.tar.gz",
			BinaryName:       "linux-amd64/helm",
			Extract:          true,
		},
		"linuxkit/linuxkit": {
			AssetName:  "linuxkit-linux-amd64",
			BinaryName: "linuxkit-linux-amd64",
			Extract:    false,
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
		"apache/cloudstack-cloudmonkey": {
			SourceOfTruthFile:     ".github/workflows/build.yml",
			GoVersionSearchString: `go-version: (1\.\d\d)`,
		},
		"aquasecurity/trivy": {
			SourceOfTruthFile:     "Dockerfile.protoc",
			GoVersionSearchString: `golang:(1\.\d\d)`,
		},
		"aws/etcdadm-bootstrap-provider": {
			SourceOfTruthFile:     "go.mod",
			GoVersionSearchString: `go (1\.\d\d)`,
		},
		"aws/etcdadm-controller": {
			SourceOfTruthFile:     "go.mod",
			GoVersionSearchString: `go (1\.\d\d)`,
		},
		"aws/rolesanywhere-credential-helper": {
			SourceOfTruthFile:     "go.mod",
			GoVersionSearchString: `go (1\.\d\d)`,
		},
		"brancz/kube-rbac-proxy": {
			SourceOfTruthFile:     ".github/workflows/build.yml",
			GoVersionSearchString: `go-version: '(1\.\d\d)'`,
		},
		"cert-manager/cert-manager": {
			SourceOfTruthFile:     "go.mod",
			GoVersionSearchString: `go (1\.\d\d)`,
		},
		"containerd/containerd": {
			SourceOfTruthFile:     ".github/workflows/release.yml",
			GoVersionSearchString: `GO_VERSION: "(1\.\d\d)"`,
		},
		"distribution/distribution": {
			SourceOfTruthFile:     "Dockerfile",
			GoVersionSearchString: `ARG GO_VERSION=(1\.\d\d)`,
		},
		"emissary-ingress/emissary": {
			SourceOfTruthFile:     "go.mod",
			GoVersionSearchString: `go (1\.\d\d)`,
		},
		"fluxcd/flux2": {
			SourceOfTruthFile:     ".github/workflows/release.yaml",
			GoVersionSearchString: `go-version: (1\.\d\d)`,
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
		"goharbor/harbor": {
			SourceOfTruthFile:     "Makefile",
			GoVersionSearchString: `GOBUILDIMAGE=golang:(1\.\d\d)`,
		},
		"goharbor/harbor-scanner-trivy": {
			SourceOfTruthFile:     "go.mod",
			GoVersionSearchString: `go (1\.\d\d)`,
		},
		"helm/helm": {
			SourceOfTruthFile:     ".github/workflows/release.yaml",
			GoVersionSearchString: `go-version: '(1\.\d\d)'`,
		},
		"linuxkit/linuxkit": {
			SourceOfTruthFile:     ".github/workflows/release.yml",
			GoVersionSearchString: `go-version: '(1\.\d\d)'`,
		},
		"kube-vip/kube-vip": {
			SourceOfTruthFile:     "Dockerfile",
			GoVersionSearchString: `golang:(1\.\d\d)`,
		},
		"kubernetes/autoscaler": {
			SourceOfTruthFile:     "cluster-autoscaler/go.mod",
			GoVersionSearchString: `go (1\.\d\d)`,
		},
		"kubernetes/cloud-provider-aws": {
			SourceOfTruthFile:     "go.mod",
			GoVersionSearchString: `go (1\.\d\d)`,
		},
		"kubernetes/cloud-provider-vsphere": {
			SourceOfTruthFile:     "go.mod",
			GoVersionSearchString: `go (1\.\d\d)`,
		},
		"kubernetes-sigs/cluster-api": {
			SourceOfTruthFile:     "Makefile",
			GoVersionSearchString: `GO_VERSION \?= (1\.\d\d)`,
		},
		"kubernetes-sigs/cluster-api-provider-cloudstack": {
			SourceOfTruthFile:     "go.mod",
			GoVersionSearchString: `go (1\.\d\d)`,
		},
		"kubernetes-sigs/cluster-api-provider-vsphere": {
			SourceOfTruthFile:     "Makefile",
			GoVersionSearchString: `GO_VERSION \?= (1\.\d\d)`,
		},
		"kubernetes-sigs/cri-tools": {
			SourceOfTruthFile:     ".github/workflows/release.yml",
			GoVersionSearchString: `go-version: '(1\.\d\d)'`,
		},
		"kubernetes-sigs/etcdadm": {
			SourceOfTruthFile:     "Makefile",
			GoVersionSearchString: `GO_IMAGE \?= golang:(1\.\d\d)`,
		},
		"kubernetes-sigs/kind": {
			SourceOfTruthFile:     ".go-version",
			GoVersionSearchString: `(1\.\d\d)`,
		},
		"metallb/metallb": {
			SourceOfTruthFile:     "controller/Dockerfile",
			GoVersionSearchString: `golang:(1\.\d\d)`,
		},
		"nutanix-cloud-native/cluster-api-provider-nutanix": {
			SourceOfTruthFile:     "go.mod",
			GoVersionSearchString: `go (1\.\d\d)`,
		},
		"nutanix-cloud-native/cloud-provider-nutanix": {
			SourceOfTruthFile:     "go.mod",
			GoVersionSearchString: `go (1\.\d\d)`,
		},
		"opencontainers/runc": {
			SourceOfTruthFile:     "Dockerfile",
			GoVersionSearchString: `ARG GO_VERSION=(1\.\d\d)`,
		},
		"prometheus/node_exporter": {
			SourceOfTruthFile:     ".promu.yml",
			GoVersionSearchString: `version: (1\.\d\d)`,
		},
		"prometheus/prometheus": {
			SourceOfTruthFile:     ".promu.yml",
			GoVersionSearchString: `version: (1\.\d\d)`,
		},
		"rancher/local-path-provisioner": {
			SourceOfTruthFile:     "Dockerfile.dapper",
			GoVersionSearchString: `golang:(1\.\d\d)`,
		},
		"replicatedhq/troubleshoot": {
			SourceOfTruthFile:     ".github/workflows/build-test-deploy.yaml",
			GoVersionSearchString: `go-version: "(1\.\d\d)"`,
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
		"tinkerbell/rufio": {
			SourceOfTruthFile:     "Dockerfile",
			GoVersionSearchString: `golang:(1\.\d\d)`,
		},
		"tinkerbell/hook": {
			SourceOfTruthFile:     "images/hook-bootkit/Dockerfile",
			GoVersionSearchString: `golang:(1\.\d\d)`,
		},
		"vmware/govmomi": {
			SourceOfTruthFile:     ".github/workflows/govmomi-release.yaml",
			GoVersionSearchString: `go-version: '(1\.\d\d)'`,
		},
	}

	DefaultProjectUpgradePRLabels  = []string{"/hold", "/area dependencies"}
	PackagesProjectUpgradePRLabels = []string{"/hold", "/area dependencies", "/sig curated-packages"}

	ProjectsWithUnconventionalUpgradeFlows = []string{
		"kubernetes-sigs/image-builder",
	}

	BottlerocketImageFormats = []string{"ami", "ova", "raw"}

	BottlerocketHostContainers = []string{"admin", "control"}

	CiliumImageDirectories = []string{"cilium", "operator-generic", "cilium-chart"}

	ADOTImageDirectories = []string{"collector"}

	ProjectsSupportingPrereleaseTags = []string{"kubernetes-sigs/cluster-api-provider-cloudstack"}

	// These projects will be upgraded only on main and won't be triggered on release branches.
	CuratedPackagesProjects = []string{
		"aquasecurity/harbor-scanner-trivy",
		"aquasecurity/trivy",
		"aws/rolesanywhere-credential-helper",
		"aws-containers/hello-eks-anywhere",
		"aws-observability/aws-otel-collector",
		"distribution/distribution",
		"emissary-ingress/emissary",
		"goharbor/harbor",
		"goharbor/harbor-scanner-trivy",
		"kubernetes/autoscaler",
		"kubernetes/cloud-provider-aws",
		"kubernetes-sigs/metrics-server",
		"metallb/metallb",
		"prometheus/node_exporter",
		"prometheus/prometheus",
		"redis/redis",
	}

	ProjectsUpgradedOnlyOnMainBranch = append(
		[]string{
			"kubernetes-sigs/cluster-api",
		},
		CuratedPackagesProjects...,
	)

	ProjectMaximumSemvers = map[string]string{
		"containerd/containerd": "v1",
		"opencontainers/runc":   "v1.1",
		"prometheus/prometheus": "v2",
	}

	ECRImageRepositories = map[string]string{
		"cilium/cilium":    CiliumImageRepository,
		"envoyproxy/envoy": EnvoyImageRepository,
	}
)
