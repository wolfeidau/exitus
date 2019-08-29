package conf

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// ErrMissingEnvironmentStage missing stage configuration
	ErrMissingEnvironmentStage = errors.New("Missing Stage ENV Variable")

	// ErrMissingEnvironmentBranch missing branch configuration
	ErrMissingEnvironmentBranch = errors.New("Missing Branch ENV Variable")
)

// Config for the environment
type Config struct {
	Debug                bool   `envconfig:"DEBUG"`
	Addr                 string `envconfig:"ADDR" default:":8080"`
	Stage                string `envconfig:"STAGE" default:"dev"`
	Branch               string `envconfig:"BRANCH"`
	PGDatasource         string `envconfig:"PGDATASOURCE"`
	OpenIDProvider       string `envconfig:"OPENID_PROVIDER_URL"`
	ClientID             string `envconfig:"OAUTH_CLIENT_ID"`
	MetricsWriteInterval int    `envconfig:"METRICS_WRITE_INTERVAL"`
	DbSecrets            string `envconfig:"DB_SECRET"`
}

type DBSecrets struct {
	Password string `json:"password,omitempty"`
	DBName   string `json:"dbname,omitempty"`
	Engine   string `json:"engine,omitempty"`
	Port     int    `json:"port,omitempty"`
	Host     string `json:"host,omitempty"`
	Username string `json:"username,omitempty"`
}

func (cfg *Config) validate() error {
	if cfg.Stage == "" {
		return ErrMissingEnvironmentStage
	}
	if cfg.Branch == "" {
		return ErrMissingEnvironmentBranch
	}

	return nil
}

// "postgresql://testing@postgres/testing?sslmode=disable&password=Tig%23fD%5BXED%2C)S%3AG%3B%3C.ruAm9"

// this extracts the JSON secret value from the environment and builds the
// datasource used to connect to postgresql
func (cfg *Config) parseDbSecrets() error {

	if cfg.DbSecrets == "" {
		log.Info().Msg("no DB secrets value provided")
		return nil
	}

	dbsecrets := &DBSecrets{}

	err := json.Unmarshal([]byte(cfg.DbSecrets), dbsecrets)
	if err != nil {
		return err
	}

	cfg.PGDatasource = fmt.Sprintf("postgres://%s@%s:%d/%s?password=%s", dbsecrets.Username, dbsecrets.Host, dbsecrets.Port, dbsecrets.DBName, url.QueryEscape(dbsecrets.Password))

	log.Debug().Str("PGDatasource", cfg.PGDatasource).Msg("configured PG datasource")

	return nil
}

func (cfg *Config) logging() error {

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if cfg.Stage == "local" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return nil
}

// NewDefaultConfig reads configuration from environment variables and validates it
func NewDefaultConfig() (*Config, error) {
	cfg := new(Config)
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse environment config")
	}
	err = cfg.parseDbSecrets()
	if err != nil {
		return nil, errors.Wrap(err, "failed parse db secret")
	}
	err = cfg.validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed validation of config")
	}
	err = cfg.logging()
	if err != nil {
		return nil, errors.Wrap(err, "failed setup logging based on config")
	}
	log.Info().Str("stage", cfg.Stage).Bool("debug", cfg.Debug).Msg("logging configured")
	log.Info().Str("stage", cfg.Stage).Str("branch", cfg.Branch).Msg("Configuration loaded")

	return cfg, nil
}
