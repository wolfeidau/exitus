package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/exitus/pkg/api"
	"github.com/wolfeidau/exitus/pkg/conf"
	"github.com/wolfeidau/exitus/pkg/db"
)

// ErrProjectNameAlreadyExists project name is already taken
var ErrProjectNameAlreadyExists = errors.New("project name is already taken")

// Projects provides a projects store
type Projects struct {
	dbconn *sql.DB
	cfg    *conf.Config
}

// NewProjects new project store
func NewProjects(dbconn *sql.DB, cfg *conf.Config) *Projects {
	return &Projects{dbconn: dbconn, cfg: cfg}
}

// Create create a project
func (ps *Projects) Create(ctx context.Context, newProj *api.NewProject, customerId string) (*api.Project, error) {

	resProj := &api.Project{
		Name:        newProj.Name,
		Description: newProj.Description,
		Labels:      newProj.Labels,
	}

	err := db.WithTransaction(ctx, ps.dbconn, func(tx db.Transaction) error {
		return tx.QueryRowContext(
			ctx,
			`INSERT INTO projects(customer_id, name, description, labels) VALUES($1, $2, $3, $4)
			RETURNING id, created_at, updated_at`,
			customerId, newProj.Name, newProj.Description, pq.Array(newProj.Labels),
		).Scan(&resProj.Id, &resProj.CreatedAt, &resProj.UpdatedAt)
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to create project")

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "projects_customer_id_name_key":
				return nil, ErrProjectNameAlreadyExists
			}
		}

		return nil, err
	}

	return resProj, nil
}
