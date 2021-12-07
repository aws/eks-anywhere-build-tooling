


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ Binary Targets
binaries: ## Build all binaries: `bottlerocket-bootstrap` for `linux/amd64 linux/arm64`
_output/bin/bottlerocket-bootstrap/linux-amd64/bottlerocket-bootstrap: ## Build `_output/bin/bottlerocket-bootstrap/linux-amd64/bottlerocket-bootstrap`
_output/bin/bottlerocket-bootstrap/linux-arm64/bottlerocket-bootstrap: ## Build `_output/bin/bottlerocket-bootstrap/linux-arm64/bottlerocket-bootstrap`

##@ Image Targets
local-images: ## Builds `bottlerocket-bootstrap/images/amd64` as oci tars for presumbit validation
images: ## Pushes `bottlerocket-bootstrap/images/push` to IMAGE_REPO
bottlerocket-bootstrap/images/amd64: ## Builds/pushes `bottlerocket-bootstrap/images/amd64`
bottlerocket-bootstrap/images/push: ## Builds/pushes `bottlerocket-bootstrap/images/push`

##@ Checksum Targets
checksums: ## Update checksums file based on currently built binaries.
validate-checksums: # Validate checksums of currently built binaries against checksums file.

##@ License Targets
gather-licenses: ## Helper to call $(GATHER_LICENSES_TARGETS) which gathers all licenses
attribution: ## Generates attribution from licenses gathered during `gather-licenses`.
attribution-pr: ## Generates PR to update attribution files for projects

##@ Clean Targets
clean: ## Removes source and _output directory

##@ Helpers
help: ## Display this help
add-generated-help-block: ## Add or update generated help block to document project make file and support shell auto completion

##@ Build Targets
build: ## Called via prow presubmit, calls `validate-checksums local-images attribution attribution-pr `
release: ## Called via prow postsubmit + release jobs, calls `validate-checksums images `
########### END GENERATED ###########################
