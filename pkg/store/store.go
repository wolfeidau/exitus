package store

import (
	"database/sql"

	"github.com/keegancsmith/sqlf"
	"github.com/wolfeidau/exitus/pkg/conf"
)

// Stores one stop for stores.
type Stores struct {
	Projects  Projects
	Customers Customers
	Issues    Issues
	Comments  Comments
}

// New create all the stores.
func New(dbconn *sql.DB, cfg *conf.Config) (*Stores, error) {
	return &Stores{
		Projects:  NewProjects(dbconn, cfg),
		Customers: NewCustomers(dbconn, cfg),
		Issues:    NewIssues(dbconn, cfg),
		Comments:  NewComments(dbconn, cfg),
	}, nil
}

// LimitOffset specifies SQL LIMIT and OFFSET counts. A pointer to it is typically embedded in other options
// structures that need to perform SQL queries with LIMIT and OFFSET.
type LimitOffset struct {
	Limit  int // SQL LIMIT count
	Offset int // SQL OFFSET count
}

// SQL returns the SQL query fragment ("LIMIT %d OFFSET %d") for use in SQL queries.
func (o *LimitOffset) SQL() *sqlf.Query {
	if o == nil {
		return &sqlf.Query{}
	}
	return sqlf.Sprintf("LIMIT %d OFFSET %d", o.Limit, o.Offset)
}

// NameLikeOptions used to query by name using like.
type NameLikeOptions struct {
	// Query specifies a search query for organizations.
	Query string
}

// ListNameLikeSQL used to search by name if query is set.
func ListNameLikeSQL(opt *NameLikeOptions) (conds []*sqlf.Query) {
	conds = []*sqlf.Query{sqlf.Sprintf("TRUE")}
	if opt.Query != "" {
		query := "%" + opt.Query + "%"
		conds = append(conds, sqlf.Sprintf("name ILIKE %s", query))
	}
	return conds
}
