package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/domain/repositories"
	"github.com/vantoan19/Petifies/server/services/relationship-service/pkg/models"
)

type redisRelationshipCacheRepository struct {
	client *redis.Client
}

func NewRedisRelationshipCacheRepository(client *redis.Client) repositories.RelationshipCacheRepository {
	return &redisRelationshipCacheRepository{client: client}
}

func (r *redisRelationshipCacheRepository) GetFollowingsInfo(ctx context.Context, userID uuid.UUID) (*models.ListFollowingsResp, error) {
	key := fmt.Sprintf("followings:%s", userID.String())
	followingsStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var followings models.ListFollowingsResp
	err = json.Unmarshal([]byte(followingsStr), &followings)
	if err != nil {
		return nil, err
	}

	return &followings, nil
}

func (r *redisRelationshipCacheRepository) SetFollowingsInfo(ctx context.Context, userID uuid.UUID, followings *models.ListFollowingsResp) error {
	key := fmt.Sprintf("followings:%s", userID.String())

	followingsStr, err := json.Marshal(followings)
	if err != nil {
		return err
	}
	err = r.client.Set(ctx, key, followingsStr, time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisRelationshipCacheRepository) ExistsFollowingsInfo(ctx context.Context, userID uuid.UUID) (bool, error) {
	key := fmt.Sprintf("followings:%s", userID.String())
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func (r *redisRelationshipCacheRepository) GetFollowersInfo(ctx context.Context, userID uuid.UUID) (*models.ListFollowersResp, error) {
	key := fmt.Sprintf("followers:%s", userID.String())
	followersStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var followers models.ListFollowersResp
	err = json.Unmarshal([]byte(followersStr), &followers)
	if err != nil {
		return nil, err
	}

	return &followers, nil
}

func (r *redisRelationshipCacheRepository) SetFollowersInfo(ctx context.Context, userID uuid.UUID, followers *models.ListFollowersResp) error {
	key := fmt.Sprintf("followers:%s", userID.String())

	followersStr, err := json.Marshal(followers)
	if err != nil {
		return err
	}
	err = r.client.Set(ctx, key, followersStr, time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisRelationshipCacheRepository) ExistsFollowersInfo(ctx context.Context, userID uuid.UUID) (bool, error) {
	key := fmt.Sprintf("followers:%s", userID.String())
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}
