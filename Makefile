MAKEFLAGS+=--no-builtin-rules --warn-undefined-variables --no-print-directory
.SUFFIXES:

BASE_DIRECTORY:=$(abspath .)
BUILD_LIB=${BASE_DIRECTORY}/build/lib
SHELL_TRACE?=false
DEFAULT_SHELL:=$(if $(filter true,$(SHELL_TRACE)),$(BUILD_LIB)/make_shell.sh trace true,bash)
SHELL:=$(DEFAULT_SHELL)
.SHELLFLAGS:=-eu -o pipefail -c

AWS_ACCOUNT_ID?=$(shell aws sts get-caller-identity --query Account --output text)
AWS_REGION?=us-west-2
IMAGE_REPO?=$(if $(AWS_ACCOUNT_ID),$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com,localhost:5000)
ECR_PUBLIC_URI?=$(shell aws ecr-public describe-registries --region us-east-1 --query 'registries[0].registryUri' --output text)
JOB_TYPE?=

RELEASE_BRANCH?=$(LATEST_EKSD_RELEASE)
GIT_HASH=$(shell git -C $(BASE_DIRECTORY) rev-parse HEAD)
ALL_PROJECTS=$(shell $(BUILD_LIB)/all_projects.sh $(BASE_DIRECTORY))

UPLOAD_ARTIFACTS_TO_S3?=false
LATEST_EKSD_RELEASE=$(shell source $(BUILD_LIB)/common.sh && build::eksd_releases::get_release_branch)

# $1 - project name using _ as separator, ex: rancher_local-path-provisoner
PROJECT_PATH_MAP=projects/$(patsubst $(firstword $(subst _, ,$(1)))_%,$(firstword $(subst _, ,$(1)))/%,$(1))

BUILDER_PLATFORM_ARCH=$(if $(filter x86_64,$(shell uname -m)),amd64,arm64)

# $1 - variable name to resolve and cache
CACHE_RESULT = $(if $(filter undefined,$(origin _cached-$1)),$(eval _cached-$1 := 1)$(eval _cache-$1 := $($1)),)$(_cache-$1)

# $1 - variable name
CACHE_VARIABLE=$(eval _old-$(1)=$(value $(1)))$(eval $(1)=$$(call CACHE_RESULT,_old-$(1)))

CACHE_VARS=ALL_PROJECTS AWS_ACCOUNT_ID BUILDER_PLATFORM_ARCH GIT_HASH LATEST_EKSD_RELEASE ECR_PUBLIC_URI
$(foreach v,$(CACHE_VARS),$(call CACHE_VARIABLE,$(v)))

.PHONY: clean-project-%
clean-project-%: PROJECT_PATH=$(call PROJECT_PATH_MAP,$*)
clean-project-%: export RELEASE_BRANCH=$(LATEST_EKSD_RELEASE)
clean-project-%:
	$(MAKE) clean -C $(PROJECT_PATH)

.PHONY: clean
clean: $(addprefix clean-project-, $(ALL_PROJECTS))
	rm -rf _output


############################## BUILD ALL ###################################

.PHONY: build-all
build-all: build-all-warning
# Build projects with dependecies first to try and validate if there are any missing
	@set -eu -o pipefail; \
	export RELEASE_BRANCH=$(RELEASE_BRANCH); \
	PROJECTS="$(foreach project,$(ALL_PROJECTS),$(call PROJECT_PATH_MAP,$(project)))"; \
	PROJS=($${PROJECTS// / }); \
  	for proj in "$${PROJS[@]}"; do \
		if  [ -n "$$($(MAKE) -C $$proj var-value-PROJECT_DEPENDENCIES)" ] && [ ! -f $$proj/eks-anywhere-full-build-complete ]; then \
			$(MAKE) $$proj/eks-anywhere-full-build-complete; \
		fi; \
	done; \
	for proj in "$${PROJS[@]}"; do \
		if  [ ! -f $$proj/eks-anywhere-full-build-complete ]; then \
			$(MAKE) $$proj/eks-anywhere-full-build-complete; \
		fi; \
	done

# Specific overrides
projects/kubernetes-sigs/kind/eks-anywhere-full-build-complete: override IMAGE_PLATFORMS=linux/amd64,linux/arm64

# tinkerbell/hook needs to be built with a public ecr repo so docker container can pull
projects/tinkerbell/hook/eks-anywhere-full-build-complete: override IMAGE_REPO_OVERRIDE=$(or $(ECR_PUBLIC_URI),$(IMAGE_REPO))
projects/tinkerbell/hook/eks-anywhere-full-build-complete: override IMAGE_PLATFORMS=linux/amd64,linux/arm64

projects/kubernetes-sigs/image-builder/eks-anywhere-full-build-complete: MAIN_TARGET=build
projects/kubernetes-sigs/image-builder/eks-anywhere-full-build-complete: export SKIP_METAL_INSTANCE_TEST=true

projects/aws/eks-a-admin-image/eks-anywhere-full-build-complete: MAIN_TARGET=build

projects/torvalds/linux/eks-anywhere-full-build-complete: export BINARY_PLATFORMS=linux/amd64 linux/arm64
projects/isc-projects/dhcp/eks-anywhere-full-build-complete: export BINARY_PLATFORMS=linux/amd64 linux/arm64
projects/tinkerbell/ipxedust/eks-anywhere-full-build-complete: export BINARY_PLATFORMS=linux/amd64 linux/arm64

# Skips
projects/aws/cluster-api-provider-aws-snow/eks-anywhere-full-build-complete:
	@echo "Skipping aws/cluster-api-provider-aws-snow: container images are pulled cross account"
	@touch $@

projects/goharbor/harbor/eks-anywhere-full-build-complete:
	@echo "Skipping /goharbor/harbor: we patch vendor directory so we skip go mod download, which can cause slight checksum differences"
	@echo "run the 'clean-go-cache' target before running harbor if you want to build with matching checksums"
	@touch $@

# Actual target
%/eks-anywhere-full-build-complete: IMAGE_PLATFORMS=linux/$(BUILDER_PLATFORM_ARCH)
%/eks-anywhere-full-build-complete: IMAGE_REPO_OVERRIDE=$(IMAGE_REPO)
# override this on the command line to true if you want to push to your own s3 bucket
%/eks-anywhere-full-build-complete: MAIN_TARGET=release
%/eks-anywhere-full-build-complete:
	@set -eu -o pipefail; \
	export RELEASE_BRANCH=$(RELEASE_BRANCH); \
	if  [ -n "$$($(MAKE) -C $(@D) var-value-PROJECT_DEPENDENCIES)" ]; then \
		PROJECT_DEPS=$$($(MAKE) -C $(@D) var-value-PROJECT_DEPENDENCIES); \
		DEPS=($${PROJECT_DEPS// / }); \
  		for dep in "$${DEPS[@]}"; do \
			if [[ "$${dep}" = *"eksa"* ]]; then \
				OVERRIDES="IMAGE_REPO=$(IMAGE_REPO_OVERRIDE) IMAGE_PLATFORMS=$(IMAGE_PLATFORMS)"; \
				DEP_RELEASE_BRANCH="$$(cut -d/ -f4 <<< $$dep)"; \
				if [ -n "$${DEP_RELEASE_BRANCH}" ]; then \
					dep="$$(dirname $$dep)"; \
					OVERRIDES+=" RELEASE_BRANCH=$$DEP_RELEASE_BRANCH"; \
				fi; \
				if [ -f projects/$${dep#"eksa/"}/eks-anywhere-full-build-complete-$$DEP_RELEASE_BRANCH ]; then \
					continue; \
				fi; \
				echo "Running make $${dep#eksa/} as dependency for $(@D)"; \
				$(MAKE) projects/$${dep#"eksa/"}/eks-anywhere-full-build-complete $$OVERRIDES; \
				if [ -n "$${DEP_RELEASE_BRANCH}" ]; then \
					mv projects/$${dep#"eksa/"}/eks-anywhere-full-build-complete projects/$${dep#"eksa/"}/eks-anywhere-full-build-complete-$$DEP_RELEASE_BRANCH; \
				fi; \
			fi; \
		done; \
	fi; \
	TARGETS="attribution $(MAIN_TARGET)"; \
	if [[ $(IMAGE_REPO_OVERRIDE) == *"ecr"* ]] && [[ -n "$$($(MAKE) -C $(@D) var-value-IMAGE_NAMES)"  ||  "$$($(MAKE) -C $(@D) var-value-HAS_HELM_CHART)" = "true" ]]; then \
		TARGETS="create-ecr-repos $${TARGETS}"; \
	fi; \
	if [ "$(UPLOAD_ARTIFACTS_TO_S3)" = "true" ]; then \
		TARGETS+=" UPLOAD_DRY_RUN=false UPLOAD_CREATE_PUBLIC_ACL=false"; \
	fi; \
	TARGETS+=" IMAGE_REPO=$(IMAGE_REPO_OVERRIDE) IMAGE_PLATFORMS=$(IMAGE_PLATFORMS) RELEASE_BRANCH=$(RELEASE_BRANCH)"; \
	echo "Running 'make -C $(@D) $${TARGETS}'"; \
	make -C $(@D) $${TARGETS}; \
	touch $@

.PHONY: build-all-warning
build-all-warning:
	@echo "*** Warning: this target is not meant to used except for specific testing situations ***"
	@echo "*** this will likely fail and either way run for a really long time ***"

#########################################################################
# to make running on mac/linux or amd/arm consistent exporting certain vars to overwrite default
.PHONY: add-generated-help-block-project-%
add-generated-help-block-project-%:
	$(eval PROJECT_PATH=$(call PROJECT_PATH_MAP,$*))
	$(MAKE) add-generated-help-block -C $(PROJECT_PATH) RELEASE_BRANCH=1-26 BUILDER_PLATFORM_ARCH=amd64 

.PHONY: add-generated-help-block
add-generated-help-block: $(addprefix add-generated-help-block-project-, $(ALL_PROJECTS))
	build/update-attribution-files/create_pr.sh

.PHONY: attribution-files-project-%
attribution-files-project-%:
	$(eval PROJECT_PATH=$(call PROJECT_PATH_MAP,$*))
	if $(MAKE) -C $(PROJECT_PATH) check-for-release-branch-skip; then \
		$(MAKE) -C $(PROJECT_PATH) all-attributions; \
	fi
	
.PHONY: attribution-files
attribution-files: $(addprefix attribution-files-project-, $(ALL_PROJECTS))
	cat _output/total_summary.txt

.PHONY: checksum-files-project-%
checksum-files-project-%:
	$(eval PROJECT_PATH=$(call PROJECT_PATH_MAP,$*))
	$(MAKE) -C $(PROJECT_PATH) all-checksums

.PHONY: update-checksum-files
update-checksum-files: $(addprefix checksum-files-project-, $(ALL_PROJECTS))
	build/lib/update_go_versions.sh
	build/update-attribution-files/create_pr.sh

.PHONY: update-attribution-files
update-attribution-files: add-generated-help-block attribution-files
	build/update-attribution-files/create_pr.sh

.PHONY: start-docker-builder
start-docker-builder: # Start long lived builder base docker container
	@$(MAKE) -C projects/aws/eks-anywhere-build-tooling start-docker-builder

.PHONY: stop-docker-builder
stop-docker-builder:
	docker rm -f -v eks-a-builder

.PHONY: run-buildkit-and-registry
run-buildkit-and-registry:
	docker run --rm -d --name buildkitd --net host --privileged moby/buildkit:v0.12.2-rootless
	docker run --rm -d --name registry  --net host registry:2

.PHONY: stop-buildkit-and-registry
stop-buildkit-and-registry:
	docker rm -v --force buildkitd
	docker rm -v --force registry

.PHONY: generate-project-list
generate-project-list: | ensure-locale
	build/lib/generate_projects_list.sh $(BASE_DIRECTORY)

.PHONY: generate-staging-buildspec
generate-staging-buildspec: export BINARY_PLATFORMS=linux/amd64 linux/arm64
generate-staging-buildspec: export IMAGE_PLATFORMS=linux/amd64 linux/arm64
generate-staging-buildspec: | ensure-locale
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "$(ALL_PROJECTS)" "$(BASE_DIRECTORY)/release/staging-build.yml" "$(BASE_DIRECTORY)/buildspec.yml" false EXCLUDE_FROM_STAGING_BUILDSPEC BUILDSPECS false
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "$(ALL_PROJECTS)" "$(BASE_DIRECTORY)/release/checksums-build.yml" "$(BASE_DIRECTORY)/buildspecs/checksums-buildspec.yml" true EXCLUDE_FROM_CHECKSUMS_BUILDSPEC CHECKSUMS_BUILDSPECS false buildspecs/checksums-pr-buildspec.yml
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "aws_bottlerocket-bootstrap" "$(BASE_DIRECTORY)/projects/aws/bottlerocket-bootstrap/buildspecs/batch-build.yml" "$(BASE_DIRECTORY)/buildspec.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "kubernetes_cloud-provider-vsphere" "$(BASE_DIRECTORY)/projects/kubernetes/cloud-provider-vsphere/buildspecs/batch-build.yml" "$(BASE_DIRECTORY)/buildspec.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "kubernetes-sigs_kind" "$(BASE_DIRECTORY)/projects/kubernetes-sigs/kind/buildspecs/batch-build.yml" "$(BASE_DIRECTORY)/buildspecs/images.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "containerd_containerd" "$(BASE_DIRECTORY)/projects/containerd/containerd/buildspecs/batch-build.yml" "$(BASE_DIRECTORY)/buildspec.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "opencontainers_runc" "$(BASE_DIRECTORY)/projects/opencontainers/runc/buildspecs/batch-build.yml" "$(BASE_DIRECTORY)/buildspec.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "torvalds_linux" "$(BASE_DIRECTORY)/projects/torvalds/linux/buildspecs/batch-build.yml" "$(BASE_DIRECTORY)/buildspec.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "isc-projects_dhcp" "$(BASE_DIRECTORY)/projects/isc-projects/dhcp/buildspecs/batch-build.yml" "$(BASE_DIRECTORY)/buildspec.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "tinkerbell_hook" "$(BASE_DIRECTORY)/projects/tinkerbell/hook/buildspecs/batch-build.yml" "$(BASE_DIRECTORY)/buildspec.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "tinkerbell_tinkerbell" "$(BASE_DIRECTORY)/projects/tinkerbell/tinkerbell/buildspecs/batch-build.yml" "$(BASE_DIRECTORY)/buildspec.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "tinkerbell_ipxedust" "$(BASE_DIRECTORY)/projects/tinkerbell/ipxedust/buildspecs/batch-build.yml" "$(BASE_DIRECTORY)/buildspec.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "tinkerbell_rufio" "$(BASE_DIRECTORY)/projects/tinkerbell/rufio/buildspecs/batch-build.yml" "$(BASE_DIRECTORY)/buildspec.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "tinkerbell_tink" "$(BASE_DIRECTORY)/projects/tinkerbell/tink/buildspecs/batch-build.yml" "$(BASE_DIRECTORY)/buildspec.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "linuxkit_linuxkit" "$(BASE_DIRECTORY)/projects/linuxkit/linuxkit/buildspecs/batch-build.yml" "$(BASE_DIRECTORY)/buildspec.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "emissary-ingress_emissary" "$(BASE_DIRECTORY)/projects/emissary-ingress/emissary/buildspecs/batch-build.yml" "$(BASE_DIRECTORY)/buildspec.yml" true "DO_NOT_EXCLUDE_FROM_BUILDSPEC"
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "distribution_distribution" "$(BASE_DIRECTORY)/projects/distribution/distribution/buildspecs/batch-build.yml" "$(BASE_DIRECTORY)/buildspec.yml" true "DO_NOT_EXCLUDE_FROM_BUILDSPEC"
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "goharbor_harbor" "$(BASE_DIRECTORY)/projects/goharbor/harbor/buildspecs/batch-build.yml" "$(BASE_DIRECTORY)/buildspec.yml" true "DO_NOT_EXCLUDE_FROM_BUILDSPEC"
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "aws_upgrader" "$(BASE_DIRECTORY)/projects/aws/upgrader/buildspecs/batch-build.yml" "$(BASE_DIRECTORY)/buildspec.yml" true "DO_NOT_EXCLUDE_FROM_BUILDSPEC"
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "$(ALL_PROJECTS)" "$(BASE_DIRECTORY)/tools/version-tracker/buildspecs/upgrade.yml" "$(BASE_DIRECTORY)/buildspecs/upgrade-buildspec.yml" true EXCLUDE_FROM_UPGRADE_BUILDSPEC UPGRADE_BUILDSPECS false "buildspecs/upgrade-eks-distro-buildspec.yml,buildspecs/upgrade-eks-distro-build-tooling-buildspec.yml" true

.PHONY: generate
generate: generate-project-list generate-staging-buildspec

.PHONY: validate-generated
validate-generated: generate validate-release-buildspecs validate-eksd-releases
	build/lib/readme_check.sh $(BASE_DIRECTORY)
	@if [ "$$(git status --porcelain -- UPSTREAM_PROJECTS.yaml release/staging-build.yml release/checksums-build.yml **/batch-build.yml **/README.md | wc -l)" -gt 0 ]; then \
		echo "Error: Generated files, UPSTREAM_PROJECTS.yaml README.md release/staging-build.yml release/checksums-build.yml batch-build.yml, do not match expected. Please run 'make generate' to update"; \
		git --no-pager diff -- UPSTREAM_PROJECTS.yaml release/staging-build.yml release/checksums-build.yml **/batch-build.yml **/README.md; \
		exit 1; \
	fi

.PHONY: check-project-path-exists
check-project-path-exists:
	@if ! stat $(PROJECT_PATH) &> /dev/null; then \
		echo "false"; \
	else \
		echo "true"; \
	fi

.PHONY: validate-release-buildspecs
validate-release-buildspecs:
	build/lib/validate_release_buildspecs.sh "$(BASE_DIRECTORY)/release/checksums-build.yml" "$(BASE_DIRECTORY)/release/staging-build.yml" "$(BASE_DIRECTORY)/tools/version-tracker/buildspecs/upgrade.yml"

.PHONY: validate-eksd-releases
validate-eksd-releases:
	build/lib/validate_eksd_releases.sh

# Locale settings impact file ordering in ls or shell file expansion. The file order is used to
# generate files that are subsequently validated by the CI. If local environments use different 
# locales to the CI we get unexpected failures that are tricky to debug without knowledge of 
# locales so we'll explicitly warn here.
# In a AL2 container image (like builder base), LANG will be empty which is equivalent to posix
# In a AL2 (or other distro) full instance the LANG will be en-us.UTF-8 which produces different sorts
# On Mac, LANG will be en-us.UTF-8 but has a fix applied to sort to avoid the difference
.PHONY: ensure-locale
ensure-locale:
	@if [ "Linux" = "$$(uname -s)" ]; then \
		LOCALE=$$(locale | grep LANG | cut -d= -f2 | tr -d '"' | tr '[:upper:]' '[:lower:]'); \
		if [[ "c.utf-8 posix" != *"$${LOCALE:-posix}"* ]]; then \
			echo WARNING: Environment locale set to $$LOCALE. On Linux systems this may create \
				non-deterministic behavior when running generation recipes. If the CI fails validation try \
				exporting LANG=C.UTF-8 to generate files instead.; \
		fi; \
	fi

.PHONY: get-default-release-branch
get-default-release-branch:
	@echo $(LATEST_EKSD_RELEASE)
