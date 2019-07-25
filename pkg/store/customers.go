package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/keegancsmith/sqlf"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/wolfeidau/exitus/pkg/api"
	"github.com/wolfeidau/exitus/pkg/conf"
	"github.com/wolfeidau/exitus/pkg/db"
)

// ErrCustomerNameAlreadyExists customer name is already taken
var ErrCustomerNameAlreadyExists = errors.New("customer name is already taken")

// CustomerNotFoundError occurs when an customer is not found.
type CustomerNotFoundError struct {
	Message string
}

func (e *CustomerNotFoundError) Error() string {
	return fmt.Sprintf("customer not found: %s", e.Message)
}

// Customers provides a customer store
type Customers interface {
	GetByID(ctx context.Context, customerId string) (*api.Customer, error)
	Create(ctx context.Context, newCustomer *api.NewCustomer) (*api.Customer, error)
	Update(ctx context.Context, updatedCustomer *api.UpdatedCustomer, customerId string) (*api.Customer, error)
	List(ctx context.Context, opt *CustomersListOptions) ([]api.Customer, error)
}

// CustomersPG provides a customer store using postgresql
type CustomersPG struct {
	dbconn *sql.DB
	cfg    *conf.Config
}

// NewCustomers new project store
func NewCustomers(dbconn *sql.DB, cfg *conf.Config) Customers {
	return &CustomersPG{dbconn: dbconn, cfg: cfg}
}

// GetByID get customer by id
func (cs *CustomersPG) GetByID(ctx context.Context, customerId string) (*api.Customer, error) {
	custs, err := cs.getBySQL(ctx, "WHERE id=$1 LIMIT 1", customerId)
	if err != nil {
		log.Error().Err(err).Msg("failed to get customer by id")
		return nil, err
	}

	if len(custs) == 0 {
		return nil, &CustomerNotFoundError{fmt.Sprintf("id %s", customerId)}
	}
	return &custs[0], nil
}

// Update update customer by id
func (cs *CustomersPG) Update(ctx context.Context, updatedCustomer *api.UpdatedCustomer, customerId string) (*api.Customer, error) {

	q := `UPDATE customers
		SET name=$2, labels=$3, updated_at=$4, description=$5
		WHERE id=$1;`

	if _, err := cs.dbconn.ExecContext(ctx, q, customerId,
		updatedCustomer.Name, pq.Array(updatedCustomer.Labels), time.Now(), updatedCustomer.Description); err != nil {
		return nil, err
	}

	resCust, err := cs.GetByID(ctx, customerId)
	if err != nil {
		log.Error().Err(err).Msg("failed to update customer")
		return nil, err
	}

	return resCust, nil
}

// Create create a customer
func (cs *CustomersPG) Create(ctx context.Context, newCustomer *api.NewCustomer) (*api.Customer, error) {

	resCust := &api.Customer{
		Name:        newCustomer.Name,
		Description: newCustomer.Description,
		Labels:      newCustomer.Labels,
	}

	err := db.WithTransaction(ctx, cs.dbconn, func(tx db.Transaction) error {
		return tx.QueryRowContext(
			ctx,
			`INSERT INTO customers(name, description, labels) VALUES($1, $2, $3)
			RETURNING id, created_at, updated_at`,
			newCustomer.Name, newCustomer.Description, pq.Array(newCustomer.Labels),
		).Scan(&resCust.Id, &resCust.CreatedAt, &resCust.UpdatedAt)
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to create customer")

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "customers_name_key":
				return nil, ErrCustomerNameAlreadyExists
			}
		}

		return nil, err
	}

	return resCust, nil
}

// CustomersListOptions specifies the options for listing customers.
type CustomersListOptions struct {
	*NameLikeOptions
	*LimitOffset
}

// NewCustomersListOptions create a new opts
func NewCustomersListOptions(query string, offset int, limit int) *CustomersListOptions {
	return &CustomersListOptions{
		NameLikeOptions: &NameLikeOptions{query},
		LimitOffset:     &LimitOffset{Limit: limit, Offset: offset},
	}
}

// List list all customers
func (cs *CustomersPG) List(ctx context.Context, opt *CustomersListOptions) ([]api.Customer, error) {
	if opt == nil {
		opt = &CustomersListOptions{}
	}

	conds := ListNameLikeSQL(opt.NameLikeOptions)

	q := sqlf.Sprintf("WHERE %s ORDER BY id ASC %s", sqlf.Join(conds, "AND"), opt.LimitOffset.SQL())

	log.Info().Str("q", q.Query(sqlf.PostgresBindVar)).Msg("Customers List getBySQL")

	return cs.getBySQL(ctx, q.Query(sqlf.PostgresBindVar), q.Args()...)
}

func (cs *CustomersPG) getBySQL(ctx context.Context, query string, args ...interface{}) ([]api.Customer, error) {
	rows, err := cs.dbconn.QueryContext(ctx, "SELECT id, name, description, labels, created_at, updated_at FROM customers "+query, args...)
	if err != nil {
		return nil, err
	}

	custs := []api.Customer{}
	defer rows.Close()
	for rows.Next() {
		cust := api.Customer{}
		err := rows.Scan(&cust.Id, &cust.Name, &cust.Description, pq.Array(&cust.Labels), &cust.CreatedAt, &cust.UpdatedAt)
		if err != nil {
			return nil, err
		}

		custs = append(custs, cust)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return custs, nil
}
