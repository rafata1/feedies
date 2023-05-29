package test

import (
	"context"
	"database/sql"
	"github.com/rafata1/feedies/model"
	"github.com/rafata1/feedies/repo/feed_repo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_UpsertFeedItem_Insert_Update_Get(t *testing.T) {
	test := NewIntegration()
	repo := feed_repo.NewRepo(test.db)

	ctx := context.Background()

	items := []model.FeedItem{
		{
			SourceID:    1,
			Title:       "title",
			Description: "description",
			Content:     "content",
			URL:         "url",
			ImageURL:    "image_url",
			Categories:  "a,b,c",
			Status:      model.FeedItemStatusDraft,
		},
	}
	err := repo.UpsertFeedItems(ctx, items)
	assert.Nil(t, err)

	actual, err := repo.GetFeedItems(ctx)
	assert.Nil(t, err)

	actual = ignoreIDFeedItems(actual)
	actual = ignoreTimestamps(actual)
	assert.Equal(t, items, actual)

	items = []model.FeedItem{
		{
			SourceID:    2,
			Title:       "title 1",
			Description: "description 1",
			Content:     "content 1",
			URL:         "url",
			ImageURL:    "image_url 1",
			Categories:  "a,b,c,d",
			Status:      model.FeedItemStatusPublished,
		},
	}

	err = repo.UpsertFeedItems(ctx, items)
	assert.Nil(t, err)

	actual, err = repo.GetFeedItems(ctx)
	assert.Nil(t, err)

	actual = ignoreIDFeedItems(actual)
	actual = ignoreTimestamps(actual)
	assert.Equal(t, items, actual)
}

func ignoreTimestamps(items []model.FeedItem) []model.FeedItem {
	var res []model.FeedItem
	for _, item := range items {
		res = append(res, ignoreTimestamp(item))
	}
	return res
}

func ignoreTimestamp(item model.FeedItem) model.FeedItem {
	item.CreatedAt = sql.NullTime{}
	item.UpdatedAt = sql.NullTime{}
	return item
}

func ignoreIDFeedItems(items []model.FeedItem) []model.FeedItem {
	var res []model.FeedItem
	for _, item := range items {
		res = append(res, ignoreIDFeedItem(item))
	}
	return res
}

func ignoreIDFeedItem(item model.FeedItem) model.FeedItem {
	item.ID = 0
	return item
}
