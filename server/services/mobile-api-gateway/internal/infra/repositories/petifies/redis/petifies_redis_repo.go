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

type redisPetifiesCacheRepository struct {
	client         *redis.Client
	petifiesClient petifiesclient.PetifiesClient
}

func NewRedisPetifiesCacheRepository(client *redis.Client, petifiesClient petifiesclient.PetifiesClient) repositories.PetifiesCacheRepository {
	return &redisPetifiesCacheRepository{client: client, petifiesClient: petifiesClient}
}

func (r *redisPetifiesCacheRepository) GetPetifies(ctx context.Context, petifiesId uuid.UUID) (*models.Petifies, error) {
	key := fmt.Sprintf("petifies:%s", petifiesId.String())
	petifiesStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		petifies, err := r.petifiesClient.GetPetifiesById(ctx, &models.GetPetifiesByIdReq{
			ID: petifiesId,
		})
		if err != nil {
			return nil, err
		}

		go func() {
			r.SetPetifies(ctx, petifiesId, petifies)
		}()

		return petifies, err
	} else if err != nil {
		return nil, err
	} else {
		var petifies models.Petifies
		err = json.Unmarshal([]byte(petifiesStr), &petifies)
		if err != nil {
			return nil, err
		}
		return &petifies, nil
	}
}

func (r *redisPetifiesCacheRepository) ListPetifies(ctx context.Context, petifiesIds []uuid.UUID) ([]*models.Petifies, error) {
	if len(petifiesIds) == 0 {
		return []*models.Petifies{}, nil
	}

	var keys []string
	for _, id := range petifiesIds {
		keys = append(keys, fmt.Sprintf("petifies:%s", id.String()))
	}

	petifiesStrs, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	results := make(map[string]*models.Petifies)
	petifiesToList := make([]uuid.UUID, 0)
	for i, key := range keys {
		petifiesStr := petifiesStrs[i]
		if petifiesStr == nil {
			petifiesToList = append(petifiesToList, petifiesIds[i])
		} else {
			var petifies models.Petifies
			err = json.Unmarshal([]byte(petifiesStr.(string)), &petifies)
			if err != nil {
				return nil, err
			}
			results[key] = &petifies
		}
	}

	if len(petifiesToList) > 0 {
		petifiesResp, err := r.petifiesClient.ListPetifiesByIds(ctx, &models.ListPetifiesByIdsReq{
			PetifiesIDs: petifiesToList,
		})
		if err != nil {
			return nil, err
		}

		for _, p := range petifiesResp.Petifies {
			key := fmt.Sprintf("petifies:%s", p.ID.String())
			results[key] = p
		}

		go func() {
			petifiesToCache := make(map[string]string)
			for _, p := range petifiesResp.Petifies {
				key := fmt.Sprintf("petifies:%s", p.ID.String())
				petifiesStr, err := json.Marshal(p)
				if err != nil {
					return
				}
				petifiesToCache[key] = string(petifiesStr)
			}
			tx := r.client.TxPipeline()
			tx.MSet(ctx, petifiesToCache)
			_, err = tx.Exec(ctx)
			if err != nil {
				return
			}
		}()
	}

	finalResults := make([]*models.Petifies, 0)
	for _, id := range petifiesIds {
		key := fmt.Sprintf("petifies:%s", id.String())
		if val, ok := results[key]; ok {
			finalResults = append(finalResults, val)
		}
	}

	return finalResults, nil
}

func (r *redisPetifiesCacheRepository) SetPetifies(ctx context.Context, petifiesId uuid.UUID, petifies *models.Petifies) error {
	key := fmt.Sprintf("petifies:%s", petifiesId.String())

	petifiesStr, err := json.Marshal(petifies)
	if err != nil {
		return err
	}

	tx := r.client.TxPipeline()

	tx.Set(ctx, key, petifiesStr, 0)

	_, err = tx.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
