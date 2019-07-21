package main

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/golang-backend-postgres/pkg/conf"
	"github.com/wolfeidau/golang-backend-postgres/pkg/db/dbconn"
)

func main() {

	// loads configuration from env and configures logger
	cfg, err := conf.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	err = dbconn.ConnectToDB("")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Docker")
	})

	log.Info().Str("addr", cfg.Addr).Msg("starting http listener")
	err = http.ListenAndServe(cfg.Addr, nil)
	log.Fatal().Err(err).Msg("Server failed")
}
