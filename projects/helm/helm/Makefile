BASE_DIRECTORY:=$(abspath ../../../)
GIT_TAG=$(shell cat GIT_TAG)
GOLANG_VERSION=$(shell cat GOLANG_VERSION)
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

include $(BASE_DIRECTORY)/Common.mk

########### DO NOT EDIT #############################
# To update call: make add-generated-help-block
# This is added to help document dynamic targets and support shell autocompletion
# Run make help for a formatted help block with all targets
include Help.mk
########### END GENERATED ###########################
