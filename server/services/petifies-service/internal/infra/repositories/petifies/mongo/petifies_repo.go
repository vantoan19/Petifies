package petifiesmongo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vantoan19/Petifies/server/libs/dbutils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/petifies-service/cmd"
	petifiesaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/db/mappers"
	"github.com/vantoan19/Petifies/server/services/petifies-service/internal/infra/db/models"
)

var logger = logging.New("PetifiesService.MongoPetifiesRepository")

var (
	ErrPetifiesNotExist = status.Errorf(codes.NotFound, "petifies does not exist")
)

type petifiesMongoRepository struct {
	client             *mongo.Client
	petifiesCollection *mongo.Collection
}

func New(client *mongo.Client) petifiesaggre.PetifiesRepository {
	return &petifiesMongoRepository{
		client:             client,
		petifiesCollection: client.Database(cmd.Conf.DatabaseName).Collection("petifies"),
	}
}

func (pr *petifiesMongoRepository) GetByUserID(ctx context.Context, userID uuid.UUID, pageSize int, afterID uuid.UUID) ([]*petifiesaggre.PetifiesAggre, error) {
	var petifies []*petifiesaggre.PetifiesAggre

	err := dbutils.ExecWithSession(ctx, pr.client, func(ssCtx mongo.SessionContext) error {
		petifies_, err := pr.GetByUserIDWithSession(ssCtx, userID, pageSize, afterID)
		if err != nil {
			return err
		}
		petifies = petifies_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return petifies, nil
}

func (pr *petifiesMongoRepository) GetByID(ctx context.Context, id uuid.UUID) (*petifiesaggre.PetifiesAggre, error) {
	var petify *petifiesaggre.PetifiesAggre

	err := dbutils.ExecWithSession(ctx, pr.client, func(ssCtx mongo.SessionContext) error {
		petify_, err := pr.GetByIDWithSession(ssCtx, id)
		if err != nil {
			return err
		}
		petify = petify_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return petify, nil
}

func (pr *petifiesMongoRepository) ListByIDs(ctx context.Context, ids []uuid.UUID) ([]*petifiesaggre.PetifiesAggre, error) {
	var petifies []*petifiesaggre.PetifiesAggre

	err := dbutils.ExecWithSession(ctx, pr.client, func(ssCtx mongo.SessionContext) error {
		petifies_, err := pr.ListByIDsWithSession(ssCtx, ids)
		if err != nil {
			return err
		}
		petifies = petifies_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return petifies, nil
}

func (pr *petifiesMongoRepository) Save(ctx context.Context, petify petifiesaggre.PetifiesAggre) (*petifiesaggre.PetifiesAggre, error) {
	var savedPetify *petifiesaggre.PetifiesAggre

	err := dbutils.ExecWithSession(ctx, pr.client, func(ssCtx mongo.SessionContext) error {
		petify_, err := pr.SaveWithSession(ssCtx, petify)
		if err != nil {
			return err
		}
		savedPetify = petify_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return savedPetify, nil
}

func (pr *petifiesMongoRepository) Update(ctx context.Context, petify petifiesaggre.PetifiesAggre) (*petifiesaggre.PetifiesAggre, error) {
	var updatedPetify *petifiesaggre.PetifiesAggre

	err := dbutils.ExecWithSession(ctx, pr.client, func(ssCtx mongo.SessionContext) error {
		petify_, err := pr.UpdateWithSession(ssCtx, petify)
		if err != nil {
			return err
		}
		updatedPetify = petify_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return updatedPetify, nil
}

func (pr *petifiesMongoRepository) GetByUserIDWithSession(ctx context.Context, userID uuid.UUID, pageSize int, afterID uuid.UUID) ([]*petifiesaggre.PetifiesAggre, error) {
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

	var result []*petifiesaggre.PetifiesAggre
	var petifies []models.Petifies
	filter := bson.D{{Key: "owner_id", Value: userID}, {Key: "created_at", Value: bson.D{{Key: "$lt", Value: createdAt}}}}
	opts := options.Find().SetLimit(int64(pageSize)).SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := pr.petifiesCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if err := cursor.All(ctx, &petifies); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	for _, p := range petifies {
		petify, err := mappers.DbModelToPetifiesAggregate(&p)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		result = append(result, petify)
	}

	return result, nil
}

func (pr *petifiesMongoRepository) GetByIDWithSession(ctx context.Context, id uuid.UUID) (*petifiesaggre.PetifiesAggre, error) {
	var petifyDB models.Petifies
	err := pr.petifiesCollection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&petifyDB)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrPetifiesNotExist
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result, err := mappers.DbModelToPetifiesAggregate(&petifyDB)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (pr *petifiesMongoRepository) ListByIDsWithSession(ctx context.Context, ids []uuid.UUID) ([]*petifiesaggre.PetifiesAggre, error) {
	var results []*petifiesaggre.PetifiesAggre

	cursor, err := pr.petifiesCollection.Find(ctx, bson.D{{Key: "id", Value: bson.D{{Key: "$in", Value: ids}}}})
	if err != nil {
		return nil, err
	}
	var petifiesDB []models.Petifies
	if err := cursor.All(ctx, &petifiesDB); err != nil {
		return nil, err
	}

	for _, p := range petifiesDB {
		petify, err := mappers.DbModelToPetifiesAggregate(&p)
		if err != nil {
			return nil, err
		}
		results = append(results, petify)
	}

	return results, nil
}

func (pr *petifiesMongoRepository) SaveWithSession(ctx context.Context, petify petifiesaggre.PetifiesAggre) (*petifiesaggre.PetifiesAggre, error) {
	petifyDB := mappers.AggregatePetifiesToDbPetifies(&petify)
	_, err := pr.petifiesCollection.InsertOne(ctx, petifyDB)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	savedPetify, err := pr.GetByIDWithSession(ctx, petify.GetID())
	if err != nil {
		return nil, err
	}
	return savedPetify, nil
}

func (pr *petifiesMongoRepository) UpdateWithSession(ctx context.Context, petify petifiesaggre.PetifiesAggre) (*petifiesaggre.PetifiesAggre, error) {
	petifyDB := mappers.AggregatePetifiesToDbPetifies(&petify)
	_, err := pr.petifiesCollection.ReplaceOne(ctx, bson.D{{Key: "id", Value: petify.GetID()}}, petifyDB)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	updatedPetify, err := pr.GetByIDWithSession(ctx, petify.GetID())
	if err != nil {
		return nil, err
	}
	return updatedPetify, nil
}
