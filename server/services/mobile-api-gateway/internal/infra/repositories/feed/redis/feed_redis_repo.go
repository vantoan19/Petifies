package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/domain/repositories"
)

type PostFeeds struct {
	PostFeedIDs []uuid.UUID `json:"post_feed_ids"`
}

type redisFeedCacheRepository struct {
	client *redis.Client
}

func NewRedisFeedCacheRepository(client *redis.Client) repositories.FeedCacheRepository {
	return &redisFeedCacheRepository{client: client}
}

func (r *redisFeedCacheRepository) GetPostFeedIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	key := fmt.Sprintf("postfeeds:%s", userID.String())
	postFeedsStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var postFeeds PostFeeds
	err = json.Unmarshal([]byte(postFeedsStr), &postFeeds)
	if err != nil {
		return nil, err
	}

	return postFeeds.PostFeedIDs, nil
}

func (r *redisFeedCacheRepository) SetPostFeedIDs(ctx context.Context, userID uuid.UUID, feedIDs []uuid.UUID) error {
	key := fmt.Sprintf("postfeeds:%s", userID.String())

	postFeeds := PostFeeds{
		PostFeedIDs: feedIDs,
	}
	postFeedsStr, err := json.Marshal(postFeeds)
	if err != nil {
		return err
	}
	err = r.client.Set(ctx, key, postFeedsStr, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisFeedCacheRepository) ExistsPostFeedIDs(ctx context.Context, userID uuid.UUID) (bool, error) {
	key := fmt.Sprintf("postfeeds:%s", userID.String())
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil
}
