package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/domain/repositories"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/models"
)

type redisUserCacheRepository struct {
	client *redis.Client
}

func NewRedisUserCacheRepository(client *redis.Client) repositories.UserCacheRepository {
	return &redisUserCacheRepository{client: client}
}

func (r *redisUserCacheRepository) GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	key := fmt.Sprintf("user:%s", userID.String())
	userStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user models.User
	err = json.Unmarshal([]byte(userStr), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *redisUserCacheRepository) SetUser(ctx context.Context, userID uuid.UUID, user models.User) error {
	key := fmt.Sprintf("user:%s", userID.String())

	userStr, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = r.client.Set(ctx, key, userStr, time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisUserCacheRepository) ExistsUser(ctx context.Context, userID uuid.UUID) (bool, error) {
	key := fmt.Sprintf("user:%s", userID.String())
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}
