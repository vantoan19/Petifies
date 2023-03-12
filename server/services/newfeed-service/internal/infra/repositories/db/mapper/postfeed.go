package mapper

import (
	"github.com/gocql/gocql"
	"github.com/google/uuid"

	postfeedaggre "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/domain/aggregates/post-feed"
	"github.com/vantoan19/Petifies/server/services/newfeed-service/internal/infra/repositories/db/models"
)

func DbPostFeedToPostFeedAggregate(postfeed *models.PostFeed) (*postfeedaggre.PostFeedAggre, error) {
	userID := uuid.UUID(postfeed.UserID)
	authorID := uuid.UUID(postfeed.AuthorID)
	postID := uuid.UUID(postfeed.PostID)

	return postfeedaggre.NewPostFeed(userID, authorID, postID, postfeed.CreatedAt)
}

func PostFeedAggregateToPostFeedDb(postfeed *postfeedaggre.PostFeedAggre) (*models.PostFeed, error) {
	return &models.PostFeed{
		UserID:    gocql.UUID(postfeed.GetUserID()),
		AuthorID:  gocql.UUID(postfeed.GetAuthorID()),
		PostID:    gocql.UUID(postfeed.GetPostID()),
		CreatedAt: postfeed.GetCreatedAt(),
	}, nil
}
