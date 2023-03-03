package postservice

import (
	"context"

	"github.com/vantoan19/Petifies/server/libs/logging-config"
	commentaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/comment"
	postaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/post"
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
