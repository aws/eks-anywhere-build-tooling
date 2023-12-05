BINARY_NAME?=version-tracker
BINARY_PATH?=bin/$(BINARY_NAME)
GO?=$(shell which go)
DRY_RUN?=false

build:
	CGO_ENABLED=0 $(GO) build -o $(BINARY_PATH) main.go

upgrade: build
	$(BINARY_PATH) upgrade --project $(PROJECT) --dry-run=$(DRY_RUN)

clean:
	rm -rf eks-anywhere-build-tooling
	git clean -fdx
