


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `flux2`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file

##@ Binary Targets
binaries: ## Build all binaries: `flux` for `linux/amd64 linux/arm64 darwin/amd64 darwin/arm64`
_output/bin/flux2/linux-amd64/flux: ## Build `_output/bin/flux2/linux-amd64/flux`
_output/bin/flux2/linux-arm64/flux: ## Build `_output/bin/flux2/linux-arm64/flux`
_output/bin/flux2/darwin-amd64/flux: ## Build `_output/bin/flux2/darwin-amd64/flux`
_output/bin/flux2/darwin-arm64/flux: ## Build `_output/bin/flux2/darwin-arm64/flux`

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

##@ Helpers
help: ## Display this help
add-generated-help-block: ## Add or update generated help block to document project make file and support shell auto completion

##@ Build Targets
build: ## Called via prow presubmit, calls `validate-checksums  attribution attribution-pr upload-artifacts`
release: ## Called via prow postsubmit + release jobs, calls `validate-checksums  upload-artifacts`
########### END GENERATED ###########################
