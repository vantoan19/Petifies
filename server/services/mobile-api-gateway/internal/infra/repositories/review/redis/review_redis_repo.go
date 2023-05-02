package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	petifiesclient "github.com/vantoan19/Petifies/server/services/grpc-clients/petifies-client"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/domain/repositories"
	"github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
)

type redisReviewCacheRepository struct {
	client         *redis.Client
	petifiesClient petifiesclient.PetifiesClient
}

func NewRedisReviewCacheRepository(client *redis.Client, petifiesClient petifiesclient.PetifiesClient) repositories.ReviewCacheRepository {
	return &redisReviewCacheRepository{client: client, petifiesClient: petifiesClient}
}

func (r *redisReviewCacheRepository) GetReview(ctx context.Context, reviewId uuid.UUID) (*models.Review, error) {
	key := fmt.Sprintf("review:%s", reviewId.String())
	reviewStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		review, err := r.petifiesClient.GetReviewById(ctx, &models.GetReviewByIdReq{
			ID: reviewId,
		})
		if err != nil {
			return nil, err
		}

		go func() {
			r.SetReview(ctx, reviewId, review)
		}()

		return review, err
	} else if err != nil {
		return nil, err
	} else {
		var review models.Review
		err = json.Unmarshal([]byte(reviewStr), &review)
		if err != nil {
			return nil, err
		}
		return &review, nil
	}
}

func (r *redisReviewCacheRepository) SetReview(ctx context.Context, reviewId uuid.UUID, review *models.Review) error {
	key := fmt.Sprintf("review:%s", reviewId.String())

	reviewStr, err := json.Marshal(review)
	if err != nil {
		return err
	}

	tx := r.client.TxPipeline()

	tx.Set(ctx, key, reviewStr, 0)

	_, err = tx.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
