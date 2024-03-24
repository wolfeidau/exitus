package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/keegancsmith/sqlf"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/wolfeidau/exitus/pkg/api"
	"github.com/wolfeidau/exitus/pkg/conf"
	"github.com/wolfeidau/exitus/pkg/db"
)

// ErrProjectNameAlreadyExists project name is already taken.
var ErrProjectNameAlreadyExists = errors.New("project name is already taken")

// ProjectNotFoundError occurs when an project is not found.
type ProjectNotFoundError struct {
	Message string
}

func (e *ProjectNotFoundError) Error() string {
	return fmt.Sprintf("project not found: %s", e.Message)
}

// Projects provides a projects store.
type Projects interface {
	GetByID(ctx context.Context, id string, customerId string) (*api.Project, error)
	Create(ctx context.Context, newProj *api.NewProject, customerId string) (*api.Project, error)
	Update(ctx context.Context, updatedProject *api.UpdatedProject, id string, customerId string) (*api.Project, error)
	List(ctx context.Context, opt *ProjectsListOptions, customerId string) ([]api.Project, error)
}

// ProjectsListOptions specifies the options for listing projects.
type ProjectsListOptions struct {
	*NameLikeOptions
	*LimitOffset
}

// NewProjectsListOptions create a new opts.
func NewProjectsListOptions(query string, offset int, limit int) *ProjectsListOptions {
	return &ProjectsListOptions{
		NameLikeOptions: &NameLikeOptions{query},
		LimitOffset:     &LimitOffset{Limit: limit, Offset: offset},
	}
}

// ProjectsPG provides a projects store for postgresql.
type ProjectsPG struct {
	dbconn *sql.DB
	cfg    *conf.Config
}

// NewProjects new project store.
func NewProjects(dbconn *sql.DB, cfg *conf.Config) Projects {
	return &ProjectsPG{dbconn: dbconn, cfg: cfg}
}

// GetByID get project by id.
func (ps *ProjectsPG) GetByID(ctx context.Context, id string, customerId string) (*api.Project, error) {
	projs, err := ps.getBySQL(ctx, "WHERE id=$1 AND customer_id=$2 LIMIT 1", id, customerId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get project by id: %s customerId: %s", id, customerId)
	}

	if len(projs) == 0 {
		return nil, &ProjectNotFoundError{fmt.Sprintf("id %s", id)}
	}
	return &projs[0], nil
}

// Update update a project.
func (ps *ProjectsPG) Update(ctx context.Context, updatedProject *api.UpdatedProject, id string, customerId string) (*api.Project, error) {
	fields := []*sqlf.Query{sqlf.Sprintf("name=%s, labels=%s, updated_at=%s", updatedProject.Name, pq.Array(updatedProject.Labels), time.Now())}

	if updatedProject.Description != nil {
		fields = append(fields, sqlf.Sprintf("description=%s", updatedProject.Description))
	}

	qry := sqlf.Sprintf("UPDATE projects SET %s WHERE id=%s AND customer_id=%s", sqlf.Join(fields, ","), id, customerId)

	if _, err := ps.dbconn.ExecContext(ctx, qry.Query(sqlf.PostgresBindVar), qry.Args()...); err != nil {
		return nil, errors.Wrapf(err, "failed to update project by id: %s customerId: %s", id, customerId)
	}

	return ps.GetByID(ctx, id, customerId)
}

// Create create a project.
func (ps *ProjectsPG) Create(ctx context.Context, newProj *api.NewProject, customerId string) (*api.Project, error) {
	resProj := api.Project{}

	qry := sqlf.Sprintf("INSERT INTO projects(customer_id, name, description, labels) VALUES(%s, %s, %s, %s)",
		customerId, newProj.Name, newProj.Description, pq.Array(newProj.Labels))

	err := db.WithTransaction(ctx, ps.dbconn, func(tx db.Transaction) error {
		return tx.QueryRowContext(
			ctx, qry.Query(sqlf.PostgresBindVar)+" RETURNING id, name, description, labels, created_at, updated_at", qry.Args()...,
		).Scan(&resProj.Id, &resProj.Name, &resProj.Description, pq.Array(&resProj.Labels), &resProj.CreatedAt, &resProj.UpdatedAt)
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "projects_customer_id_name_key":
				return nil, ErrProjectNameAlreadyExists
			}
		}
		return nil, errors.Wrapf(err, "failed to create project with name: %s customerId: %s", newProj.Name, customerId)
	}

	return &resProj, nil
}

// List list all projects.
func (ps *ProjectsPG) List(ctx context.Context, opt *ProjectsListOptions, customerId string) ([]api.Project, error) {
	if opt == nil {
		opt = &ProjectsListOptions{}
	}

	conds := ListNameLikeSQL(opt.NameLikeOptions)
	conds = append(conds, sqlf.Sprintf("customer_id = %s", customerId))

	qry := sqlf.Sprintf("WHERE %s ORDER BY id ASC %s", sqlf.Join(conds, "AND"), opt.LimitOffset.SQL())

	return ps.getBySQL(ctx, qry.Query(sqlf.PostgresBindVar), qry.Args()...)
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
