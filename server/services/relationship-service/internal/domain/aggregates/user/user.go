package useraggre

import (
	"errors"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/relationship-service/internal/domain/aggregates/user/entities"
)

var (
	ErrRelationshipAlreadyExists = errors.New("relationship already exist")
	ErrRelationshipNotExist      = errors.New("relationship does not exist")
	ErrExceedRelationshipLimit   = errors.New("exceeds relationship limit")
)

// UserAggregate represents the aggregate for a user and their relationships
type UserAggregate struct {
	user       *entities.User
	followers  []uuid.UUID
	followings []uuid.UUID
}

func NewUserAggregate(user entities.User) (*UserAggregate, error) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	return &UserAggregate{
		user:       &user,
		followers:  make([]uuid.UUID, 0),
		followings: make([]uuid.UUID, 0),
	}, nil
}

func (u *UserAggregate) Follow(userID uuid.UUID) error {
	if len(u.followings) > 2000 {
		return ErrExceedRelationshipLimit
	}
	for _, id := range u.followings {
		if id == userID {
			return ErrRelationshipAlreadyExists
		}
	}

	u.followings = append(u.followings, userID)
	return nil
}

func (u *UserAggregate) Unfollow(userID uuid.UUID) error {
	for i, r := range u.followings {
		if r == userID {
			u.followings = append(u.followings[:i], u.followings[i+1:]...)
			return nil
		}
	}
	return ErrRelationshipNotExist
}

func (u *UserAggregate) GetFollowings() []uuid.UUID {
	return append(u.followings, u.user.ID)
}

func (u *UserAggregate) AddFollower(userID uuid.UUID) {
	u.followers = append(u.followers, userID)
}

func (u *UserAggregate) GetFollowers() []uuid.UUID {
	return append(u.followers, u.user.ID)
}

// ========= Aggregate Root Getter =========

func (u *UserAggregate) GetID() uuid.UUID {
	return u.user.ID
}

func (u *UserAggregate) GetEmail() string {
	return u.user.Email
}
