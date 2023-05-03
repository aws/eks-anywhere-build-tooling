


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `cluster-api`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file
patch-repo: ## Patch upstream repo with patches in patches directory

##@ Binary Targets
binaries: ## Build all binaries: `clusterctl manager kubeadm-bootstrap-manager kubeadm-control-plane-manager cluster-api-provider-docker-manager` for `linux/amd64 linux/arm64`
_output/bin/cluster-api/linux-amd64/clusterctl: ## Build `_output/bin/cluster-api/linux-amd64/clusterctl`
_output/bin/cluster-api/linux-amd64/manager: ## Build `_output/bin/cluster-api/linux-amd64/manager`
_output/bin/cluster-api/linux-amd64/kubeadm-bootstrap-manager: ## Build `_output/bin/cluster-api/linux-amd64/kubeadm-bootstrap-manager`
_output/bin/cluster-api/linux-amd64/kubeadm-control-plane-manager: ## Build `_output/bin/cluster-api/linux-amd64/kubeadm-control-plane-manager`
_output/bin/cluster-api/linux-amd64/cluster-api-provider-docker-manager: ## Build `_output/bin/cluster-api/linux-amd64/cluster-api-provider-docker-manager`
_output/bin/cluster-api/linux-arm64/clusterctl: ## Build `_output/bin/cluster-api/linux-arm64/clusterctl`
_output/bin/cluster-api/linux-arm64/manager: ## Build `_output/bin/cluster-api/linux-arm64/manager`
_output/bin/cluster-api/linux-arm64/kubeadm-bootstrap-manager: ## Build `_output/bin/cluster-api/linux-arm64/kubeadm-bootstrap-manager`
_output/bin/cluster-api/linux-arm64/kubeadm-control-plane-manager: ## Build `_output/bin/cluster-api/linux-arm64/kubeadm-control-plane-manager`
_output/bin/cluster-api/linux-arm64/cluster-api-provider-docker-manager: ## Build `_output/bin/cluster-api/linux-arm64/cluster-api-provider-docker-manager`

##@ Image Targets
local-images: ## Builds `cluster-api-controller/images/amd64 kubeadm-bootstrap-controller/images/amd64 kubeadm-control-plane-controller/images/amd64 cluster-api-docker-controller/images/amd64` as oci tars for presumbit validation
images: ## Pushes `cluster-api-controller/images/push kubeadm-bootstrap-controller/images/push kubeadm-control-plane-controller/images/push cluster-api-docker-controller/images/push` to IMAGE_REPO
cluster-api-controller/images/amd64: ## Builds/pushes `cluster-api-controller/images/amd64`
kubeadm-bootstrap-controller/images/amd64: ## Builds/pushes `kubeadm-bootstrap-controller/images/amd64`
kubeadm-control-plane-controller/images/amd64: ## Builds/pushes `kubeadm-control-plane-controller/images/amd64`
cluster-api-docker-controller/images/amd64: ## Builds/pushes `cluster-api-docker-controller/images/amd64`
cluster-api-controller/images/push: ## Builds/pushes `cluster-api-controller/images/push`
kubeadm-bootstrap-controller/images/push: ## Builds/pushes `kubeadm-bootstrap-controller/images/push`
kubeadm-control-plane-controller/images/push: ## Builds/pushes `kubeadm-control-plane-controller/images/push`
cluster-api-docker-controller/images/push: ## Builds/pushes `cluster-api-docker-controller/images/push`

##@ Fetch Binary Targets
_output/dependencies/linux-amd64/eksd/kubernetes/client: ## Fetch `_output/dependencies/linux-amd64/eksd/kubernetes/client`
_output/dependencies/linux-arm64/eksd/kubernetes/client: ## Fetch `_output/dependencies/linux-arm64/eksd/kubernetes/client`

##@ Checksum Targets
checksums: ## Update checksums file based on currently built binaries.
validate-checksums: # Validate checksums of currently built binaries against checksums file.
all-checksums: ## Update checksums files for all RELEASE_BRANCHes.

##@ Run in Docker Targets
run-all-attributions-in-docker: ## Run `all-attributions` in docker builder container
run-all-attributions-checksums-in-docker: ## Run `all-attributions-checksums` in docker builder container
run-all-checksums-in-docker: ## Run `all-checksums` in docker builder container
run-attribution-in-docker: ## Run `attribution` in docker builder container
run-attribution-checksums-in-docker: ## Run `attribution-checksums` in docker builder container
run-binaries-in-docker: ## Run `binaries` in docker builder container
run-checksums-in-docker: ## Run `checksums` in docker builder container
run-clean-in-docker: ## Run `clean` in docker builder container
run-clean-go-cache-in-docker: ## Run `clean-go-cache` in docker builder container

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

##@ Helpers
help: ## Display this help
add-generated-help-block: ## Add or update generated help block to document project make file and support shell auto completion

##@Update Helpers
run-target-in-docker: ## Run `MAKE_TARGET` using builder base docker container
stop-docker-builder: ## Clean up builder base docker container
generate: ## Update UPSTREAM_PROJECTS.yaml
update-go-mods: ## Update locally checked-in go sum to assist in vuln scanning
update-vendor-for-dep-patch: ## After bumping dep in go.mod file, uses generic vendor update script or one provided from upstream project
patch-for-dep-update: ## After bumping dep in go.mod file and updating vendor, generates patch
create-ecr-repos: ## Create repos in ECR for project images for local testing

##@ Build Targets
build: ## Called via prow presubmit, calls `github-rate-limit-pre validate-checksums attribution local-images  upload-artifacts attribution-pr github-rate-limit-post`
release: ## Called via prow postsubmit + release jobs, calls `validate-checksums images  upload-artifacts`
########### END GENERATED ###########################
