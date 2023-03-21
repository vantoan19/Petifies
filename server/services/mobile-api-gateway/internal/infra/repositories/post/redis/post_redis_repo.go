package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	err = r.client.Set(ctx, key, postContentStr, time.Minute).Err()
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

func (r *redisPostCacheRepository) GetPostLoveCount(ctx context.Context, postID uuid.UUID) (int, error) {
	key := fmt.Sprintf("post:%s:loveCount", postID.String())
	count, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return -1, nil
	}

	return count, nil
}

func (r *redisPostCacheRepository) SetPostLoveCount(ctx context.Context, postID uuid.UUID, loveCount int) error {
	key := fmt.Sprintf("post:%s:loveCount", postID.String())
	err := r.client.Set(ctx, key, loveCount, time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisPostCacheRepository) ExistsPostLoveCount(ctx context.Context, postID uuid.UUID) (bool, error) {
	key := fmt.Sprintf("post:%s:loveCount", postID.String())
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func (r *redisPostCacheRepository) GetPostCommentCount(ctx context.Context, postID uuid.UUID) (int, error) {
	key := fmt.Sprintf("post:%s:commentCount", postID.String())
	commentCount, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return -1, nil
	}

	return commentCount, nil
}

func (r *redisPostCacheRepository) SetPostCommentCount(ctx context.Context, postID uuid.UUID, commentCount int) error {
	key := fmt.Sprintf("post:%s:commentCount", postID.String())
	err := r.client.Set(ctx, key, commentCount, time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisPostCacheRepository) ExistsPostCommentCount(ctx context.Context, postID uuid.UUID) (bool, error) {
	key := fmt.Sprintf("post:%s:commentCount", postID.String())
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}
