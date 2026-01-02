


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ Binary Targets
binaries: ## Build all binaries: `upgrader` for `linux/amd64 linux/arm64`
_output/bin/upgrader/linux-amd64/upgrader: ## Build `_output/bin/upgrader/linux-amd64/upgrader`
_output/bin/upgrader/linux-arm64/upgrader: ## Build `_output/bin/upgrader/linux-arm64/upgrader`

##@ Image Targets
local-images: ## Builds `upgrader/images/amd64` as oci tars for presumbit validation
images: ## Pushes `upgrader/images/push` to IMAGE_REPO
upgrader/images/amd64: ## Builds/pushes `upgrader/images/amd64`
upgrader/images/push: ## Builds/pushes `upgrader/images/push`

##@ Fetch Binary Targets
_output/1-26/dependencies/linux-amd64/eksd/kubernetes/client: ## Fetch `_output/1-26/dependencies/linux-amd64/eksd/kubernetes/client`
_output/1-26/dependencies/linux-arm64/eksd/kubernetes/client: ## Fetch `_output/1-26/dependencies/linux-arm64/eksd/kubernetes/client`
_output/1-26/dependencies/linux-amd64/eksd/kubernetes/server: ## Fetch `_output/1-26/dependencies/linux-amd64/eksd/kubernetes/server`
_output/1-26/dependencies/linux-arm64/eksd/kubernetes/server: ## Fetch `_output/1-26/dependencies/linux-arm64/eksd/kubernetes/server`
_output/1-26/dependencies/linux-amd64/eksd/cni-plugins: ## Fetch `_output/1-26/dependencies/linux-amd64/eksd/cni-plugins`
_output/1-26/dependencies/linux-arm64/eksd/cni-plugins: ## Fetch `_output/1-26/dependencies/linux-arm64/eksd/cni-plugins`
_output/1-26/dependencies/linux-amd64/eksa/containerd/containerd: ## Fetch `_output/1-26/dependencies/linux-amd64/eksa/containerd/containerd`
_output/1-26/dependencies/linux-arm64/eksa/containerd/containerd: ## Fetch `_output/1-26/dependencies/linux-arm64/eksa/containerd/containerd`
_output/1-26/dependencies/linux-amd64/eksa/kubernetes-sigs/cri-tools: ## Fetch `_output/1-26/dependencies/linux-amd64/eksa/kubernetes-sigs/cri-tools`
_output/1-26/dependencies/linux-arm64/eksa/kubernetes-sigs/cri-tools: ## Fetch `_output/1-26/dependencies/linux-arm64/eksa/kubernetes-sigs/cri-tools`

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
run-in-docker/upgrader/../eks-anywhere-go-mod-download: ## Run `upgrader/../eks-anywhere-go-mod-download` in docker builder container
run-in-docker/_output/bin/upgrader/linux-amd64/upgrader: ## Run `_output/bin/upgrader/linux-amd64/upgrader` in docker builder container
run-in-docker/_output/bin/upgrader/linux-arm64/upgrader: ## Run `_output/bin/upgrader/linux-arm64/upgrader` in docker builder container
run-in-docker/_output/attribution/go-license.csv: ## Run `_output/attribution/go-license.csv` in docker builder container

##@ License Targets
gather-licenses: ## Helper to call $(GATHER_LICENSES_TARGETS) which gathers all licenses
attribution: ## Generates attribution from licenses gathered during `gather-licenses`.
attribution-pr: ## Generates PR to update attribution files for projects
attribution-checksums: ## Update attribution and checksums files.
all-attributions: ## Update attribution files for all RELEASE_BRANCHes.
all-attributions-checksums: ## Update attribution and checksums files for all RELEASE_BRANCHes.

##@ Clean Targets
clean: ## Removes source and _output directory
clean-go-cache: ## Removes the GOMODCACHE AND GOCACHE folders

##@Fetch Binary Targets
handle-dependencies: ## Download and extract TARs for each dependency listed in PROJECT_DEPENDENCIES

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
build: ## Called via prow presubmit, calls `github-rate-limit-pre validate-checksums attribution local-images   attribution-pr github-rate-limit-post`
release: ## Called via prow postsubmit + release jobs, calls `validate-checksums images  `
########### END GENERATED ###########################
