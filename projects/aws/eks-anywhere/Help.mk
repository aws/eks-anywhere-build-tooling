


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ Image Targets
local-images: ## Builds `eks-anywhere-diagnostic-collector/images/amd64` as oci tars for presumbit validation
images: ## Pushes `eks-anywhere-diagnostic-collector/images/push` to IMAGE_REPO
eks-anywhere-diagnostic-collector/images/amd64: ## Builds/pushes `eks-anywhere-diagnostic-collector/images/amd64`
eks-anywhere-diagnostic-collector/images/push: ## Builds/pushes `eks-anywhere-diagnostic-collector/images/push`

##@ Fetch Binary Targets
_output/dependencies/linux-amd64/eksd/kubernetes/client: ## Fetch `_output/dependencies/linux-amd64/eksd/kubernetes/client`
_output/dependencies/linux-arm64/eksd/kubernetes/client: ## Fetch `_output/dependencies/linux-arm64/eksd/kubernetes/client`

##@ Clean Targets
clean: ## Removes source and _output directory

##@ Helpers
help: ## Display this help
add-generated-help-block: ## Add or update generated help block to document project make file and support shell auto completion

##@ Build Targets
build: ## Called via prow presubmit, calls `local-images`
release: ## Called via prow postsubmit + release jobs, calls `images`
########### END GENERATED ###########################
