package store_test

import (
	"os"
	"testing"

	"github.com/wolfeidau/exitus/pkg/db"
)

func TestMigrations_Down_Up(t *testing.T) {
	if os.Getenv("SKIP_MIGRATION_TEST") != "" {
		t.Skip()
	}

	// get testing context to ensure we can connect to the DB
	_ = db.TestContext(t)

	m := db.NewMigrate(db.Global)
	// Run all down migrations then up migrations again to ensure there are no SQL errors.
	if err := m.Down(); err != nil {
		t.Errorf("error running down migrations: %s", err)
	}
	if err := db.DoMigrate(m); err != nil {
		t.Errorf("error running up migrations: %s", err)
	}
}
