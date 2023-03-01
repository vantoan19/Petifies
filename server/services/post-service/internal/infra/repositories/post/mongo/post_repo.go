package mongo_post

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/post-service/cmd"
	postaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/post"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/infra/db/mapper"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/infra/db/models"
	mongo_comment "github.com/vantoan19/Petifies/server/services/post-service/internal/infra/repositories/comment/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrPostNotExist = errors.New("post does not exist")
)

type PostRepository struct {
	client            *mongo.Client
	postCollection    *mongo.Collection
	commentCollection *mongo.Collection
	loveCollection    *mongo.Collection
}

func New(client *mongo.Client) *PostRepository {
	return &PostRepository{
		client:            client,
		postCollection:    client.Database(cmd.Conf.DatabaseName).Collection("posts"),
		commentCollection: client.Database(cmd.Conf.DatabaseName).Collection("comments"),
		loveCollection:    client.Database(cmd.Conf.DatabaseName).Collection("loves"),
	}
}

func (pr *PostRepository) GetByUUID(ctx context.Context, id uuid.UUID) (*postaggre.Post, error) {
	var post *postaggre.Post

	err := pr.execSession(ctx, func(ss mongo.Session) error {
		post_, err := pr.GetByUUIDWithSession(ctx, id, ss)
		if err != nil {
			return err
		}
		post = post_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (pr *PostRepository) SavePost(ctx context.Context, post postaggre.Post) (*postaggre.Post, error) {
	var savedPost *postaggre.Post
	err := pr.execSession(ctx, func(ss mongo.Session) error {
		post_, err := pr.SavePostWithSession(ctx, &post, ss)
		if err != nil {
			return err
		}
		savedPost = post_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return savedPost, nil
}

func (pr *PostRepository) UpdatePost(ctx context.Context, post postaggre.Post) (*postaggre.Post, error) {
	var updatedPost *postaggre.Post
	err := pr.execSession(ctx, func(ss mongo.Session) error {
		post_, err := pr.UpdatePostWithSession(ctx, &post, ss)
		if err != nil {
			return err
		}
		updatedPost = post_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return updatedPost, nil
}

func (pr *PostRepository) DeleteByUUID(ctx context.Context, id uuid.UUID) (*postaggre.Post, error) {
	var post *postaggre.Post
	err := pr.execSession(ctx, func(ss mongo.Session) error {
		post_, err := pr.DeleteByUUIDWithSession(ctx, id, ss)
		if err != nil {
			return err
		}
		post = post_
		return nil
	})
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (pr *PostRepository) GetByUUIDWithSession(ctx context.Context, id uuid.UUID, ss mongo.Session) (*postaggre.Post, error) {
	postDb, err := ss.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		var result models.Post
		err := pr.postCollection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&result)
		return &result, err
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrPostNotExist
		}
		return nil, err
	}
	postDb_, _ := postDb.(*models.Post)

	loves, err := ss.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		var results []models.Love
		cursor, err := pr.loveCollection.Find(ctx, bson.D{{Key: "post_id", Value: id}})
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
		cursor, err := pr.commentCollection.Find(ctx, bson.D{{Key: "parent_id", Value: id}})
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

	post, err := mapper.DbModelsToPostAggregate(postDb_, loves_, comments_)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (pr *PostRepository) SavePostWithSession(ctx context.Context, post *postaggre.Post, ss mongo.Session) (*postaggre.Post, error) {
	_, err := ss.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		postEntity := post.GetPostEntity()
		result, err := pr.postCollection.InsertOne(ctx, mapper.EntityPostToDbPost(&postEntity))
		return result, err
	})
	if err != nil {
		return nil, err
	}

	_, err = ss.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		var loves []interface{}
		for _, l := range post.GetLoves() {
			loves = append(loves, *mapper.EntityLoveToDbLove(&l))
		}
		result, err := pr.loveCollection.InsertMany(ctx, loves)
		return result, err
	})
	if err != nil {
		return nil, err
	}

	post_, err := pr.GetByUUIDWithSession(ctx, post.GetPostEntity().ID, ss)
	if err != nil {
		return nil, err
	}
	return post_, nil
}

func (pr *PostRepository) UpdatePostWithSession(ctx context.Context, post *postaggre.Post, ss mongo.Session) (*postaggre.Post, error) {
	_, err := ss.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		postEntity := post.GetPostEntity()
		result, err := pr.postCollection.ReplaceOne(ctx, bson.D{{Key: "id", Value: postEntity.ID}}, mapper.EntityPostToDbPost(&postEntity))
		return result, err
	})
	if err != nil {
		return nil, err
	}

	_, err = ss.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		var loves []mongo.WriteModel
		for _, l := range post.GetLoves() {
			loves = append(loves, mongo.NewReplaceOneModel().SetFilter(bson.D{{Key: "id", Value: l.ID}}).SetUpsert(true).SetReplacement(mapper.EntityLoveToDbLove(&l)))
		}
		result, err := pr.loveCollection.BulkWrite(ctx, loves)
		return result, err
	})
	if err != nil {
		return nil, err
	}

	post_, err := pr.GetByUUIDWithSession(ctx, post.GetPostEntity().ID, ss)
	if err != nil {
		return nil, err
	}
	return post_, nil
}

func (pr *PostRepository) DeleteByUUIDWithSession(ctx context.Context, id uuid.UUID, ss mongo.Session) (*postaggre.Post, error) {
	post, err := pr.GetByUUIDWithSession(ctx, id, ss)
	if err != nil {
		return nil, err
	}

	_, err = ss.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		result, err := pr.postCollection.DeleteOne(ctx, bson.D{{Key: "id", Value: post.GetPostEntity().ID}})
		return result, err
	})
	if err != nil {
		return nil, err
	}

	_, err = ss.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		var loves []mongo.WriteModel
		for _, l := range post.GetLoves() {
			loves = append(loves, mongo.NewDeleteOneModel().SetFilter(bson.D{{Key: "id", Value: l.ID}}))
		}
		result, err := pr.loveCollection.BulkWrite(ctx, loves)
		return result, err
	})
	if err != nil {
		return nil, err
	}

	cr := mongo_comment.New(pr.client)
	for _, commentID := range post.GetComments() {
		_, err := cr.DeleteByUUIDWithSession(ctx, commentID, ss)
		if err != nil {
			return nil, err
		}
	}

	return post, nil
}

func (pr *PostRepository) execSession(ctx context.Context, fn func(ss mongo.Session) error) error {
	session, err := pr.client.StartSession()
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
