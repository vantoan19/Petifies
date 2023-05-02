package locationaggre

import (
	"context"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/location-service/internal/domain/aggregates/location/valueobjects"
)

type LocationRepository interface {
	FindNearbyLocationsByEntityType(ctx context.Context, longitude, latitude float64, maxDistance float64, entityType valueobjects.EntityType,
		pageSize int, offset int) ([]*LocationAggre, error)
	GetByID(ctx context.Context, id uuid.UUID) (*LocationAggre, error)
	GetByEntityID(ctx context.Context, entityID uuid.UUID) (*LocationAggre, error)
	Save(ctx context.Context, location LocationAggre) (*LocationAggre, error)
	Update(ctx context.Context, location LocationAggre) (*LocationAggre, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	DeleteByEntityID(ctx context.Context, entityID uuid.UUID) error
}
