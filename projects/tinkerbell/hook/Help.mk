


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `hook`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file
patch-repo: ## Patch upstream repo with patches in patches directory

##@ Binary Targets
binaries: ## Build all binaries: `hook-bootkit hook-docker` for `linux/amd64 linux/arm64`
_output/bin/hook/linux-amd64/hook-bootkit: ## Build `_output/bin/hook/linux-amd64/hook-bootkit`
_output/bin/hook/linux-amd64/hook-docker: ## Build `_output/bin/hook/linux-amd64/hook-docker`
_output/bin/hook/linux-arm64/hook-bootkit: ## Build `_output/bin/hook/linux-arm64/hook-bootkit`
_output/bin/hook/linux-arm64/hook-docker: ## Build `_output/bin/hook/linux-arm64/hook-docker`

##@ Image Targets
local-images: ## Builds `hook-bootkit/images/amd64 hook-docker/images/amd64 hook-runc/images/amd64 hook-containerd/images/amd64 kernel/images/amd64 hook-dind/images/amd64 hook-embedded/images/amd64 hook-udev/images/amd64 hook-acpid/images/amd64` as oci tars for presumbit validation
images: ## Pushes `hook-bootkit/images/push hook-docker/images/push hook-runc/images/push hook-containerd/images/push kernel/images/push hook-dind/images/push hook-embedded/images/push hook-udev/images/push hook-acpid/images/push` to IMAGE_REPO
hook-bootkit/images/amd64: ## Builds/pushes `hook-bootkit/images/amd64`
hook-docker/images/amd64: ## Builds/pushes `hook-docker/images/amd64`
hook-runc/images/amd64: ## Builds/pushes `hook-runc/images/amd64`
hook-containerd/images/amd64: ## Builds/pushes `hook-containerd/images/amd64`
kernel/images/amd64: ## Builds/pushes `kernel/images/amd64`
hook-dind/images/amd64: ## Builds/pushes `hook-dind/images/amd64`
hook-embedded/images/amd64: ## Builds/pushes `hook-embedded/images/amd64`
hook-udev/images/amd64: ## Builds/pushes `hook-udev/images/amd64`
hook-acpid/images/amd64: ## Builds/pushes `hook-acpid/images/amd64`
hook-bootkit/images/push: ## Builds/pushes `hook-bootkit/images/push`
hook-docker/images/push: ## Builds/pushes `hook-docker/images/push`
hook-runc/images/push: ## Builds/pushes `hook-runc/images/push`
hook-containerd/images/push: ## Builds/pushes `hook-containerd/images/push`
kernel/images/push: ## Builds/pushes `kernel/images/push`
hook-dind/images/push: ## Builds/pushes `hook-dind/images/push`
hook-embedded/images/push: ## Builds/pushes `hook-embedded/images/push`
hook-udev/images/push: ## Builds/pushes `hook-udev/images/push`
hook-acpid/images/push: ## Builds/pushes `hook-acpid/images/push`

##@ Fetch Binary Targets
_output/dependencies/linux-amd64/eksa/containerd/containerd: ## Fetch `_output/dependencies/linux-amd64/eksa/containerd/containerd`
_output/dependencies/linux-arm64/eksa/containerd/containerd: ## Fetch `_output/dependencies/linux-arm64/eksa/containerd/containerd`
_output/dependencies/linux-amd64/eksa/tinkerbell/tinkerbell: ## Fetch `_output/dependencies/linux-amd64/eksa/tinkerbell/tinkerbell`
_output/dependencies/linux-arm64/eksa/tinkerbell/tinkerbell: ## Fetch `_output/dependencies/linux-arm64/eksa/tinkerbell/tinkerbell`
_output/dependencies/linux-amd64/eksa/linuxkit/linuxkit: ## Fetch `_output/dependencies/linux-amd64/eksa/linuxkit/linuxkit`
_output/dependencies/linux-arm64/eksa/linuxkit/linuxkit: ## Fetch `_output/dependencies/linux-arm64/eksa/linuxkit/linuxkit`

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
run-in-docker/hook/images/hook-bootkit/eks-anywhere-go-mod-download: ## Run `hook/images/hook-bootkit/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/hook/images/hook-docker/eks-anywhere-go-mod-download: ## Run `hook/images/hook-docker/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/_output/bin/hook/linux-amd64/hook-bootkit: ## Run `_output/bin/hook/linux-amd64/hook-bootkit` in docker builder container
run-in-docker/_output/bin/hook/linux-amd64/hook-docker: ## Run `_output/bin/hook/linux-amd64/hook-docker` in docker builder container
run-in-docker/_output/bin/hook/linux-arm64/hook-bootkit: ## Run `_output/bin/hook/linux-arm64/hook-bootkit` in docker builder container
run-in-docker/_output/bin/hook/linux-arm64/hook-docker: ## Run `_output/bin/hook/linux-arm64/hook-docker` in docker builder container
run-in-docker/_output/hook-bootkit/attribution/go-license.csv: ## Run `_output/hook-bootkit/attribution/go-license.csv` in docker builder container
run-in-docker/_output/hook-docker/attribution/go-license.csv: ## Run `_output/hook-docker/attribution/go-license.csv` in docker builder container

##@ Artifact Targets
tarballs: ## Create tarballs by calling build/lib/simple_create_tarballs.sh unless SIMPLE_CREATE_TARBALLS=false, then tarballs must be defined in project Makefile
s3-artifacts: # Prepare ARTIFACTS_PATH folder structure with tarballs/manifests/other items to be uploaded to s3
upload-artifacts: # Upload tarballs and other artifacts from ARTIFACTS_PATH to S3

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
build: ## Called via prow presubmit, calls `github-rate-limit-pre validate-checksums attribution local-images  upload-artifacts attribution-pr github-rate-limit-post`
release: ## Called via prow postsubmit + release jobs, calls `validate-checksums images  upload-artifacts`
########### END GENERATED ###########################
