package mongo_comment

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/vantoan19/Petifies/server/services/post-service/cmd"
	commentaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/comment"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/infra/db/mapper"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/infra/db/models"
)

var (
	ErrCommentNoExist = errors.New("comment does not exist")
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
	err := cr.execSession(ctx, func(ss mongo.Session) error {
		comment_, err := cr.GetByUUIDWithSession(ctx, id, ss)
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
	err := cr.execSession(ctx, func(ss mongo.Session) error {
		comment_, err := cr.SavePostWithSession(ctx, &comment, ss)
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
	err := cr.execSession(ctx, func(ss mongo.Session) error {
		comment_, err := cr.UpdatePostWithSession(ctx, &comment, ss)
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
	err := cr.execSession(ctx, func(ss mongo.Session) error {
		comment_, err := cr.DeleteByUUIDWithSession(ctx, id, ss)
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

func (cr *CommentRepository) GetByUUIDWithSession(ctx context.Context, id uuid.UUID, ss mongo.Session) (*commentaggre.Comment, error) {
	commentDb, err := ss.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		var result models.Comment
		err := cr.commentCollection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&result)
		return &result, err
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrCommentNoExist
		}
		return nil, err
	}
	commentDb_, _ := commentDb.(*models.Comment)

	loves, err := ss.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		var results []models.Love
		cursor, err := cr.loveCollection.Find(ctx, bson.D{{Key: "comment_id", Value: id}})
		if err != nil {
			return nil, err
		}
		if err := cursor.All(ctx, &results); err != nil {
			return nil, err
		}
		return &results, nil
	})
	if err != nil {
		return nil, err
	}
	loves_, _ := loves.(*[]models.Love)

	comments, err := ss.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		var results []models.Comment
		cursor, err := cr.commentCollection.Find(ctx, bson.D{{Key: "parent_id", Value: id}})
		if err != nil {
			return nil, err
		}
		if err := cursor.All(ctx, &results); err != nil {
			return nil, err
		}
		return &results, err
	})
	if err != nil {
		return nil, err
	}
	comments_, _ := comments.(*[]models.Comment)

	comment, err := mapper.DbModelsToCommentAggregate(commentDb_, loves_, comments_)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (cr *CommentRepository) SavePostWithSession(ctx context.Context, comment *commentaggre.Comment, ss mongo.Session) (*commentaggre.Comment, error) {
	_, err := ss.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		commentEntity := comment.GetCommentEntity()
		result, err := cr.commentCollection.InsertOne(ctx, mapper.EntityCommentToDbComment(&commentEntity))
		return result, err
	})
	if err != nil {
		return nil, err
	}

	_, err = ss.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		var loves []interface{}
		for _, l := range comment.GetLoves() {
			loves = append(loves, *mapper.EntityLoveToDbLove(&l))
		}
		result, err := cr.loveCollection.InsertMany(ctx, loves)
		return result, err
	})
	if err != nil {
		return nil, err
	}

	comment_, err := cr.GetByUUIDWithSession(ctx, comment.GetCommentEntity().ID, ss)
	if err != nil {
		return nil, err
	}
	return comment_, nil
}

func (cr *CommentRepository) UpdatePostWithSession(ctx context.Context, comment *commentaggre.Comment, ss mongo.Session) (*commentaggre.Comment, error) {
	_, err := ss.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		commentEntity := comment.GetCommentEntity()
		result, err := cr.commentCollection.ReplaceOne(ctx, bson.D{{Key: "id", Value: commentEntity.ID}}, mapper.EntityCommentToDbComment(&commentEntity))
		return result, err
	})
	if err != nil {
		return nil, err
	}

	_, err = ss.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		var loves []mongo.WriteModel
		for _, l := range comment.GetLoves() {
			loves = append(loves, mongo.NewReplaceOneModel().SetFilter(bson.D{{Key: "id", Value: l.ID}}).SetUpsert(true).SetReplacement(mapper.EntityLoveToDbLove(&l)))
		}
		result, err := cr.loveCollection.BulkWrite(ctx, loves)
		return result, err
	})
	if err != nil {
		return nil, err
	}

	comment_, err := cr.GetByUUIDWithSession(ctx, comment.GetCommentEntity().ID, ss)
	if err != nil {
		return nil, err
	}
	return comment_, nil
}

func (cr *CommentRepository) DeleteByUUIDWithSession(ctx context.Context, id uuid.UUID, ss mongo.Session) (*commentaggre.Comment, error) {
	comment, err := cr.GetByUUIDWithSession(ctx, id, ss)
	if err != nil {
		return nil, err
	}

	_, err = ss.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		result, err := cr.commentCollection.DeleteOne(ctx, bson.D{{Key: "id", Value: comment.GetCommentEntity().ID}})
		return result, err
	})
	if err != nil {
		return nil, err
	}

	_, err = ss.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		var loves []mongo.WriteModel
		for _, l := range comment.GetLoves() {
			loves = append(loves, mongo.NewDeleteOneModel().SetFilter(bson.D{{Key: "id", Value: l.ID}}))
		}
		result, err := cr.loveCollection.BulkWrite(ctx, loves)
		return result, err
	})
	if err != nil {
		return nil, err
	}

	for _, subcommentID := range comment.GetSubcommentsID() {
		_, err = cr.DeleteByUUIDWithSession(ctx, subcommentID, ss)
		if err != nil {
			return nil, err
		}
	}

	return comment, nil
}

func (cr *CommentRepository) execSession(ctx context.Context, fn func(ss mongo.Session) error) error {
	session, err := cr.client.StartSession()
	defer session.EndSession(ctx)
	if err != nil {
		return err
	}

	if err = fn(session); err != nil {
		if abErr := session.AbortTransaction(ctx); abErr != nil {
			return fmt.Errorf("session err: %v, abort err: %v", err, abErr)
		}
		return err
	}

	return session.CommitTransaction(ctx)
}
