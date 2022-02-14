//go:build go1.16
// +build go1.16

// Package migrations contains the migration scripts for the DB.
package migrations

import "embed"

//go:embed *.sql
var MigrationsFs embed.FS
