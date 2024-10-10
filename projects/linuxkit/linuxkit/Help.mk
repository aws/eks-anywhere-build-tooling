


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Repo Targets
clone-repo:  ## Clone upstream `linuxkit`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file
patch-repo: ## Patch upstream repo with patches in patches directory

##@ Binary Targets
binaries: ## Build all binaries: `linuxkit rngd sysfs sysctl init rc.init service` for `linux/amd64`
_output/bin/linuxkit/linux-amd64/linuxkit: ## Build `_output/bin/linuxkit/linux-amd64/linuxkit`
_output/bin/linuxkit/linux-amd64/rngd: ## Build `_output/bin/linuxkit/linux-amd64/rngd`
_output/bin/linuxkit/linux-amd64/sysfs: ## Build `_output/bin/linuxkit/linux-amd64/sysfs`
_output/bin/linuxkit/linux-amd64/sysctl: ## Build `_output/bin/linuxkit/linux-amd64/sysctl`
_output/bin/linuxkit/linux-amd64/init: ## Build `_output/bin/linuxkit/linux-amd64/init`
_output/bin/linuxkit/linux-amd64/rc.init: ## Build `_output/bin/linuxkit/linux-amd64/rc.init`
_output/bin/linuxkit/linux-amd64/service: ## Build `_output/bin/linuxkit/linux-amd64/service`

##@ Image Targets
local-images: ## Builds `init/images/amd64 ca-certificates/images/amd64 firmware/images/amd64 rngd/images/amd64 sysctl/images/amd64 sysfs/images/amd64 modprobe/images/amd64 dhcpcd/images/amd64 openntpd/images/amd64 getty/images/amd64` as oci tars for presumbit validation
images: ## Pushes `init/images/push ca-certificates/images/push firmware/images/push rngd/images/push sysctl/images/push sysfs/images/push modprobe/images/push dhcpcd/images/push openntpd/images/push getty/images/push` to IMAGE_REPO
init/images/amd64: ## Builds/pushes `init/images/amd64`
ca-certificates/images/amd64: ## Builds/pushes `ca-certificates/images/amd64`
firmware/images/amd64: ## Builds/pushes `firmware/images/amd64`
rngd/images/amd64: ## Builds/pushes `rngd/images/amd64`
sysctl/images/amd64: ## Builds/pushes `sysctl/images/amd64`
sysfs/images/amd64: ## Builds/pushes `sysfs/images/amd64`
modprobe/images/amd64: ## Builds/pushes `modprobe/images/amd64`
dhcpcd/images/amd64: ## Builds/pushes `dhcpcd/images/amd64`
openntpd/images/amd64: ## Builds/pushes `openntpd/images/amd64`
getty/images/amd64: ## Builds/pushes `getty/images/amd64`
init/images/push: ## Builds/pushes `init/images/push`
ca-certificates/images/push: ## Builds/pushes `ca-certificates/images/push`
firmware/images/push: ## Builds/pushes `firmware/images/push`
rngd/images/push: ## Builds/pushes `rngd/images/push`
sysctl/images/push: ## Builds/pushes `sysctl/images/push`
sysfs/images/push: ## Builds/pushes `sysfs/images/push`
modprobe/images/push: ## Builds/pushes `modprobe/images/push`
dhcpcd/images/push: ## Builds/pushes `dhcpcd/images/push`
openntpd/images/push: ## Builds/pushes `openntpd/images/push`
getty/images/push: ## Builds/pushes `getty/images/push`

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
run-in-docker/linuxkit/src/cmd/linuxkit/eks-anywhere-go-mod-download: ## Run `linuxkit/src/cmd/linuxkit/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/linuxkit/pkg/rngd/eks-anywhere-go-mod-download: ## Run `linuxkit/pkg/rngd/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/linuxkit/pkg/sysfs/eks-anywhere-go-mod-download: ## Run `linuxkit/pkg/sysfs/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/linuxkit/pkg/sysctl/eks-anywhere-go-mod-download: ## Run `linuxkit/pkg/sysctl/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/linuxkit/pkg/init/eks-anywhere-go-mod-download: ## Run `linuxkit/pkg/init/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/_output/bin/linuxkit/linux-amd64/linuxkit: ## Run `_output/bin/linuxkit/linux-amd64/linuxkit` in docker builder container
run-in-docker/_output/bin/linuxkit/linux-amd64/rngd: ## Run `_output/bin/linuxkit/linux-amd64/rngd` in docker builder container
run-in-docker/_output/bin/linuxkit/linux-amd64/sysfs: ## Run `_output/bin/linuxkit/linux-amd64/sysfs` in docker builder container
run-in-docker/_output/bin/linuxkit/linux-amd64/sysctl: ## Run `_output/bin/linuxkit/linux-amd64/sysctl` in docker builder container
run-in-docker/_output/bin/linuxkit/linux-amd64/init: ## Run `_output/bin/linuxkit/linux-amd64/init` in docker builder container
run-in-docker/_output/bin/linuxkit/linux-amd64/rc.init: ## Run `_output/bin/linuxkit/linux-amd64/rc.init` in docker builder container
run-in-docker/_output/bin/linuxkit/linux-amd64/service: ## Run `_output/bin/linuxkit/linux-amd64/service` in docker builder container
run-in-docker/_output/attribution/go-license.csv: ## Run `_output/attribution/go-license.csv` in docker builder container
run-in-docker/_output/rngd/attribution/go-license.csv: ## Run `_output/rngd/attribution/go-license.csv` in docker builder container
run-in-docker/_output/sysfs/attribution/go-license.csv: ## Run `_output/sysfs/attribution/go-license.csv` in docker builder container
run-in-docker/_output/sysctl/attribution/go-license.csv: ## Run `_output/sysctl/attribution/go-license.csv` in docker builder container
run-in-docker/_output/init/attribution/go-license.csv: ## Run `_output/init/attribution/go-license.csv` in docker builder container
run-in-docker/_output/rc.init/attribution/go-license.csv: ## Run `_output/rc.init/attribution/go-license.csv` in docker builder container
run-in-docker/_output/service/attribution/go-license.csv: ## Run `_output/service/attribution/go-license.csv` in docker builder container

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
