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
	testIssueId = "3b5d27e3-3524-4c34-a189-2c0cc30765f9"
	testAuthor  = "3b5d27e3-3524-4c34-a189-2c0cc30765f9"
)

func TestComments_CreateGetUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	assert := require.New(t)
	ctx := db.TestContext(t)

	cfg, err := conf.NewDefaultConfig()
	if err != nil {
		t.Fatal("failed to load config")
	}

	cstore := store.NewComments(db.Global, cfg)

	newComment, err := cstore.Create(ctx, &api.NewComment{
		Content: "test issue",
	}, testIssueId, testProjectId, testCustomerId, testAuthor)
	if err != nil {
		t.Fatal("failed to create a comment")
	}

	assert.NotEmpty(newComment.Id)
	assert.NotEmpty(newComment.UpdatedAt)
	assert.NotEmpty(newComment.CreatedAt)

	getComment, err := cstore.GetByID(ctx, newComment.Id, testIssueId, testProjectId, testCustomerId)
	if err != nil {
		t.Fatal("failed to get comment by id")
	}

	assert.Equal(newComment, getComment)

	newComment, err = cstore.Update(ctx, &api.UpdatedComment{
		NewComment: api.NewComment{
			Content: "updated test comment",
		},
	}, newComment.Id, testIssueId, testProjectId, testCustomerId)

	assert.Equal("updated test comment", newComment.Content)

	listComment, err := cstore.List(ctx, store.NewCommentListOptions("test", 0, 100), testIssueId, testProjectId, testCustomerId)
	if err != nil {
		t.Fatal("failed to get comments")
	}

	assert.Len(listComment, 1)
	assert.Equal(newComment, &listComment[0])
}
