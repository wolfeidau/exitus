package store_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wolfeidau/exitus/pkg/api"
	"github.com/wolfeidau/exitus/pkg/conf"
	"github.com/wolfeidau/exitus/pkg/db"
	"github.com/wolfeidau/exitus/pkg/store"
)

func TestCustomers_CreateGetUpdateList(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	assert := require.New(t)
	ctx := db.TestContext(t)

	cfg, err := conf.NewDefaultConfig()
	if err != nil {
		t.Fatal("failed to load config")
	}

	cstore := store.NewCustomers(db.Global, cfg)

	newCust, err := cstore.Create(ctx, &api.NewCustomer{
		Name:   "test customer",
		Labels: []string{"test"},
	})
	if err != nil {
		t.Fatal("failed to create customer")
	}

	assert.NotEmpty(newCust.Id)
	assert.NotEmpty(newCust.UpdatedAt)
	assert.NotEmpty(newCust.CreatedAt)

	getCust, err := cstore.GetByID(ctx, newCust.Id)
	if err != nil {
		t.Fatal("failed to get customer by id")
	}

	assert.Equal(newCust, getCust)

	newCust, err = cstore.Update(ctx, &api.UpdatedCustomer{
		NewCustomer: api.NewCustomer{
			Name:   "updated test customer",
			Labels: []string{"test", "update"},
		},
	}, getCust.Id)

	assert.Equal("updated test customer", newCust.Name)

	listCust, err := cstore.List(ctx, store.NewCustomersListOptions("test", 0, 100))
	if err != nil {
		t.Fatal("failed to get customer by id")
	}

	assert.Len(listCust, 1)
	assert.Equal(newCust, &listCust[0])
}
