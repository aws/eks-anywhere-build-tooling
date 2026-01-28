



########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ GIT/Coverage

clone-repo:  ## Clone upstream `tinkerbell`
checkout-repo: ## Checkout upstream tag based on value in GIT_TAG file
patch-repo: ## Patch upstream repo with patches in patches directory

##@ Binary Coverage

binaries: ## Build all binaries: `tinkerbell tink-agent` for `linux/amd64 linux/arm64`
_output/bin/tinkerbell/linux-amd64/tinkerbell: ## Build `_output/bin/tinkerbell/linux-amd64/tinkerbell`
_output/bin/tinkerbell/linux-amd64/tink-agent: ## Build `_output/bin/tinkerbell/linux-amd64/tink-agent`
_output/bin/tinkerbell/linux-arm64/tinkerbell: ## Build `_output/bin/tinkerbell/linux-arm64/tinkerbell`
_output/bin/tinkerbell/linux-arm64/tink-agent: ## Build `_output/bin/tinkerbell/linux-arm64/tink-agent`

##@ Image Targets

local-images: ## Builds `tinkerbell/images/amd64 tink-agent/images/amd64` as oci tars for local testing
images: ## Pushes `tinkerbell/images/push tink-agent/images/push` to IMAGE_REPO
tinkerbell/images/amd64: ## Builds/pushes `tinkerbell/images/amd64`
tink-agent/images/amd64: ## Builds/pushes `tink-agent/images/amd64`
tinkerbell/images/push: ## Builds/pushes `tinkerbell/images/push`
tink-agent/images/push: ## Builds/pushes `tink-agent/images/push`

##@ Fetch Binary Targets

_output/dependencies/linux-amd64/eksa/tinkerbell/tinkerbell: ## Fetch `_output/dependencies/linux-amd64/eksa/tinkerbell/tinkerbell`
_output/dependencies/linux-arm64/eksa/tinkerbell/tinkerbell: ## Fetch `_output/dependencies/linux-arm64/eksa/tinkerbell/tinkerbell`

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
run-in-docker/tinkerbell/eks-anywhere-go-mod-download: ## Run `tinkerbell/eks-anywhere-go-mod-download` in docker builder container
run-in-docker/_output/bin/tinkerbell/linux-amd64/tinkerbell: ## Run `_output/bin/tinkerbell/linux-amd64/tinkerbell` in docker builder container
run-in-docker/_output/bin/tinkerbell/linux-amd64/tink-agent: ## Run `_output/bin/tinkerbell/linux-amd64/tink-agent` in docker builder container
run-in-docker/_output/bin/tinkerbell/linux-arm64/tinkerbell: ## Run `_output/bin/tinkerbell/linux-arm64/tinkerbell` in docker builder container
run-in-docker/_output/bin/tinkerbell/linux-arm64/tink-agent: ## Run `_output/bin/tinkerbell/linux-arm64/tink-agent` in docker builder container
run-in-docker/_output/attribution/go-license.csv: ## Run `_output/attribution/go-license.csv` in docker builder container

##@ License Targets

gather-licenses: ## Helper to call $(GATHER_LICENSES_TARGETS) which gathers all licenses
attribution: ## Generates attribution from licenses gathered during `gather-licenses`.
attribution-pr: ## Generates PR to update attribution files for projects

##@ Clean Targets

clean: ## Removes source and samples directories
clean-go-cache: ## Removes the GOMODCACHE AND GOCACHE folders
clean-repo: ## Removes source directory

##@Fetch Licenses Targets

_output/attribution/go-license.csv: ## Fetch licenses for tinkerbell project

##@ Helpers

help: ## Display this help
add-hierarchical-help-block: ## Add hierarchical help block to Help.mk
add-generated-help-block: ## Add generated help block to Help.mk

##@Update Coverage

run-target-in-docker: ## Run `MAKE_TARGET` using builder base docker container
update-attribution-checksums-docker: ## Update attribution and checksums using the builder base docker container
stop-docker-builder: ## Clean up builder base docker container
generate: ## Update UPSTREAM_PROJECTS.yaml
update-go-mods: ## Update locally checked-in go sum to assist in vuln scanning
update-vendor-for-dep-patch: ## After bumping dep in go.mod file, uses generic vendor update script or one provided from upstream project
patch-for-dep-update: ## After bumping dep in go.mod file and updating vendor, generates patch
create-hierarchical-ecr-repos: ## Create ECR repos for project images for local testing

##@ Build Coverage

build: ## Called via prow presubmit, calls `github-rate-limit-pre validate-checksums attribution local-images upload-artifacts attribution-pr github-rate-limit-post`
release: ## Called via prow postsubmit + release phase, calls `validate-checksums images `
########### END GENERATED ###########################
