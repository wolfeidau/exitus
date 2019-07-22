package store

import (
	"database/sql"

	"github.com/wolfeidau/exitus/pkg/conf"
)

// Stores one stop for stores
type Stores struct {
	Projects *Projects
}

// New create all the stores
func New(dbconn *sql.DB, cfg *conf.Config) (*Stores, error) {
	return &Stores{
		Projects: NewProjects(dbconn, cfg),
	}, nil
}
