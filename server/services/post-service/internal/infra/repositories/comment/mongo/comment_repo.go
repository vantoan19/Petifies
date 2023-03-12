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

	utils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/post-service/cmd"
	commentaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/comment"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/infra/db/mapper"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/infra/db/models"
)

var (
	ErrCommentNoExist = status.Errorf(codes.NotFound, "comment does not exist")
	wc                = writeconcern.New(writeconcern.WMajority())
	rc                = readconcern.Snapshot()
	transOpts         = options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
)

type CommentRepository struct {
	client            *mongo.Client
	postCollection    *mongo.Collection
	commentCollection *mongo.Collection
	loveCollection    *mongo.Collection
}

func New(client *mongo.Client) *CommentRepository {
	return &CommentRepository{
		client:            client,
		postCollection:    client.Database(cmd.Conf.DatabaseName).Collection("posts"),
		commentCollection: client.Database(cmd.Conf.DatabaseName).Collection("comments"),
		loveCollection:    client.Database(cmd.Conf.DatabaseName).Collection("loves"),
	}
}

func (cr *CommentRepository) GetByUUID(ctx context.Context, id uuid.UUID) (*commentaggre.Comment, error) {
	var comment *commentaggre.Comment
	err := cr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		comment_, err := cr.GetByUUIDWithSession(ssCtx, id)
		if err != nil {
			return err
		}
		comment = comment_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (cr *CommentRepository) SaveComment(ctx context.Context, comment commentaggre.Comment) (*commentaggre.Comment, error) {
	var savedComment *commentaggre.Comment
	err := cr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		comment_, err := cr.SavePostWithSession(ssCtx, &comment)
		if err != nil {
			return err
		}
		savedComment = comment_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return savedComment, nil
}

func (cr *CommentRepository) UpdateComment(ctx context.Context, comment commentaggre.Comment) (*commentaggre.Comment, error) {
	var updatedComment *commentaggre.Comment
	err := cr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		comment_, err := cr.UpdatePostWithSession(ssCtx, &comment)
		if err != nil {
			return err
		}
		updatedComment = comment_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return updatedComment, nil
}

func (cr *CommentRepository) DeleteByUUID(ctx context.Context, id uuid.UUID) (*commentaggre.Comment, error) {
	var comment *commentaggre.Comment
	err := cr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		comment_, err := cr.DeleteByUUIDWithSession(ssCtx, id)
		if err != nil {
			return err
		}
		comment = comment_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (cr *CommentRepository) GetByUUIDWithSession(ctx context.Context, id uuid.UUID) (*commentaggre.Comment, error) {
	var comment models.Comment
	err := cr.commentCollection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&comment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrCommentNoExist
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var loves []models.Love
	cursor, err := cr.loveCollection.Find(ctx, bson.D{{Key: "comment_id", Value: id}})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if err := cursor.All(ctx, &loves); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var comments []models.Comment
	cursor, err = cr.commentCollection.Find(ctx, bson.D{{Key: "parent_id", Value: id}})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result, err := mapper.DbModelsToCommentAggregate(&comment, &loves, &comments)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (cr *CommentRepository) SavePostWithSession(ctx context.Context, comment *commentaggre.Comment) (*commentaggre.Comment, error) {
	commentEntity := comment.GetCommentEntity()
	_, err := cr.commentCollection.InsertOne(ctx, mapper.EntityCommentToDbComment(&commentEntity))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	loves := utils.Map2(comment.GetLoves(), func(l entities.Love) mongo.WriteModel {
		return mongo.NewInsertOneModel().SetDocument(mapper.EntityLoveToDbLove(&l))
	})
	if len(loves) > 0 {
		_, err = cr.loveCollection.BulkWrite(ctx, loves)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}

	comment_, err := cr.GetByUUIDWithSession(ctx, comment.GetCommentEntity().ID)
	if err != nil {
		return nil, err
	}
	return comment_, nil
}

func (cr *CommentRepository) UpdatePostWithSession(ctx context.Context, comment *commentaggre.Comment) (*commentaggre.Comment, error) {
	commentEntity := comment.GetCommentEntity()
	_, err := cr.commentCollection.ReplaceOne(ctx, bson.D{{Key: "id", Value: commentEntity.ID}}, mapper.EntityCommentToDbComment(&commentEntity))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	loves := utils.Map2(comment.GetLoves(), func(l entities.Love) mongo.WriteModel {
		return mongo.NewReplaceOneModel().SetFilter(bson.D{{Key: "id", Value: l.ID}}).SetUpsert(true).SetReplacement(mapper.EntityLoveToDbLove(&l))
	})
	if len(loves) > 0 {
		_, err = cr.loveCollection.BulkWrite(ctx, loves)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}

	comment_, err := cr.GetByUUIDWithSession(ctx, comment.GetCommentEntity().ID)
	if err != nil {
		return nil, err
	}
	return comment_, nil
}

func (cr *CommentRepository) DeleteByUUIDWithSession(ctx context.Context, id uuid.UUID) (*commentaggre.Comment, error) {
	comment, err := cr.GetByUUIDWithSession(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = cr.commentCollection.DeleteOne(ctx, bson.D{{Key: "id", Value: comment.GetCommentEntity().ID}})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	loves := utils.Map2(comment.GetLoves(), func(l entities.Love) mongo.WriteModel {
		return mongo.NewDeleteOneModel().SetFilter(bson.D{{Key: "id", Value: l.ID}})
	})
	if len(loves) > 0 {
		_, err = cr.loveCollection.BulkWrite(ctx, loves)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}

	for _, subcommentID := range comment.GetSubcommentsID() {
		_, err = cr.DeleteByUUIDWithSession(ctx, subcommentID)
		if err != nil {
			return nil, err
		}
	}

	return comment, nil
}

func (cr *CommentRepository) execSession(ctx context.Context, fn func(ssCtx mongo.SessionContext) error) error {
	session, err := cr.client.StartSession()
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
		return status.Errorf(codes.Internal, err.Error())
	}

	return session.CommitTransaction(ctx)
}
