


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `autoscaler`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file
patch-repo: ## Patch upstream repo with patches in patches directory

##@ Binary Targets
binaries: ## Build all binaries: `cluster-autoscaler` for `linux/amd64 linux/arm64`
_output/1-26/bin/autoscaler/linux-amd64/cluster-autoscaler: ## Build `_output/1-26/bin/autoscaler/linux-amd64/cluster-autoscaler`
_output/1-26/bin/autoscaler/linux-arm64/cluster-autoscaler: ## Build `_output/1-26/bin/autoscaler/linux-arm64/cluster-autoscaler`

##@ Image Targets
local-images: ## Builds `autoscaler/images/amd64` as oci tars for presumbit validation
images: ## Pushes `autoscaler/images/push` to IMAGE_REPO
autoscaler/images/amd64: ## Builds/pushes `autoscaler/images/amd64`
autoscaler/images/push: ## Builds/pushes `autoscaler/images/push`

##@ Helm Targets
helm/build: ## Build helm chart
helm/push: ## Build helm chart and push to registry defined in IMAGE_REPO.

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
run-in-docker/autoscaler/cluster-autoscaler/eks-anywhere-go-mod-download: ## Run `autoscaler/cluster-autoscaler/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/_output/1-26/bin/autoscaler/linux-amd64/cluster-autoscaler: ## Run `_output/1-26/bin/autoscaler/linux-amd64/cluster-autoscaler` in docker builder container
run-in-docker/_output/1-26/bin/autoscaler/linux-arm64/cluster-autoscaler: ## Run `_output/1-26/bin/autoscaler/linux-arm64/cluster-autoscaler` in docker builder container
run-in-docker/_output/1-26/attribution/go-license.csv: ## Run `_output/1-26/attribution/go-license.csv` in docker builder container

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
clean-repo: ## Removes source directory

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
build: ## Called via prow presubmit, calls `github-rate-limit-pre validate-checksums attribution local-images helm/build  attribution-pr github-rate-limit-post`
release: ## Called via prow postsubmit + release jobs, calls `validate-checksums images helm/push `
########### END GENERATED ###########################
