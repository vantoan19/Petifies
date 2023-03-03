package postservice

import (
	"context"
	"errors"
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

var logger = logging.New("PostService.PostSvc")

type postService struct {
	postRepo    postaggre.PostRepository
	commentRepo commentaggre.CommentRepository
}

type PostConfiguration func(ps *postService) error

type PostService interface {
	CreatePost(ctx context.Context, post *models.CreatePostReq) (*postaggre.Post, error)
	LoveReactPost(ctx context.Context, req *models.LoveReactReq) (*entities.Love, error)
	EditPost(ctx context.Context, post *models.EditPostReq) (*postaggre.Post, error)
	ListPosts(ctx context.Context, req *models.ListPostsReq) ([]*postaggre.Post, error)
}

func NewPostService(cfgs ...PostConfiguration) (PostService, error) {
	ps := &postService{}
	for _, cfg := range cfgs {
		if err := cfg(ps); err != nil {
			return nil, err
		}
	}
	return ps, nil
}

func WithMongoPostRepository(client *mongo.Client) PostConfiguration {
	return func(ps *postService) error {
		repo := mongo_post.New(client)
		ps.postRepo = repo
		return nil
	}
}

func WithMongoCommentRepository(client *mongo.Client) PostConfiguration {
	return func(ps *postService) error {
		repo := mongo_comment.New(client)
		ps.commentRepo = repo
		return nil
	}
}

func (ps *postService) CreatePost(ctx context.Context, post *models.CreatePostReq) (*postaggre.Post, error) {
	logger.Info("Start CreatePost")

	newPost, err := postaggre.NewPost(post)
	if err != nil {
		logger.ErrorData("Finish CreatePost: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}
	createdPost, err := ps.postRepo.SavePost(ctx, *newPost)
	if err != nil {
		logger.ErrorData("Finish CreatePost: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish CreatePost: Successful")
	return createdPost, nil
}

func (ps *postService) LoveReactPost(ctx context.Context, req *models.LoveReactReq) (*entities.Love, error) {
	logger.Info("Start LoveReactPost")

	post, err := ps.postRepo.GetByUUID(ctx, req.TargetID)
	if err != nil {
		logger.ErrorData("Finish LoveReactPost: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}
	err = post.AddLoveByAuthorID(req.AuthorID)
	if err != nil {
		logger.ErrorData("Finish LoveReactPost: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	updatedPost, err := ps.postRepo.UpdatePost(ctx, *post)
	if err != nil {
		logger.ErrorData("Finish LoveReactPost: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}
	love := updatedPost.GetLovesByAuthorID(req.AuthorID)
	if love.AuthorID == uuid.Nil {
		logger.ErrorData("Finish LoveReactPost: Failed", logging.Data{"error": err.Error()})
		return nil, errors.New("failed to react post")
	}

	logger.Info("Finish LoveReactPost: Successful")
	return &love, nil
}

func (ps *postService) EditPost(ctx context.Context, req *models.EditPostReq) (*postaggre.Post, error) {
	logger.Info("Start EditPost")

	post, err := ps.postRepo.GetByUUID(ctx, req.ID)
	if err != nil {
		logger.ErrorData("Finish EditPost: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	post.RemoveAllImages()
	post.RemoveAllVideos()

	post.UpdateTextContent(valueobjects.NewTextContent(req.Content))
	for _, i := range req.Images {
		err = post.AddNewImage(valueobjects.NewImageContent(i.URL, i.Description))
		if err != nil {
			logger.ErrorData("Finish EditPost: Failed", logging.Data{"error": err.Error()})
			return nil, err
		}
	}
	for _, v := range req.Videos {
		err = post.AddNewVideo(valueobjects.NewVideoContent(v.URL, v.Description))
		if err != nil {
			logger.ErrorData("Finish EditPost: Failed", logging.Data{"error": err.Error()})
			return nil, err
		}
	}

	updatedPost, err := ps.postRepo.UpdatePost(ctx, *post)
	if err != nil {
		logger.ErrorData("Finish EditPost: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish EditPost: Successful")
	return updatedPost, err
}

func (ps *postService) ListPosts(ctx context.Context, req *models.ListPostsReq) ([]*postaggre.Post, error) {
	logger.Info("Start ListPosts")

	var wg sync.WaitGroup
	resultsChan := make(chan *postaggre.Post, len(req.PostIDs))
	errsChan := make(chan error, len(req.PostIDs))

	for _, id := range req.PostIDs {
		wg.Add(1)
		go func(id uuid.UUID) {
			defer wg.Done()
			post, err := ps.postRepo.GetByUUID(ctx, id)
			if err != nil {
				errsChan <- err
				return
			}
			resultsChan <- post
		}(id)
	}

	wg.Wait()
	close(errsChan)
	close(resultsChan)
	errs := utils.ToSlice(errsChan)
	results := utils.ToSlice(resultsChan)
	if len(errs) > 0 {
		logger.ErrorData("Finish ListPosts: Failed", logging.Data{"error": errs[0].Error()})
		return nil, errs[0]
	}

	logger.Info("Finish ListPosts: Successful")
	return results, nil
}
