


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `eks-anywhere-packages`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file

##@ Binary Targets
binaries: ## Build all binaries: `package-manager` for `linux/amd64 linux/arm64`
_output/bin/eks-anywhere-packages/linux-amd64/package-manager: ## Build `_output/bin/eks-anywhere-packages/linux-amd64/package-manager`
_output/bin/eks-anywhere-packages/linux-arm64/package-manager: ## Build `_output/bin/eks-anywhere-packages/linux-arm64/package-manager`

##@ Image Targets
local-images: ## Builds `eks-anywhere-packages/images/amd64 helm/build` as oci tars for presumbit validation
images: ## Pushes `eks-anywhere-packages/images/push helm/push` to IMAGE_REPO
eks-anywhere-packages/images/amd64: ## Builds/pushes `eks-anywhere-packages/images/amd64`
helm/build: ## Builds/pushes `helm/build`
eks-anywhere-packages/images/push: ## Builds/pushes `eks-anywhere-packages/images/push`
helm/push: ## Builds/pushes `helm/push`

##@ Helm Targets
helm/build: ## Build helm chart
helm/push: ## Build helm chart and push to registry defined in IMAGE_REPO.

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

##@Update Helpers
run-target-in-docker: ## Run `MAKE_TARGET` using builder base docker container
update-attribution-checksums-docker: ## Update attribution and checksums using the builder base docker container
stop-docker-builder: ## Clean up builder base docker container
generate: ## Update UPSTREAM_PROJECTS.yaml
create-ecr-repos: ## Create repos in ECR for project images for local testing

##@ Build Targets
build: ## Called via prow presubmit, calls `handle-dependencies validate-checksums attribution local-images  attribution-pr`
release: ## Called via prow postsubmit + release jobs, calls `handle-dependencies validate-checksums images `
########### END GENERATED ###########################
