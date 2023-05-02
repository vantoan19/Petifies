package mongo_post

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

	"github.com/vantoan19/Petifies/server/services/post-service/cmd"
	postaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/post"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/infra/db/mapper"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/infra/db/models"
	mongo_comment "github.com/vantoan19/Petifies/server/services/post-service/internal/infra/repositories/comment/mongo"
)

var (
	ErrPostNotExist = status.Errorf(codes.NotFound, "post does not exist")
	wc              = writeconcern.New(writeconcern.WMajority())
	rc              = readconcern.Snapshot()
	transOpts       = options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
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

	err := pr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		post_, err := pr.GetByUUIDWithSession(ssCtx, id)
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
	err := pr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		post_, err := pr.SavePostWithSession(ssCtx, &post)
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
	err := pr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		post_, err := pr.UpdatePostWithSession(ssCtx, &post)
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
	err := pr.execSession(ctx, func(ssCtx mongo.SessionContext) error {
		post_, err := pr.DeleteByUUIDWithSession(ssCtx, id)
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

func (pr *PostRepository) GetByUUIDWithSession(ctx context.Context, id uuid.UUID) (*postaggre.Post, error) {
	var post models.Post
	err := pr.postCollection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrPostNotExist
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result, err := mapper.DbModelsToPostAggregate(&post)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (pr *PostRepository) SavePostWithSession(ctx context.Context, post *postaggre.Post) (*postaggre.Post, error) {
	postEntity := post.GetPostEntity()
	_, err := pr.postCollection.InsertOne(ctx, mapper.EntityPostToDbPost(&postEntity))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	post_, err := pr.GetByUUIDWithSession(ctx, post.GetPostEntity().ID)
	if err != nil {
		return nil, err
	}
	return post_, nil
}

func (pr *PostRepository) UpdatePostWithSession(ctx context.Context, post *postaggre.Post) (*postaggre.Post, error) {
	postEntity := post.GetPostEntity()
	_, err := pr.postCollection.ReplaceOne(ctx, bson.D{{Key: "id", Value: postEntity.ID}}, mapper.EntityPostToDbPost(&postEntity))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	post_, err := pr.GetByUUIDWithSession(ctx, post.GetPostEntity().ID)
	if err != nil {
		return nil, err
	}
	return post_, nil
}

func (pr *PostRepository) DeleteByUUIDWithSession(ctx context.Context, id uuid.UUID) (*postaggre.Post, error) {
	post, err := pr.GetByUUIDWithSession(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = pr.postCollection.DeleteOne(ctx, bson.D{{Key: "id", Value: post.GetPostEntity().ID}})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	cr := mongo_comment.New(pr.client)
	for _, commentID := range post.GetComments() {
		_, err := cr.DeleteByUUIDWithSession(ctx, commentID)
		if err != nil {
			return nil, err
		}
	}

	return post, nil
}

func (pr *PostRepository) execSession(ctx context.Context, fn func(ssCtx mongo.SessionContext) error) error {
	session, err := pr.client.StartSession()
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
		if err == ErrPostNotExist {
			return err
		}
		return status.Errorf(codes.Internal, err.Error())
	}
	return session.CommitTransaction(ctx)
}
