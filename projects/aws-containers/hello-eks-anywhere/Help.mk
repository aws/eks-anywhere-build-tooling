


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `hello-eks-anywhere`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file

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
build: ## Called via prow presubmit, calls `validate-checksums attribution local-images  attribution-pr`
release: ## Called via prow postsubmit + release jobs, calls `validate-checksums images `
########### END GENERATED ###########################
