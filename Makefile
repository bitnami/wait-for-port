.PHONY: all build clean download get-build-deps vet lint test cover

TOOL_NAME := wait-for-port

# Load relative to the common.mk file
include $(dir $(lastword $(MAKEFILE_LIST)))/vars.mk

include ./vars.mk

all:
	@$(MAKE) get-build-deps
	$(MAKE) download
	@$(MAKE) vet
	@$(MAKE) lint
	@$(MAKE) cover
	@$(MAKE) build

build/%:
	@echo "Building GOARCH=$(*F)"
	@GOARCH=$(*F) go build -ldflags=$(LDFLAGS) -o $(TOOL_PATH)
	@echo "*** Binary created under $(TOOL_PATH) ***"

build: build/amd64

clean:
	@rm -rf $(BUILD_DIR)

download:
	$(GO_MOD) download

get-build-deps:
	@echo "+ Downloading build dependencies"
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest

vet:
	@echo "+ Vet"
	@go vet ./...

lint:
	@echo "+ Linting package"
	@staticcheck ./...
	$(call fmtcheck, .)

get-deps:
	@echo "+ Downloading dependencies"
	@go get -d -t ./...

test:
	@echo "+ Testing package"
	$(GO_TEST) .

cover: test
	@echo "+ Tests Coverage"
	@mkdir -p $(BUILD_DIR)
	@touch $(BUILD_DIR)/cover.out
	@go test -coverprofile=$(BUILD_DIR)/cover.out
	@go tool cover -html=$(BUILD_DIR)/cover.out -o=$(BUILD_DIR)/coverage.html
