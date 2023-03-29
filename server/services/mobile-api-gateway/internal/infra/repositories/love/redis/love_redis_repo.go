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

type redisLoveCacheRepository struct {
	client *redis.Client
}

func NewRedisLoveCacheRepository(client *redis.Client) repositories.LoveCacheRepository {
	return &redisLoveCacheRepository{client: client}
}

func (r *redisLoveCacheRepository) GetLove(ctx context.Context, authorID, targetID uuid.UUID) (*models.Love, error) {
	key := fmt.Sprintf("love:%s_%s", authorID.String(), targetID.String())
	loveStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var love models.Love
	err = json.Unmarshal([]byte(loveStr), &love)
	if err != nil {
		return nil, err
	}

	return &love, nil
}

func (r *redisLoveCacheRepository) SetLove(ctx context.Context, authorID, targetID uuid.UUID, love models.Love) error {
	key := fmt.Sprintf("love:%s_%s", authorID.String(), targetID.String())

	loveStr, err := json.Marshal(love)
	if err != nil {
		return err
	}
	err = r.client.Set(ctx, key, loveStr, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisLoveCacheRepository) RemoveLove(ctx context.Context, authorID, targetID uuid.UUID) error {
	if exists, err := r.ExistsLove(ctx, authorID, targetID); err != nil {
		return err
	} else if exists {
		key := fmt.Sprintf("love:%s_%s", authorID.String(), targetID.String())
		err := r.client.Del(ctx, key).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *redisLoveCacheRepository) ExistsLove(ctx context.Context, authorID, targetID uuid.UUID) (bool, error) {
	key := fmt.Sprintf("love:%s_%s", authorID.String(), targetID.String())
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}
