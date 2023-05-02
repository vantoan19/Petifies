package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
)

type PetifiesSessionCacheRepository interface {
	GetPetifiesSession(ctx context.Context, sessionId uuid.UUID) (*models.PetifiesSession, error)
	SetPetifiesSession(ctx context.Context, sessionId uuid.UUID, petifiesSession *models.PetifiesSession) error
}
