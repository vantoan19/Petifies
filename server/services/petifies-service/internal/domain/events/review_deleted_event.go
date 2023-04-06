package events

import (
	"time"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	reviewaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/reviews"
)

type ReviewDeletedEvent models.ReviewEvent

func NewReviewDeletedEvent(review *reviewaggre.ReviewAggre) ReviewDeletedEvent {
	return ReviewDeletedEvent{
		ID:         review.GetID(),
		PetifiesID: review.GetPetifiesID(),
		AuthorID:   review.GetAuthorID(),
		EventType:  models.REVIEW_DELETED,
		CreatedAt:  time.Now(),
	}
}
