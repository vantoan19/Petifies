package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	postclient "github.com/vantoan19/Petifies/server/services/grpc-clients/post-client"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/domain/repositories"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

type redisLoveCacheRepository struct {
	client     *redis.Client
	postClient postclient.PostClient
}

func NewRedisLoveCacheRepository(client *redis.Client, postClient postclient.PostClient) repositories.LoveCacheRepository {
	return &redisLoveCacheRepository{client: client, postClient: postClient}
}

func (r *redisLoveCacheRepository) GetLove(ctx context.Context, authorID, targetID uuid.UUID) (*models.Love, error) {
	key := fmt.Sprintf("love:%s_%s", authorID.String(), targetID.String())
	loveStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		love, err := r.postClient.GetLove(ctx, &models.GetLoveReq{
			AuthorID: authorID,
			TargetID: targetID,
		})
		if err != nil {
			return nil, err
		}

		go func() {
			r.SetLove(ctx, authorID, targetID, *love)
		}()

		return love, err
	} else if err != nil {
		return nil, err
	} else {
		var love models.Love
		err = json.Unmarshal([]byte(loveStr), &love)
		if err != nil {
			return nil, err
		}
		return &love, nil
	}
}

func (r *redisLoveCacheRepository) SetLove(ctx context.Context, authorID, targetID uuid.UUID, love models.Love) error {
	key := fmt.Sprintf("love:%s_%s", authorID.String(), targetID.String())

	loveStr, err := json.Marshal(love)
	if err != nil {
		return err
	}

	tx := r.client.TxPipeline()

	tx.Set(ctx, key, loveStr, 0)

	_, err = tx.Exec(ctx)
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
