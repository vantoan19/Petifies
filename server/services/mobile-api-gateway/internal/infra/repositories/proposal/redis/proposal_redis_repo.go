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

type redisPetifiesProposalCacheRepository struct {
	client         *redis.Client
	petifiesClient petifiesclient.PetifiesClient
}

func NewRedisPetifiesProposalCacheRepository(client *redis.Client, petifiesClient petifiesclient.PetifiesClient) repositories.PetifiesProposalCacheRepository {
	return &redisPetifiesProposalCacheRepository{client: client, petifiesClient: petifiesClient}
}

func (r *redisPetifiesProposalCacheRepository) GetPetifiesProposal(ctx context.Context, proposalId uuid.UUID) (*models.PetifiesProposal, error) {
	key := fmt.Sprintf("petifiesproposal:%s", proposalId.String())
	proposalStr, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		proposal, err := r.petifiesClient.GetPetifiesProposalById(ctx, &models.GetProposalByIdReq{
			ID: proposalId,
		})
		if err != nil {
			return nil, err
		}

		go func() {
			r.SetPetifiesProposal(ctx, proposalId, proposal)
		}()

		return proposal, err
	} else if err != nil {
		return nil, err
	} else {
		var proposal models.PetifiesProposal
		err = json.Unmarshal([]byte(proposalStr), &proposal)
		if err != nil {
			return nil, err
		}
		return &proposal, nil
	}
}

func (r *redisPetifiesProposalCacheRepository) SetPetifiesProposal(ctx context.Context, proposalId uuid.UUID, proposal *models.PetifiesProposal) error {
	key := fmt.Sprintf("petifiesproposal:%s", proposalId.String())

	proposalStr, err := json.Marshal(proposal)
	if err != nil {
		return err
	}

	tx := r.client.TxPipeline()

	tx.Set(ctx, key, proposalStr, 0)

	_, err = tx.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
