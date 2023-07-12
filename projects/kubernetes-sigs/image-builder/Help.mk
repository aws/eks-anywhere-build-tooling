


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
build: ## Called via prow presubmit, calls `build-raw-ubuntu-2004 fake-ubuntu-2004-raw.gz upload-artifacts-raw-ubuntu-2004 build-raw-ubuntu-2204 fake-ubuntu-2204-raw.gz upload-artifacts-raw-ubuntu-2204 build-raw-redhat-8 fake-redhat-8-raw.gz upload-artifacts-raw-redhat-8 release-raw-ubuntu-2004 upload-bottlerocket-1-raw build-ami-ubuntu-2004 build-ami-ubuntu-2204 upload-bottlerocket-1-ami packer/ova/vsphere.json build-ova-ubuntu-2004 fake-ubuntu-2004-ova.ova upload-artifacts-ova-ubuntu-2004 build-ova-ubuntu-2204 fake-ubuntu-2204-ova.ova upload-artifacts-ova-ubuntu-2204 build-ova-redhat-8 fake-redhat-8-ova.ova upload-artifacts-ova-redhat-8 build-ova-ubuntu-2004-efi fake-ubuntu-2004-ova-efi.ova upload-artifacts-ova-ubuntu-2004-efi build-ova-ubuntu-2204-efi fake-ubuntu-2204-ova-efi.ova upload-artifacts-ova-ubuntu-2204-efi upload-bottlerocket-1-ova build-cloudstack-redhat-8 fake-redhat-8-cloudstack.qcow2 upload-artifacts-cloudstack-redhat-8 build-nutanix-ubuntu-2004 build-nutanix-ubuntu-2204`
release: ## Called via prow postsubmit + release jobs, calls `validate-supported-image-all release-ova-ubuntu-2004 upload-artifacts-ova-ubuntu-2004`
########### END GENERATED ###########################
