package user

import (
	"context"
	"errors"

	"github.com/google/uuid"

	postfeedaggre "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/domain/aggregates/post-feed"
	storyfeedaggre "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/domain/aggregates/story-feed"
	"github.com/vantoan19/Petifies/server/services/newfeed-service/internal/domain/aggregates/user/entities"
)

var (
	ErrPostFeedAlreadyExists   = errors.New("post feed already exists")
	ErrStoryFeedAlreadyExists  = errors.New("story feed already exists")
	ErrPetifyFeedAlreadyExists = errors.New("petify feed already exists")
)

type UserAggre struct {
	user       *entities.User
	postFeeds  []uuid.UUID
	storyFeeds []uuid.UUID
}

func NewUserAggregate(id uuid.UUID, email string) (*UserAggre, error) {
	user := &entities.User{
		ID:    id,
		Email: email,
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}

	return &UserAggre{
		user:       user,
		postFeeds:  make([]uuid.UUID, 0),
		storyFeeds: make([]uuid.UUID, 0),
	}, nil
}

func (u *UserAggre) AddPostFeed(post postfeedaggre.PostFeedAggre, repo postfeedaggre.PostFeedRepository) (*postfeedaggre.PostFeedAggre, error) {
	exists, err := repo.ExistsPostFeed(context.Background(), post.GetUserID(), post.GetPostID())
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPostFeedAlreadyExists
	}

	return repo.Save(context.Background(), post)
}

func (u *UserAggre) AddStoryFeed(story storyfeedaggre.StoryFeedAggre, repo storyfeedaggre.StoryFeedRepository) (*storyfeedaggre.StoryFeedAggre, error) {
	exists, err := repo.ExistsStoryFeed(context.Background(), story.GetUserID(), story.GetStoryID())
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrStoryFeedAlreadyExists
	}

	return repo.Save(context.Background(), story)
}

// ========== Aggregate Root Getters ============

func (u *UserAggre) GetID() uuid.UUID {
	return u.user.ID
}

func (u *UserAggre) GetEmail() string {
	return u.user.Email
}
