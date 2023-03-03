package commentservice

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"

	utils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	commentaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/comment"
	postaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/post"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/valueobjects"
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
	LoveReactComment(ctx context.Context, req *models.LoveReactReq) (*entities.Love, error)
	EditComment(ctx context.Context, req *models.EditCommentReq) (*commentaggre.Comment, error)
	ListComments(ctx context.Context, req *models.ListCommentsReq) ([]*commentaggre.Comment, error)
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

func (cs *commentService) LoveReactComment(ctx context.Context, req *models.LoveReactReq) (*entities.Love, error) {
	logger.Info("Start LoveReactComment")

	comment, err := cs.commentRepo.GetByUUID(ctx, req.TargetID)
	if err != nil {
		logger.ErrorData("Finish LoveReactComment: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}
	err = comment.AddLoveByAuthorID(req.AuthorID)
	if err != nil {
		logger.ErrorData("Finish LoveReactComment: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	updatedComment, err := cs.commentRepo.UpdateComment(ctx, *comment)
	if err != nil {
		logger.ErrorData("Finish LoveReactComment: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}
	love := updatedComment.GetLovesByAuthorID(req.AuthorID)
	if love.AuthorID == uuid.Nil {
		logger.ErrorData("Finish LoveReactComment: Failed", logging.Data{"error": err.Error()})
		return nil, errors.New("failed to react comment")
	}

	logger.Info("Finish LoveReactComment: Successful")
	return &love, nil
}

func (cs *commentService) EditComment(ctx context.Context, req *models.EditCommentReq) (*commentaggre.Comment, error) {
	logger.Info("Start EditComment")

	comment, err := cs.commentRepo.GetByUUID(ctx, req.ID)
	if err != nil {
		logger.ErrorData("Finish EditComment: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	err = comment.SetCommentEntity(entities.Comment{
		ID:           comment.GetID(),
		PostID:       comment.GetPostID(),
		AuthorID:     comment.GetAuthorID(),
		ParentID:     comment.GetParentID(),
		IsPostParent: comment.GetIsPostParent(),
		Content:      valueobjects.NewTextContent(req.Content),
		ImageContent: valueobjects.NewImageContent(req.Image.URL, req.Image.Description),
		VideoContent: valueobjects.NewVideoContent(req.Video.URL, req.Video.Description),
		CreatedAt:    comment.GetCreatedAt(),
		UpdatedAt:    comment.GetUpdatedAt(),
	})
	if err != nil {
		logger.ErrorData("Finish EditComment: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	updatedComment, err := cs.commentRepo.UpdateComment(ctx, *comment)
	if err != nil {
		logger.ErrorData("Finish EditPost: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	return updatedComment, err
}

func (cs *commentService) ListComments(ctx context.Context, req *models.ListCommentsReq) ([]*commentaggre.Comment, error) {
	logger.Info("Start ListComments")

	var wg sync.WaitGroup
	resultsChan := make(chan *commentaggre.Comment, len(req.CommentIDs))
	errsChan := make(chan error, len(req.CommentIDs))

	for _, id := range req.CommentIDs {
		wg.Add(1)
		go func(id uuid.UUID) {
			defer wg.Done()
			fmt.Println(id)
			comment, err := cs.commentRepo.GetByUUID(ctx, id)
			if err != nil {
				errsChan <- err
				return
			}
			resultsChan <- comment
		}(id)
	}

	wg.Wait()

	close(errsChan)
	close(resultsChan)
	errs := utils.ToSlice(errsChan)
	results := utils.ToSlice(resultsChan)
	if len(errs) > 0 {
		logger.ErrorData("Finish ListComments: Failed", logging.Data{"error": errs[0].Error()})
		return nil, errs[0]
	}

	logger.Info("Finish ListComments: Successful")
	return results, nil
}
