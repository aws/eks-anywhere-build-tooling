


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `tink`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file
patch-repo: ## Patch upstream repo with patches in patches directory

##@ Binary Targets
binaries: ## Build all binaries: `tink-cli tink-server tink-worker` for `linux/amd64 linux/arm64`
_output/bin/tink/linux-amd64/tink-cli: ## Build `_output/bin/tink/linux-amd64/tink-cli`
_output/bin/tink/linux-amd64/tink-server: ## Build `_output/bin/tink/linux-amd64/tink-server`
_output/bin/tink/linux-amd64/tink-worker: ## Build `_output/bin/tink/linux-amd64/tink-worker`
_output/bin/tink/linux-arm64/tink-cli: ## Build `_output/bin/tink/linux-arm64/tink-cli`
_output/bin/tink/linux-arm64/tink-server: ## Build `_output/bin/tink/linux-arm64/tink-server`
_output/bin/tink/linux-arm64/tink-worker: ## Build `_output/bin/tink/linux-arm64/tink-worker`

##@ Image Targets
local-images: ## Builds `tink-cli/images/amd64 tink-server/images/amd64 tink-worker/images/amd64` as oci tars for presumbit validation
images: ## Pushes `tink-cli/images/push tink-server/images/push tink-worker/images/push` to IMAGE_REPO
tink-cli/images/amd64: ## Builds/pushes `tink-cli/images/amd64`
tink-server/images/amd64: ## Builds/pushes `tink-server/images/amd64`
tink-worker/images/amd64: ## Builds/pushes `tink-worker/images/amd64`
tink-cli/images/push: ## Builds/pushes `tink-cli/images/push`
tink-server/images/push: ## Builds/pushes `tink-server/images/push`
tink-worker/images/push: ## Builds/pushes `tink-worker/images/push`

##@ Fetch Binary Targets
_output/dependencies/linux-amd64/eksa/cloudflare/cfssl: ## Fetch `_output/dependencies/linux-amd64/eksa/cloudflare/cfssl`
_output/dependencies/linux-arm64/eksa/cloudflare/cfssl: ## Fetch `_output/dependencies/linux-arm64/eksa/cloudflare/cfssl`

##@ Checksum Targets
checksums: ## Update checksums file based on currently built binaries.
validate-checksums: # Validate checksums of currently built binaries against checksums file.

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
build: ## Called via prow presubmit, calls `validate-checksums local-images attribution  attribution-pr`
release: ## Called via prow postsubmit + release jobs, calls `validate-checksums images `
########### END GENERATED ###########################
