


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `cloudstack-cloudmonkey`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file

##@ Binary Targets
binaries: ## Build all binaries: `cmk` for `linux/amd64 linux/arm64`
_output/bin/cloudstack-cloudmonkey/linux-amd64/cmk: ## Build `_output/bin/cloudstack-cloudmonkey/linux-amd64/cmk`
_output/bin/cloudstack-cloudmonkey/linux-arm64/cmk: ## Build `_output/bin/cloudstack-cloudmonkey/linux-arm64/cmk`

##@ Image Targets
local-images: ## Builds `` as oci tars for presumbit validation
images: ## Pushes `` to IMAGE_REPO

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
helm/build: ## Build helm chart
helm/push: ## Build helm chart and push to registry defined in IMAGE_REPO.

##@ Helpers
help: ## Display this help
add-generated-help-block: ## Add or update generated help block to document project make file and support shell auto completion

##@ Build Targets
build: ## Called via prow presubmit, calls `validate-checksums  attribution attribution-pr upload-artifacts`
release: ## Called via prow postsubmit + release jobs, calls `validate-checksums  upload-artifacts`
########### END GENERATED ###########################
