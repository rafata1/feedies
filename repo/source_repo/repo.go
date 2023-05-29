package source_repo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/rafata1/feedies/model"
)

type IRepo interface {
	UpsertSources(ctx context.Context, sources []model.Source) error
	GetSources(ctx context.Context) ([]model.Source, error)
	GetSourcesByIDs(ctx context.Context, ids []int64) ([]model.Source, error)
}

type repo struct {
	db *sqlx.DB
}

var upsertSourcesQuery = `INSERT INTO source (id, url, logo_url, status)
VALUES (:id, :url, :logo_url, :status)
ON DUPLICATE KEY UPDATE id=VALUES(id), url=VALUES(url), logo_url=VALUES(logo_url), status=VALUES(status)`

func (r *repo) UpsertSources(ctx context.Context, sources []model.Source) error {
	_, err := r.db.NamedExecContext(ctx, upsertSourcesQuery, sources)
	return err
}

var getSourcesByIDsQuery = `SELECT id, url, logo_url, status FROM source WHERE id IN (?)`

func (r *repo) GetSourcesByIDs(ctx context.Context, ids []int64) ([]model.Source, error) {
	var res []model.Source
	query, args, _ := sqlx.In(getSourcesByIDsQuery, ids)
	err := r.db.SelectContext(ctx, &res, query, args...)
	return res, err
}

var getSourcesQuery = `SELECT id, url, logo_url, status FROM source`

func (r *repo) GetSources(ctx context.Context) ([]model.Source, error) {
	var res []model.Source
	err := r.db.SelectContext(ctx, &res, getSourcesQuery)
	return res, err
}

func NewRepo(db *sqlx.DB) IRepo {
	return &repo{db: db}
}
