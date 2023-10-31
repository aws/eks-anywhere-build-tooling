


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `kind`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file
patch-repo: ## Patch upstream repo with patches in patches directory

##@ Binary Targets
binaries: ## Build all binaries: `kind kindnetd` for `linux/amd64 linux/arm64`
_output/bin/kind/linux-amd64/kind: ## Build `_output/bin/kind/linux-amd64/kind`
_output/bin/kind/linux-amd64/kindnetd: ## Build `_output/bin/kind/linux-amd64/kindnetd`
_output/bin/kind/linux-arm64/kind: ## Build `_output/bin/kind/linux-arm64/kind`
_output/bin/kind/linux-arm64/kindnetd: ## Build `_output/bin/kind/linux-arm64/kindnetd`

##@ Image Targets
local-images: ## Builds `haproxy/images/amd64 kindnetd/images/amd64 kind-base/images/amd64` as oci tars for presumbit validation
images: ## Pushes `haproxy/images/push kindnetd/images/push kind-base/images/push` to IMAGE_REPO
haproxy/images/amd64: ## Builds/pushes `haproxy/images/amd64`
kindnetd/images/amd64: ## Builds/pushes `kindnetd/images/amd64`
kind-base/images/amd64: ## Builds/pushes `kind-base/images/amd64`
haproxy/images/push: ## Builds/pushes `haproxy/images/push`
kindnetd/images/push: ## Builds/pushes `kindnetd/images/push`
kind-base/images/push: ## Builds/pushes `kind-base/images/push`

##@ Fetch Binary Targets
_output/1-26/dependencies/linux-amd64/eksd/kubernetes/client: ## Fetch `_output/1-26/dependencies/linux-amd64/eksd/kubernetes/client`
_output/1-26/dependencies/linux-arm64/eksd/kubernetes/client: ## Fetch `_output/1-26/dependencies/linux-arm64/eksd/kubernetes/client`
_output/1-26/dependencies/linux-amd64/eksa/kubernetes-sigs/etcdadm: ## Fetch `_output/1-26/dependencies/linux-amd64/eksa/kubernetes-sigs/etcdadm`
_output/1-26/dependencies/linux-arm64/eksa/kubernetes-sigs/etcdadm: ## Fetch `_output/1-26/dependencies/linux-arm64/eksa/kubernetes-sigs/etcdadm`
_output/1-26/dependencies/linux-amd64/eksd/cni-plugins: ## Fetch `_output/1-26/dependencies/linux-amd64/eksd/cni-plugins`
_output/1-26/dependencies/linux-arm64/eksd/cni-plugins: ## Fetch `_output/1-26/dependencies/linux-arm64/eksd/cni-plugins`
_output/1-26/dependencies/linux-amd64/eksa/kubernetes-sigs/cri-tools: ## Fetch `_output/1-26/dependencies/linux-amd64/eksa/kubernetes-sigs/cri-tools`
_output/1-26/dependencies/linux-arm64/eksa/kubernetes-sigs/cri-tools: ## Fetch `_output/1-26/dependencies/linux-arm64/eksa/kubernetes-sigs/cri-tools`
_output/1-26/dependencies/linux-amd64/eksd/etcd/etcd.tar.gz: ## Fetch `_output/1-26/dependencies/linux-amd64/eksd/etcd/etcd.tar.gz`
_output/1-26/dependencies/linux-arm64/eksd/etcd/etcd.tar.gz: ## Fetch `_output/1-26/dependencies/linux-arm64/eksd/etcd/etcd.tar.gz`

##@ Checksum Targets
checksums: ## Update checksums file based on currently built binaries.
validate-checksums: # Validate checksums of currently built binaries against checksums file.
all-checksums: ## Update checksums files for all RELEASE_BRANCHes.

##@ Run in Docker Targets
run-in-docker/all-attributions: ## Run `all-attributions` in docker builder container
run-in-docker/all-attributions-checksums: ## Run `all-attributions-checksums` in docker builder container
run-in-docker/all-checksums: ## Run `all-checksums` in docker builder container
run-in-docker/attribution: ## Run `attribution` in docker builder container
run-in-docker/attribution-checksums: ## Run `attribution-checksums` in docker builder container
run-in-docker/binaries: ## Run `binaries` in docker builder container
run-in-docker/checksums: ## Run `checksums` in docker builder container
run-in-docker/clean: ## Run `clean` in docker builder container
run-in-docker/clean-go-cache: ## Run `clean-go-cache` in docker builder container
run-in-docker/validate-checksums: ## Run `validate-checksums` in docker builder container
run-in-docker/kind/eks-anywhere-go-mod-download: ## Run `kind/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/kind/./images/kindnetd/eks-anywhere-go-mod-download: ## Run `kind/./images/kindnetd/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/_output/bin/kind/linux-amd64/kind: ## Run `_output/bin/kind/linux-amd64/kind` in docker builder container
run-in-docker/_output/bin/kind/linux-amd64/kindnetd: ## Run `_output/bin/kind/linux-amd64/kindnetd` in docker builder container
run-in-docker/_output/bin/kind/linux-arm64/kind: ## Run `_output/bin/kind/linux-arm64/kind` in docker builder container
run-in-docker/_output/bin/kind/linux-arm64/kindnetd: ## Run `_output/bin/kind/linux-arm64/kindnetd` in docker builder container
run-in-docker/_output/attribution/go-license.csv: ## Run `_output/attribution/go-license.csv` in docker builder container
run-in-docker/_output/kindnetd/attribution/go-license.csv: ## Run `_output/kindnetd/attribution/go-license.csv` in docker builder container

##@ Artifact Targets
tarballs: ## Create tarballs by calling build/lib/simple_create_tarballs.sh unless SIMPLE_CREATE_TARBALLS=false, then tarballs must be defined in project Makefile
s3-artifacts: # Prepare ARTIFACTS_PATH folder structure with tarballs/manifests/other items to be uploaded to s3
upload-artifacts: # Upload tarballs and other artifacts from ARTIFACTS_PATH to S3

##@ License Targets
gather-licenses: ## Helper to call $(GATHER_LICENSES_TARGETS) which gathers all licenses
attribution: ## Generates attribution from licenses gathered during `gather-licenses`.
attribution-pr: ## Generates PR to update attribution files for projects
attribution-checksums: ## Update attribution and checksums files.
all-attributions: ## Update attribution files for all RELEASE_BRANCHes.
all-attributions-checksums: ## Update attribution and checksums files for all RELEASE_BRANCHes.

##@ Clean Targets
clean: ## Removes source and _output directory
clean-repo: ## Removes source directory
create-kind-cluster-amd64: ## Create local kind cluster using built amd64 image
create-kind-cluster-arm64: ## Create local kind cluster using built arm64 image

##@ Helpers
help: ## Display this help
add-generated-help-block: ## Add or update generated help block to document project make file and support shell auto completion

##@Update Helpers
start-docker-builder: ## Start long lived builder base docker container
stop-docker-builder: ## Clean up builder base docker container
run-buildkit-and-registry: ## Run buildkitd and a local docker registry as containers
stop-buildkit-and-registry: ## Stop the buildkitd and a local docker registry containers
generate: ## Update UPSTREAM_PROJECTS.yaml
update-go-mods: ## Update locally checked-in go sum to assist in vuln scanning
update-vendor-for-dep-patch: ## After bumping dep in go.mod file, uses generic vendor update script or one provided from upstream project
patch-for-dep-update: ## After bumping dep in go.mod file and updating vendor, generates patch
create-ecr-repos: ## Create repos in ECR for project images for local testing

##@ Build Targets
build: ## Called via prow presubmit, calls `github-rate-limit-pre validate-checksums attribution local-images  upload-artifacts attribution-pr github-rate-limit-post`
release: ## Called via prow postsubmit + release jobs, calls `validate-checksums images  upload-artifacts`
########### END GENERATED ###########################
