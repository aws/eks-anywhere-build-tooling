


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `image-builder`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file
patch-repo: ## Patch upstream repo with patches in patches directory

##@ Fetch Binary Targets
_output/1-26/dependencies/linux-amd64/eksa/kubernetes-sigs/etcdadm: ## Fetch `_output/1-26/dependencies/linux-amd64/eksa/kubernetes-sigs/etcdadm`
_output/1-26/dependencies/linux-arm64/eksa/kubernetes-sigs/etcdadm: ## Fetch `_output/1-26/dependencies/linux-arm64/eksa/kubernetes-sigs/etcdadm`
_output/1-26/dependencies/linux-amd64/eksa/kubernetes-sigs/cri-tools: ## Fetch `_output/1-26/dependencies/linux-amd64/eksa/kubernetes-sigs/cri-tools`
_output/1-26/dependencies/linux-arm64/eksa/kubernetes-sigs/cri-tools: ## Fetch `_output/1-26/dependencies/linux-arm64/eksa/kubernetes-sigs/cri-tools`

##@ Run in Docker Targets
run-in-docker/all-attributions: ## Run `all-attributions` in docker builder container
run-in-docker/all-attributions-checksums: ## Run `all-attributions-checksums` in docker builder container
run-in-docker/all-checksums: ## Run `all-checksums` in docker builder container
run-in-docker/attribution: ## Run `attribution` in docker builder container
run-in-docker/attribution-checksums: ## Run `attribution-checksums` in docker builder container
run-in-docker/binaries: ## Run `binaries` in docker builder container
run-in-docker/checksums: ## Run `checksums` in docker builder container
run-in-docker/clean: ## Run `clean` in docker builder container
run-in-docker/clean-go-cache: ## Run `clean-go-cache` in docker builder container
run-in-docker/validate-checksums: ## Run `validate-checksums` in docker builder container

##@ Clean Targets
clean: ## Removes source and _output directory
clean-go-cache: ## Removes the GOMODCACHE AND GOCACHE folders
clean-repo: ## Removes source directory

##@Fetch Binary Targets
handle-dependencies: ## Download and extract TARs for each dependency listed in PROJECT_DEPENDENCIES

##@ Helpers
help: ## Display this help
add-generated-help-block: ## Add or update generated help block to document project make file and support shell auto completion

##@Update Helpers
start-docker-builder: ## Start long lived builder base docker container
stop-docker-builder: ## Clean up builder base docker container
run-buildkit-and-registry: ## Run buildkitd and a local docker registry as containers
stop-buildkit-and-registry: ## Stop the buildkitd and a local docker registry containers
generate: ## Update UPSTREAM_PROJECTS.yaml
update-go-mods: ## Update locally checked-in go sum to assist in vuln scanning
update-vendor-for-dep-patch: ## After bumping dep in go.mod file, uses generic vendor update script or one provided from upstream project
patch-for-dep-update: ## After bumping dep in go.mod file and updating vendor, generates patch
create-ecr-repos: ## Create repos in ECR for project images for local testing

##@ Build Targets
build: ## Called via prow presubmit, calls ` build-raw-ubuntu-2004 fake-ubuntu-2004-raw.gz upload-artifacts-raw-ubuntu-2004 build-raw-ubuntu-2204 fake-ubuntu-2204-raw.gz upload-artifacts-raw-ubuntu-2204 build-raw-ubuntu-2404 fake-ubuntu-2404-raw.gz upload-artifacts-raw-ubuntu-2404 build-raw-redhat-8 fake-redhat-8-raw.gz upload-artifacts-raw-redhat-8 build-raw-redhat-9 fake-redhat-9-raw.gz upload-artifacts-raw-redhat-9 metal-instance-test build-ami-ubuntu-2004 build-ami-ubuntu-2204 packer/ova/vsphere.json build-ova-ubuntu-2004 fake-ubuntu-2004-ova.ova upload-artifacts-ova-ubuntu-2004 build-ova-ubuntu-2204 fake-ubuntu-2204-ova.ova upload-artifacts-ova-ubuntu-2204 build-ova-redhat-8 fake-redhat-8-ova.ova upload-artifacts-ova-redhat-8 build-ova-ubuntu-2004-efi fake-ubuntu-2004-ova-efi.ova upload-artifacts-ova-ubuntu-2004-efi build-ova-ubuntu-2204-efi fake-ubuntu-2204-ova-efi.ova upload-artifacts-ova-ubuntu-2204-efi build-cloudstack-redhat-8 fake-redhat-8-cloudstack.qcow2 upload-artifacts-cloudstack-redhat-8 build-cloudstack-redhat-9 fake-redhat-9-cloudstack.qcow2 upload-artifacts-cloudstack-redhat-9 packer/nutanix/nutanix.json build-nutanix-ubuntu-2004 fake-ubuntu-2004-nutanix.img upload-artifacts-nutanix-ubuntu-2004 build-nutanix-ubuntu-2204 fake-ubuntu-2204-nutanix.img upload-artifacts-nutanix-ubuntu-2204 build-nutanix-ubuntu-2404 fake-ubuntu-2404-nutanix.img upload-artifacts-nutanix-ubuntu-2404 build-nutanix-redhat-8 fake-redhat-8-nutanix.img upload-artifacts-nutanix-redhat-8 build-nutanix-redhat-9 fake-redhat-9-nutanix.img upload-artifacts-nutanix-redhat-9`
release: ## Called via prow postsubmit + release jobs, calls `validate-supported-image-all release-ova-ubuntu-2004 upload-artifacts-ova-ubuntu-2004`
########### END GENERATED ###########################
