


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `cert-manager`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file

##@ Binary Targets
binaries: ## Build all binaries: `cert-manager-acmesolver cert-manager-cainjector cert-manager-controller cert-manager-webhook cert-manager-ctl` for `linux/amd64 linux/arm64`
_output/bin/cert-manager/linux-amd64/cert-manager-acmesolver: ## Build `_output/bin/cert-manager/linux-amd64/cert-manager-acmesolver`
_output/bin/cert-manager/linux-amd64/cert-manager-cainjector: ## Build `_output/bin/cert-manager/linux-amd64/cert-manager-cainjector`
_output/bin/cert-manager/linux-amd64/cert-manager-controller: ## Build `_output/bin/cert-manager/linux-amd64/cert-manager-controller`
_output/bin/cert-manager/linux-amd64/cert-manager-webhook: ## Build `_output/bin/cert-manager/linux-amd64/cert-manager-webhook`
_output/bin/cert-manager/linux-amd64/cert-manager-ctl: ## Build `_output/bin/cert-manager/linux-amd64/cert-manager-ctl`
_output/bin/cert-manager/linux-arm64/cert-manager-acmesolver: ## Build `_output/bin/cert-manager/linux-arm64/cert-manager-acmesolver`
_output/bin/cert-manager/linux-arm64/cert-manager-cainjector: ## Build `_output/bin/cert-manager/linux-arm64/cert-manager-cainjector`
_output/bin/cert-manager/linux-arm64/cert-manager-controller: ## Build `_output/bin/cert-manager/linux-arm64/cert-manager-controller`
_output/bin/cert-manager/linux-arm64/cert-manager-webhook: ## Build `_output/bin/cert-manager/linux-arm64/cert-manager-webhook`
_output/bin/cert-manager/linux-arm64/cert-manager-ctl: ## Build `_output/bin/cert-manager/linux-arm64/cert-manager-ctl`

##@ Image Targets
local-images: ## Builds `cert-manager-acmesolver/images/amd64 cert-manager-cainjector/images/amd64 cert-manager-controller/images/amd64 cert-manager-webhook/images/amd64 cert-manager-ctl/images/amd64` as oci tars for presumbit validation
images: ## Pushes `cert-manager-acmesolver/images/push cert-manager-cainjector/images/push cert-manager-controller/images/push cert-manager-webhook/images/push cert-manager-ctl/images/push` to IMAGE_REPO
cert-manager-acmesolver/images/amd64: ## Builds/pushes `cert-manager-acmesolver/images/amd64`
cert-manager-cainjector/images/amd64: ## Builds/pushes `cert-manager-cainjector/images/amd64`
cert-manager-controller/images/amd64: ## Builds/pushes `cert-manager-controller/images/amd64`
cert-manager-webhook/images/amd64: ## Builds/pushes `cert-manager-webhook/images/amd64`
cert-manager-ctl/images/amd64: ## Builds/pushes `cert-manager-ctl/images/amd64`
cert-manager-acmesolver/images/push: ## Builds/pushes `cert-manager-acmesolver/images/push`
cert-manager-cainjector/images/push: ## Builds/pushes `cert-manager-cainjector/images/push`
cert-manager-controller/images/push: ## Builds/pushes `cert-manager-controller/images/push`
cert-manager-webhook/images/push: ## Builds/pushes `cert-manager-webhook/images/push`
cert-manager-ctl/images/push: ## Builds/pushes `cert-manager-ctl/images/push`

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
run-in-docker/cert-manager/cmd/acmesolver/eks-anywhere-go-mod-download: ## Run `cert-manager/cmd/acmesolver/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/cert-manager/cmd/cainjector/eks-anywhere-go-mod-download: ## Run `cert-manager/cmd/cainjector/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/cert-manager/cmd/controller/eks-anywhere-go-mod-download: ## Run `cert-manager/cmd/controller/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/cert-manager/cmd/webhook/eks-anywhere-go-mod-download: ## Run `cert-manager/cmd/webhook/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/cert-manager/cmd/ctl/eks-anywhere-go-mod-download: ## Run `cert-manager/cmd/ctl/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/_output/bin/cert-manager/linux-amd64/cert-manager-acmesolver: ## Run `_output/bin/cert-manager/linux-amd64/cert-manager-acmesolver` in docker builder container
run-in-docker/_output/bin/cert-manager/linux-amd64/cert-manager-cainjector: ## Run `_output/bin/cert-manager/linux-amd64/cert-manager-cainjector` in docker builder container
run-in-docker/_output/bin/cert-manager/linux-amd64/cert-manager-controller: ## Run `_output/bin/cert-manager/linux-amd64/cert-manager-controller` in docker builder container
run-in-docker/_output/bin/cert-manager/linux-amd64/cert-manager-webhook: ## Run `_output/bin/cert-manager/linux-amd64/cert-manager-webhook` in docker builder container
run-in-docker/_output/bin/cert-manager/linux-amd64/cert-manager-ctl: ## Run `_output/bin/cert-manager/linux-amd64/cert-manager-ctl` in docker builder container
run-in-docker/_output/bin/cert-manager/linux-arm64/cert-manager-acmesolver: ## Run `_output/bin/cert-manager/linux-arm64/cert-manager-acmesolver` in docker builder container
run-in-docker/_output/bin/cert-manager/linux-arm64/cert-manager-cainjector: ## Run `_output/bin/cert-manager/linux-arm64/cert-manager-cainjector` in docker builder container
run-in-docker/_output/bin/cert-manager/linux-arm64/cert-manager-controller: ## Run `_output/bin/cert-manager/linux-arm64/cert-manager-controller` in docker builder container
run-in-docker/_output/bin/cert-manager/linux-arm64/cert-manager-webhook: ## Run `_output/bin/cert-manager/linux-arm64/cert-manager-webhook` in docker builder container
run-in-docker/_output/bin/cert-manager/linux-arm64/cert-manager-ctl: ## Run `_output/bin/cert-manager/linux-arm64/cert-manager-ctl` in docker builder container
run-in-docker/_output/cert-manager-acmesolver/attribution/go-license.csv: ## Run `_output/cert-manager-acmesolver/attribution/go-license.csv` in docker builder container
run-in-docker/_output/cert-manager-cainjector/attribution/go-license.csv: ## Run `_output/cert-manager-cainjector/attribution/go-license.csv` in docker builder container
run-in-docker/_output/cert-manager-controller/attribution/go-license.csv: ## Run `_output/cert-manager-controller/attribution/go-license.csv` in docker builder container
run-in-docker/_output/cert-manager-webhook/attribution/go-license.csv: ## Run `_output/cert-manager-webhook/attribution/go-license.csv` in docker builder container
run-in-docker/_output/cert-manager-ctl/attribution/go-license.csv: ## Run `_output/cert-manager-ctl/attribution/go-license.csv` in docker builder container

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
clean-go-cache: ## Removes the GOMODCACHE AND GOCACHE folders
clean-repo: ## Removes source directory

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
build: ## Called via prow presubmit, calls `github-rate-limit-pre validate-checksums attribution local-images helm/build upload-artifacts attribution-pr github-rate-limit-post`
release: ## Called via prow postsubmit + release jobs, calls `validate-checksums images helm/push upload-artifacts`
########### END GENERATED ###########################
