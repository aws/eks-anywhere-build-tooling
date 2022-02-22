


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `kube-rbac-proxy`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file

##@ Binary Targets
binaries: ## Build all binaries: `kube-rbac-proxy` for `linux/amd64 linux/arm64`
_output/bin/kube-rbac-proxy/linux-amd64/kube-rbac-proxy: ## Build `_output/bin/kube-rbac-proxy/linux-amd64/kube-rbac-proxy`
_output/bin/kube-rbac-proxy/linux-arm64/kube-rbac-proxy: ## Build `_output/bin/kube-rbac-proxy/linux-arm64/kube-rbac-proxy`

##@ Image Targets
local-images: ## Builds `kube-rbac-proxy/images/amd64` as oci tars for presumbit validation
images: ## Pushes `kube-rbac-proxy/images/push` to IMAGE_REPO
kube-rbac-proxy/images/amd64: ## Builds/pushes `kube-rbac-proxy/images/amd64`
kube-rbac-proxy/images/push: ## Builds/pushes `kube-rbac-proxy/images/push`

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
build: ## Called via prow presubmit, calls `checksums local-images attribution  attribution-pr`
release: ## Called via prow postsubmit + release jobs, calls `checksums images `
########### END GENERATED ###########################
