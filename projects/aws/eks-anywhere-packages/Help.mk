


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `eks-anywhere-packages`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file

##@ Binary Targets
binaries: ## Build all binaries: `package-manager ecrtokenrefresher credentialprovider` for `linux/amd64 linux/arm64`
_output/bin/eks-anywhere-packages/linux-amd64/package-manager: ## Build `_output/bin/eks-anywhere-packages/linux-amd64/package-manager`
_output/bin/eks-anywhere-packages/linux-amd64/ecrtokenrefresher: ## Build `_output/bin/eks-anywhere-packages/linux-amd64/ecrtokenrefresher`
_output/bin/eks-anywhere-packages/linux-amd64/credentialprovider: ## Build `_output/bin/eks-anywhere-packages/linux-amd64/credentialprovider`
_output/bin/eks-anywhere-packages/linux-arm64/package-manager: ## Build `_output/bin/eks-anywhere-packages/linux-arm64/package-manager`
_output/bin/eks-anywhere-packages/linux-arm64/ecrtokenrefresher: ## Build `_output/bin/eks-anywhere-packages/linux-arm64/ecrtokenrefresher`
_output/bin/eks-anywhere-packages/linux-arm64/credentialprovider: ## Build `_output/bin/eks-anywhere-packages/linux-arm64/credentialprovider`

##@ Image Targets
local-images: ## Builds `eks-anywhere-packages/images/amd64 ecrtokenrefresher/images/amd64 credentialprovider/images/amd64` as oci tars for presumbit validation
images: ## Pushes `eks-anywhere-packages/images/push ecrtokenrefresher/images/push credentialprovider/images/push` to IMAGE_REPO
eks-anywhere-packages/images/amd64: ## Builds/pushes `eks-anywhere-packages/images/amd64`
ecrtokenrefresher/images/amd64: ## Builds/pushes `ecrtokenrefresher/images/amd64`
credentialprovider/images/amd64: ## Builds/pushes `credentialprovider/images/amd64`
eks-anywhere-packages/images/push: ## Builds/pushes `eks-anywhere-packages/images/push`
ecrtokenrefresher/images/push: ## Builds/pushes `ecrtokenrefresher/images/push`
credentialprovider/images/push: ## Builds/pushes `credentialprovider/images/push`

##@ Helm Targets
helm/build: ## Build helm chart
helm/push: ## Build helm chart and push to registry defined in IMAGE_REPO.

##@ Fetch Binary Targets
_output/dependencies/linux-amd64/eksd/kubernetes/client: ## Fetch `_output/dependencies/linux-amd64/eksd/kubernetes/client`
_output/dependencies/linux-arm64/eksd/kubernetes/client: ## Fetch `_output/dependencies/linux-arm64/eksd/kubernetes/client`
_output/dependencies/linux-amd64/eksa/aws/rolesanywhere-credential-helper: ## Fetch `_output/dependencies/linux-amd64/eksa/aws/rolesanywhere-credential-helper`
_output/dependencies/linux-arm64/eksa/aws/rolesanywhere-credential-helper: ## Fetch `_output/dependencies/linux-arm64/eksa/aws/rolesanywhere-credential-helper`
_output/dependencies/linux-amd64/eksa/kubernetes/cloud-provider-aws/1-29: ## Fetch `_output/dependencies/linux-amd64/eksa/kubernetes/cloud-provider-aws/1-29`
_output/dependencies/linux-arm64/eksa/kubernetes/cloud-provider-aws/1-29: ## Fetch `_output/dependencies/linux-arm64/eksa/kubernetes/cloud-provider-aws/1-29`

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
run-in-docker/eks-anywhere-packages/eks-anywhere-go-mod-download: ## Run `eks-anywhere-packages/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/eks-anywhere-packages/ecrtokenrefresher/eks-anywhere-go-mod-download: ## Run `eks-anywhere-packages/ecrtokenrefresher/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/eks-anywhere-packages/credentialproviderpackage/eks-anywhere-go-mod-download: ## Run `eks-anywhere-packages/credentialproviderpackage/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/_output/bin/eks-anywhere-packages/linux-amd64/package-manager: ## Run `_output/bin/eks-anywhere-packages/linux-amd64/package-manager` in docker builder container
run-in-docker/_output/bin/eks-anywhere-packages/linux-amd64/ecrtokenrefresher: ## Run `_output/bin/eks-anywhere-packages/linux-amd64/ecrtokenrefresher` in docker builder container
run-in-docker/_output/bin/eks-anywhere-packages/linux-amd64/credentialprovider: ## Run `_output/bin/eks-anywhere-packages/linux-amd64/credentialprovider` in docker builder container
run-in-docker/_output/bin/eks-anywhere-packages/linux-arm64/package-manager: ## Run `_output/bin/eks-anywhere-packages/linux-arm64/package-manager` in docker builder container
run-in-docker/_output/bin/eks-anywhere-packages/linux-arm64/ecrtokenrefresher: ## Run `_output/bin/eks-anywhere-packages/linux-arm64/ecrtokenrefresher` in docker builder container
run-in-docker/_output/bin/eks-anywhere-packages/linux-arm64/credentialprovider: ## Run `_output/bin/eks-anywhere-packages/linux-arm64/credentialprovider` in docker builder container
run-in-docker/_output/attribution/go-license.csv: ## Run `_output/attribution/go-license.csv` in docker builder container
run-in-docker/_output/ecrtokenrefresher/attribution/go-license.csv: ## Run `_output/ecrtokenrefresher/attribution/go-license.csv` in docker builder container
run-in-docker/_output/credentialprovider/attribution/go-license.csv: ## Run `_output/credentialprovider/attribution/go-license.csv` in docker builder container

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
