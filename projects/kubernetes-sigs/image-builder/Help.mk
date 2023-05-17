


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `image-builder`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file
patch-repo: ## Patch upstream repo with patches in patches directory

##@ Image Targets
local-images: ## Builds `image-builder/images/amd64` as oci tars for presumbit validation
images: ## Pushes `image-builder/images/push` to IMAGE_REPO
image-builder/images/amd64: ## Builds/pushes `image-builder/images/amd64`
image-builder/images/push: ## Builds/pushes `image-builder/images/push`

##@ Fetch Binary Targets
_output/1-26/dependencies/linux-amd64/eksa/kubernetes-sigs/etcdadm: ## Fetch `_output/1-26/dependencies/linux-amd64/eksa/kubernetes-sigs/etcdadm`
_output/1-26/dependencies/linux-arm64/eksa/kubernetes-sigs/etcdadm: ## Fetch `_output/1-26/dependencies/linux-arm64/eksa/kubernetes-sigs/etcdadm`
_output/1-26/dependencies/linux-amd64/eksa/kubernetes-sigs/cri-tools: ## Fetch `_output/1-26/dependencies/linux-amd64/eksa/kubernetes-sigs/cri-tools`
_output/1-26/dependencies/linux-arm64/eksa/kubernetes-sigs/cri-tools: ## Fetch `_output/1-26/dependencies/linux-arm64/eksa/kubernetes-sigs/cri-tools`

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
clean-repo: ## Removes source directory

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
build: ## Called via prow presubmit, calls `setup-packer-configs-raw release-image-build-on-metal-ubuntu fake-ubuntu.gz _output/tar/1-26/raw/ubuntu/ubuntu.gz upload-artifacts-raw download-bottlerocket-raw _output/tar/1-26/raw/bottlerocket/bottlerocket.img.gz upload-bottlerocket-raw build-ami-ubuntu-2004 download-bottlerocket-ami _output/tar/1-26/ami/bottlerocket/bottlerocket.img.gz upload-bottlerocket-ami setup-packer-configs-ova  fake-ubuntu.ova _output/tar/1-26/ova/ubuntu/ubuntu.ova  upload-artifacts-ova download-bottlerocket-ova _output/tar/1-26/ova/bottlerocket/bottlerocket.ova upload-bottlerocket-ova`
release: ## Called via prow postsubmit + release jobs, calls `release-image-build-on-metal-ubuntu upload-artifacts-raw`
########### END GENERATED ###########################
