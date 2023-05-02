package loveaggre

import (
	"time"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

type Love struct {
	love *entities.Love
}

func NewLove(info *models.Love) (*Love, error) {
	loveEntity := &entities.Love{
		ID:           info.ID,
		TargetID:     info.TargetID,
		IsPostTarget: info.IsPostTarget,
		AuthorID:     info.AuthorID,
		CreatedAt:    info.CreatedAt,
	}

	if errs := loveEntity.Validate(); errs.Exist() {
		return nil, errs[0]
	}

	return &Love{
		love: loveEntity,
	}, nil
}

// ========== Aggregate Root Getters ===========

func (l *Love) GetID() uuid.UUID {
	return l.love.ID
}

func (l *Love) GetTargetID() uuid.UUID {
	return l.love.TargetID
}

func (l *Love) GetAuthorID() uuid.UUID {
	return l.love.AuthorID
}

func (l *Love) GetIsPostTarget() bool {
	return l.love.IsPostTarget
}

func (l *Love) GetCreatedAt() time.Time {
	return l.love.CreatedAt
}
