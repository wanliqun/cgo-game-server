# See https://tech.davis-hansson.com/p/make/
SHELL := bash
.DELETE_ON_ERROR:
.SHELLFLAGS := -eu -o pipefail -c
#.DEFAULT_GOAL := all
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-print-directory
BIN := .tmp/bin

CXX_SOURCE_DIR = ./cgo/cpp
CXX_INCLUDE_DIR = ./submodules/name-generator/dasmig
CXX_OUTPUT_LIB_DIR = /usr/local/lib

# Set to use a different compiler. For example, `GO=go1.18rc1 make test`.
GO ?= go
ARGS ?=

.PHONY: help
help: ## Describe useful make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-15s %s\n", $$1, $$2}'

.PHONY: clean
clean: ## Delete intermediate build artifacts
	@# -X only removes untracked files, -d recurses into directories, -f actually removes files/dirs
	git clean -Xdf

.PHONY: test
test: generate cgo ## Run all unit tests
	$(GO) test -race -cover ./...

.PHONY: lint
lint: lint-proto lint-go  ## Lint code and protos

.PHONY: lint-go
lint-go: $(BIN)/golangci-lint
	$(BIN)/golangci-lint run --modules-download-mode=readonly --timeout=3m0s ./...

.PHONY: lint-go-fix
lint-go-fix: $(BIN)/golangci-lint
	$(BIN)/golangci-lint run --fix --modules-download-mode=readonly --timeout=3m0s ./...

.PHONY: lint-proto
lint-proto: $(BIN)/buf
	$(BIN)/buf lint
	$(BIN)/buf breaking --against '.git#branch=main'

.PHONY: generate
generate: generate-proto ## Generate protobuf Go codes

.PHONY: generate-proto
generate-proto: $(BIN)/buf
	rm -rf proto/*.pb.go
	$(BIN)/buf generate
	go mod tidy -v

.PHONY: checkgenerate
checkgenerate: generate
	@# Used in CI to verify that `make generate` doesn't produce a diff.
	test -z "$$(git status --porcelain | tee /dev/stderr)"

.PHONY: upgrade-go
upgrade-go:
	$(GO) get -u -t ./... && go mod tidy -v

$(BIN):
	@mkdir -p $(BIN)

$(BIN)/buf: $(BIN) Makefile
	GOBIN=$(abspath $(@D)) $(GO) install github.com/bufbuild/buf/cmd/buf@latest

$(BIN)/golangci-lint: $(BIN) Makefile
	GOBIN=$(abspath $(@D)) $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

.PHONY: cgo
cgo: ## Generate C++ dynamic link library
	clang++ -o $(CXX_OUTPUT_LIB_DIR)/libnamegen.so $(CXX_SOURCE_DIR)/lib-bridge.cpp \
		-I$(CXX_INCLUDE_DIR) \
		-std=c++17 -O3 -Wall -Wextra -fPIC -shared
