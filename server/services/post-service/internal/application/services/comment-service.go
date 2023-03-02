package services

import (
	"context"

	"github.com/vantoan19/Petifies/server/libs/logging-config"
	commentaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/comment"
	mongo_comment "github.com/vantoan19/Petifies/server/services/post-service/internal/infra/repositories/comment/mongo"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type commentService struct {
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

func (cs *commentService) CreateComment(ctx context.Context, comment *models.CreateCommentReq) (*commentaggre.Comment, error) {
	logger.Info("Start CommentService.CreateComment")

	newComment, err := commentaggre.New(comment)
	if err != nil {
		logger.ErrorData("Finish CommentService.CreateComment: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}
	createdComment, err := cs.commentRepo.SaveComment(ctx, *newComment)
	if err != nil {
		logger.ErrorData("Finish CommentService.CreateComment: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish CommentService.CreateComment: Successful")
	return createdComment, nil
}
