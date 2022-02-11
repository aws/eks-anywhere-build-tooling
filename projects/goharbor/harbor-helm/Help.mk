


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `harbor-helm`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file

##@ Image Targets
local-images: ## Builds `harbor-helm/images/amd64 helm/build` as oci tars for presumbit validation
images: ## Pushes `harbor-helm/images/push helm/push` to IMAGE_REPO
harbor-helm/images/amd64: ## Builds/pushes `harbor-helm/images/amd64`
helm/build: ## Builds/pushes `helm/build`
harbor-helm/images/push: ## Builds/pushes `harbor-helm/images/push`
helm/push: ## Builds/pushes `helm/push`

##@ Helm Targets
helm/build: ## Build helm chart
helm/push: ## Build helm chart and push to registry defined in IMAGE_REPO.

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
