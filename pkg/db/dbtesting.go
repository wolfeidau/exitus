package db

import (
	"context"
	"database/sql"
	"strings"
	"sync"
	"testing"

	"github.com/wolfeidau/exitus/pkg/conf"
)

var (
	Global      *sql.DB
	connectOnce sync.Once
)

func TestContext(t testing.TB) context.Context {
	connectOnce.Do(func() {
		// loads configuration from env and configures logger
		cfg, err := conf.NewDefaultConfig()
		if err != nil {
			t.Fatal("failed to load config")
		}

		Global, err = NewDB(cfg)
		if err != nil {
			t.Fatal("failed to load config")
		}
	})

	emptyDBPreserveSchema(t, Global)

	return context.TODO()
}

func emptyDBPreserveSchema(t testing.TB, d *sql.DB) {
	_, err := d.Exec(`SELECT * FROM schema_migrations`)
	if err != nil {
		t.Fatalf("Table schema_migrations not found: %v", err)
	}

	rows, err := d.Query("SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE' AND table_name != 'schema_migrations'")
	if err != nil {
		t.Fatal(err)
	}
	var tables []string
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			t.Fatal(err)
		}
		tables = append(tables, table)
	}
	if err = rows.Close(); err != nil {
		t.Fatal(err)
	}
	if err = rows.Err(); err != nil {
		t.Fatal(err)
	}
	if testing.Verbose() {
		t.Logf("Truncating all %d tables", len(tables))
	}
	_, err = d.Exec("TRUNCATE " + strings.Join(tables, ", ") + " RESTART IDENTITY")
	if err != nil {
		t.Fatal(err)
	}
}
