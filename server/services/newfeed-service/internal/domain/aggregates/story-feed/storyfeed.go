package storyfeedaggre

import (
	"time"

	"github.com/google/uuid"

	"github.com/vantoan19/Petifies/server/services/newfeed-service/internal/domain/aggregates/story-feed/entities"
)

type StoryFeedAggre struct {
	story *entities.StoryFeed
}

func NewStoryFeed(id, userID, authorID, storyID uuid.UUID, createdAt time.Time) (*StoryFeedAggre, error) {
	story := &entities.StoryFeed{
		UserID:    userID,
		AuthorID:  authorID,
		StoryID:   storyID,
		CreatedAt: createdAt,
	}
	if errs := story.Validate(); errs.Exist() {
		return nil, errs[0]
	}

	return &StoryFeedAggre{
		story: story,
	}, nil
}

// ========== Aggregate Root Getters ============

func (p *StoryFeedAggre) GetUserID() uuid.UUID {
	return p.story.UserID
}

func (p *StoryFeedAggre) GetAuthorID() uuid.UUID {
	return p.story.AuthorID
}

func (p *StoryFeedAggre) GetStoryID() uuid.UUID {
	return p.story.StoryID
}

func (p *StoryFeedAggre) GetCreatedAt() time.Time {
	return p.story.CreatedAt
}
