


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ Image Targets
local-images: ## Builds `eks-anywhere-cli-tools/images/amd64` as oci tars for presumbit validation
images: ## Pushes `eks-anywhere-cli-tools/images/push` to IMAGE_REPO
eks-anywhere-cli-tools/images/amd64: ## Builds/pushes `eks-anywhere-cli-tools/images/amd64`
eks-anywhere-cli-tools/images/push: ## Builds/pushes `eks-anywhere-cli-tools/images/push`

##@ Fetch Binary Targets
_output/dependencies/linux-amd64/eksa/fluxcd/flux2: ## Fetch `_output/dependencies/linux-amd64/eksa/fluxcd/flux2`
_output/dependencies/linux-amd64/eksa/kubernetes-sigs/cluster-api: ## Fetch `_output/dependencies/linux-amd64/eksa/kubernetes-sigs/cluster-api`
_output/dependencies/linux-amd64/eksa/kubernetes-sigs/cluster-api-provider-aws: ## Fetch `_output/dependencies/linux-amd64/eksa/kubernetes-sigs/cluster-api-provider-aws`
_output/dependencies/linux-amd64/eksa/kubernetes-sigs/kind: ## Fetch `_output/dependencies/linux-amd64/eksa/kubernetes-sigs/kind`
_output/dependencies/linux-amd64/eksa/replicatedhq/troubleshoot: ## Fetch `_output/dependencies/linux-amd64/eksa/replicatedhq/troubleshoot`
_output/dependencies/linux-amd64/eksa/vmware/govmomi: ## Fetch `_output/dependencies/linux-amd64/eksa/vmware/govmomi`
_output/dependencies/linux-amd64/eksd/kubernetes/client: ## Fetch `_output/dependencies/linux-amd64/eksd/kubernetes/client`
_output/dependencies/linux-amd64/eksa/helm/helm: ## Fetch `_output/dependencies/linux-amd64/eksa/helm/helm`
_output/dependencies/linux-arm64/eksa/fluxcd/flux2: ## Fetch `_output/dependencies/linux-arm64/eksa/fluxcd/flux2`
_output/dependencies/linux-arm64/eksa/kubernetes-sigs/cluster-api: ## Fetch `_output/dependencies/linux-arm64/eksa/kubernetes-sigs/cluster-api`
_output/dependencies/linux-arm64/eksa/kubernetes-sigs/cluster-api-provider-aws: ## Fetch `_output/dependencies/linux-arm64/eksa/kubernetes-sigs/cluster-api-provider-aws`
_output/dependencies/linux-arm64/eksa/kubernetes-sigs/kind: ## Fetch `_output/dependencies/linux-arm64/eksa/kubernetes-sigs/kind`
_output/dependencies/linux-arm64/eksa/replicatedhq/troubleshoot: ## Fetch `_output/dependencies/linux-arm64/eksa/replicatedhq/troubleshoot`
_output/dependencies/linux-arm64/eksa/vmware/govmomi: ## Fetch `_output/dependencies/linux-arm64/eksa/vmware/govmomi`
_output/dependencies/linux-arm64/eksd/kubernetes/client: ## Fetch `_output/dependencies/linux-arm64/eksd/kubernetes/client`
_output/dependencies/linux-arm64/eksa/helm/helm: ## Fetch `_output/dependencies/linux-arm64/eksa/helm/helm`

##@ Clean Targets
clean: ## Removes source and _output directory

##@ Helpers
help: ## Display this help
add-generated-help-block: ## Add or update generated help block to document project make file and support shell auto completion

##@ Build Targets
build: ## Called via prow presubmit, calls `local-images`
release: ## Called via prow postsubmit + release jobs, calls `images`
########### END GENERATED ###########################
