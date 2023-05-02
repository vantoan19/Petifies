package events

import (
	"time"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	reviewaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/reviews"
)

type ReviewCreatedEvent models.ReviewEvent

func NewReviewCreatedEvent(review *reviewaggre.ReviewAggre) ReviewCreatedEvent {
	return ReviewCreatedEvent{
		ID:         review.GetID(),
		PetifiesID: review.GetPetifiesID(),
		AuthorID:   review.GetAuthorID(),
		EventType:  models.REVIEW_CREATED,
		CreatedAt:  time.Now(),
	}
}
