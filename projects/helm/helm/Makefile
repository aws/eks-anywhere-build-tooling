BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)
GIT_TAG:=$(shell cat GIT_TAG)
GOLANG_VERSION?="1.19"

REPO=helm
REPO_OWNER=helm

BINARY_TARGET_FILES=helm
SOURCE_PATTERNS=./cmd/helm

# Match upstream makefile - https://github.com/helm/helm/blob/main/Makefile
K8S_MODULES_VER=$(subst ., ,$(subst v,,$(shell cd $(REPO) && go list -f '{{.Version}}' -m k8s.io/client-go)))
K8S_MODULES_MAJOR_VER=$(shell echo $$(($(firstword $(K8S_MODULES_VER)) + 1)))
K8S_MODULES_MINOR_VER=$(word 2,$(K8S_MODULES_VER))

EXTRA_GO_LDFLAGS=-X helm.sh/helm/v3/internal/version.version=$(GIT_TAG) \
	-X helm.sh/helm/v3/internal/version.metadata= \
	-X helm.sh/helm/v3/internal/version.gitCommit=$(shell git -C $(REPO) rev-list -n 1  "${GIT_TAG}") \
	-X helm.sh/helm/v3/internal/version.gitTreeState=clean \
	-X helm.sh/helm/v3/pkg/lint/rules.k8sVersionMajor=$(K8S_MODULES_MAJOR_VER) \
	-X helm.sh/helm/v3/pkg/lint/rules.k8sVersionMinor=$(K8S_MODULES_MINOR_VER) \
	-X helm.sh/helm/v3/pkg/chartutil.k8sVersionMajor=$(K8S_MODULES_MAJOR_VER) \
	-X helm.sh/helm/v3/pkg/chartutil.k8sVersionMinor=$(K8S_MODULES_MINOR_VER)

HAS_S3_ARTIFACTS=true

IMAGE_NAMES=

FIX_LICENSES_XEIPUUV_TARGET=$(REPO)/vendor/github.com/xeipuuv/gojsonpointer/LICENSE.txt

include $(BASE_DIRECTORY)/Common.mk


$(GATHER_LICENSES_TARGETS): $(FIX_LICENSES_XEIPUUV_TARGET)

$(FIX_LICENSES_XEIPUUV_TARGET): | $(GO_MOD_DOWNLOAD_TARGETS)
# The xeipuuv dependency github repos all have licenses however they all do not have go.mod files
# checked in to the repo. Hence we need to manually download licenses from Github for each of them 
# and place them in the respective folders under vendor directory so that they is available for 
# go-licenses to pick up	
	for package in gojsonpointer gojsonreference gojsonschema ; do \
		wget --retry-connrefused -q https://raw.githubusercontent.com/xeipuuv/$$package/master/LICENSE-APACHE-2.0.txt -O \
			$(REPO)/vendor/github.com/xeipuuv/$$package/LICENSE.txt; \
	done;


########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################
