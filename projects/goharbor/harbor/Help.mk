


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `harbor`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file

##@ Image Targets
local-images: ## Builds `harbor-chartserver/images/amd64` as oci tars for presumbit validation
images: ## Pushes `harbor-chartserver/images/push` to IMAGE_REPO
harbor-chartserver/images/amd64: ## Builds/pushes `harbor-chartserver/images/amd64`
harbor-chartserver/images/push: ## Builds/pushes `harbor-chartserver/images/push`

##@ Fetch Binary Targets
_output/dependencies/linux-amd64/eksa/helm/chartmuseum: ## Fetch `_output/dependencies/linux-amd64/eksa/helm/chartmuseum`
_output/dependencies/linux-arm64/eksa/helm/chartmuseum: ## Fetch `_output/dependencies/linux-arm64/eksa/helm/chartmuseum`

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
