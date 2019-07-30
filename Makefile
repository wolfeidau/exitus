APPNAME  := golang-backend
export APPNAME

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

test: deps
	$(shell go env GOPATH)/bin/migrate.linux-amd64 -database "postgresql://testing:${POSTGRES_PASSWORD}@postgres/testing?sslmode=disable" -path ./migrations up
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
