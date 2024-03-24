package main

import (
	"io/ioutil"

	"github.com/labstack/echo/v4"
	echolog "github.com/labstack/gommon/log"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/exitus/pkg/api"
	"github.com/wolfeidau/exitus/pkg/conf"
	"github.com/wolfeidau/exitus/pkg/db"
	"github.com/wolfeidau/exitus/pkg/healthz"
	"github.com/wolfeidau/exitus/pkg/metrics"
	"github.com/wolfeidau/exitus/pkg/middleware"
	"github.com/wolfeidau/exitus/pkg/server"
	"github.com/wolfeidau/exitus/pkg/store"
)

func main() {
	// loads configuration from env and configures logger
	cfg, err := conf.NewDefaultConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	// configure metrics writer
	mr := metrics.New(cfg)
	go mr.Start()

	dbconn, err := db.NewDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
	}

	dbm := metrics.NewDBMonitor(dbconn)
	go dbm.Start()

	stores, err := store.New(dbconn, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
	}

	svr, err := server.NewServer(cfg, stores)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to bind api")
	}

	e := echo.New()
	// shut up
	e.Logger.SetOutput(ioutil.Discard)
	e.Logger.SetLevel(echolog.OFF)

	e.GET("/healthz", echo.WrapHandler(healthz.Handler()))

	// add a version to the api
	g := e.Group("/v1")

	g.Use(middleware.RequestID)
	g.Use(middleware.ErrorLog)
	g.Use(middleware.RequestLog)
	g.Use(middleware.JWTWithConfig(&middleware.JWTConfig{
		ProviderURL: cfg.OpenIDProvider,
		ClientID:    cfg.ClientID,
	}))

	api.RegisterHandlers(g, svr)

	log.Info().Str("addr", cfg.Addr).Msg("starting http listener")
	err = e.Start(cfg.Addr)
	log.Fatal().Err(err).Msg("Server failed")
}
