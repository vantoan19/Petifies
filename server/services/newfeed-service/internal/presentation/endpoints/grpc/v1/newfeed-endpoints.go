package v1

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"

	commmonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	postfeedservice "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/application/services/postfeed-service"
	postfeedaggre "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/domain/aggregates/post-feed"
	"github.com/vantoan19/Petifies/server/services/newfeed-service/pkg/models"
)

type NewfeedEndpoints struct {
	ListPostFeeds  endpoint.Endpoint
	ListStoryFeeds endpoint.Endpoint
}

func NewNewfeedEndpoints(ps postfeedservice.PostfeedService) NewfeedEndpoints {
	return NewfeedEndpoints{
		ListPostFeeds:  makeListPostFeedsEndpoint(ps),
		ListStoryFeeds: makeListStoryFeedsEndpoint(),
	}
}

func makeListPostFeedsEndpoint(s postfeedservice.PostfeedService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListPostFeedsReq)
		result, err := s.ListPostFeeds(ctx, req)
		if err != nil {
			return nil, err
		}

		return &models.ListPostFeedsResp{
			PostIDs: commmonutils.Map2(result, func(p *postfeedaggre.PostFeedAggre) uuid.UUID { return p.GetPostID() }),
		}, nil
	}
}

func makeListStoryFeedsEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// req := request.(*models.ListStoryFeedsReq)
		return &models.ListStoryFeedsResp{}, nil
	}
}
