package petifiesproposalmongo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/libs/dbutils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/petifies-service/cmd"
	petifiesproposalaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_proposal"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/db/mappers"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var logger = logging.New("PetifiesService.MongoPetifiesProposalRepository")

var (
	ErrProposalNotExist = status.Errorf(codes.NotFound, "proposal does not exist")
)

type petifiesProposalMongoRepository struct {
	client             *mongo.Client
	proposalCollection *mongo.Collection
}

func New(client *mongo.Client) petifiesproposalaggre.PetifiesProposalRepository {
	return &petifiesProposalMongoRepository{
		client:             client,
		proposalCollection: client.Database(cmd.Conf.DatabaseName).Collection("petifies_proposals"),
	}
}

func (pr *petifiesProposalMongoRepository) GetBySessionAndUserID(ctx context.Context, sessionID, userID uuid.UUID) (*petifiesproposalaggre.PetifiesProposalAggre, error) {
	var proposal *petifiesproposalaggre.PetifiesProposalAggre

	err := dbutils.ExecWithSession(ctx, pr.client, func(ssCtx mongo.SessionContext) error {
		proposal_, err := pr.GetBySessionAndUserIDWithSession(ssCtx, sessionID, userID)
		if err != nil {
			return err
		}
		proposal = proposal_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return proposal, nil
}

func (pr *petifiesProposalMongoRepository) GetByID(ctx context.Context, id uuid.UUID) (*petifiesproposalaggre.PetifiesProposalAggre, error) {
	var proposal *petifiesproposalaggre.PetifiesProposalAggre

	err := dbutils.ExecWithSession(ctx, pr.client, func(ssCtx mongo.SessionContext) error {
		proposal_, err := pr.GetByIDWithSession(ssCtx, id)
		if err != nil {
			return err
		}
		proposal = proposal_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return proposal, nil
}

func (pr *petifiesProposalMongoRepository) GetBySessionID(ctx context.Context, sessionID uuid.UUID, pageSize int, afterID uuid.UUID) ([]*petifiesproposalaggre.PetifiesProposalAggre, error) {
	var proposals []*petifiesproposalaggre.PetifiesProposalAggre

	err := dbutils.ExecWithSession(ctx, pr.client, func(ssCtx mongo.SessionContext) error {
		proposals_, err := pr.GetBySessionIDWithSession(ssCtx, sessionID, pageSize, afterID)
		if err != nil {
			return err
		}
		proposals = proposals_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return proposals, nil
}

func (pr *petifiesProposalMongoRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*petifiesproposalaggre.PetifiesProposalAggre, error) {
	var proposals []*petifiesproposalaggre.PetifiesProposalAggre

	err := dbutils.ExecWithSession(ctx, pr.client, func(ssCtx mongo.SessionContext) error {
		proposals_, err := pr.GetByUserIDWithSession(ssCtx, userID)
		if err != nil {
			return err
		}
		proposals = proposals_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return proposals, nil
}

func (pr *petifiesProposalMongoRepository) ListByIDs(ctx context.Context, ids []uuid.UUID) ([]*petifiesproposalaggre.PetifiesProposalAggre, error) {
	var proposals []*petifiesproposalaggre.PetifiesProposalAggre

	err := dbutils.ExecWithSession(ctx, pr.client, func(ssCtx mongo.SessionContext) error {
		proposals_, err := pr.ListByIDsWithSession(ssCtx, ids)
		if err != nil {
			return err
		}
		proposals = proposals_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return proposals, nil
}

func (pr *petifiesProposalMongoRepository) ExistsBySessionAndUserID(ctx context.Context, sessionID, userID uuid.UUID) (bool, error) {
	var exists bool

	err := dbutils.ExecWithSession(ctx, pr.client, func(ssCtx mongo.SessionContext) error {
		exists_, err := pr.ExistsBySessionAndUserIDWithSession(ssCtx, sessionID, userID)
		if err != nil {
			return err
		}
		exists = exists_
		return nil
	})
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (pr *petifiesProposalMongoRepository) Save(ctx context.Context, proposal petifiesproposalaggre.PetifiesProposalAggre) (*petifiesproposalaggre.PetifiesProposalAggre, error) {
	var savedProposal *petifiesproposalaggre.PetifiesProposalAggre

	err := dbutils.ExecWithSession(ctx, pr.client, func(ssCtx mongo.SessionContext) error {
		proposal_, err := pr.SaveWithSession(ssCtx, proposal)
		if err != nil {
			return err
		}
		savedProposal = proposal_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return savedProposal, nil
}

func (pr *petifiesProposalMongoRepository) Update(ctx context.Context, proposal petifiesproposalaggre.PetifiesProposalAggre) (*petifiesproposalaggre.PetifiesProposalAggre, error) {
	var updatedProposal *petifiesproposalaggre.PetifiesProposalAggre

	err := dbutils.ExecWithSession(ctx, pr.client, func(ssCtx mongo.SessionContext) error {
		proposal_, err := pr.UpdateWithSession(ssCtx, proposal)
		if err != nil {
			return err
		}
		updatedProposal = proposal_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return updatedProposal, nil
}

func (pr *petifiesProposalMongoRepository) GetBySessionAndUserIDWithSession(ctx context.Context, sessionID, userID uuid.UUID) (*petifiesproposalaggre.PetifiesProposalAggre, error) {
	var proposal models.PetifiesProposal
	err := pr.proposalCollection.FindOne(ctx, bson.D{{Key: "petifies_session_id", Value: sessionID}, {Key: "user_id", Value: userID}}).Decode(&proposal)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrProposalNotExist
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result, err := mappers.DbModelToPetifiesProposalAggregate(&proposal)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (pr *petifiesProposalMongoRepository) GetByIDWithSession(ctx context.Context, id uuid.UUID) (*petifiesproposalaggre.PetifiesProposalAggre, error) {
	var proposal models.PetifiesProposal
	err := pr.proposalCollection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&proposal)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrProposalNotExist
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result, err := mappers.DbModelToPetifiesProposalAggregate(&proposal)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (pr *petifiesProposalMongoRepository) GetBySessionIDWithSession(ctx context.Context, sessionID uuid.UUID, pageSize int, afterID uuid.UUID) ([]*petifiesproposalaggre.PetifiesProposalAggre, error) {
	var createdAt time.Time

	if afterID == uuid.Nil {
		createdAt = time.Now()
	} else {
		petifies, err := pr.GetByIDWithSession(ctx, afterID)
		if err != nil {
			return nil, err
		}
		createdAt = petifies.GetCreatedAt()
	}

	var result []*petifiesproposalaggre.PetifiesProposalAggre
	var proposals []models.PetifiesProposal
	filter := bson.D{{Key: "petifies_session_id", Value: sessionID}, {Key: "created_at", Value: bson.D{{Key: "$lt", Value: createdAt}}}}
	opts := options.Find().SetLimit(int64(pageSize)).SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := pr.proposalCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if err := cursor.All(ctx, &proposals); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	for _, p := range proposals {
		proposal, err := mappers.DbModelToPetifiesProposalAggregate(&p)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		result = append(result, proposal)
	}

	return result, nil
}

func (pr *petifiesProposalMongoRepository) GetByUserIDWithSession(ctx context.Context, userID uuid.UUID) ([]*petifiesproposalaggre.PetifiesProposalAggre, error) {
	var result []*petifiesproposalaggre.PetifiesProposalAggre
	var proposals []models.PetifiesProposal
	cursor, err := pr.proposalCollection.Find(
		ctx,
		bson.D{{Key: "user_id", Value: userID}},
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if err := cursor.All(ctx, &proposals); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	for _, p := range proposals {
		proposal, err := mappers.DbModelToPetifiesProposalAggregate(&p)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		result = append(result, proposal)
	}

	return result, nil
}

func (pr *petifiesProposalMongoRepository) ListByIDsWithSession(ctx context.Context, ids []uuid.UUID) ([]*petifiesproposalaggre.PetifiesProposalAggre, error) {
	var results []*petifiesproposalaggre.PetifiesProposalAggre

	cursor, err := pr.proposalCollection.Find(ctx, bson.D{{Key: "id", Value: bson.D{{Key: "$in", Value: ids}}}})
	if err != nil {
		return nil, err
	}
	var proposalsDB []models.PetifiesProposal
	if err := cursor.All(ctx, proposalsDB); err != nil {
		return nil, err
	}

	for _, p := range proposalsDB {
		proposal, err := mappers.DbModelToPetifiesProposalAggregate(&p)
		if err != nil {
			return nil, err
		}
		results = append(results, proposal)
	}

	return results, nil
}

func (pr *petifiesProposalMongoRepository) ExistsBySessionAndUserIDWithSession(ctx context.Context, sessionID, userID uuid.UUID) (bool, error) {
	count, err := pr.proposalCollection.CountDocuments(ctx, bson.D{{Key: "petifies_session_id", Value: sessionID}, {Key: "user_id", Value: userID}})
	if err != nil {
		return false, err
	}

	return (count == 1), nil
}

func (pr *petifiesProposalMongoRepository) SaveWithSession(ctx context.Context, proposal petifiesproposalaggre.PetifiesProposalAggre) (*petifiesproposalaggre.PetifiesProposalAggre, error) {
	proposalDB := mappers.AggregatePetifiesProposalToDbPetifiesProposal(&proposal)
	_, err := pr.proposalCollection.InsertOne(ctx, proposalDB)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	savedProposal, err := pr.GetByIDWithSession(ctx, proposal.GetID())
	if err != nil {
		return nil, err
	}
	return savedProposal, nil
}

func (pr *petifiesProposalMongoRepository) UpdateWithSession(ctx context.Context, proposal petifiesproposalaggre.PetifiesProposalAggre) (*petifiesproposalaggre.PetifiesProposalAggre, error) {
	proposalDB := mappers.AggregatePetifiesProposalToDbPetifiesProposal(&proposal)
	_, err := pr.proposalCollection.ReplaceOne(ctx, bson.D{{Key: "id", Value: proposal.GetID()}}, proposalDB)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	updatedProposal, err := pr.GetByIDWithSession(ctx, proposal.GetID())
	if err != nil {
		return nil, err
	}
	return updatedProposal, nil
}
