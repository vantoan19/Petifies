package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/domain/repositories"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

type redisCommentCacheRepository struct {
	client *redis.Client
}

func NewRedisCommentCacheRepository(client *redis.Client) repositories.CommentCacheRepository {
	return &redisCommentCacheRepository{client: client}
}

func (r *redisCommentCacheRepository) GetCommentContent(ctx context.Context, commentID uuid.UUID) (*models.Comment, error) {
	key := fmt.Sprintf("comment:%s:content", commentID.String())
	commentContentStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var commentContent models.Comment
	err = json.Unmarshal([]byte(commentContentStr), &commentContent)
	if err != nil {
		return nil, err
	}

	return &commentContent, nil
}

func (r *redisCommentCacheRepository) SetCommentContent(ctx context.Context, commentID uuid.UUID, comment models.Comment) error {
	key := fmt.Sprintf("comment:%s:content", commentID.String())

	commentContentStr, err := json.Marshal(comment)
	if err != nil {
		return err
	}
	err = r.client.Set(ctx, key, commentContentStr, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisCommentCacheRepository) ExistsCommentContent(ctx context.Context, commentID uuid.UUID) (bool, error) {
	key := fmt.Sprintf("comment:%s:content", commentID.String())
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}
