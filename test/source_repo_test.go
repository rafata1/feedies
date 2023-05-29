package test

import (
	"context"
	"github.com/rafata1/feedies/model"
	"github.com/rafata1/feedies/repo/source_repo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_UpsertSources_Insert_Update_Get(t *testing.T) {
	test := NewIntegration()
	repo := source_repo.NewRepo(test.db)

	sources := []model.Source{
		{
			URL:     "url",
			LogoURL: "image_url",
			Status:  model.SourceStatusInactive,
		},
	}

	ctx := context.Background()
	err := repo.UpsertSources(ctx, sources)
	assert.Nil(t, err)

	actual, err := repo.GetSources(ctx)
	assert.Nil(t, err)

	actual = ignoreIDSources(actual)
	assert.Equal(t, sources, actual)

	sources = []model.Source{
		{
			URL:     "url",
			LogoURL: "image_url 1",
			Status:  model.SourceStatusActive,
		},
	}

	err = repo.UpsertSources(ctx, sources)
	assert.Nil(t, err)

	actual, err = repo.GetSources(ctx)
	assert.Nil(t, err)

	actual = ignoreIDSources(actual)
	assert.Equal(t, sources, actual)
}

func ignoreIDSources(sources []model.Source) []model.Source {
	var res []model.Source
	for _, source := range sources {
		res = append(res, ignoreIDSource(source))
	}
	return res
}

func ignoreIDSource(source model.Source) model.Source {
	source.ID = 0
	return source
}
