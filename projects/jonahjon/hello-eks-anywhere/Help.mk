


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ Image Targets
local-images: ## Builds `hello-eks-anywhere/images/amd64 helm/build` as oci tars for presumbit validation
images: ## Pushes `hello-eks-anywhere/images/push helm/push` to IMAGE_REPO
hello-eks-anywhere/images/amd64: ## Builds/pushes `hello-eks-anywhere/images/amd64`
helm/build: ## Builds/pushes `helm/build`
hello-eks-anywhere/images/push: ## Builds/pushes `hello-eks-anywhere/images/push`
helm/push: ## Builds/pushes `helm/push`

##@ Helm Targets
helm/build: ## Build helm chart
helm/push: ## Build helm chart and push to registry defined in IMAGE_REPO.

##@ Clean Targets
clean: ## Removes source and _output directory

##@ Helpers
help: ## Display this help
add-generated-help-block: ## Add or update generated help block to document project make file and support shell auto completion

##@ Build Targets
build: ## Called via prow presubmit, calls `local-images`
release: ## Called via prow postsubmit + release jobs, calls `images`
########### END GENERATED ###########################
