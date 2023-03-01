package mapper

import (
	commentaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/comment"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/valueobjects"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/infra/db/models"
)

func DbCommentToEntityComment(c *models.Comment) *entities.Comment {
	return &entities.Comment{
		ID:           c.ID,
		PostID:       c.PostID,
		AuthorID:     c.AuthorID,
		ParentID:     c.ParentID,
		IsPostParent: c.IsPostParent,
		Content:      valueobjects.NewTextContent(c.Content),
		ImageContent: valueobjects.NewImageContent(c.ImageContent.URL, c.ImageContent.Description),
		VideoContent: valueobjects.NewVideoContent(c.VideoContent.URL, c.VideoContent.Description),
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

func EntityCommentToDbComment(c *entities.Comment) *models.Comment {
	return &models.Comment{
		ID:           c.ID,
		PostID:       c.PostID,
		AuthorID:     c.AuthorID,
		ParentID:     c.ParentID,
		IsPostParent: c.IsPostParent,
		Content:      c.Content.Content(),
		ImageContent: models.Image{
			URL:         c.ImageContent.URL(),
			Description: c.ImageContent.Description(),
		},
		VideoContent: models.Video{
			URL:         c.VideoContent.URL(),
			Description: c.VideoContent.Description(),
		},
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func DbModelsToCommentAggregate(c *models.Comment, ls *[]models.Love, subcs *[]models.Comment) (*commentaggre.Comment, error) {
	comment := &commentaggre.Comment{}

	if err := comment.SetCommentEntity(*DbCommentToEntityComment(c)); err != nil {
		return nil, err
	}

	for _, l := range *ls {
		if err := comment.AddLoveByEntity(*DbLoveToEntityLove(&l)); err != nil {
			return nil, err
		}
	}
	for _, subc := range *subcs {
		if err := comment.AddSubcommentByEntity(*DbCommentToEntityComment(&subc)); err != nil {
			return nil, err
		}
	}

	return comment, nil
}
