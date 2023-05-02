package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	relationshipclient "github.com/vantoan19/Petifies/server/services/grpc-clients/relationship-client"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/domain/repositories"
	"github.com/vantoan19/Petifies/server/services/relationship-service/pkg/models"
)

type redisRelationshipCacheRepository struct {
	client             *redis.Client
	relationshipClient relationshipclient.RelationshipClient
}

func NewRedisRelationshipCacheRepository(client *redis.Client, relationshipClient relationshipclient.RelationshipClient) repositories.RelationshipCacheRepository {
	return &redisRelationshipCacheRepository{client: client, relationshipClient: relationshipClient}
}

func (r *redisRelationshipCacheRepository) GetFollowingsInfo(ctx context.Context, userID uuid.UUID) (*models.ListFollowingsResp, error) {
	key := fmt.Sprintf("followings:%s", userID.String())
	followingsStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		followings, err := r.relationshipClient.ListFollowings(ctx, userID)
		if err != nil {
			return nil, err
		}

		go func() {
			r.SetFollowingsInfo(ctx, userID, followings)
		}()

		return followings, nil
	} else if err != nil {
		return nil, err
	} else {
		var followings models.ListFollowingsResp
		err = json.Unmarshal([]byte(followingsStr), &followings)
		if err != nil {
			return nil, err
		}

		return &followings, nil
	}
}

func (r *redisRelationshipCacheRepository) SetFollowingsInfo(ctx context.Context, userID uuid.UUID, followings *models.ListFollowingsResp) error {
	key := fmt.Sprintf("followings:%s", userID.String())

	followingsStr, err := json.Marshal(followings)
	if err != nil {
		return err
	}

	tx := r.client.TxPipeline()

	tx.Set(ctx, key, followingsStr, 0)

	_, err = tx.Exec(ctx)
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
		followers, err := r.relationshipClient.ListFollowers(ctx, userID)
		if err != nil {
			return nil, err
		}

		go func() {
			r.SetFollowersInfo(ctx, userID, followers)
		}()

		return followers, nil
	} else if err != nil {
		return nil, err
	} else {
		var followers models.ListFollowersResp
		err = json.Unmarshal([]byte(followersStr), &followers)
		if err != nil {
			return nil, err
		}

		return &followers, nil
	}

}

func (r *redisRelationshipCacheRepository) SetFollowersInfo(ctx context.Context, userID uuid.UUID, followers *models.ListFollowersResp) error {
	key := fmt.Sprintf("followers:%s", userID.String())

	followersStr, err := json.Marshal(followers)
	if err != nil {
		return err
	}

	tx := r.client.TxPipeline()

	tx.Set(ctx, key, followersStr, 0)

	_, err = tx.Exec(ctx)
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
