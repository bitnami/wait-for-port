.PHONY: dep-ensure test cover all lint

TOOL_NAME := wait-for-port

# Load relative to the common.mk file
include $(dir $(lastword $(MAKEFILE_LIST)))/vars.mk

include ./vars.mk

all:
	@$(MAKE) get-build-deps
	$(DEP_ENSURE) -vendor-only
	@$(MAKE) vet
	@$(MAKE) lint
	@$(MAKE) cover
	@$(MAKE) build

build:
	@go build -ldflags="-s -w" -o $(TOOL_PATH)

clean:
	@rm -rf $(BUILD_DIR)

dep-ensure:
	$(DEP_ENSURE)

get-build-deps:
	@echo "+ Downloading build dependencies"
	@go get golang.org/x/tools/cmd/goimports
	@go get github.com/golang/lint/golint
	@go get -u github.com/golang/dep/cmd/dep

vet:
	@echo "+ Vet"
	@go vet ./...

lint:
	@echo "+ Linting package"
	@golint .
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
