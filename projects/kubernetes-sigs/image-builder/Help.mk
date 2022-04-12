


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

##@ Build Targets
build: ## Called via prow presubmit, calls `release-raw-ubuntu-2004-efi image-builder/images/capi/output/fake-ubuntu.gz /home/prow/go/src/github.com/aws/eks-anywhere-build-tooling/projects/kubernetes-sigs/image-builder/_output/tar/1-21/raw/ubuntu/ubuntu.gz upload-artifacts-raw build-ami-ubuntu-2004 setup-packer-configs-ova download-ova-bottlerocket image-builder/images/capi/output/fake-ubuntu.ova /home/prow/go/src/github.com/aws/eks-anywhere-build-tooling/projects/kubernetes-sigs/image-builder/_output/tar/1-21/ova/ubuntu/ubuntu.ova /home/prow/go/src/github.com/aws/eks-anywhere-build-tooling/projects/kubernetes-sigs/image-builder/_output/tar/1-21/ova/bottlerocket/bottlerocket.ova upload-artifacts-ova upload-ova-bottlerocket`
release: ## Called via prow postsubmit + release jobs, calls `validate-checksums images `
########### END GENERATED ###########################
