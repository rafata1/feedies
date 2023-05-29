package feed_service

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/mmcdole/gofeed"
	"github.com/rafata1/feedies/model"
	"github.com/rafata1/feedies/repo/feed_repo"
	"github.com/rafata1/feedies/repo/source_repo"
	"github.com/rafata1/feedies/service/common"
	"log"
	"strings"
)

type IService interface {
	ConsumeRSS(ctx context.Context)
	GetNews(ctx context.Context) (GetNewsOutput, error)
}

type service struct {
	feedRepo   feed_repo.IRepo
	sourceRepo source_repo.IRepo
	parser     *gofeed.Parser
}

func (s *service) GetNews(ctx context.Context) (GetNewsOutput, error) {
	feedItems, err := s.feedRepo.GetFeedItems(ctx)
	if err != nil {
		return GetNewsOutput{}, common.WrapErrQueryDB(err)
	}

	sources, err := s.sourceRepo.GetSourcesByIDs(ctx, getSourceIDs(feedItems))
	if err != nil {
		return GetNewsOutput{}, common.WrapErrSaveDB(err)
	}

	return GetNewsOutput{
		News: getNews(feedItems, sources),
	}, nil
}

func getNews(feedItems []model.FeedItem, sources []model.Source) []New {
	sourceMp := make(map[int64]model.Source)
	for _, source := range sources {
		sourceMp[source.ID] = source
	}
	var res []New
	for _, item := range feedItems {
		source, ok := sourceMp[item.SourceID]
		if !ok {
			continue
		}
		res = append(res, New{
			SourceURL:     source.URL,
			SourceLogoURL: source.LogoURL,
			Title:         item.Title,
			ImageURL:      item.ImageURL,
			CreatedAt:     item.CreatedAt.Time.Format("January 2"),
		})
	}
	return res
}

func getSourceIDs(feedItems []model.FeedItem) []int64 {
	var res []int64
	for _, item := range feedItems {
		res = append(res, item.SourceID)
	}
	return res
}

func (s *service) ConsumeRSS(ctx context.Context) {
	sources, err := s.sourceRepo.GetSources(ctx)
	if err != nil {
		log.Printf("Error getting rss sources: %s\n", err.Error())
		return
	}

	for _, source := range sources {
		nums, err := s.consumeFeedItems(ctx, source)
		if err != nil {
			log.Printf("Error consume feed items of %s: %s\n", source.URL, err.Error())
		}
		log.Printf("Consumed %d feed items from %s\n", nums, source.URL)
	}
	return
}

func (s *service) consumeFeedItems(ctx context.Context, source model.Source) (int, error) {
	feed, err := s.parser.ParseURL(source.URL)
	if err != nil {
		return 0, err
	}

	feedItems := getFeedItems(feed, source.ID)
	return len(feedItems), s.feedRepo.UpsertFeedItems(ctx, feedItems)
}

func getFeedItems(feed *gofeed.Feed, sourceID int64) []model.FeedItem {
	var res []model.FeedItem
	for _, item := range feed.Items {
		res = append(res, toFeedItem(item, sourceID))
	}
	return res
}

func toFeedItem(item *gofeed.Item, sourceID int64) model.FeedItem {
	res := model.FeedItem{
		SourceID:    sourceID,
		Title:       item.Title,
		Description: item.Description,
		Content:     item.Content,
		URL:         item.Link,
		Categories:  strings.Join(item.Categories, ","),
		Status:      model.FeedItemStatusPublished,
	}
	if item.Image != nil {
		res.ImageURL = item.Image.URL
	}
	return res
}

func NewService(feedRepo feed_repo.IRepo, sourceRepo source_repo.IRepo) IService {
	return &service{
		feedRepo:   feedRepo,
		sourceRepo: sourceRepo,
		parser:     gofeed.NewParser(),
	}
}

func Init(db *sqlx.DB) IService {
	feedRepo := feed_repo.NewRepo(db)
	sourceRepo := source_repo.NewRepo(db)
	return NewService(feedRepo, sourceRepo)
}
