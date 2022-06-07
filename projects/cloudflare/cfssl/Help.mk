


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `cfssl`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file

##@ Binary Targets
binaries: ## Build all binaries: `cfssl cfssl-bundle cfssl-certinfo cfssl-newkey cfssl-scan cfssljson mkbundle multirootca` for `linux/amd64 linux/arm64`
_output/bin/cfssl/linux-amd64/cfssl: ## Build `_output/bin/cfssl/linux-amd64/cfssl`
_output/bin/cfssl/linux-amd64/cfssl-bundle: ## Build `_output/bin/cfssl/linux-amd64/cfssl-bundle`
_output/bin/cfssl/linux-amd64/cfssl-certinfo: ## Build `_output/bin/cfssl/linux-amd64/cfssl-certinfo`
_output/bin/cfssl/linux-amd64/cfssl-newkey: ## Build `_output/bin/cfssl/linux-amd64/cfssl-newkey`
_output/bin/cfssl/linux-amd64/cfssl-scan: ## Build `_output/bin/cfssl/linux-amd64/cfssl-scan`
_output/bin/cfssl/linux-amd64/cfssljson: ## Build `_output/bin/cfssl/linux-amd64/cfssljson`
_output/bin/cfssl/linux-amd64/mkbundle: ## Build `_output/bin/cfssl/linux-amd64/mkbundle`
_output/bin/cfssl/linux-amd64/multirootca: ## Build `_output/bin/cfssl/linux-amd64/multirootca`
_output/bin/cfssl/linux-arm64/cfssl: ## Build `_output/bin/cfssl/linux-arm64/cfssl`
_output/bin/cfssl/linux-arm64/cfssl-bundle: ## Build `_output/bin/cfssl/linux-arm64/cfssl-bundle`
_output/bin/cfssl/linux-arm64/cfssl-certinfo: ## Build `_output/bin/cfssl/linux-arm64/cfssl-certinfo`
_output/bin/cfssl/linux-arm64/cfssl-newkey: ## Build `_output/bin/cfssl/linux-arm64/cfssl-newkey`
_output/bin/cfssl/linux-arm64/cfssl-scan: ## Build `_output/bin/cfssl/linux-arm64/cfssl-scan`
_output/bin/cfssl/linux-arm64/cfssljson: ## Build `_output/bin/cfssl/linux-arm64/cfssljson`
_output/bin/cfssl/linux-arm64/mkbundle: ## Build `_output/bin/cfssl/linux-arm64/mkbundle`
_output/bin/cfssl/linux-arm64/multirootca: ## Build `_output/bin/cfssl/linux-arm64/multirootca`

##@ Image Targets
local-images: ## Builds `cfssl/images/amd64` as oci tars for presumbit validation
images: ## Pushes `cfssl/images/push` to IMAGE_REPO
cfssl/images/amd64: ## Builds/pushes `cfssl/images/amd64`
cfssl/images/push: ## Builds/pushes `cfssl/images/push`

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
create-ecr-repos: ## Create repos in ECR for project images for local testing

##@ Build Targets
build: ## Called via prow presubmit, calls `handle-dependencies validate-checksums attribution local-images upload-artifacts attribution-pr" `
release: ## Called via prow postsubmit + release jobs, calls `handle-dependencies validate-checksums images upload-artifacts`
########### END GENERATED ###########################
