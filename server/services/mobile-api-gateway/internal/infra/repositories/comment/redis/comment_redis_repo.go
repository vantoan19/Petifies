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
	err = r.client.Set(ctx, key, commentContentStr, time.Minute).Err()
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

func (r *redisCommentCacheRepository) GetCommentLoveCount(ctx context.Context, commentID uuid.UUID) (int, error) {
	key := fmt.Sprintf("comment:%s:loveCount", commentID.String())
	count, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return -1, nil
	}

	return count, nil
}

func (r *redisCommentCacheRepository) SetCommentLoveCount(ctx context.Context, commentID uuid.UUID, loveCount int) error {
	key := fmt.Sprintf("comment:%s:loveCount", commentID.String())
	err := r.client.Set(ctx, key, loveCount, time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisCommentCacheRepository) ExistsCommentLoveCount(ctx context.Context, commentID uuid.UUID) (bool, error) {
	key := fmt.Sprintf("comment:%s:loveCount", commentID.String())
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func (r *redisCommentCacheRepository) GetCommentSubCommentCount(ctx context.Context, commentID uuid.UUID) (int, error) {
	key := fmt.Sprintf("comment:%s:subcommentCount", commentID.String())
	commentCount, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return -1, nil
	}

	return commentCount, nil
}

func (r *redisCommentCacheRepository) SetCommentSubCommentCount(ctx context.Context, commentID uuid.UUID, commentCount int) error {
	key := fmt.Sprintf("comment:%s:subcommentCount", commentID.String())
	err := r.client.Set(ctx, key, commentCount, time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisCommentCacheRepository) ExistsCommentSubCommentCount(ctx context.Context, commentID uuid.UUID) (bool, error) {
	key := fmt.Sprintf("comment:%s:subcommentCount", commentID.String())
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}
