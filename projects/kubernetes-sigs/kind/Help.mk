


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `kind`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file
patch-repo: ## Patch upstream repo with patches in patches directory

##@ Binary Targets
binaries: ## Build all binaries: `` for `linux/amd64 linux/arm64 darwin/amd64 darwin/arm64`
_output/bin/kind/linux-amd64/kind: ## Build `_output/bin/kind/linux-amd64/kind`
_output/bin/kind/linux-arm64/kind: ## Build `_output/bin/kind/linux-arm64/kind`
_output/bin/kind/darwin-amd64/kind: ## Build `_output/bin/kind/darwin-amd64/kind`
_output/bin/kind/darwin-arm64/kind: ## Build `_output/bin/kind/darwin-arm64/kind`
_output/bin/kind/linux-amd64/kindnetd: ## Build `_output/bin/kind/linux-amd64/kindnetd`
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
_output/1-21/dependencies/linux-amd64/eksd/kubernetes/client: ## Fetch `_output/1-21/dependencies/linux-amd64/eksd/kubernetes/client`
_output/1-21/dependencies/linux-amd64/eksa/kubernetes-sigs/etcdadm: ## Fetch `_output/1-21/dependencies/linux-amd64/eksa/kubernetes-sigs/etcdadm`
_output/1-21/dependencies/linux-amd64/eksd/cni-plugins: ## Fetch `_output/1-21/dependencies/linux-amd64/eksd/cni-plugins`
_output/1-21/dependencies/linux-amd64/eksa/kubernetes-sigs/cri-tools: ## Fetch `_output/1-21/dependencies/linux-amd64/eksa/kubernetes-sigs/cri-tools`
_output/1-21/dependencies/linux-amd64/eksd/etcd/etcd.tar.gz: ## Fetch `_output/1-21/dependencies/linux-amd64/eksd/etcd/etcd.tar.gz`
_output/1-21/dependencies/linux-arm64/eksd/kubernetes/client: ## Fetch `_output/1-21/dependencies/linux-arm64/eksd/kubernetes/client`
_output/1-21/dependencies/linux-arm64/eksa/kubernetes-sigs/etcdadm: ## Fetch `_output/1-21/dependencies/linux-arm64/eksa/kubernetes-sigs/etcdadm`
_output/1-21/dependencies/linux-arm64/eksd/cni-plugins: ## Fetch `_output/1-21/dependencies/linux-arm64/eksd/cni-plugins`
_output/1-21/dependencies/linux-arm64/eksa/kubernetes-sigs/cri-tools: ## Fetch `_output/1-21/dependencies/linux-arm64/eksa/kubernetes-sigs/cri-tools`
_output/1-21/dependencies/linux-arm64/eksd/etcd/etcd.tar.gz: ## Fetch `_output/1-21/dependencies/linux-arm64/eksd/etcd/etcd.tar.gz`

##@ Checksum Targets
checksums: ## Update checksums file based on currently built binaries.
validate-checksums: # Validate checksums of currently built binaries against checksums file.

##@ Artifact Targets
tarballs: ## Create tarballs by calling build/lib/simple_create_tarballs.sh unless SIMPLE_CREATE_TARBALLS=false, then tarballs must be defined in project Makefile
s3-artifacts: # Prepare ARTIFACTS_PATH folder structure with tarballs/manifests/other items to be uploaded to s3
upload-artifacts: # Upload tarballs and other artifacts from ARTIFACTS_PATH to S3

##@ License Targets
gather-licenses: ## Helper to call $(GATHER_LICENSES_TARGETS) which gathers all licenses
attribution: ## Generates attribution from licenses gathered during `gather-licenses`.
attribution-pr: ## Generates PR to update attribution files for projects

##@ Clean Targets
clean: ## Removes source and _output directory
clean-repo: ## Removes source directory

##@ Helpers
help: ## Display this help
add-generated-help-block: ## Add or update generated help block to document project make file and support shell auto completion

##@ Build Targets
build: ## Called via prow presubmit, calls `validate-checksums local-images attribution upload-artifacts attribution-pr`
release: ## Called via prow postsubmit + release jobs, calls `validate-checksums images upload-artifacts`
########### END GENERATED ###########################
