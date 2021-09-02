BASE_DIRECTORY=$(shell git rev-parse --show-toplevel)
ARTIFACTS_BUCKET?=my-s3-bucket
GIT_TAG?=$(shell cat GIT_TAG)
HELM_REPO_VERSION?=$(shell cat HELM_REPO_VERSION)

ifdef CODEBUILD_SRC_DIR
	ARTIFACTS_PATH?=$(CODEBUILD_SRC_DIR)/$(PROJECT_PATH)/$(CODEBUILD_BUILD_NUMBER)-$(CODEBUILD_RESOLVED_SOURCE_VERSION)/artifacts
else
	ARTIFACTS_PATH?="_output/artifacts"
endif

.PHONY: upload-artifacts
upload-artifacts:
	mkdir -p $(ARTIFACTS_PATH)
	cp -rf _output/manifests/ $(ARTIFACTS_PATH)
	build/upload_artifacts.sh $(ARTIFACTS_PATH) $(ARTIFACTS_BUCKET) $(PROJECT_PATH) $(CODEBUILD_BUILD_NUMBER) $(CODEBUILD_RESOLVED_SOURCE_VERSION)

.PHONY: manifests
manifests:
	build/create_manifests.sh $(GIT_TAG) $(HELM_REPO_VERSION)

.PHONY: build
build: manifests
	echo "Done $(COMPONENT)"

.PHONY: release
release: build
	echo "Done $(COMPONENT)"

.PHONY: all
all: release

.PHONY: clean
clean:
	rm -rf $(REPO)
	rm -rf "_output"
