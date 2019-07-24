package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/keegancsmith/sqlf"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/exitus/pkg/api"
	"github.com/wolfeidau/exitus/pkg/conf"
	"github.com/wolfeidau/exitus/pkg/db"
)

// ErrProjectNameAlreadyExists project name is already taken
var ErrProjectNameAlreadyExists = errors.New("project name is already taken")

// ProjectNotFoundError occurs when an project is not found.
type ProjectNotFoundError struct {
	Message string
}

func (e *ProjectNotFoundError) Error() string {
	return fmt.Sprintf("project not found: %s", e.Message)
}

// Projects provides a projects store
type Projects interface {
	GetByID(ctx context.Context, projectId string) (*api.Project, error)
	Create(ctx context.Context, newProj *api.NewProject, customerId string) (*api.Project, error)
	List(ctx context.Context, opt *ProjectsListOptions) ([]api.Project, error)
}

// ProjectsPG provides a projects store for postgresql
type ProjectsPG struct {
	dbconn *sql.DB
	cfg    *conf.Config
}

// NewProjects new project store
func NewProjects(dbconn *sql.DB, cfg *conf.Config) Projects {
	return &ProjectsPG{dbconn: dbconn, cfg: cfg}
}

// GetByID get project by id
func (ps *ProjectsPG) GetByID(ctx context.Context, projectId string) (*api.Project, error) {
	projs, err := ps.getBySQL(ctx, "WHERE id=$1 LIMIT 1", projectId)
	if err != nil {
		return nil, err
	}

	if len(projs) == 0 {
		return nil, &ProjectNotFoundError{fmt.Sprintf("id %s", projectId)}
	}
	return &projs[0], nil
}

// Create create a project
func (ps *ProjectsPG) Create(ctx context.Context, newProj *api.NewProject, customerId string) (*api.Project, error) {

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

// ProjectsListOptions specifies the options for listing projects.
type ProjectsListOptions struct {
	*NameLikeOptions
	*LimitOffset
}

// NewProjectsListOptions create a new opts
func NewProjectsListOptions(query string, offset int, limit int) *ProjectsListOptions {
	return &ProjectsListOptions{
		NameLikeOptions: &NameLikeOptions{query},
		LimitOffset:     &LimitOffset{Limit: limit, Offset: offset},
	}
}

// List list all projects
func (ps *ProjectsPG) List(ctx context.Context, opt *ProjectsListOptions) ([]api.Project, error) {
	if opt == nil {
		opt = &ProjectsListOptions{}
	}

	conds := ListNameLikeSQL(opt.NameLikeOptions)

	q := sqlf.Sprintf("WHERE %s ORDER BY id ASC %s", sqlf.Join(conds, "AND"), opt.LimitOffset.SQL())

	log.Info().Str("q", q.Query(sqlf.PostgresBindVar)).Msg("Projects List getBySQL")

	return ps.getBySQL(ctx, q.Query(sqlf.PostgresBindVar), q.Args()...)
}

func (ps *ProjectsPG) getBySQL(ctx context.Context, query string, args ...interface{}) ([]api.Project, error) {
	rows, err := ps.dbconn.QueryContext(ctx, "SELECT id, name, description, labels, created_at, updated_at FROM projects "+query, args...)
	if err != nil {
		return nil, err
	}

	projs := []api.Project{}
	defer rows.Close()
	for rows.Next() {
		proj := api.Project{}
		err := rows.Scan(&proj.Id, &proj.Name, &proj.Description, pq.Array(&proj.Labels), &proj.CreatedAt, &proj.UpdatedAt)
		if err != nil {
			return nil, err
		}

		projs = append(projs, proj)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return projs, nil
}
