package petifiessessionmongo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/libs/dbutils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/petifies-service/cmd"
	petifiessessionaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_session"
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
	ErrSessionNotExist = status.Errorf(codes.NotFound, "session does not exist")
)

type petifiesSessionMongoRepository struct {
	client            *mongo.Client
	sessionCollection *mongo.Collection
}

func New(client *mongo.Client) petifiessessionaggre.PetifiesSessionRepository {
	return &petifiesSessionMongoRepository{
		client:            client,
		sessionCollection: client.Database(cmd.Conf.DatabaseName).Collection("petifies_sessions"),
	}
}

func (pr *petifiesSessionMongoRepository) GetByPetifiesID(ctx context.Context, petifiesID uuid.UUID, pageSize int, afterID uuid.UUID) ([]*petifiessessionaggre.PetifiesSessionAggre, error) {
	var sessions []*petifiessessionaggre.PetifiesSessionAggre

	err := dbutils.ExecWithSession(ctx, pr.client, func(ssCtx mongo.SessionContext) error {
		sessions_, err := pr.GetByPetifiesIDWithSession(ssCtx, petifiesID, pageSize, afterID)
		if err != nil {
			return err
		}
		sessions = sessions_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (pr *petifiesSessionMongoRepository) GetByID(ctx context.Context, id uuid.UUID) (*petifiessessionaggre.PetifiesSessionAggre, error) {
	var session *petifiessessionaggre.PetifiesSessionAggre

	err := dbutils.ExecWithSession(ctx, pr.client, func(ssCtx mongo.SessionContext) error {
		session_, err := pr.GetByIDWithSession(ssCtx, id)
		if err != nil {
			return err
		}
		session = session_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (pr *petifiesSessionMongoRepository) ListByIds(ctx context.Context, ids []uuid.UUID) ([]*petifiessessionaggre.PetifiesSessionAggre, error) {
	var sessions []*petifiessessionaggre.PetifiesSessionAggre

	err := dbutils.ExecWithSession(ctx, pr.client, func(ssCtx mongo.SessionContext) error {
		sessions_, err := pr.ListByIdsWithSession(ssCtx, ids)
		if err != nil {
			return err
		}
		sessions = sessions_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (pr *petifiesSessionMongoRepository) Save(ctx context.Context, session petifiessessionaggre.PetifiesSessionAggre) (*petifiessessionaggre.PetifiesSessionAggre, error) {
	var savedSession *petifiessessionaggre.PetifiesSessionAggre

	err := dbutils.ExecWithSession(ctx, pr.client, func(ssCtx mongo.SessionContext) error {
		session_, err := pr.SaveWithSession(ssCtx, session)
		if err != nil {
			return err
		}
		savedSession = session_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return savedSession, nil
}

func (pr *petifiesSessionMongoRepository) Update(ctx context.Context, session petifiessessionaggre.PetifiesSessionAggre) (*petifiessessionaggre.PetifiesSessionAggre, error) {
	var updatedSession *petifiessessionaggre.PetifiesSessionAggre

	err := dbutils.ExecWithSession(ctx, pr.client, func(ssCtx mongo.SessionContext) error {
		session_, err := pr.UpdateWithSession(ssCtx, session)
		if err != nil {
			return err
		}
		updatedSession = session_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return updatedSession, nil
}

func (pr *petifiesSessionMongoRepository) GetByPetifiesIDWithSession(
	ctx context.Context,
	petifiesID uuid.UUID,
	pageSize int,
	afterID uuid.UUID,
) ([]*petifiessessionaggre.PetifiesSessionAggre, error) {
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

	var result []*petifiessessionaggre.PetifiesSessionAggre
	var sessions []models.PetifiesSession
	filter := bson.D{{Key: "petifies_id", Value: petifiesID}, {Key: "created_at", Value: bson.D{{Key: "$lt", Value: createdAt}}}}
	opts := options.Find().SetLimit(int64(pageSize)).SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := pr.sessionCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	for _, s := range sessions {
		petify, err := mappers.DbModelToPetifiesSessionAggregate(&s)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		result = append(result, petify)
	}

	return result, nil
}

func (pr *petifiesSessionMongoRepository) GetByIDWithSession(ctx context.Context, id uuid.UUID) (*petifiessessionaggre.PetifiesSessionAggre, error) {
	var petifySessionDB models.PetifiesSession
	err := pr.sessionCollection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&petifySessionDB)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrSessionNotExist
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result, err := mappers.DbModelToPetifiesSessionAggregate(&petifySessionDB)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (pr *petifiesSessionMongoRepository) ListByIdsWithSession(ctx context.Context, ids []uuid.UUID) ([]*petifiessessionaggre.PetifiesSessionAggre, error) {
	var results []*petifiessessionaggre.PetifiesSessionAggre

	cursor, err := pr.sessionCollection.Find(ctx, bson.D{{Key: "id", Value: bson.D{{Key: "$in", Value: ids}}}})
	if err != nil {
		return nil, err
	}
	var sessionsDB []models.PetifiesSession
	if err := cursor.All(ctx, sessionsDB); err != nil {
		return nil, err
	}

	for _, s := range sessionsDB {
		session, err := mappers.DbModelToPetifiesSessionAggregate(&s)
		if err != nil {
			return nil, err
		}
		results = append(results, session)
	}

	return results, nil
}

func (pr *petifiesSessionMongoRepository) SaveWithSession(ctx context.Context, session petifiessessionaggre.PetifiesSessionAggre) (*petifiessessionaggre.PetifiesSessionAggre, error) {
	petifySessionDB := mappers.AggregatePetifiesSessionToDbPetifiesSession(&session)
	_, err := pr.sessionCollection.InsertOne(ctx, petifySessionDB)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	savedSession, err := pr.GetByIDWithSession(ctx, session.GetID())
	if err != nil {
		return nil, err
	}
	return savedSession, nil
}

func (pr *petifiesSessionMongoRepository) UpdateWithSession(ctx context.Context, session petifiessessionaggre.PetifiesSessionAggre) (*petifiessessionaggre.PetifiesSessionAggre, error) {
	petifySessionDB := mappers.AggregatePetifiesSessionToDbPetifiesSession(&session)
	_, err := pr.sessionCollection.ReplaceOne(ctx, bson.D{{Key: "id", Value: session.GetID()}}, petifySessionDB)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	updatedSession, err := pr.GetByIDWithSession(ctx, session.GetID())
	if err != nil {
		return nil, err
	}
	return updatedSession, nil
}
