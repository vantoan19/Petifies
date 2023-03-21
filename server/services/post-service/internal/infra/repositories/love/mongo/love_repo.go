package mongo_comment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vantoan19/Petifies/server/libs/logging-config"
	paginateutils "github.com/vantoan19/Petifies/server/libs/paginate-utils"
	"github.com/vantoan19/Petifies/server/services/post-service/cmd"
	loveaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/love"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/infra/db/mapper"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/infra/db/models"
)

var logger = logging.New("PostService.MongoLoveRepository")

var (
	ErrLoveNotExist = status.Errorf(codes.NotFound, "love does not exist")
	wc              = writeconcern.New(writeconcern.WMajority())
	rc              = readconcern.Snapshot()
	transOpts       = options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
)

type LoveRepository struct {
	client            *mongo.Client
	postCollection    *mongo.Collection
	commentCollection *mongo.Collection
	loveCollection    *mongo.Collection
}

func New(client *mongo.Client) *LoveRepository {
	return &LoveRepository{
		client:            client,
		postCollection:    client.Database(cmd.Conf.DatabaseName).Collection("posts"),
		commentCollection: client.Database(cmd.Conf.DatabaseName).Collection("comments"),
		loveCollection:    client.Database(cmd.Conf.DatabaseName).Collection("loves"),
	}
}

func (lr *LoveRepository) GetByTargetID(ctx context.Context, targetID uuid.UUID, pageToken paginateutils.PageToken) ([]*loveaggre.Love, error) {
	logger.Info("Start GetByTargetID")
	var loves []*loveaggre.Love

	err := lr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		loves_, err := lr.GetByTargetIDWithSession(ssCtx, targetID, pageToken)
		if err != nil {
			return err
		}
		loves = loves_

		return nil
	})
	if err != nil {
		logger.ErrorData("Finish GetByTargetID: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetByTargetID: Successful")
	return loves, nil
}

func (lr *LoveRepository) GetByTargetIDAndAuthorID(ctx context.Context, authorID uuid.UUID, targetID uuid.UUID) (*loveaggre.Love, error) {
	logger.Info("Start GetByTargetIDAndAuthorID")
	var love *loveaggre.Love

	err := lr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		love_, err := lr.GetByTargetIDAndAuthorIDWithSession(ssCtx, authorID, targetID)
		if err != nil {
			return err
		}
		love = love_
		return nil
	})
	if err != nil {
		if err == ErrLoveNotExist {
			return nil, nil
		}
		logger.ErrorData("Finish GetByTargetIDAndAuthorID: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetByTargetIDAndAuthorID: Successful")
	return love, nil
}

func (lr *LoveRepository) CountLoveByTargetID(ctx context.Context, targetID uuid.UUID) (int, error) {
	logger.Info("Start CountLoveByTargetID")
	var result int
	err := lr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		count, err := lr.CountLoveByTargetIDWithSession(ssCtx, targetID)
		if err != nil {
			return err
		}
		result = int(count)
		return nil
	})
	if err != nil {
		logger.ErrorData("Finish CountLoveByTargetID: Failed", logging.Data{"error": err.Error()})
		return 0, err
	}

	logger.Info("Finish CountLoveByTargetID: Successful")
	return result, nil
}

func (lr *LoveRepository) SaveLove(ctx context.Context, love loveaggre.Love) (*loveaggre.Love, error) {
	logger.Info("Start SaveLove")

	var savedLove *loveaggre.Love
	err := lr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		love_, err := lr.SaveLoveWithSession(ssCtx, &love)
		if err != nil {
			return err
		}
		savedLove = love_
		return nil
	})
	if err != nil {
		logger.ErrorData("Finish SaveLove: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish SaveLove: Successful")
	return savedLove, nil
}

func (lr *LoveRepository) DeleteByUUID(ctx context.Context, id uuid.UUID) error {
	logger.Info("Start DeleteByUUID")

	err := lr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		_, err := lr.DeleteByUUIDWithSession(ssCtx, id)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.ErrorData("Finish DeleteByUUID: Failed", logging.Data{"error": err.Error()})
		return err
	}

	logger.Info("Finish DeleteByUUID: Successful")
	return nil
}

func (lr *LoveRepository) DeleteByTargetID(ctx context.Context, targetID uuid.UUID) error {
	logger.Info("Start DeleteByTargetID")

	err := lr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		err := lr.DeleteByTargetIDWithSession(ssCtx, targetID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.ErrorData("Finish DeleteByTargetID: Failed", logging.Data{"error": err.Error()})
		return err
	}

	logger.Info("Finish DeleteByTargetID: Successful")
	return nil
}

func (lr *LoveRepository) GetByUUIDWithSession(ctx context.Context, id uuid.UUID) (*loveaggre.Love, error) {
	var love models.Love
	err := lr.loveCollection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&love)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrLoveNotExist
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result, err := mapper.DbModelsToLoveAggregate(&love)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (lr *LoveRepository) GetByTargetIDAndAuthorIDWithSession(ctx context.Context, authorID, targetID uuid.UUID) (*loveaggre.Love, error) {
	var love models.Love
	err := lr.loveCollection.FindOne(ctx, bson.D{{Key: "target_id", Value: targetID}, {Key: "author_id", Value: authorID}}).Decode(&love)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrLoveNotExist
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result, err := mapper.DbModelsToLoveAggregate(&love)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (lr *LoveRepository) GetByTargetIDWithSession(ctx context.Context, targetID uuid.UUID, pageToken paginateutils.PageToken) ([]*loveaggre.Love, error) {
	var result []*loveaggre.Love
	var loves []models.Love
	cursor, err := lr.loveCollection.Find(
		ctx,
		bson.D{{Key: "target_id", Value: targetID}},
		options.Find().SetSkip(pageToken.Offset).SetLimit(int64(pageToken.PageSize)),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if err := cursor.All(ctx, &loves); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	for _, l := range loves {
		love, err := mapper.DbModelsToLoveAggregate(&l)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		result = append(result, love)
	}

	return result, nil
}

func (lr *LoveRepository) CountLoveByTargetIDWithSession(ctx context.Context, targetID uuid.UUID) (int64, error) {
	return lr.loveCollection.CountDocuments(ctx, bson.D{{Key: "target_id", Value: targetID}})
}

func (lr *LoveRepository) SaveLoveWithSession(ctx context.Context, love *loveaggre.Love) (*loveaggre.Love, error) {
	loveDBModel := models.Love{
		ID:           love.GetID(),
		TargetID:     love.GetTargetID(),
		IsPostTarget: love.GetIsPostTarget(),
		AuthorID:     love.GetAuthorID(),
		CreatedAt:    love.GetCreatedAt(),
	}
	_, err := lr.loveCollection.InsertOne(ctx, &loveDBModel)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	love_, err := lr.GetByUUIDWithSession(ctx, love.GetID())
	if err != nil {
		return nil, err
	}
	return love_, nil
}

func (lr *LoveRepository) DeleteByUUIDWithSession(ctx context.Context, id uuid.UUID) (*loveaggre.Love, error) {
	love, err := lr.GetByUUIDWithSession(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = lr.loveCollection.DeleteOne(ctx, bson.D{{Key: "id", Value: love.GetID()}})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return love, nil
}

func (lr *LoveRepository) DeleteByTargetIDWithSession(ctx context.Context, targetID uuid.UUID) error {
	_, err := lr.loveCollection.DeleteMany(ctx, bson.D{{Key: "target_id", Value: targetID}})
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	return nil
}

func (lr *LoveRepository) execSession(ctx context.Context, fn func(ssCtx mongo.SessionContext) error) error {
	session, err := lr.client.StartSession()
	defer session.EndSession(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	err = session.StartTransaction(transOpts)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	if err = fn(mongo.NewSessionContext(ctx, session)); err != nil {
		if abErr := session.AbortTransaction(ctx); abErr != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("session err: %v, abort err: %v", err, abErr))
		}
		if err == ErrLoveNotExist {
			return err
		}
		return status.Errorf(codes.Internal, err.Error())
	}

	return session.CommitTransaction(ctx)
}
