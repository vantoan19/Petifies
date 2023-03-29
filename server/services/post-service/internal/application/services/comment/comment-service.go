package commentservice

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	utils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	commentaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/comment"
	loveaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/love"
	postaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/post"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/entities"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/valueobjects"
	mongo_comment "github.com/vantoan19/Petifies/server/services/post-service/internal/infra/repositories/comment/mongo"
	mongo_love "github.com/vantoan19/Petifies/server/services/post-service/internal/infra/repositories/love/mongo"
	mongo_post "github.com/vantoan19/Petifies/server/services/post-service/internal/infra/repositories/post/mongo"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var logger = logging.New("PostService.CommentSvc")

type commentService struct {
	postRepo    postaggre.PostRepository
	commentRepo commentaggre.CommentRepository
	loveRepo    loveaggre.LoveRepository
}

type CommentConfiguration func(cs *commentService) error

type CommentService interface {
	CreateComment(ctx context.Context, comment *models.CreateCommentReq) (*commentaggre.Comment, error)
	LoveReactComment(ctx context.Context, req *models.LoveReactReq) (*loveaggre.Love, error)
	EditComment(ctx context.Context, req *models.EditCommentReq) (*commentaggre.Comment, error)
	ListComments(ctx context.Context, req *models.ListCommentsReq) ([]*commentaggre.Comment, error)
	GetLoveCount(ctx context.Context, commentID uuid.UUID) (int, error)
	GetSubcommentCount(ctx context.Context, commentID uuid.UUID) (int, error)
	GetComment(ctx context.Context, commentID uuid.UUID) (*commentaggre.Comment, error)
	RemoveLoveReactComment(ctx context.Context, req *models.RemoveLoveReactReq) error
	ListCommentsByParentID(ctx context.Context, parentID uuid.UUID, pageSize int, afterCommentID uuid.UUID) ([]*commentaggre.Comment, error)
	ListCommentAncestors(ctx context.Context, commentID uuid.UUID) ([]*commentaggre.Comment, error)
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

func WithMongoLoveRepository(client *mongo.Client) CommentConfiguration {
	return func(cs *commentService) error {
		repo := mongo_love.New(client)
		cs.loveRepo = repo
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

func (cs *commentService) LoveReactComment(ctx context.Context, req *models.LoveReactReq) (*loveaggre.Love, error) {
	logger.Info("Start LoveReactComment")

	comment, err := cs.commentRepo.GetByUUID(ctx, req.TargetID)
	if err != nil {
		logger.ErrorData("Finish LoveReactComment: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}
	err = comment.AddLoveByAuthorIDAndSave(req.AuthorID, cs.loveRepo)
	if err != nil {
		logger.ErrorData("Finish LoveReactComment: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}
	love, err := cs.loveRepo.GetByTargetIDAndAuthorID(ctx, req.AuthorID, req.TargetID)
	if err != nil {
		logger.ErrorData("Finish LoveReactComment: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish LoveReactComment: Successful")
	return love, nil
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
			comment, err := cs.commentRepo.GetByUUID(ctx, id)
			if err == mongo_comment.ErrCommentNoExist {
				return
			}
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

func (cs *commentService) GetLoveCount(ctx context.Context, commentID uuid.UUID) (int, error) {
	logger.Info("Start GetLoveCount")

	count, err := cs.loveRepo.CountLoveByTargetID(ctx, commentID)
	if err != nil {
		logger.ErrorData("Finish GetLoveCount: Failed", logging.Data{"error": err.Error()})
		return 0, err
	}

	logger.Info("Finish LoveReactComment: Successful")
	return count, nil
}

func (cs *commentService) GetSubcommentCount(ctx context.Context, commentID uuid.UUID) (int, error) {
	logger.Info("Start GetSubcommentCount")

	count, err := cs.commentRepo.CountCommentByParentID(ctx, commentID)
	if err != nil {
		logger.ErrorData("Finish GetSubcommentCount: Failed", logging.Data{"error": err.Error()})
		return 0, err
	}

	logger.Info("Finish GetSubcommentCount: Successful")
	return count, nil
}

func (cs *commentService) GetComment(ctx context.Context, commentID uuid.UUID) (*commentaggre.Comment, error) {
	logger.Info("Start GetComment")

	comment, err := cs.commentRepo.GetByUUID(ctx, commentID)
	if err != nil {
		logger.ErrorData("Finish GetComment: Failed", logging.Data{"error": err.Error()})
	}

	logger.Info("Finish GetComment: Successful")
	return comment, nil
}

func (cs *commentService) RemoveLoveReactComment(ctx context.Context, req *models.RemoveLoveReactReq) error {
	logger.Info("RemoveLoveReactComment")

	comment, err := cs.commentRepo.GetByUUID(ctx, req.TargetID)
	if err != nil {
		logger.ErrorData("Finish RemoveLoveReactComment: Failed", logging.Data{"error": err.Error()})
		return err
	}
	err = comment.RemoveLoveByAuthorIDAndDelete(req.AuthorID, cs.loveRepo)
	if err != nil {
		logger.ErrorData("Finish RemoveLoveReactComment: Failed", logging.Data{"error": err.Error()})
		return err
	}

	if exists, err := cs.loveRepo.ExistsLoveByTargetIDAndAuthorID(ctx, req.AuthorID, req.TargetID); err != nil {
		logger.ErrorData("Finish RemoveLoveReactComment: Failed", logging.Data{"error": err.Error()})
		return err
	} else if exists {
		logger.ErrorData("Finish RemoveLoveReactComment: Failed", logging.Data{"error": "faiiled to remove react"})
		return status.Errorf(codes.Internal, "failed to remove react")
	}

	logger.Info("Finish RemoveLoveReactComment: Successful")
	return nil
}

func (cs *commentService) ListCommentsByParentID(ctx context.Context, parentID uuid.UUID, pageSize int, afterCommentID uuid.UUID) ([]*commentaggre.Comment, error) {
	logger.Info("ListCommentsByParentID")

	comments, err := cs.commentRepo.GetByParentID(ctx, parentID, pageSize, afterCommentID)
	if err != nil {
		logger.ErrorData("Finish RemoveLoveReactComment: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListCommentsByParentID: Successful")
	return comments, nil
}

func (cs *commentService) ListCommentAncestors(ctx context.Context, commentID uuid.UUID) ([]*commentaggre.Comment, error) {
	logger.Info("ListCommentAncestors")

	comments, err := cs.commentRepo.GetCommentAncestors(ctx, commentID)
	if err != nil {
		logger.ErrorData("Finish ListCommentAncestors: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListCommentAncestors: Successful")
	return comments, nil
}
