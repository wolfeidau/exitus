package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/keegancsmith/sqlf"
	"github.com/pkg/errors"
	"github.com/wolfeidau/exitus/pkg/api"
	"github.com/wolfeidau/exitus/pkg/conf"
	"github.com/wolfeidau/exitus/pkg/db"
)

// CommentNotFoundError occurs when an comment is not found.
type CommentNotFoundError struct {
	Message string
}

func (e *CommentNotFoundError) Error() string {
	return fmt.Sprintf("comment not found: %s", e.Message)
}

// Comments provides a comments store.
type Comments interface {
	GetByID(ctx context.Context, id, issueId, projectId, customerId string) (*api.Comment, error)
	Create(ctx context.Context, newComment *api.NewComment, issueId, projectId, customerId, author string) (*api.Comment, error)
	Update(ctx context.Context, updatedComment *api.UpdatedComment, id, issueId, projectId, customerId string) (*api.Comment, error)
	List(ctx context.Context, opt *CommentListOptions, issueId, projectId, customerId string) ([]api.Comment, error)
}

// CommentListOptions specifies the options for listing comments.
type CommentListOptions struct {
	*ContentLikeOptions
	*LimitOffset
}

// NewCommentListOptions create a new opts.
func NewCommentListOptions(query string, offset int, limit int) *CommentListOptions {
	return &CommentListOptions{
		ContentLikeOptions: &ContentLikeOptions{query},
		LimitOffset:        &LimitOffset{Limit: limit, Offset: offset},
	}
}

// ContentLikeOptions used to query by content using like.
type ContentLikeOptions struct {
	// Query specifies a search query for organizations.
	Query string
}

// ListContentLikeSQL used to search by content if query is set.
func ListContentLikeSQL(opt *ContentLikeOptions) (conds []*sqlf.Query) {
	conds = []*sqlf.Query{sqlf.Sprintf("TRUE")}
	if opt.Query != "" {
		query := "%" + opt.Query + "%"
		conds = append(conds, sqlf.Sprintf("content ILIKE %s", query))
	}
	return conds
}

// CommentsPG provides a comments store for postgresql.
type CommentsPG struct {
	dbconn *sql.DB
	cfg    *conf.Config
}

// NewComments new comments store.
func NewComments(dbconn *sql.DB, cfg *conf.Config) Comments {
	return &CommentsPG{dbconn: dbconn, cfg: cfg}
}

// GetById get comment by id.
func (cs *CommentsPG) GetByID(ctx context.Context, id, issueId, projectId, customerId string) (*api.Comment, error) {
	comments, err := cs.getBySQL(ctx, "WHERE id=$1 AND issue_id=$2 AND project_id=$3 AND customer_id=$4 LIMIT 1", id, issueId, projectId, customerId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get comment by id: %s issueId: %s projectId: %s customerId: %s", id, issueId, projectId, customerId)
	}

	if len(comments) == 0 {
		return nil, &CommentNotFoundError{fmt.Sprintf("id %s issueId: %s project_id %s", id, issueId, projectId)}
	}

	return &comments[0], nil
}

// Create create new comment.
func (cs *CommentsPG) Create(ctx context.Context, newComment *api.NewComment, issueId, projectId, customerId, author string) (*api.Comment, error) {
	comment := api.Comment{}

	qry := sqlf.Sprintf("INSERT INTO comments(issue_id, project_id, customer_id, author, content) VALUES(%s, %s, %s, %s, %s)",
		issueId, projectId, customerId, author, newComment.Content)

	err := db.WithTransaction(ctx, cs.dbconn, func(tx db.Transaction) error {
		return tx.QueryRowContext(
			ctx, qry.Query(sqlf.PostgresBindVar)+" RETURNING id, content, created_at, updated_at", qry.Args()...,
		).Scan(&comment.Id, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create comment with subject: %s issueId: %s projectId: %s customerId: %s", newComment.Content, issueId, projectId, customerId)
	}

	return &comment, nil
}

// Update update an comment.
func (cs *CommentsPG) Update(ctx context.Context, updatedComment *api.UpdatedComment, id, issueId, projectId, customerId string) (*api.Comment, error) {
	fields := []*sqlf.Query{sqlf.Sprintf("content=%s, updated_at=%s", updatedComment.Content, time.Now())}

	qry := sqlf.Sprintf("UPDATE comments SET %s WHERE id=%s AND customer_id=%s", sqlf.Join(fields, ","), id, customerId)

	if _, err := cs.dbconn.ExecContext(ctx, qry.Query(sqlf.PostgresBindVar), qry.Args()...); err != nil {
		return nil, errors.Wrapf(err, "failed to update issue by id: %s customerId: %s", id, customerId)
	}

	return cs.GetByID(ctx, id, issueId, projectId, customerId)
}

// List list comments.
func (cs *CommentsPG) List(ctx context.Context, opt *CommentListOptions, issueId, projectId, customerId string) ([]api.Comment, error) {
	if opt == nil {
		opt = &CommentListOptions{}
	}

	conds := ListContentLikeSQL(opt.ContentLikeOptions)
	conds = append(conds, sqlf.Sprintf("issue_id = %s", issueId))
	conds = append(conds, sqlf.Sprintf("project_id = %s", projectId))
	conds = append(conds, sqlf.Sprintf("customer_id = %s", customerId))

	qry := sqlf.Sprintf("WHERE %s ORDER BY id ASC %s", sqlf.Join(conds, "AND"), opt.LimitOffset.SQL())

	return cs.getBySQL(ctx, qry.Query(sqlf.PostgresBindVar), qry.Args()...)
}

func (cs *CommentsPG) getBySQL(ctx context.Context, query string, args ...interface{}) ([]api.Comment, error) {
	rows, err := cs.dbconn.QueryContext(ctx, "SELECT id, content, created_at, updated_at FROM comments "+query, args...)
	if err != nil {
		return nil, err
	}

	comments := []api.Comment{}
	defer rows.Close()
	for rows.Next() {
		comment := api.Comment{}
		err := rows.Scan(&comment.Id, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
