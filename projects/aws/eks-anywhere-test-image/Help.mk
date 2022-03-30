


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ Image Targets
local-images: ## Builds `eks-anywhere-test/images/amd64 helm/build` as oci tars for presumbit validation
images: ## Pushes `eks-anywhere-test/images/push helm/push` to IMAGE_REPO
eks-anywhere-test/images/amd64: ## Builds/pushes `eks-anywhere-test/images/amd64`
helm/build: ## Builds/pushes `helm/build`
eks-anywhere-test/images/push: ## Builds/pushes `eks-anywhere-test/images/push`
helm/push: ## Builds/pushes `helm/push`

##@ Clean Targets
clean: ## Removes source and _output directory

##@ Helpers
help: ## Display this help
add-generated-help-block: ## Add or update generated help block to document project make file and support shell auto completion

##@ Build Targets
build: ## Called via prow presubmit, calls `local-images`
release: ## Called via prow postsubmit + release jobs, calls `images`
########### END GENERATED ###########################
