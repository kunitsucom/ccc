SHELL    := /usr/bin/env bash -Eeu -o pipefail
GITROOT  := $(shell git rev-parse --show-toplevel || pwd || echo '.')
PRE_PUSH := ${GITROOT}/.git/hooks/pre-push

.DEFAULT_GOAL := help
.PHONY: help
help: githooks ## display this help documents.
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'

.PHONY: githooks
githooks:  ## githooks をインストールします。
	@[[ -f "${PRE_PUSH}" ]] || cp -ai "${GITROOT}/.githooks/pre-push" "${PRE_PUSH}"

.PHONY: lint
lint:  ## go mod tidy の後に golangci-lint を実行します。
	# tidy
	go mod tidy
	git diff --exit-code go.mod go.sum
	# lint
	# cf. https://golangci-lint.run/usage/linters/
	./.bin/golangci-lint run --fix --sort-results
	git diff --exit-code

.PHONY: test
test: githooks ## go test を実行し coverage を出力します。
	# test
	go test -v -race -p=4 -parallel=8 -timeout=300s -cover -coverprofile=./coverage.txt ./...
	go tool cover -func=./coverage.txt

.PHONY: ci
ci: lint test ## CI 上で実行する lint や test のコマンドセットです。

.PHONY: credits
credits:  ## CREDITS ファイルを生成します。
	command -v gocredits || go install github.com/Songmu/gocredits/cmd/gocredits@latest
	gocredits . > CREDITS
