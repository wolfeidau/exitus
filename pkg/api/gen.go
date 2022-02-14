// Package api contains the REST interfaces.
package api

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -generate client,types,server,spec -package api -o exitus.gen.go exitus.yml
//go:generate gofmt -s -w exitus.gen.go
