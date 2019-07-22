APPNAME  := golang-backend
export APPNAME

deps:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.17.1
.PHONY: deps

lint:
	golangci-lint run ./...
.PHONY: lint

generate:
	go generate ./migrations/
	go generate ./pkg/api/
.PHONY: generate

local: generate
	docker-compose up
.PHONY: local

watchexec:
	watchexec --restart --exts "go" --watch . "docker-compose restart backend"
.PHONY: watchexec
