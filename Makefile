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

# $1 - project name using _ as seperator, ex: rancher_local-path-provisoner
PROJECT_PATH_MAP=projects/$(patsubst $(firstword $(subst _, ,$(1)))_%,$(firstword $(subst _, ,$(1)))/%,$(1))

# Locale settings impact file ordering in ls or shell file expansion. The file order is used to
# generate files that are subsequently validated by the CI. If local environments use different 
# locales to the CI we get unexpected failures that are tricky to debug without knowledge of 
# locales so we'll explicitly warn here.
TO_LOWER = $(shell echo $(1) | tr '[:upper:]' '[:lower:]')
ifeq ($(shell uname -s),Linux)
  LOCALE := $(call TO_LOWER,$(shell locale | grep LANG | cut -d= -f2 | tr -d '"'))
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
	$(MAKE) add-generated-help-block -C $(PROJECT_PATH) RELEASE_BRANCH=1-21

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

.PHONY: run-target-in-docker
run-target-in-docker:
	build/lib/run_target_docker.sh $(PROJECT) $(MAKE_TARGET) $(IMAGE_REPO) "$(RELEASE_BRANCH)" $(ARTIFACTS_BUCKET)

.PHONY: update-attribution-checksums-docker
update-attribution-checksums-docker:
	build/lib/update_checksum_docker.sh $(PROJECT) $(IMAGE_REPO) $(RELEASE_BRANCH)

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
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "aws_bottlerocket-bootstrap" "$(BASE_DIRECTORY)/projects/aws/bottlerocket-bootstrap/buildspecs/batch-build.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "kubernetes_cloud-provider-vsphere" "$(BASE_DIRECTORY)/projects/kubernetes/cloud-provider-vsphere/buildspecs/batch-build.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "kubernetes-sigs_kind" "$(BASE_DIRECTORY)/projects/kubernetes-sigs/kind/buildspecs/batch-build.yml" true
	build/lib/generate_staging_buildspec.sh $(BASE_DIRECTORY) "fluxcd_source-controller" "$(BASE_DIRECTORY)/projects/fluxcd/source-controller/buildspecs/batch-build.yml" false

.PHONY: generate
generate: generate-project-list generate-staging-buildspec

.PHONY: validate-generated
validate-generated: generate
	@if [ "$$(git status --porcelain -- UPSTREAM_PROJECTS.yaml release/staging-build.yml **/batch-build.yml | wc -l)" -gt 0 ]; then \
		echo "Error: Generated files, UPSTREAM_PROJECTS.yaml release/staging-build.yml, do not match expected. Please run `make generate` to update"; \
		git diff -- UPSTREAM_PROJECTS.yaml release/staging-build.yml **/batch-build.yml; \
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
