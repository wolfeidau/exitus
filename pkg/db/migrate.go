package db

import (
	"database/sql"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/exitus/migrations"
)

// NewMigrate load migrations.
func NewMigrate(db *sql.DB) *migrate.Migrate {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load driver")
	}

	d, err := iofs.New(migrations.MigrationsFs, ".")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read assets from iofs")
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create migration from go-bindata")
	}
	// In case another process was faster and runs migrations, we will wait
	// this long
	m.LockTimeout = 5 * time.Minute

	return m
}

// DoMigrate do sql migrations.
func DoMigrate(m *migrate.Migrate) (err error) {
	err = m.Up()
	if err == nil || err == migrate.ErrNoChange {
		return nil
	}

	if os.IsNotExist(err) {
		// This should only happen if the DB is ahead of the migrations available
		version, dirty, verr := m.Version()
		if verr != nil {
			return verr
		}
		if dirty { // this shouldn't happen, but checking anyways
			return err
		}
		log.Warn().Uint("db_version", version).Msg("WARNING: Detected an old version of database.")
		return nil
	}
	return err
}
