package listener

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	relationshipclient "github.com/vantoan19/Petifies/server/services/grpc-clients/relationship-client"
	postfeedaggre "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/domain/aggregates/post-feed"
	cassandraPostRepo "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/infra/repositories/post/cassandra"
	cassandraUserRepo "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/infra/repositories/user/cassandra"
)

var logger = logging.New("NewFeedService.PostEventListener")

type PostEventListener struct {
	postRepo           *cassandraPostRepo.PostRepository
	userRepo           *cassandraUserRepo.UserRepository
	relationshipClient relationshipclient.RelationshipClient
}

func NewPostEventListener(userRepo *cassandraUserRepo.UserRepository, postRepo *cassandraPostRepo.PostRepository, relationshipClient relationshipclient.RelationshipClient) *PostEventListener {
	return &PostEventListener{
		userRepo:           userRepo,
		postRepo:           postRepo,
		relationshipClient: relationshipClient,
	}
}

func (pl *PostEventListener) PostCreated(ctx context.Context, event models.PostEvent) ([]*postfeedaggre.PostFeedAggre, error) {
	logger.Info("Start PostCreated")

	results := make([]*postfeedaggre.PostFeedAggre, 0)

	// Check whether author exists or not
	author, err := pl.userRepo.GetByID(ctx, event.AuthorID)
	if err != nil {
		logger.ErrorData("Finish PostCreated: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}
	if author == nil {
		logger.ErrorData("Finish PostCreated: FAILED", logging.Data{"error": "author does not exist"})
		return nil, status.Errorf(codes.NotFound, "author does not exist")
	}

	// Get follower list of the author
	resp, err := pl.relationshipClient.ListFollowers(ctx, event.AuthorID)
	if err != nil {
		logger.ErrorData("Finish PostCreated: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}
	for _, followerID := range resp.FollowerIDs {
		// Get follower
		follower, err := pl.userRepo.GetByID(ctx, followerID)
		if err != nil {
			logger.ErrorData("Finish PostCreated: FAILED", logging.Data{"error": err.Error()})
			return nil, err
		}
		if follower == nil {
			logger.ErrorData("Finish PostCreated: FAILED", logging.Data{"error": "follower does not exist"})
			return nil, status.Errorf(codes.NotFound, "follower does not exist")
		}

		// Add post feed to the DB by using User Aggregate as the entry
		postfeed, err := postfeedaggre.NewPostFeed(followerID, event.AuthorID, event.ID, event.CreatedAt)
		if err != nil {
			logger.ErrorData("Finish PostCreated: FAILED", logging.Data{"error": err.Error()})
			return nil, err
		}
		savedPostFeed, err := follower.AddPostFeed(*postfeed, pl.postRepo)
		if err != nil {
			logger.ErrorData("Finish PostCreated: FAILED", logging.Data{"error": err.Error()})
			return nil, err
		}

		results = append(results, savedPostFeed)
	}

	logger.Info("Finish PostCreated: Success")
	return results, nil
}
