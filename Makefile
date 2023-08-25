SHELL           := /usr/bin/env bash -Eeu -o pipefail
REPO_ROOT       := $(shell git rev-parse --show-toplevel || pwd || echo '.')
REPO_LOCAL_DIR  := ${REPO_ROOT}/.local
PRE_PUSH        := ${REPO_ROOT}/.git/hooks/pre-push
GO_MODULE_NAME  := github.com/kunitsucom/ccc
BUILD_VERSION   := $(shell git describe --tags --exact-match HEAD 2>/dev/null || git rev-parse --short HEAD)
BUILD_REVISION  := $(shell git rev-parse HEAD)
BUILD_BRANCH    := $(shell git rev-parse --abbrev-ref HEAD | tr / -)
BUILD_TIMESTAMP := $(shell git log -n 1 --format='%cI')

export PATH := ${REPO_LOCAL_DIR}/bin:${REPO_ROOT}/.bin:${PATH}

.DEFAULT_GOAL := help
.PHONY: help
help: githooks ## display this help documents
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'

.PHONY: setup
setup: githooks ## Setup tools for development
	# == SETUP =====================================================
	# versenv
	make versenv
	# --------------------------------------------------------------

.PHONY: versenv
versenv:
	# direnv
	direnv allow .
	# golangci-lint
	golangci-lint --version

.PHONY: githooks
githooks:
	@[[ -f "${PRE_PUSH}" ]] || cp -ai "${REPO_ROOT}/.githooks/pre-push" "${PRE_PUSH}"

clean:  ## Clean up chace, etc
	go clean -x -cache -testcache -modcache -fuzzcache
	golangci-lint cache clean

.PHONY: lint
lint:  ## Run secretlint, go mod tidy, golangci-lint
	# ref. https://github.com/secretlint/secretlint
	docker run -v `pwd`:`pwd` -w `pwd` --rm secretlint/secretlint secretlint "**/*"
	# tidy
	go mod tidy
	git diff --exit-code go.mod go.sum
	# lint
	# ref. https://golangci-lint.run/usage/linters/
	./.bin/golangci-lint run --fix --sort-results
	git diff --exit-code

.PHONY: credits
credits:  ## Generate CREDITS file
	command -v gocredits || go install github.com/Songmu/gocredits/cmd/gocredits@latest
	gocredits -skip-missing . > CREDITS
	git diff --exit-code

.PHONY: test
test: githooks ## Run go test and display coverage
	# test
	go test -v -race -p=4 -parallel=8 -timeout=300s -cover -coverprofile=./coverage.txt ./...
	go tool cover -func=./coverage.txt

.PHONY: ci
ci: lint credits test ## CI command set

.PHONY: act-check
act-check:
	@if ! command -v act >/dev/null 2>&1; then \
		printf "\033[31;1m%s\033[0m\n" "act is not installed: brew install act" 1>&2; \
		exit 1; \
	fi

.PHONY: act-go-lint
act-go-lint: act-check
	act pull_request --container-architecture linux/amd64 -P ubuntu-latest=catthehacker/ubuntu:act-latest -W .github/workflows/go-lint.yml

.PHONY: act-go-test
act-go-test: act-check
	act pull_request --container-architecture linux/amd64 -P ubuntu-latest=catthehacker/ubuntu:act-latest -W .github/workflows/go-test.yml

.PHONY: goxz
goxz: ci ## Run goxz for release
	command -v goxz || go install github.com/Songmu/goxz/cmd/goxz@latest
	goxz -d ./.tmp -os=linux,darwin,windows -arch=amd64,arm64 -pv ${BUILD_VERSION} -build-ldflags "-X ${GO_MODULE_NAME}/pkg/config.version=${BUILD_VERSION} -X ${GO_MODULE_NAME}/pkg/config.revision=${BUILD_REVISION} -X ${GO_MODULE_NAME}/pkg/config.branch=${BUILD_BRANCH} -X ${GO_MODULE_NAME}/pkg/config.timestamp=${BUILD_TIMESTAMP}" ./cmd/ccc
