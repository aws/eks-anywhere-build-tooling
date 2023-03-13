


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `hook`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file
patch-repo: ## Patch upstream repo with patches in patches directory

##@ Binary Targets
binaries: ## Build all binaries: `hook-bootkit hook-docker` for `linux/amd64 linux/arm64`
_output/bin/hook/linux-amd64/hook-bootkit: ## Build `_output/bin/hook/linux-amd64/hook-bootkit`
_output/bin/hook/linux-amd64/hook-docker: ## Build `_output/bin/hook/linux-amd64/hook-docker`
_output/bin/hook/linux-arm64/hook-bootkit: ## Build `_output/bin/hook/linux-arm64/hook-bootkit`
_output/bin/hook/linux-arm64/hook-docker: ## Build `_output/bin/hook/linux-arm64/hook-docker`

##@ Image Targets
local-images: ## Builds `hook-bootkit/images/amd64 hook-docker/images/amd64 kernel/images/amd64` as oci tars for presumbit validation
images: ## Pushes `hook-bootkit/images/push hook-docker/images/push kernel/images/push` to IMAGE_REPO
hook-bootkit/images/amd64: ## Builds/pushes `hook-bootkit/images/amd64`
hook-docker/images/amd64: ## Builds/pushes `hook-docker/images/amd64`
kernel/images/amd64: ## Builds/pushes `kernel/images/amd64`
hook-bootkit/images/push: ## Builds/pushes `hook-bootkit/images/push`
hook-docker/images/push: ## Builds/pushes `hook-docker/images/push`
kernel/images/push: ## Builds/pushes `kernel/images/push`

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

##@Update Helpers
run-target-in-docker: ## Run `MAKE_TARGET` using builder base docker container
update-attribution-checksums-docker: ## Update attribution and checksums using the builder base docker container
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
