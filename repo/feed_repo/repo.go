package feed_repo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/rafata1/feedies/model"
)

type IRepo interface {
	UpsertFeedItems(ctx context.Context, items []model.FeedItem) error
	GetFeedItems(ctx context.Context) ([]model.FeedItem, error)
}

type repo struct {
	db *sqlx.DB
}

var upsertFeedItemsQuery = `INSERT INTO feed_item
(id, source_id, title, description, content, url, image_url, categories, status)
VALUES (:id, :source_id, :title, :description, :content, :url, :image_url, :categories, :status)
ON DUPLICATE KEY UPDATE
id = VALUES(id),
source_id = VALUES(source_id),
title = VALUES(title),
description = VALUES(description),
content = VALUES(content),
url = VALUES(url),
image_url = VALUES(image_url),
categories = VALUES(categories),
status = VALUES(status)`

func (r repo) UpsertFeedItems(ctx context.Context, items []model.FeedItem) error {
	_, err := r.db.NamedExecContext(ctx, upsertFeedItemsQuery, items)
	return err
}

var getFeedItemsQuery = `SELECT
id, source_id, title, description, content, url, image_url, categories, status, created_at
FROM feed_item ORDER BY created_at desc`

func (r repo) GetFeedItems(ctx context.Context) ([]model.FeedItem, error) {
	var res []model.FeedItem
	err := r.db.SelectContext(ctx, &res, getFeedItemsQuery)
	return res, err
}

func NewRepo(db *sqlx.DB) IRepo {
	return &repo{
		db: db,
	}
}
