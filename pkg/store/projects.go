package store

import (
	"context"
	"database/sql"
	"errors"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/exitus/pkg/api"
	"github.com/wolfeidau/exitus/pkg/conf"
)

// ErrProjectNameAlreadyExists project name is already taken
var ErrProjectNameAlreadyExists = errors.New("project name is already taken")

type Projects struct {
	dbconn *sql.DB
	cfg    *conf.Config
}

func NewProjects(dbconn *sql.DB, cfg *conf.Config) *Projects {
	return &Projects{dbconn: dbconn, cfg: cfg}
}

func (ps *Projects) Create(ctx context.Context, newProj *api.NewProject, ownerId string) (*api.Project, error) {

	resProj := &api.Project{
		Name:        newProj.Name,
		Description: newProj.Description,
		Labels:      newProj.Labels,
	}

	tx, err := ps.dbconn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			rollErr := tx.Rollback()
			if rollErr != nil {
				err = multierror.Append(err, rollErr)
			}
			return
		}
		err = tx.Commit()
	}()

	err = tx.QueryRowContext(
		ctx,
		"INSERT INTO projects(name, description, owner_id, labels) VALUES($1, $2, $3, $4) RETURNING id, created_at, updated_at",
		newProj.Name, newProj.Description, ownerId, pq.Array(newProj.Labels),
	).Scan(&resProj.Id, &resProj.CreatedAt, &resProj.UpdatedAt)
	if err != nil {
		log.Error().Err(err).Msg("failed to create project")

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "projects_name_key":
				return nil, ErrProjectNameAlreadyExists
			}
		}

		return nil, err
	}

	return resProj, nil
}
