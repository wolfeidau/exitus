package dbconn

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gchaincl/sqlhooks"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/golang-backend-postgres/pkg/env"
)

// DefaultMaxOpenConnections the default value for max open connections in the
// PostgreSQL connection pool
const DefaultMaxOpenConnections = 30

var (
	// Global is the global DB connection.
	// Only use this after a call to ConnectToDB.
	Global *sql.DB

	defaultDataSource = env.Get("PGDATASOURCE", "")

	registerOnce sync.Once
)

// ConnectToDB connects to the given DB and stores the handle globally.
//
// Note: github.com/lib/pq parses the environment as well. This function will
// also use the value of PGDATASOURCE if supplied and dataSource is the empty
// string.
func ConnectToDB(dataSource string) error {
	if dataSource == "" {
		dataSource = defaultDataSource
	}

	// Force PostgreSQL session timezone to UTC.
	if v, ok := os.LookupEnv("PGTZ"); ok && v != "UTC" && v != "utc" {
		log.Warn().Str("ignoredPGTZ", v).Msg("Ignoring PGTZ environment variable; using PGTZ=UTC.")
	}
	if err := os.Setenv("PGTZ", "UTC"); err != nil {
		return errors.Wrap(err, "Error setting PGTZ=UTC")
	}

	var err error
	Global, err = openDBWithStartupWait(dataSource)
	if err != nil {
		return errors.Wrap(err, "DB not available")
	}
	configureConnectionPool(Global)

	if err := DoMigrate(NewMigrate(Global)); err != nil {
		return errors.Wrap(err, "Failed to migrate the DB.")
	}

	return nil
}

var startupTimeout = func() time.Duration {
	str := env.Get("DB_STARTUP_TIMEOUT", "10s")
	d, err := time.ParseDuration(str)
	if err != nil {
		log.Fatal().Err(err).Msg("db startup timed out")
	}
	return d
}()

func openDBWithStartupWait(dataSource string) (db *sql.DB, err error) {
	// Allow the DB to take up to 10s while it reports "pq: the database system is starting up".
	startupDeadline := time.Now().Add(startupTimeout)
	for {
		if time.Now().After(startupDeadline) {
			return nil, fmt.Errorf("database did not start up within %s (%v)", startupTimeout, err)
		}
		db, err = Open(dataSource)
		if err == nil {
			err = db.Ping()
		}
		if err != nil && isDatabaseLikelyStartingUp(err) {
			time.Sleep(startupTimeout / 10)
			continue
		}
		return db, err
	}
}

// isDatabaseLikelyStartingUp returns whether the err likely just means the PostgreSQL database is
// starting up, and it should not be treated as a fatal error during program initialization.
func isDatabaseLikelyStartingUp(err error) bool {
	if strings.Contains(err.Error(), "pq: the database system is starting up") {
		// Wait for DB to start up.
		return true
	}
	if e, ok := errors.Cause(err).(net.Error); ok && strings.Contains(e.Error(), "connection refused") {
		// Wait for DB to start listening.
		return true
	}
	return false
}

// Open creates a new DB handle with the given schema by connecting to
// the database identified by dataSource (e.g., "dbname=mypgdb" or
// blank to use the PG* env vars).
//
// Open assumes that the database already exists.
func Open(dataSource string) (*sql.DB, error) {
	registerOnce.Do(func() {
		sql.Register("postgres-proxy", sqlhooks.Wrap(&pq.Driver{}, &hook{}))
	})
	db, err := sql.Open("postgres-proxy", dataSource)
	if err != nil {
		return nil, errors.Wrap(err, "postgresql open")
	}
	return db, nil
}

// Ping attempts to contact the database and returns a non-nil error upon failure. It is intended to
// be used by health checks.
func Ping(ctx context.Context) error { return Global.PingContext(ctx) }

// configureConnectionPool sets reasonable sizes on the built in DB queue. By
// default the connection pool is unbounded, which leads to the error `pq:
// sorry too many clients already`.
func configureConnectionPool(db *sql.DB) {
	var err error
	maxOpen := DefaultMaxOpenConnections
	if e := os.Getenv("SRC_PGSQL_MAX_OPEN"); e != "" {
		maxOpen, err = strconv.Atoi(e)
		if err != nil {
			log.Fatal().Err(err).Msg("SRC_PGSQL_MAX_OPEN is not an int")
		}
	}
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxOpen)
	db.SetConnMaxLifetime(time.Minute)
}
