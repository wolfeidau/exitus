package store_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wolfeidau/exitus/pkg/api"
	"github.com/wolfeidau/exitus/pkg/conf"
	"github.com/wolfeidau/exitus/pkg/db"
	"github.com/wolfeidau/exitus/pkg/store"
)

const testCustomerId = "3b5d27e3-3524-4c34-a189-2c0cc30765f9"

func TestProjects_CreateGetUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	assert := require.New(t)
	ctx := db.TestContext(t)

	cfg, err := conf.NewDefaultConfig()
	if err != nil {
		t.Fatal("failed to create project")
	}

	pstore := store.NewProjects(db.Global, cfg)

	newProj, err := pstore.Create(ctx, &api.NewProject{
		Name:   "test project",
		Labels: []string{"test"},
	}, testCustomerId)
	if err != nil {
		t.Fatal("failed to load config")
	}

	assert.NotEmpty(newProj.Id)
	assert.NotEmpty(newProj.UpdatedAt)
	assert.NotEmpty(newProj.CreatedAt)

	getProj, err := pstore.GetByID(ctx, newProj.Id, "3b5d27e3-3524-4c34-a189-2c0cc30765f9")
	if err != nil {
		t.Fatal("failed to get project by id")
	}

	assert.Equal(newProj, getProj)

	newProj, err = pstore.Update(ctx, &api.UpdatedProject{
		NewProject: api.NewProject{
			Name:   "updated test project",
			Labels: []string{"test", "update"},
		},
	}, newProj.Id, testCustomerId)
	if err != nil {
		t.Fatal("failed to update project by id")
	}

	assert.Equal("updated test project", newProj.Name)

	listProj, err := pstore.List(ctx, store.NewProjectsListOptions("test", 0, 100), testCustomerId)
	if err != nil {
		t.Fatal("failed to list projects")
	}

	assert.Len(listProj, 1)
	assert.Equal(newProj, &listProj[0])
}
