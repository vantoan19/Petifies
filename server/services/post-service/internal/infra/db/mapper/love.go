package mapper

import (
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/infra/db/models"
)

func DbLoveToEntityLove(l *models.Love) *entities.Love {
	return &entities.Love{
		ID:        l.ID,
		PostID:    l.PostID,
		CommentID: l.CommentID,
		AuthorID:  l.AuthorID,
		CreatedAt: l.CreatedAt,
	}
}

func EntityLoveToDbLove(l *entities.Love) *models.Love {
	return &models.Love{
		ID:        l.ID,
		PostID:    l.PostID,
		CommentID: l.CommentID,
		AuthorID:  l.AuthorID,
		CreatedAt: l.CreatedAt,
	}
}
