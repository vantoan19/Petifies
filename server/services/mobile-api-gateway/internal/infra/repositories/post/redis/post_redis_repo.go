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

type redisPostCacheRepository struct {
	client *redis.Client
}

func NewRedisPostCacheRepository(client *redis.Client) repositories.PostCacheRepository {
	return &redisPostCacheRepository{client: client}
}

func (r *redisPostCacheRepository) GetPostContent(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	key := fmt.Sprintf("post:%s:content", postID.String())
	postContentStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var postContent models.Post
	err = json.Unmarshal([]byte(postContentStr), &postContent)
	if err != nil {
		return nil, err
	}

	return &postContent, nil
}

func (r *redisPostCacheRepository) SetPostContent(ctx context.Context, postID uuid.UUID, post models.Post) error {
	key := fmt.Sprintf("post:%s:content", postID.String())

	postContentStr, err := json.Marshal(post)
	if err != nil {
		return err
	}
	err = r.client.Set(ctx, key, postContentStr, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisPostCacheRepository) ExistsPostContent(ctx context.Context, postID uuid.UUID) (bool, error) {
	key := fmt.Sprintf("post:%s:content", postID.String())
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}
