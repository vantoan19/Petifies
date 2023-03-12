package postfeedaggre

import (
	"time"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/newfeed-service/internal/domain/aggregates/post-feed/entities"
)

type PostFeedAggre struct {
	post *entities.PostFeed
}

func NewPostFeed(userID, authorID, postID uuid.UUID, createdAt time.Time) (*PostFeedAggre, error) {
	post := &entities.PostFeed{
		UserID:    userID,
		AuthorID:  authorID,
		PostID:    postID,
		CreatedAt: createdAt,
	}
	if errs := post.Validate(); errs.Exist() {
		return nil, errs[0]
	}

	return &PostFeedAggre{
		post: post,
	}, nil
}

// ========== Aggregate Root Getters ============

func (p *PostFeedAggre) GetUserID() uuid.UUID {
	return p.post.UserID
}

func (p *PostFeedAggre) GetAuthorID() uuid.UUID {
	return p.post.AuthorID
}

func (p *PostFeedAggre) GetPostID() uuid.UUID {
	return p.post.PostID
}

func (p *PostFeedAggre) GetCreatedAt() time.Time {
	return p.post.CreatedAt
}
