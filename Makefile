SHELL              := /usr/bin/env bash -Eeu -o pipefail
REPO_ROOT          := $(shell git rev-parse --show-toplevel || pwd || echo '.')
REPO_LOCAL_DIR     := ${REPO_ROOT}/.local
REPO_TMP_DIR       := ${REPO_ROOT}/.tmp
PRE_PUSH           := ${REPO_ROOT}/.git/hooks/pre-push
GIT_TAG_LATEST     := $(shell git describe --tags --abbrev=0)
GIT_BRANCH_CURRENT := $(shell git rev-parse --abbrev-ref HEAD)
GO_MODULE_NAME     := github.com/kunitsucom/ccc
BUILD_VERSION      := $(shell git describe --tags --exact-match HEAD 2>/dev/null || git rev-parse --short HEAD)
BUILD_REVISION     := $(shell git rev-parse HEAD)
BUILD_BRANCH       := $(shell git rev-parse --abbrev-ref HEAD | tr / -)
BUILD_TIMESTAMP    := $(shell git log -n 1 --format='%cI')

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

.PHONY: release
release: ci ## Run goxz and gh release upload
	@command -v goxz >/dev/null || go install github.com/Songmu/goxz/cmd/goxz@latest
	git checkout main
	git checkout "${GIT_TAG_LATEST}"
	-goxz -d "${REPO_TMP_DIR}" -os=linux,darwin,windows -arch=amd64,arm64 -pv "`git describe --tags --abbrev=0`" -trimpath -build-ldflags "-s -w -X ${GO_MODULE_NAME}/pkg/config.version=`git describe --tags --abbrev=0` -X ${GO_MODULE_NAME}/pkg/config.revision=`git rev-parse HEAD` -X ${GO_MODULE_NAME}/pkg/config.branch=`git rev-parse --abbrev-ref HEAD` -X ${GO_MODULE_NAME}/pkg/config.timestamp=`git log -n 1 --format='%cI'`" ./cmd/ccc
	-gh release upload "`git describe --tags --abbrev=0`" "${REPO_TMP_DIR}"/*"`git describe --tags --abbrev=0`"*
	git checkout "${GIT_BRANCH_CURRENT}"
