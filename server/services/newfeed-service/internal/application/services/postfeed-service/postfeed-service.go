package postfeedservice

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	postfeedaggre "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/domain/aggregates/post-feed"
	"github.com/vantoan19/Petifies/server/services/newfeed-service/internal/infra/repositories/post/cassandra"
	"github.com/vantoan19/Petifies/server/services/newfeed-service/pkg/models"
)

var logger = logging.New("NewfeedService.Service")

type PostfeedConfiguration func(ps *postfeedService) error

type postfeedService struct {
	postFeedRepository postfeedaggre.PostFeedRepository
}

type PostfeedService interface {
	ListPostFeeds(ctx context.Context, req *models.ListPostFeedsReq) ([]*postfeedaggre.PostFeedAggre, error)
}

func NewPostFeedService(cfgs ...PostfeedConfiguration) (PostfeedService, error) {
	ps := &postfeedService{}
	for _, cfg := range cfgs {
		err := cfg(ps)
		if err != nil {
			return nil, err
		}
	}
	return ps, nil
}

func WithCassandraPostfeedRepository(session *gocql.Session) PostfeedConfiguration {
	return func(ps *postfeedService) error {
		postfeedRepo, err := cassandra.NewCassandraPostRepository(session)
		if err != nil {
			return err
		}
		ps.postFeedRepository = postfeedRepo
		return nil
	}
}

func (s *postfeedService) ListPostFeeds(ctx context.Context, req *models.ListPostFeedsReq) ([]*postfeedaggre.PostFeedAggre, error) {
	logger.Info("Start ListPostFeeds")

	postFeeds, err := s.postFeedRepository.GetByUserID(ctx, req.UserID, req.PageSize, req.BeforeTime)
	if err != nil {
		logger.ErrorData("Finish ListPostFeeds: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListPostFeeds: Successful")
	return postFeeds, err
}
