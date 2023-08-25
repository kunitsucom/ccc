SHELL     := /usr/bin/env bash -Eeu -o pipefail
GITROOT   := $(shell git rev-parse --show-toplevel || pwd || echo '.')
PRE_PUSH  := ${GITROOT}/.git/hooks/pre-push
GOMODULE  := github.com/kunitsucom/ccc
VERSION   := $(shell git describe --tags --abbrev=0 --always)
REVISION  := $(shell git log -1 --format='%H')
BRANCH    := $(shell git rev-parse --abbrev-ref HEAD)
TIMESTAMP := $(shell git log -1 --format='%cI')

.DEFAULT_GOAL := help
.PHONY: help
help: githooks ## display this help documents
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'

.PHONY: setup
setup: githooks ## Setup tools for development
	# direnv
	./.bin/direnv allow
	# golangci-lint
	./.bin/golangci-lint --version

.PHONY: githooks
githooks:
	@[[ -f "${PRE_PUSH}" ]] || cp -ai "${GITROOT}/.githooks/pre-push" "${PRE_PUSH}"

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

.PHONY: goxz
goxz: ci ## Run goxz for release
	command -v goxz || go install github.com/Songmu/goxz/cmd/goxz@latest
	goxz -d ./.tmp -os=linux,darwin,windows -arch=amd64,arm64 -pv ${VERSION} -build-ldflags "-X ${GOMODULE}/pkg/config.version=${VERSION} -X ${GOMODULE}/pkg/config.revision=${REVISION} -X ${GOMODULE}/pkg/config.branch=${BRANCH} -X ${GOMODULE}/pkg/config.timestamp=${TIMESTAMP}" ./cmd/ccc
