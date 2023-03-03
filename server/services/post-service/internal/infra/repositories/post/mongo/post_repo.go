package mongo_post

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	utils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/post-service/cmd"
	postaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/post"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/infra/db/mapper"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/infra/db/models"
	mongo_comment "github.com/vantoan19/Petifies/server/services/post-service/internal/infra/repositories/comment/mongo"
)

var (
	ErrPostNotExist = errors.New("post does not exist")
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
		return nil, err
	}

	var loves []models.Love
	cursor, err := pr.loveCollection.Find(ctx, bson.D{{Key: "post_id", Value: id}})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &loves); err != nil {
		return nil, err
	}

	var comments []models.Comment
	cursor, err = pr.commentCollection.Find(ctx, bson.D{{Key: "parent_id", Value: id}})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, err
	}

	result, err := mapper.DbModelsToPostAggregate(&post, &loves, &comments)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (pr *PostRepository) SavePostWithSession(ctx context.Context, post *postaggre.Post) (*postaggre.Post, error) {
	postEntity := post.GetPostEntity()
	_, err := pr.postCollection.InsertOne(ctx, mapper.EntityPostToDbPost(&postEntity))
	if err != nil {
		return nil, err
	}

	loves := utils.Map2(post.GetLoves(), func(l entities.Love) mongo.WriteModel {
		return mongo.NewInsertOneModel().SetDocument(mapper.EntityLoveToDbLove(&l))
	})
	if len(loves) > 0 {
		_, err = pr.loveCollection.BulkWrite(ctx, loves)
		if err != nil {
			return nil, err
		}
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
		return nil, err
	}

	// Get all loves
	var allLoves []models.Love
	cursor, err := pr.loveCollection.Find(ctx, bson.D{{Key: "post_id", Value: post.GetPostID()}})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &allLoves); err != nil {
		return nil, err
	}

	// Mark loves which currently exist in the aggregate
	existence := make(map[uuid.UUID]bool)
	for _, l := range post.GetLoves() {
		existence[l.ID] = true
	}
	lovesToDelete := utils.Filter(allLoves, func(l models.Love) bool { return !existence[l.ID] })

	operations := utils.Map2(post.GetLoves(), func(l entities.Love) mongo.WriteModel {
		return mongo.NewReplaceOneModel().SetFilter(bson.D{{Key: "id", Value: l.ID}}).SetUpsert(true).SetReplacement(mapper.EntityLoveToDbLove(&l))
	})
	operations = append(operations, utils.Map2(lovesToDelete, func(l models.Love) mongo.WriteModel {
		return mongo.NewDeleteOneModel().SetFilter(bson.D{{Key: "id", Value: l.ID}})
	})...)

	if len(operations) > 0 {
		_, err = pr.loveCollection.BulkWrite(ctx, operations)
		if err != nil {
			return nil, err
		}
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
		return nil, err
	}

	loves := utils.Map2(post.GetLoves(), func(l entities.Love) mongo.WriteModel {
		return mongo.NewDeleteOneModel().SetFilter(bson.D{{Key: "id", Value: l.ID}})
	})
	if len(loves) > 0 {
		_, err = pr.loveCollection.BulkWrite(ctx, loves)
		if err != nil {
			return nil, err
		}
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
		return err
	}
	err = session.StartTransaction(transOpts)
	if err != nil {
		return err
	}

	if err = fn(mongo.NewSessionContext(ctx, session)); err != nil {
		if abErr := session.AbortTransaction(ctx); abErr != nil {
			return fmt.Errorf("session err: %v, abort err: %v", err, abErr)
		}
		return err
	}
	return session.CommitTransaction(ctx)
}
