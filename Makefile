APPNAME  := golang-backend
export APPNAME

LDFLAGS := -ldflags="-s -w"
DEPLOY := infra/deploy

deps:
	mkdir -p $(shell go env GOPATH)/bin
	if [ ! -f "$(shell go env GOPATH)/bin/golangci-lint" ]; then curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.17.1; fi
	if [ ! -f "$(shell go env GOPATH)/bin/migrate.linux-amd64" ]; then curl -sfL https://github.com/golang-migrate/migrate/releases/download/v4.5.0/migrate.linux-amd64.tar.gz | tar xvz -C $(shell go env GOPATH)/bin; fi
.PHONY: deps

lint:
	golangci-lint run ./...
.PHONY: lint

generate:
	go generate ./migrations/
	go generate ./pkg/api/
.PHONY: generate

docker-build:
	@echo "--- build all the things"
	@go mod download
	@docker run --rm \
		-v $$(pwd):/src/$$(basename $$(pwd)) \
		-v $$(go env GOPATH)/pkg/mod:/go/pkg/mod \
		-w /src/$$(basename $$(pwd)) -it golang make linux
.PHONY: build-docker

linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o infra/deploy/exitus ./cmd/backend
.PHONY: linux

test: deps
	$(shell go env GOPATH)/bin/migrate.linux-amd64 -database "postgresql://testing@postgres/testing?sslmode=disable&password=${POSTGRES_PASSWORD}" -path ./migrations up
	go test -cover -v ./...
.PHONY: test

docker-compose-test:
	docker-compose -f dev/docker-compose.yml up  --exit-code-from exitus_testing
.PHONY: docker-compose-test

local: generate
	docker-compose up
.PHONY: local

watchexec:
	watchexec --restart --exts "go" --watch . "docker-compose restart backend"
.PHONY: watchexec
