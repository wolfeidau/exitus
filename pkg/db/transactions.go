package db

import (
	"context"
	"database/sql"

	multierror "github.com/hashicorp/go-multierror"
)

//
// code was adapted from https://pseudomuto.com/2018/01/clean-sql-transactions-in-golang/
//

// Transaction is an interface that models the standard transaction in
// `database/sql`.
//
// To ensure `TxFn` funcs cannot commit or rollback a transaction (which is
// handled by `WithTransaction`), those methods are not included here.
type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// TxFn is a function that will be called with an initialised `Transaction` object
// that can be used for executing statements and queries against a database.
type TxFn func(Transaction) error

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn`.
func WithTransaction(ctx context.Context, db *sql.DB, fn TxFn) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return
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

	err = fn(tx)
	return
}
