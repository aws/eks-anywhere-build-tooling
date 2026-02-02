


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `harbor`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file
patch-repo: ## Patch upstream repo with patches in patches directory

##@ Binary Targets
binaries: ## Build all binaries: `harbor-core harbor-jobservice harbor-registryctl harbor-migrate harbor-exporter` for `linux/amd64 linux/arm64`
_output/bin/harbor/linux-amd64/harbor-core: ## Build `_output/bin/harbor/linux-amd64/harbor-core`
_output/bin/harbor/linux-amd64/harbor-jobservice: ## Build `_output/bin/harbor/linux-amd64/harbor-jobservice`
_output/bin/harbor/linux-amd64/harbor-registryctl: ## Build `_output/bin/harbor/linux-amd64/harbor-registryctl`
_output/bin/harbor/linux-amd64/harbor-migrate: ## Build `_output/bin/harbor/linux-amd64/harbor-migrate`
_output/bin/harbor/linux-amd64/harbor-exporter: ## Build `_output/bin/harbor/linux-amd64/harbor-exporter`
_output/bin/harbor/linux-arm64/harbor-core: ## Build `_output/bin/harbor/linux-arm64/harbor-core`
_output/bin/harbor/linux-arm64/harbor-jobservice: ## Build `_output/bin/harbor/linux-arm64/harbor-jobservice`
_output/bin/harbor/linux-arm64/harbor-registryctl: ## Build `_output/bin/harbor/linux-arm64/harbor-registryctl`
_output/bin/harbor/linux-arm64/harbor-migrate: ## Build `_output/bin/harbor/linux-arm64/harbor-migrate`
_output/bin/harbor/linux-arm64/harbor-exporter: ## Build `_output/bin/harbor/linux-arm64/harbor-exporter`

##@ Image Targets
local-images: ## Builds `harbor-db/images/amd64 harbor-portal/images/amd64 harbor-core/images/amd64 harbor-log/images/amd64 harbor-nginx/images/amd64 harbor-jobservice/images/amd64 harbor-registry/images/amd64 harbor-registryctl/images/amd64 harbor-redis/images/amd64 harbor-exporter/images/amd64 harbor-trivy/images/amd64` as oci tars for presumbit validation
images: ## Pushes `harbor-db/images/push harbor-portal/images/push harbor-core/images/push harbor-log/images/push harbor-nginx/images/push harbor-jobservice/images/push harbor-registry/images/push harbor-registryctl/images/push harbor-redis/images/push harbor-exporter/images/push harbor-trivy/images/push` to IMAGE_REPO
harbor-db/images/amd64: ## Builds/pushes `harbor-db/images/amd64`
harbor-portal/images/amd64: ## Builds/pushes `harbor-portal/images/amd64`
harbor-core/images/amd64: ## Builds/pushes `harbor-core/images/amd64`
harbor-log/images/amd64: ## Builds/pushes `harbor-log/images/amd64`
harbor-nginx/images/amd64: ## Builds/pushes `harbor-nginx/images/amd64`
harbor-jobservice/images/amd64: ## Builds/pushes `harbor-jobservice/images/amd64`
harbor-registry/images/amd64: ## Builds/pushes `harbor-registry/images/amd64`
harbor-registryctl/images/amd64: ## Builds/pushes `harbor-registryctl/images/amd64`
harbor-redis/images/amd64: ## Builds/pushes `harbor-redis/images/amd64`
harbor-exporter/images/amd64: ## Builds/pushes `harbor-exporter/images/amd64`
harbor-trivy/images/amd64: ## Builds/pushes `harbor-trivy/images/amd64`
harbor-db/images/push: ## Builds/pushes `harbor-db/images/push`
harbor-portal/images/push: ## Builds/pushes `harbor-portal/images/push`
harbor-core/images/push: ## Builds/pushes `harbor-core/images/push`
harbor-log/images/push: ## Builds/pushes `harbor-log/images/push`
harbor-nginx/images/push: ## Builds/pushes `harbor-nginx/images/push`
harbor-jobservice/images/push: ## Builds/pushes `harbor-jobservice/images/push`
harbor-registry/images/push: ## Builds/pushes `harbor-registry/images/push`
harbor-registryctl/images/push: ## Builds/pushes `harbor-registryctl/images/push`
harbor-redis/images/push: ## Builds/pushes `harbor-redis/images/push`
harbor-exporter/images/push: ## Builds/pushes `harbor-exporter/images/push`
harbor-trivy/images/push: ## Builds/pushes `harbor-trivy/images/push`

##@ Helm Targets
helm/build: ## Build helm chart
helm/push: ## Build helm chart and push to registry defined in IMAGE_REPO.

##@ Fetch Binary Targets
_output/dependencies/linux-amd64/eksa/goharbor/distribution: ## Fetch `_output/dependencies/linux-amd64/eksa/goharbor/distribution`
_output/dependencies/linux-arm64/eksa/goharbor/distribution: ## Fetch `_output/dependencies/linux-arm64/eksa/goharbor/distribution`
_output/dependencies/linux-amd64/eksa/aquasecurity/trivy: ## Fetch `_output/dependencies/linux-amd64/eksa/aquasecurity/trivy`
_output/dependencies/linux-arm64/eksa/aquasecurity/trivy: ## Fetch `_output/dependencies/linux-arm64/eksa/aquasecurity/trivy`
_output/dependencies/linux-amd64/eksa/goharbor/harbor-scanner-trivy: ## Fetch `_output/dependencies/linux-amd64/eksa/goharbor/harbor-scanner-trivy`
_output/dependencies/linux-arm64/eksa/goharbor/harbor-scanner-trivy: ## Fetch `_output/dependencies/linux-arm64/eksa/goharbor/harbor-scanner-trivy`

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
run-in-docker/harbor/src/eks-anywhere-go-mod-download: ## Run `harbor/src/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/_output/bin/harbor/linux-amd64/harbor-core: ## Run `_output/bin/harbor/linux-amd64/harbor-core` in docker builder container
run-in-docker/_output/bin/harbor/linux-amd64/harbor-jobservice: ## Run `_output/bin/harbor/linux-amd64/harbor-jobservice` in docker builder container
run-in-docker/_output/bin/harbor/linux-amd64/harbor-registryctl: ## Run `_output/bin/harbor/linux-amd64/harbor-registryctl` in docker builder container
run-in-docker/_output/bin/harbor/linux-amd64/harbor-migrate: ## Run `_output/bin/harbor/linux-amd64/harbor-migrate` in docker builder container
run-in-docker/_output/bin/harbor/linux-amd64/harbor-exporter: ## Run `_output/bin/harbor/linux-amd64/harbor-exporter` in docker builder container
run-in-docker/_output/bin/harbor/linux-arm64/harbor-core: ## Run `_output/bin/harbor/linux-arm64/harbor-core` in docker builder container
run-in-docker/_output/bin/harbor/linux-arm64/harbor-jobservice: ## Run `_output/bin/harbor/linux-arm64/harbor-jobservice` in docker builder container
run-in-docker/_output/bin/harbor/linux-arm64/harbor-registryctl: ## Run `_output/bin/harbor/linux-arm64/harbor-registryctl` in docker builder container
run-in-docker/_output/bin/harbor/linux-arm64/harbor-migrate: ## Run `_output/bin/harbor/linux-arm64/harbor-migrate` in docker builder container
run-in-docker/_output/bin/harbor/linux-arm64/harbor-exporter: ## Run `_output/bin/harbor/linux-arm64/harbor-exporter` in docker builder container
run-in-docker/_output/attribution/go-license.csv: ## Run `_output/attribution/go-license.csv` in docker builder container
run-in-docker/_output/attribution/go-license.csv: ## Run `_output/attribution/go-license.csv` in docker builder container
run-in-docker/_output/attribution/go-license.csv: ## Run `_output/attribution/go-license.csv` in docker builder container
run-in-docker/_output/attribution/go-license.csv: ## Run `_output/attribution/go-license.csv` in docker builder container
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
