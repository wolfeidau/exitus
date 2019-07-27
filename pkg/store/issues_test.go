package store_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wolfeidau/exitus/pkg/api"
	"github.com/wolfeidau/exitus/pkg/conf"
	"github.com/wolfeidau/exitus/pkg/db"
	"github.com/wolfeidau/exitus/pkg/store"
)

const (
	testProjectId = "3b5d27e3-3524-4c34-a189-2c0cc30765f9"
	testReporter  = "3b5d27e3-3524-4c34-a189-2c0cc30765f9"
)

func TestIssues_CreateGetUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	assert := require.New(t)
	ctx := db.TestContext(t)

	cfg, err := conf.NewDefaultConfig()
	if err != nil {
		t.Fatal("failed to create issue")
	}

	istore := store.NewIssues(db.Global, cfg)

	newIssue, err := istore.Create(ctx, &api.NewIssue{
		Subject: "test issue",
		Labels:  []string{"test"},
	}, testProjectId, testCustomerId, testReporter)
	if err != nil {
		t.Fatal("failed to load config")
	}

	assert.NotEmpty(newIssue.Id)
	assert.NotEmpty(newIssue.UpdatedAt)
	assert.NotEmpty(newIssue.CreatedAt)

	getIssue, err := istore.GetByID(ctx, newIssue.Id, testProjectId, testCustomerId)
	if err != nil {
		t.Fatal("failed to get issue by id")
	}

	assert.Equal(newIssue, getIssue)

	newIssue, err = istore.Update(ctx, &api.UpdatedIssue{
		NewIssue: api.NewIssue{
			Subject: "updated test issue",
			Labels:  []string{"test", "updated"},
		},
	}, newIssue.Id, testProjectId, testCustomerId)

	assert.Equal("updated test issue", newIssue.Subject)

	listIssue, err := istore.List(ctx, store.NewIssueListOptions("test", 0, 100), testProjectId, testCustomerId)
	if err != nil {
		t.Fatal("failed to get issue by id")
	}

	assert.Len(listIssue, 1)
	assert.Equal(newIssue, &listIssue[0])
}
