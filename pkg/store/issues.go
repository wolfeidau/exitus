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

// IssueNotFoundError occurs when an issue is not found.
type IssueNotFoundError struct {
	Message string
}

func (e *IssueNotFoundError) Error() string {
	return fmt.Sprintf("issue not found: %s", e.Message)
}

// Issues provides a issues store.
type Issues interface {
	GetByID(ctx context.Context, id, projectId, customerId string) (*api.Issue, error)
	Create(ctx context.Context, newProj *api.NewIssue, projectId, customerId, reporter string) (*api.Issue, error)
	Update(ctx context.Context, updatedIssue *api.UpdatedIssue, id, projectId, customerId string) (*api.Issue, error)
	List(ctx context.Context, opt *IssueListOptions, projectId, customerId string) ([]api.Issue, error)
}

// IssueListOptions specifies the options for listing issues.
type IssueListOptions struct {
	*SubjectLikeOptions
	*LimitOffset
}

// NewIssueListOptions create a new opts.
func NewIssueListOptions(query string, offset int, limit int) *IssueListOptions {
	return &IssueListOptions{
		SubjectLikeOptions: &SubjectLikeOptions{query},
		LimitOffset:        &LimitOffset{Limit: limit, Offset: offset},
	}
}

// SubjectLikeOptions used to query by subject using like.
type SubjectLikeOptions struct {
	// Query specifies a search query for organizations.
	Query string
}

// ListSubjectLikeSQL used to search by subject if query is set.
func ListSubjectLikeSQL(opt *SubjectLikeOptions) (conds []*sqlf.Query) {
	conds = []*sqlf.Query{sqlf.Sprintf("TRUE")}
	if opt.Query != "" {
		query := "%" + opt.Query + "%"
		conds = append(conds, sqlf.Sprintf("subject ILIKE %s", query))
	}
	return conds
}

// IssuesPG provides a issues store for postgresql.
type IssuesPG struct {
	dbconn *sql.DB
	cfg    *conf.Config
}

// NewIssues new issues store.
func NewIssues(dbconn *sql.DB, cfg *conf.Config) Issues {
	return &IssuesPG{dbconn: dbconn, cfg: cfg}
}

// GetByID get issue by id.
func (is *IssuesPG) GetByID(ctx context.Context, id, projectId, customerId string) (*api.Issue, error) {
	issues, err := is.getBySQL(ctx, "WHERE id=$1 AND project_id=$2 AND customer_id=$3 LIMIT 1", id, projectId, customerId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get issue by id: %s projectId: %s customerId: %s", id, projectId, customerId)
	}

	if len(issues) == 0 {
		return nil, &IssueNotFoundError{fmt.Sprintf("id %s project_id %s", id, projectId)}
	}

	return &issues[0], nil
}

// Create create new issue.
func (is *IssuesPG) Create(ctx context.Context, newIssue *api.NewIssue, projectId, customerId, reporter string) (*api.Issue, error) {
	issue := api.Issue{}

	qry := sqlf.Sprintf("INSERT INTO issues(project_id, customer_id, reporter, subject, state, severity, category, labels, content) VALUES(%s, %s, %s, %s, %s, %s, %s, %s, %s)",
		projectId, customerId, reporter, newIssue.Subject, "created", newIssue.Severity, newIssue.Category, pq.Array(newIssue.Labels), newIssue.Content)

	err := db.WithTransaction(ctx, is.dbconn, func(tx db.Transaction) error {
		return tx.QueryRowContext(
			ctx, qry.Query(sqlf.PostgresBindVar)+" RETURNING id, subject, state, severity, category, labels, content, created_at, updated_at", qry.Args()...,
		).Scan(&issue.Id, &issue.Subject, &issue.State, &issue.Severity, &issue.Category, pq.Array(&issue.Labels), &issue.Content, &issue.CreatedAt, &issue.UpdatedAt)
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create issue with subject: %s, customer_id: %s", newIssue.Subject, customerId)
	}

	return &issue, nil
}

func (is *IssuesPG) Update(ctx context.Context, updatedIssue *api.UpdatedIssue, id, projectId, customerId string) (*api.Issue, error) {
	fields := []*sqlf.Query{sqlf.Sprintf("subject=%s, content=%s, severity=%s, category=%s, labels=%s, updated_at=%s", updatedIssue.Subject, updatedIssue.Content, updatedIssue.Severity, updatedIssue.Category, pq.Array(updatedIssue.Labels), time.Now())}

	qry := sqlf.Sprintf("UPDATE issues SET %s WHERE id=%s AND customer_id=%s", sqlf.Join(fields, ","), id, customerId)

	if _, err := is.dbconn.ExecContext(ctx, qry.Query(sqlf.PostgresBindVar), qry.Args()...); err != nil {
		return nil, errors.Wrapf(err, "failed to update issue by id: %s customerId: %s", id, customerId)
	}

	return is.GetByID(ctx, id, projectId, customerId)
}

// List list issues.
func (is *IssuesPG) List(ctx context.Context, opt *IssueListOptions, projectId, customerId string) ([]api.Issue, error) {
	if opt == nil {
		opt = &IssueListOptions{}
	}

	conds := ListSubjectLikeSQL(opt.SubjectLikeOptions)
	conds = append(conds, sqlf.Sprintf("project_id = %s", projectId))
	conds = append(conds, sqlf.Sprintf("customer_id = %s", customerId))

	qry := sqlf.Sprintf("WHERE %s ORDER BY id ASC %s", sqlf.Join(conds, "AND"), opt.LimitOffset.SQL())

	return is.getBySQL(ctx, qry.Query(sqlf.PostgresBindVar), qry.Args()...)
}

func (is *IssuesPG) getBySQL(ctx context.Context, query string, args ...interface{}) ([]api.Issue, error) {
	rows, err := is.dbconn.QueryContext(ctx, "SELECT id, subject, state, severity, category, labels, content, created_at, updated_at FROM issues "+query, args...)
	if err != nil {
		return nil, err
	}

	issues := []api.Issue{}
	defer rows.Close()
	for rows.Next() {
		issue := api.Issue{}
		err := rows.Scan(&issue.Id, &issue.Subject, &issue.State, &issue.Severity, &issue.Category, pq.Array(&issue.Labels), &issue.Content, &issue.CreatedAt, &issue.UpdatedAt)
		if err != nil {
			return nil, err
		}

		issues = append(issues, issue)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return issues, nil
}
