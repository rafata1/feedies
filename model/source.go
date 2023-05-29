package model

import "database/sql"

type SourceStatus int

const (
	SourceStatusUnspecified SourceStatus = 0
	SourceStatusActive      SourceStatus = 1
	SourceStatusInactive    SourceStatus = 2
)

type Source struct {
	ID        int64        `db:"id"`
	URL       string       `db:"url"`
	LogoURL   string       `db:"logo_url"`
	Status    SourceStatus `db:"status"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
