SHELL := bash
.SHELLFLAGS:=-eux -o pipefail -c
BASE_DIRECTORY:=$(shell git rev-parse --show-toplevel)

BINARY_NAME?=version-tracker
BINARY_PATH?=bin/$(BINARY_NAME)
GOLANG_VERSION=1.24
GO ?= $(shell source $(BASE_DIRECTORY)/build/lib/common.sh && build::common::get_go_path $(GOLANG_VERSION))/go
DRY_RUN?=false
VERBOSITY?=0

build:
	CGO_ENABLED=0 $(GO) build -o $(BINARY_PATH) main.go

.PHONY: test
test:
	$(GO) test . ./cmd ./pkg/...
	python3 -m unittest discover -s ../../projects/kubernetes/autoscaler -p 'test_*.py'

upgrade: build
	$(BINARY_PATH) upgrade --project $(PROJECT) --dry-run=$(DRY_RUN) --verbosity $(VERBOSITY)

clean:
	rm -rf eks-anywhere-build-tooling
	git clean -fdx
