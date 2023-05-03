


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
_output/dependencies/linux-arm64/eksa/fluxcd/flux2: ## Fetch `_output/dependencies/linux-arm64/eksa/fluxcd/flux2`
_output/dependencies/linux-amd64/eksa/kubernetes-sigs/cluster-api: ## Fetch `_output/dependencies/linux-amd64/eksa/kubernetes-sigs/cluster-api`
_output/dependencies/linux-arm64/eksa/kubernetes-sigs/cluster-api: ## Fetch `_output/dependencies/linux-arm64/eksa/kubernetes-sigs/cluster-api`
_output/dependencies/linux-amd64/eksa/kubernetes-sigs/kind: ## Fetch `_output/dependencies/linux-amd64/eksa/kubernetes-sigs/kind`
_output/dependencies/linux-arm64/eksa/kubernetes-sigs/kind: ## Fetch `_output/dependencies/linux-arm64/eksa/kubernetes-sigs/kind`
_output/dependencies/linux-amd64/eksa/replicatedhq/troubleshoot: ## Fetch `_output/dependencies/linux-amd64/eksa/replicatedhq/troubleshoot`
_output/dependencies/linux-arm64/eksa/replicatedhq/troubleshoot: ## Fetch `_output/dependencies/linux-arm64/eksa/replicatedhq/troubleshoot`
_output/dependencies/linux-amd64/eksa/vmware/govmomi: ## Fetch `_output/dependencies/linux-amd64/eksa/vmware/govmomi`
_output/dependencies/linux-arm64/eksa/vmware/govmomi: ## Fetch `_output/dependencies/linux-arm64/eksa/vmware/govmomi`
_output/dependencies/linux-amd64/eksd/kubernetes/client: ## Fetch `_output/dependencies/linux-amd64/eksd/kubernetes/client`
_output/dependencies/linux-arm64/eksd/kubernetes/client: ## Fetch `_output/dependencies/linux-arm64/eksd/kubernetes/client`
_output/dependencies/linux-amd64/eksa/helm/helm: ## Fetch `_output/dependencies/linux-amd64/eksa/helm/helm`
_output/dependencies/linux-arm64/eksa/helm/helm: ## Fetch `_output/dependencies/linux-arm64/eksa/helm/helm`
_output/dependencies/linux-amd64/eksa/apache/cloudstack-cloudmonkey: ## Fetch `_output/dependencies/linux-amd64/eksa/apache/cloudstack-cloudmonkey`
_output/dependencies/linux-arm64/eksa/apache/cloudstack-cloudmonkey: ## Fetch `_output/dependencies/linux-arm64/eksa/apache/cloudstack-cloudmonkey`

##@ Run in Docker Targets
run-all-attributions-in-docker: ## Run `all-attributions` in docker builder container
run-all-attributions-checksums-in-docker: ## Run `all-attributions-checksums` in docker builder container
run-all-checksums-in-docker: ## Run `all-checksums` in docker builder container
run-attribution-in-docker: ## Run `attribution` in docker builder container
run-attribution-checksums-in-docker: ## Run `attribution-checksums` in docker builder container
run-binaries-in-docker: ## Run `binaries` in docker builder container
run-checksums-in-docker: ## Run `checksums` in docker builder container
run-clean-in-docker: ## Run `clean` in docker builder container
run-clean-go-cache-in-docker: ## Run `clean-go-cache` in docker builder container

##@ Clean Targets
clean: ## Removes source and _output directory

##@ Helpers
help: ## Display this help
add-generated-help-block: ## Add or update generated help block to document project make file and support shell auto completion

##@Update Helpers
run-target-in-docker: ## Run `MAKE_TARGET` using builder base docker container
stop-docker-builder: ## Clean up builder base docker container
generate: ## Update UPSTREAM_PROJECTS.yaml
update-go-mods: ## Update locally checked-in go sum to assist in vuln scanning
update-vendor-for-dep-patch: ## After bumping dep in go.mod file, uses generic vendor update script or one provided from upstream project
patch-for-dep-update: ## After bumping dep in go.mod file and updating vendor, generates patch
create-ecr-repos: ## Create repos in ECR for project images for local testing

##@ Build Targets
build: ## Called via prow presubmit, calls `local-images`
release: ## Called via prow postsubmit + release jobs, calls `images`
########### END GENERATED ###########################
