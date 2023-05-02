package feedservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	newfeedclient "github.com/vantoan19/Petifies/server/services/grpc-clients/newfeed-client"
	postservice "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/application/services/post"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/domain/repositories"
	redisFeedCache "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/infra/repositories/feed/redis"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/pkg/models"
	newfeedModels "github.com/vantoan19/Petifies/server/services/newfeed-service/pkg/models"
)

var logger = logging.New("MobileGateway.LoveSvc")

type FeedConfiguration func(rs *feedService) error

type feedService struct {
	newfeedClient       newfeedclient.NewfeedClient
	feedCacheRepository repositories.FeedCacheRepository
	postService         postservice.PostService
}

type FeedService interface {
	ListBatchPostFeeds(ctx context.Context, userID uuid.UUID, afterPostID uuid.UUID) ([]*models.PostWithUserInfo, error)
}

func NewFeedService(newfeedClientConn *grpc.ClientConn, postService postservice.PostService, cfgs ...FeedConfiguration) (FeedService, error) {
	fs := &feedService{
		newfeedClient: newfeedclient.New(newfeedClientConn),
		postService:   postService,
	}
	for _, cfg := range cfgs {
		err := cfg(fs)
		if err != nil {
			return nil, err
		}
	}
	return fs, nil
}

func WithRedisFeedCacheRepository(client *redis.Client) FeedConfiguration {
	return func(fs *feedService) error {
		repo := redisFeedCache.NewRedisFeedCacheRepository(client)
		fs.feedCacheRepository = repo
		return nil
	}
}

func (rs *feedService) ListBatchPostFeeds(ctx context.Context, userID uuid.UUID, afterPostID uuid.UUID) ([]*models.PostWithUserInfo, error) {
	logger.Info("Start ListBatchPostFeeds")

	postFeedIds, err := rs.FetchPostFeedIDs(ctx, userID, afterPostID)
	if err != nil {
		logger.ErrorData("Finished ListBatchPostFeeds: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	var willBeReturnedPostFeedIds []uuid.UUID
	for i := 0; i < 40 && i < len(postFeedIds); i++ {
		willBeReturnedPostFeedIds = append(willBeReturnedPostFeedIds, postFeedIds[i])
	}
	posts, err := rs.postService.ListPostsWithUserInfos(ctx, willBeReturnedPostFeedIds)
	if err != nil {
		logger.ErrorData("Finished ListBatchPostFeeds: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListBatchPostFeeds: Successful")
	return posts, nil
}

func (rs *feedService) FetchPostFeedIDs(ctx context.Context, userID uuid.UUID, afterPostID uuid.UUID) ([]uuid.UUID, error) {
	if exists, _ := rs.feedCacheRepository.ExistsPostFeedIDs(ctx, userID); exists {
		logger.Info("Executing FetchPostFeedIDs: getting post feed ids from cache")
		feeds, err := rs.feedCacheRepository.GetPostFeedIDs(ctx, userID)
		if err != nil {
			logger.WarningData("Executing FetchPostFeedIDs: failed to get feeds from cache", logging.Data{"error": err.Error()})
		} else {
			lastPostIdx := commonutils.FindFirst(feeds, func(id uuid.UUID) bool { return id == afterPostID })

			if lastPostIdx != -1 && lastPostIdx < len(feeds)-1 {
				logger.Info("Finish FetchPostFeedIDs: Successful")
				return feeds[lastPostIdx+1:], nil
			}
		}
	}

	logger.Info("Executing FetchPostFeedIDs: getting post feed ids from NewFeed Service")
	feedIds, err := rs.newfeedClient.ListPostFeeds(ctx, &newfeedModels.ListPostFeedsReq{
		UserID:      userID,
		PageSize:    400,
		AfterPostID: afterPostID,
	})
	if err != nil {
		logger.ErrorData("Finish FetchPostFeedIDs: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Executing FetchPostFeedIDs: Populating posts' info from Post Service")
	posts, err := rs.postService.ListPostsWithUserInfos(ctx, feedIds.PostIDs)
	if err != nil {
		logger.ErrorData("Finished FetchPostFeedIDs: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}
	validPosts := commonutils.Filter(posts, func(p *models.PostWithUserInfo) bool { return p.Visibility != "private" })
	validPostIDs := commonutils.Map2(validPosts, func(p *models.PostWithUserInfo) uuid.UUID { return p.ID })

	go func() {
		err := rs.feedCacheRepository.SetPostFeedIDs(context.Background(), userID, validPostIDs)
		if err != nil {
			logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
		}
	}()

	return validPostIDs, nil
}
