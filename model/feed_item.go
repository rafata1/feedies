package model

import "database/sql"

type FeedItemStatus int

const (
	FeedItemStatusUnspecified FeedItemStatus = 0
	FeedItemStatusPublished   FeedItemStatus = 1
	FeedItemStatusDraft       FeedItemStatus = 2
)

type FeedItem struct {
	ID          int64          `db:"id"`
	SourceID    int64          `db:"source_id"`
	Title       string         `db:"title"`
	Description string         `db:"description"`
	Content     string         `db:"content"`
	URL         string         `db:"url"`
	ImageURL    string         `db:"image_url"`
	Categories  string         `db:"categories"`
	Status      FeedItemStatus `db:"status"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}
