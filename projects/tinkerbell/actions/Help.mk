


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `actions`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file
patch-repo: ## Patch upstream repo with patches in patches directory

##@ Binary Targets
binaries: ## Build all binaries: `cexec kexec image2disk oci2disk writefile` for `linux/amd64 linux/arm64`
_output/bin/actions/linux-amd64/cexec: ## Build `_output/bin/actions/linux-amd64/cexec`
_output/bin/actions/linux-amd64/kexec: ## Build `_output/bin/actions/linux-amd64/kexec`
_output/bin/actions/linux-amd64/image2disk: ## Build `_output/bin/actions/linux-amd64/image2disk`
_output/bin/actions/linux-amd64/oci2disk: ## Build `_output/bin/actions/linux-amd64/oci2disk`
_output/bin/actions/linux-amd64/writefile: ## Build `_output/bin/actions/linux-amd64/writefile`
_output/bin/actions/linux-arm64/cexec: ## Build `_output/bin/actions/linux-arm64/cexec`
_output/bin/actions/linux-arm64/kexec: ## Build `_output/bin/actions/linux-arm64/kexec`
_output/bin/actions/linux-arm64/image2disk: ## Build `_output/bin/actions/linux-arm64/image2disk`
_output/bin/actions/linux-arm64/oci2disk: ## Build `_output/bin/actions/linux-arm64/oci2disk`
_output/bin/actions/linux-arm64/writefile: ## Build `_output/bin/actions/linux-arm64/writefile`

##@ Image Targets
local-images: ## Builds `cexec/images/amd64 kexec/images/amd64 image2disk/images/amd64 oci2disk/images/amd64 writefile/images/amd64 reboot/images/amd64` as oci tars for presumbit validation
images: ## Pushes `cexec/images/push kexec/images/push image2disk/images/push oci2disk/images/push writefile/images/push reboot/images/push` to IMAGE_REPO
cexec/images/amd64: ## Builds/pushes `cexec/images/amd64`
kexec/images/amd64: ## Builds/pushes `kexec/images/amd64`
image2disk/images/amd64: ## Builds/pushes `image2disk/images/amd64`
oci2disk/images/amd64: ## Builds/pushes `oci2disk/images/amd64`
writefile/images/amd64: ## Builds/pushes `writefile/images/amd64`
reboot/images/amd64: ## Builds/pushes `reboot/images/amd64`
cexec/images/push: ## Builds/pushes `cexec/images/push`
kexec/images/push: ## Builds/pushes `kexec/images/push`
image2disk/images/push: ## Builds/pushes `image2disk/images/push`
oci2disk/images/push: ## Builds/pushes `oci2disk/images/push`
writefile/images/push: ## Builds/pushes `writefile/images/push`
reboot/images/push: ## Builds/pushes `reboot/images/push`

##@ Fetch Binary Targets
_output/dependencies/linux-amd64/eksa/torvalds/linux: ## Fetch `_output/dependencies/linux-amd64/eksa/torvalds/linux`
_output/dependencies/linux-arm64/eksa/torvalds/linux: ## Fetch `_output/dependencies/linux-arm64/eksa/torvalds/linux`

##@ Checksum Targets
checksums: ## Update checksums file based on currently built binaries.
validate-checksums: # Validate checksums of currently built binaries against checksums file.
all-checksums: ## Update checksums files for all RELEASE_BRANCHes.

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
run-in-docker/actions/eks-anywhere-go-mod-download: ## Run `actions/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/_output/bin/actions/linux-amd64/cexec: ## Run `_output/bin/actions/linux-amd64/cexec` in docker builder container
run-in-docker/_output/bin/actions/linux-amd64/kexec: ## Run `_output/bin/actions/linux-amd64/kexec` in docker builder container
run-in-docker/_output/bin/actions/linux-amd64/image2disk: ## Run `_output/bin/actions/linux-amd64/image2disk` in docker builder container
run-in-docker/_output/bin/actions/linux-amd64/oci2disk: ## Run `_output/bin/actions/linux-amd64/oci2disk` in docker builder container
run-in-docker/_output/bin/actions/linux-amd64/writefile: ## Run `_output/bin/actions/linux-amd64/writefile` in docker builder container
run-in-docker/_output/bin/actions/linux-arm64/cexec: ## Run `_output/bin/actions/linux-arm64/cexec` in docker builder container
run-in-docker/_output/bin/actions/linux-arm64/kexec: ## Run `_output/bin/actions/linux-arm64/kexec` in docker builder container
run-in-docker/_output/bin/actions/linux-arm64/image2disk: ## Run `_output/bin/actions/linux-arm64/image2disk` in docker builder container
run-in-docker/_output/bin/actions/linux-arm64/oci2disk: ## Run `_output/bin/actions/linux-arm64/oci2disk` in docker builder container
run-in-docker/_output/bin/actions/linux-arm64/writefile: ## Run `_output/bin/actions/linux-arm64/writefile` in docker builder container
run-in-docker/_output/attribution/go-license.csv: ## Run `_output/attribution/go-license.csv` in docker builder container
run-in-docker/_output/attribution/go-license.csv: ## Run `_output/attribution/go-license.csv` in docker builder container
run-in-docker/_output/attribution/go-license.csv: ## Run `_output/attribution/go-license.csv` in docker builder container
run-in-docker/_output/attribution/go-license.csv: ## Run `_output/attribution/go-license.csv` in docker builder container
run-in-docker/_output/attribution/go-license.csv: ## Run `_output/attribution/go-license.csv` in docker builder container

##@ License Targets
gather-licenses: ## Helper to call $(GATHER_LICENSES_TARGETS) which gathers all licenses
attribution: ## Generates attribution from licenses gathered during `gather-licenses`.
attribution-pr: ## Generates PR to update attribution files for projects
attribution-checksums: ## Update attribution and checksums files.
all-attributions: ## Update attribution files for all RELEASE_BRANCHes.
all-attributions-checksums: ## Update attribution and checksums files for all RELEASE_BRANCHes.

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
build: ## Called via prow presubmit, calls `github-rate-limit-pre validate-checksums attribution local-images   attribution-pr github-rate-limit-post`
release: ## Called via prow postsubmit + release jobs, calls `validate-checksums images  `
########### END GENERATED ###########################
