package mapper

import (
	loveaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/love"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/infra/db/models"
	reqModels "github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

func DbLoveToEntityLove(l *models.Love) *entities.Love {
	return &entities.Love{
		ID:           l.ID,
		TargetID:     l.TargetID,
		IsPostTarget: l.IsPostTarget,
		AuthorID:     l.AuthorID,
		CreatedAt:    l.CreatedAt,
	}
}

func EntityLoveToDbLove(l *entities.Love) *models.Love {
	return &models.Love{
		ID:           l.ID,
		TargetID:     l.TargetID,
		IsPostTarget: l.IsPostTarget,
		AuthorID:     l.AuthorID,
		CreatedAt:    l.CreatedAt,
	}
}

func DbModelsToLoveAggregate(l *models.Love) (*loveaggre.Love, error) {
	return loveaggre.NewLove(&reqModels.Love{
		ID:           l.ID,
		TargetID:     l.TargetID,
		IsPostTarget: l.IsPostTarget,
		AuthorID:     l.AuthorID,
		CreatedAt:    l.CreatedAt,
	})
}
