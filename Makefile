BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
BUILD_LIB=${BASE_DIRECTORY}/build/lib
AWS_ACCOUNT_ID?=$(shell aws sts get-caller-identity --query Account --output text)
AWS_REGION?=us-west-2
IMAGE_REPO?=$(if $(AWS_ACCOUNT_ID),$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com,localhost:5000)
ECR_PUBLIC_URI?=$(shell aws ecr-public describe-registries --region us-east-1 --query 'registries[0].registryUri' --output text)
JOB_TYPE?=

RELEASE_BRANCH?=
GIT_HASH:=$(shell git -C $(BASE_DIRECTORY) rev-parse HEAD)
ALL_PROJECTS=$(shell $(BUILD_LIB)/all_projects.sh $(BASE_DIRECTORY))

# $1 - project name using _ as separator, ex: rancher_local-path-provisoner
PROJECT_PATH_MAP=projects/$(patsubst $(firstword $(subst _, ,$(1)))_%,$(firstword $(subst _, ,$(1)))/%,$(1))

# Locale settings impact file ordering in ls or shell file expansion. The file order is used to
# generate files that are subsequently validated by the CI. If local environments use different 
# locales to the CI we get unexpected failures that are tricky to debug without knowledge of 
# locales so we'll explicitly warn here.
# In a AL2 container image (like builder base), LANG will be empty which is equivalent to posix
# In a AL2 (or other distro) full instance the LANG will be en-us.UTF-8 which produces different sorts
# On Mac, LANG will be en-us.UTF-8 but has a fix applied to sort to avoid the difference
ifeq ($(shell uname -s),Linux)
  LOCALE:=$(call TO_LOWER,$(shell locale | grep LANG | cut -d= -f2 | tr -d '"'))
  LOCALE:=$(if $(LOCALE),$(LOCALE),posix)
  ifeq ($(filter c.utf-8 posix,$(LOCALE)),)
    $(warning WARNING: Environment locale set to $(LANG). On Linux systems this may create \
	non-deterministic behavior when running generation recipes. If the CI fails validation try \
	`LANG=C.UTF-8 make <recipe>` to generate files instead.)
  endif
endif

.PHONY: clean-project-%
clean-project-%:
	$(eval PROJECT_PATH=$(call PROJECT_PATH_MAP,$*))
	$(MAKE) clean -C $(PROJECT_PATH)

.PHONY: clean
clean: $(addprefix clean-project-, $(ALL_PROJECTS))
	rm -rf _output

.PHONY: add-generated-help-block-project-%
add-generated-help-block-project-%:
	$(eval PROJECT_PATH=$(call PROJECT_PATH_MAP,$*))
	$(MAKE) add-generated-help-block -C $(PROJECT_PATH) RELEASE_BRANCH=1-26

.PHONY: add-generated-help-block
add-generated-help-block: $(addprefix add-generated-help-block-project-, $(ALL_PROJECTS))
	build/update-attribution-files/create_pr.sh

.PHONY: attribution-files-project-%
attribution-files-project-%:
	$(eval PROJECT_PATH=$(call PROJECT_PATH_MAP,$*))
	$(MAKE) -C $(PROJECT_PATH) all-attributions

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

.PHONY: stop-docker-builder
stop-docker-builder:
	docker rm -f -v eks-a-builder

.PHONY: run-buildkit-and-registry
run-buildkit-and-registry:
	docker run -d --name buildkitd --net host --privileged moby/buildkit:v0.10.6-rootless
	docker run -d --name registry  --net host registry:2

.PHONY: stop-buildkit-and-registry
stop-buildkit-and-registry:
	docker rm -v --force buildkitd
	docker rm -v --force registry

.PHONY: generate-project-list
generate-project-list:
	build/lib/generate_projects_list.sh $(BASE_DIRECTORY)

.PHONY: generate-staging-buildspec
generate-staging-buildspec:
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "$(ALL_PROJECTS)" "$(BASE_DIRECTORY)/release/staging-build.yml"
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "$(ALL_PROJECTS)" "$(BASE_DIRECTORY)/release/checksums-build.yml" true EXCLUDE_FROM_CHECKSUMS_BUILDSPEC CHECKSUMS_BUILDSPECS false buildspecs/checksums-pr-buildspec.yml
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "aws_bottlerocket-bootstrap" "$(BASE_DIRECTORY)/projects/aws/bottlerocket-bootstrap/buildspecs/batch-build.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "kubernetes_cloud-provider-vsphere" "$(BASE_DIRECTORY)/projects/kubernetes/cloud-provider-vsphere/buildspecs/batch-build.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "kubernetes-sigs_kind" "$(BASE_DIRECTORY)/projects/kubernetes-sigs/kind/buildspecs/batch-build.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "containerd_containerd" "$(BASE_DIRECTORY)/projects/containerd/containerd/buildspecs/batch-build.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "opencontainers_runc" "$(BASE_DIRECTORY)/projects/opencontainers/runc/buildspecs/batch-build.yml" true

.PHONY: generate
generate: generate-project-list generate-staging-buildspec

.PHONY: validate-generated
validate-generated: generate
	@if [ "$$(git status --porcelain -- UPSTREAM_PROJECTS.yaml release/staging-build.yml release/checksums-build.yml **/batch-build.yml | wc -l)" -gt 0 ]; then \
		echo "Error: Generated files, UPSTREAM_PROJECTS.yaml release/staging-build.yml release/checksums-build.yml batch-build.yml, do not match expected. Please run `make generate` to update"; \
		git diff -- UPSTREAM_PROJECTS.yaml release/staging-build.yml release/checksums-build.yml **/batch-build.yml; \
		exit 1; \
	fi
	build/lib/readme_check.sh

.PHONY: check-project-path-exists
check-project-path-exists:
	@if ! stat $(PROJECT_PATH) &> /dev/null; then \
		echo "false"; \
	else \
		echo "true"; \
	fi
