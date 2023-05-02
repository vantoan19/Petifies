package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
)

type PetifiesCacheRepository interface {
	GetPetifies(ctx context.Context, petifiesId uuid.UUID) (*models.Petifies, error)
	ListPetifies(ctx context.Context, petifiesIds []uuid.UUID) ([]*models.Petifies, error)
	SetPetifies(ctx context.Context, petifiesId uuid.UUID, petifies *models.Petifies) error
}
