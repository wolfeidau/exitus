package store

import (
	"database/sql"

	"github.com/wolfeidau/exitus/pkg/conf"
)

type Stores struct {
	Projects *Projects
}

func New(dbconn *sql.DB, cfg *conf.Config) (*Stores, error) {
	return &Stores{
		Projects: NewProjects(dbconn, cfg),
	}, nil
}
