


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion


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

##@ Build Targets
build: ## Called via prow presubmit, calls `validate`
release: ## Called via prow postsubmit + release jobs, calls `build-raw-image`
########### END GENERATED ###########################
