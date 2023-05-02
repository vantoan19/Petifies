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

type redisPetifiesSessionCacheRepository struct {
	client         *redis.Client
	petifiesClient petifiesclient.PetifiesClient
}

func NewRedisPetifiesSessionCacheRepository(client *redis.Client, petifiesClient petifiesclient.PetifiesClient) repositories.PetifiesSessionCacheRepository {
	return &redisPetifiesSessionCacheRepository{client: client, petifiesClient: petifiesClient}
}

func (r *redisPetifiesSessionCacheRepository) GetPetifiesSession(ctx context.Context, sessionId uuid.UUID) (*models.PetifiesSession, error) {
	key := fmt.Sprintf("petifiessession:%s", sessionId.String())
	sessionStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		session, err := r.petifiesClient.GetPetifiesSessionById(ctx, &models.GetSessionByIdReq{
			ID: sessionId,
		})
		if err != nil {
			return nil, err
		}

		go func() {
			r.SetPetifiesSession(ctx, sessionId, session)
		}()

		return session, err
	} else if err != nil {
		return nil, err
	} else {
		var session models.PetifiesSession
		err = json.Unmarshal([]byte(sessionStr), &session)
		if err != nil {
			return nil, err
		}
		return &session, nil
	}
}

func (r *redisPetifiesSessionCacheRepository) SetPetifiesSession(ctx context.Context, sessionId uuid.UUID, petifiesSession *models.PetifiesSession) error {
	key := fmt.Sprintf("petifiessession:%s", sessionId.String())

	sessionStr, err := json.Marshal(petifiesSession)
	if err != nil {
		return err
	}

	tx := r.client.TxPipeline()

	tx.Set(ctx, key, sessionStr, 0)

	_, err = tx.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
