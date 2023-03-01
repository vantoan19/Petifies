package mapper

import (
	postaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/post"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/valueobjects"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/infra/db/models"
)

func DbPostToEntityPost(p *models.Post) *entities.Post {
	images := make([]valueobjects.ImageContent, 0)
	videos := make([]valueobjects.VideoContent, 0)

	for _, image := range p.Images {
		images = append(images, valueobjects.NewImageContent(image.URL, image.Description))
	}
	for _, video := range p.Videos {
		videos = append(videos, valueobjects.NewVideoContent(video.URL, video.Description))
	}

	return &entities.Post{
		ID:          p.ID,
		AuthorID:    p.AuthorID,
		TextContent: valueobjects.NewTextContent(p.TextContent),
		Images:      images,
		Videos:      videos,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func DbModelsToPostAggregate(p *models.Post, ls *[]models.Love, cs *[]models.Comment) (*postaggre.Post, error) {
	post := &postaggre.Post{}

	if err := post.SetPostEntity(*DbPostToEntityPost(p)); err != nil {
		return nil, err
	}

	for _, l := range *ls {
		if err := post.AddLoveByEntity(*DbLoveToEntityLove(&l)); err != nil {
			return nil, err
		}
	}
	for _, c := range *cs {
		if err := post.AddCommentByEntity(*DbCommentToEntityComment(&c)); err != nil {
			return nil, err
		}
	}

	return post, nil
}
