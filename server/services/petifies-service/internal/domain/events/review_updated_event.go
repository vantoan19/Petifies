package events

import (
	"time"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	reviewaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/reviews"
)

type ReviewUpdatedEvent models.ReviewEvent

func NewReviewUpdatedEvent(review *reviewaggre.ReviewAggre) ReviewUpdatedEvent {
	return ReviewUpdatedEvent{
		ID:         review.GetID(),
		PetifiesID: review.GetPetifiesID(),
		AuthorID:   review.GetAuthorID(),
		EventType:  models.REVIEW_UPDATED,
		CreatedAt:  time.Now(),
	}
}
