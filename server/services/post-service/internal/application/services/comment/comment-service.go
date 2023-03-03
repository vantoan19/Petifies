package commentservice

import (
	"context"
	"fmt"

	"github.com/vantoan19/Petifies/server/libs/logging-config"
	commentaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/comment"
	postaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/post"
	mongo_comment "github.com/vantoan19/Petifies/server/services/post-service/internal/infra/repositories/comment/mongo"
	mongo_post "github.com/vantoan19/Petifies/server/services/post-service/internal/infra/repositories/post/mongo"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var logger = logging.New("PostService.CommentSvc")

type commentService struct {
	postRepo    postaggre.PostRepository
	commentRepo commentaggre.CommentRepository
}

type CommentConfiguration func(cs *commentService) error

type CommentService interface {
	CreateComment(ctx context.Context, comment *models.CreateCommentReq) (*commentaggre.Comment, error)
}

func NewCommentService(cfgs ...CommentConfiguration) (CommentService, error) {
	cs := &commentService{}
	for _, cfg := range cfgs {
		if err := cfg(cs); err != nil {
			return nil, err
		}
	}
	return cs, nil
}

func WithMongoCommentRepository(client *mongo.Client) CommentConfiguration {
	return func(cs *commentService) error {
		repo := mongo_comment.New(client)
		cs.commentRepo = repo
		return nil
	}
}

func WithMongoPostRepository(client *mongo.Client) CommentConfiguration {
	return func(cs *commentService) error {
		repo := mongo_post.New(client)
		cs.postRepo = repo
		return nil
	}
}

func (cs *commentService) CreateComment(ctx context.Context, comment *models.CreateCommentReq) (*commentaggre.Comment, error) {
	logger.Info("Start CreateComment")

	newComment, err := commentaggre.New(comment)
	if err != nil {
		return nil, err
	}

	if comment.IsPostParent {
		fmt.Println(comment.ParentID)
		post, err := cs.postRepo.GetByUUID(ctx, comment.ParentID)
		if err != nil {
			logger.ErrorData("Finish CreateComment: Failed", logging.Data{"error": err.Error()})
			return nil, err
		}

		err = post.AddCommentAndSave(newComment, cs.commentRepo)
		if err != nil {
			logger.ErrorData("Finish CreateComment: Failed", logging.Data{"error": err.Error()})
			return nil, err
		}
	} else {
		parentComment, err := cs.commentRepo.GetByUUID(ctx, comment.ParentID)
		if err != nil {
			logger.ErrorData("Finish CreateComment: Failed", logging.Data{"error": err.Error()})
			return nil, err
		}

		err = parentComment.AddSubcommentAndSave(newComment, cs.commentRepo)
		if err != nil {
			logger.ErrorData("Finish CreateComment: Failed", logging.Data{"error": err.Error()})
			return nil, err
		}
	}

	createdComment, err := cs.commentRepo.GetByUUID(ctx, newComment.GetID())
	if err != nil {
		logger.ErrorData("Finish CreateComment: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish CreateComment: Successful")
	return createdComment, nil
}
