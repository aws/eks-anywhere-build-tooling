


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ Binary Targets
binaries: ## Build all binaries: `bottlerocket-bootstrap bottlerocket-bootstrap-snow` for `linux/amd64 linux/arm64`
_output/bin/bottlerocket-bootstrap/linux-amd64/bottlerocket-bootstrap: ## Build `_output/bin/bottlerocket-bootstrap/linux-amd64/bottlerocket-bootstrap`
_output/bin/bottlerocket-bootstrap/linux-amd64/bottlerocket-bootstrap-snow: ## Build `_output/bin/bottlerocket-bootstrap/linux-amd64/bottlerocket-bootstrap-snow`
_output/bin/bottlerocket-bootstrap/linux-arm64/bottlerocket-bootstrap: ## Build `_output/bin/bottlerocket-bootstrap/linux-arm64/bottlerocket-bootstrap`
_output/bin/bottlerocket-bootstrap/linux-arm64/bottlerocket-bootstrap-snow: ## Build `_output/bin/bottlerocket-bootstrap/linux-arm64/bottlerocket-bootstrap-snow`

##@ Image Targets
local-images: ## Builds `bottlerocket-bootstrap/images/amd64 bottlerocket-bootstrap-snow/images/amd64` as oci tars for presumbit validation
images: ## Pushes `bottlerocket-bootstrap/images/push bottlerocket-bootstrap-snow/images/push` to IMAGE_REPO
bottlerocket-bootstrap/images/amd64: ## Builds/pushes `bottlerocket-bootstrap/images/amd64`
bottlerocket-bootstrap-snow/images/amd64: ## Builds/pushes `bottlerocket-bootstrap-snow/images/amd64`
bottlerocket-bootstrap/images/push: ## Builds/pushes `bottlerocket-bootstrap/images/push`
bottlerocket-bootstrap-snow/images/push: ## Builds/pushes `bottlerocket-bootstrap-snow/images/push`

##@ Fetch Binary Targets
_output/1-21/dependencies/linux-amd64/eksd/kubernetes/client: ## Fetch `_output/1-21/dependencies/linux-amd64/eksd/kubernetes/client`
_output/1-21/dependencies/linux-arm64/eksd/kubernetes/client: ## Fetch `_output/1-21/dependencies/linux-arm64/eksd/kubernetes/client`
_output/1-21/dependencies/linux-amd64/eksd/kubernetes/server: ## Fetch `_output/1-21/dependencies/linux-amd64/eksd/kubernetes/server`
_output/1-21/dependencies/linux-arm64/eksd/kubernetes/server: ## Fetch `_output/1-21/dependencies/linux-arm64/eksd/kubernetes/server`
_output/1-21/dependencies/linux-amd64/eksa/kubernetes-sigs/etcdadm: ## Fetch `_output/1-21/dependencies/linux-amd64/eksa/kubernetes-sigs/etcdadm`
_output/1-21/dependencies/linux-arm64/eksa/kubernetes-sigs/etcdadm: ## Fetch `_output/1-21/dependencies/linux-arm64/eksa/kubernetes-sigs/etcdadm`

##@ Checksum Targets
checksums: ## Update checksums file based on currently built binaries.
validate-checksums: # Validate checksums of currently built binaries against checksums file.

##@ License Targets
gather-licenses: ## Helper to call $(GATHER_LICENSES_TARGETS) which gathers all licenses
attribution: ## Generates attribution from licenses gathered during `gather-licenses`.
attribution-pr: ## Generates PR to update attribution files for projects

##@ Clean Targets
clean: ## Removes source and _output directory

##@ Helpers
help: ## Display this help
add-generated-help-block: ## Add or update generated help block to document project make file and support shell auto completion

##@Update Helpers
run-target-in-docker: ## Run `MAKE_TARGET` using builder base docker container
update-attribution-checksums-docker: ## Update attribution and checksums using the builder base docker container
stop-docker-builder: ## Clean up builder base docker container
generate: ## Update UPSTREAM_PROJECTS.yaml
update-go-mods: ## Update locally checked-in go sum to assist in vuln scanning
update-vendor-for-dep-patch: ## After bumping dep in go.mod file, uses generic vendor update script or one provided from upstream project
patch-for-dep-update: ## After bumping dep in go.mod file and updating vendor, generates patch
create-ecr-repos: ## Create repos in ECR for project images for local testing

##@ Build Targets
build: ## Called via prow presubmit, calls `validate-checksums attribution local-images   attribution-pr`
release: ## Called via prow postsubmit + release jobs, calls `validate-checksums images  `
########### END GENERATED ###########################
