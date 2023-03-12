package petifyfeedaggre

import (
	"time"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/newfeed-service/internal/domain/aggregates/petify-feed/entities"
)

type PetifyFeedAggre struct {
	petify *entities.PetifyFeed
}

func NewPetifyFeed(userID, authorID, petifyID uuid.UUID, createdAt time.Time) (*PetifyFeedAggre, error) {
	petify := &entities.PetifyFeed{
		UserID:    userID,
		PetifyID:  petifyID,
		AuthorID:  authorID,
		CreatedAt: createdAt,
	}
	if errs := petify.Validate(); errs.Exist() {
		return nil, errs[0]
	}

	return &PetifyFeedAggre{
		petify: petify,
	}, nil
}

// ========== Aggregate Root Getters ============

func (p *PetifyFeedAggre) GetUserID() uuid.UUID {
	return p.petify.UserID
}

func (p *PetifyFeedAggre) GetAuthorID() uuid.UUID {
	return p.petify.AuthorID
}

func (p *PetifyFeedAggre) GetPetifyID() uuid.UUID {
	return p.petify.PetifyID
}

func (p *PetifyFeedAggre) GetCreatedAt() time.Time {
	return p.petify.CreatedAt
}
