


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `vsphere-csi-driver`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file

##@ Binary Targets
binaries: ## Build all binaries: `vsphere-csi-driver vsphere-csi-syncer` for `linux/amd64 linux/arm64`
_output/bin/vsphere-csi-driver/linux-amd64/vsphere-csi-driver: ## Build `_output/bin/vsphere-csi-driver/linux-amd64/vsphere-csi-driver`
_output/bin/vsphere-csi-driver/linux-amd64/vsphere-csi-syncer: ## Build `_output/bin/vsphere-csi-driver/linux-amd64/vsphere-csi-syncer`
_output/bin/vsphere-csi-driver/linux-arm64/vsphere-csi-driver: ## Build `_output/bin/vsphere-csi-driver/linux-arm64/vsphere-csi-driver`
_output/bin/vsphere-csi-driver/linux-arm64/vsphere-csi-syncer: ## Build `_output/bin/vsphere-csi-driver/linux-arm64/vsphere-csi-syncer`

##@ Image Targets
local-images: ## Builds `csi-driver/images/amd64 csi-syncer/images/amd64` as oci tars for presumbit validation
images: ## Pushes `csi-driver/images/push csi-syncer/images/push` to IMAGE_REPO
csi-driver/images/amd64: ## Builds/pushes `csi-driver/images/amd64`
csi-syncer/images/amd64: ## Builds/pushes `csi-syncer/images/amd64`
csi-driver/images/push: ## Builds/pushes `csi-driver/images/push`
csi-syncer/images/push: ## Builds/pushes `csi-syncer/images/push`

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
build: ## Called via prow presubmit, calls `validate-checksums local-images attribution attribution-pr `
release: ## Called via prow postsubmit + release jobs, calls `validate-checksums images `
########### END GENERATED ###########################
