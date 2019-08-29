// Package api contains the REST interfaces.
package api

//go:generate env GOBIN=$PWD/.bin GO111MODULE=on go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen
//go:generate $PWD/.bin/oapi-codegen -generate client,types,server,spec -package api -o exitus.gen.go exitus.yml
//go:generate gofumpt -s -w exitus.gen.go
