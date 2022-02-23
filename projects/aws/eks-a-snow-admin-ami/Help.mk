


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


##@ Binary Targets
binaries: ## Build all binaries: `eks-a-snow-admin-ami` for `darwin/amd64`
_output/bin/eks-a-snow-admin-ami/darwin-amd64/eks-a-snow-admin-ami: ## Build `_output/bin/eks-a-snow-admin-ami/darwin-amd64/eks-a-snow-admin-ami`

##@ Checksum Targets
checksums: ## Update checksums file based on currently built binaries.
validate-checksums: # Validate checksums of currently built binaries against checksums file.

##@ Clean Targets
clean: ## Removes source and _output directory

##@ Helpers
help: ## Display this help
add-generated-help-block: ## Add or update generated help block to document project make file and support shell auto completion

##@ Build Targets
build: ## Called via prow presubmit, calls `binaries`
release: ## Called via prow postsubmit + release jobs, calls `binaries`
########### END GENERATED ###########################
