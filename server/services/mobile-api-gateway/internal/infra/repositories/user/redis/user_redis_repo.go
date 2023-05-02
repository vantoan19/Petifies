package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	userclient "github.com/vantoan19/Petifies/server/services/grpc-clients/user-client"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/domain/repositories"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/models"
)

type redisUserCacheRepository struct {
	client     *redis.Client
	userClient userclient.UserClient
}

func NewRedisUserCacheRepository(client *redis.Client, userClient userclient.UserClient) repositories.UserCacheRepository {
	return &redisUserCacheRepository{client: client, userClient: userClient}
}

func (r *redisUserCacheRepository) GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	key := fmt.Sprintf("user:%s", userID.String())
	userStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		user, err := r.userClient.GetUser(ctx, userID)
		if err != nil {
			return nil, err
		}

		go func() {
			r.SetUser(ctx, userID, *user)
		}()

		return user, nil
	} else if err != nil {
		return nil, err
	} else {
		var user models.User
		err = json.Unmarshal([]byte(userStr), &user)
		if err != nil {
			return nil, err
		}

		return &user, nil
	}
}

func (r *redisUserCacheRepository) ListUsers(ctx context.Context, userIds []uuid.UUID) ([]*models.User, error) {
	if len(userIds) == 0 {
		return []*models.User{}, nil
	}

	var keys []string
	for _, id := range userIds {
		keys = append(keys, fmt.Sprintf("user:%s", id.String()))
	}

	userStrs, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	results := make(map[string]*models.User)
	usersToList := make([]uuid.UUID, 0)
	for i, key := range keys {
		userStr := userStrs[i]
		if userStr == nil {
			usersToList = append(usersToList, userIds[i])
		} else {
			var user models.User
			err = json.Unmarshal([]byte(userStr.(string)), &user)
			if err != nil {
				return nil, err
			}
			results[key] = &user
		}
	}

	if len(usersToList) > 0 {
		users, err := r.userClient.ListUsersByIds(ctx, usersToList)
		if err != nil {
			return nil, err
		}

		for _, u := range users {
			key := fmt.Sprintf("user:%s", u.ID.String())
			results[key] = u
		}

		go func() {
			usersToCache := make(map[string]string)
			for _, p := range users {
				key := fmt.Sprintf("user:%s", p.ID.String())
				userStr, err := json.Marshal(p)
				if err != nil {
					return
				}
				usersToCache[key] = string(userStr)
			}
			tx := r.client.TxPipeline()
			tx.MSet(ctx, usersToCache)
			_, err = tx.Exec(ctx)
			if err != nil {
				return
			}
		}()
	}

	finalResults := make([]*models.User, 0)
	for _, id := range userIds {
		key := fmt.Sprintf("user:%s", id.String())
		if val, ok := results[key]; ok {
			finalResults = append(finalResults, val)
		}
	}

	return finalResults, nil
}

func (r *redisUserCacheRepository) SetUser(ctx context.Context, userID uuid.UUID, user models.User) error {
	key := fmt.Sprintf("user:%s", userID.String())

	userStr, err := json.Marshal(user)
	if err != nil {
		return err
	}

	tx := r.client.TxPipeline()

	tx.Set(ctx, key, userStr, 0)

	_, err = tx.Exec(ctx)
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
