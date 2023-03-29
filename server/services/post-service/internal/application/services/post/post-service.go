package postservice

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	kafkamodels "github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/producer"
	outbox_repo "github.com/vantoan19/Petifies/server/infrastructure/outbox/repository"
	utils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	commentaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/comment"
	loveaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/love"
	postaggre "github.com/vantoan19/Petifies/server/services/post-service/internal/domain/aggregates/post"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/domain/common/valueobjects"
	"github.com/vantoan19/Petifies/server/services/post-service/internal/infra/publishers/kafka"
	mongo_comment "github.com/vantoan19/Petifies/server/services/post-service/internal/infra/repositories/comment/mongo"
	mongo_love "github.com/vantoan19/Petifies/server/services/post-service/internal/infra/repositories/love/mongo"
	mongo_post "github.com/vantoan19/Petifies/server/services/post-service/internal/infra/repositories/post/mongo"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

var logger = logging.New("PostService.PostSvc")

type postService struct {
	postRepo           postaggre.PostRepository
	commentRepo        commentaggre.CommentRepository
	loveRepo           loveaggre.LoveRepository
	postEventPublisher *kafka.PostEventPublisher
}

type PostConfiguration func(ps *postService) error

type PostService interface {
	CreatePost(ctx context.Context, post *models.CreatePostReq) (*postaggre.Post, error)
	LoveReactPost(ctx context.Context, req *models.LoveReactReq) (*loveaggre.Love, error)
	RemoveLoveReactPost(ctx context.Context, req *models.RemoveLoveReactReq) error
	EditPost(ctx context.Context, post *models.EditPostReq) (*postaggre.Post, error)
	ListPosts(ctx context.Context, req *models.ListPostsReq) ([]*postaggre.Post, error)
	GetLoveCount(ctx context.Context, postID uuid.UUID) (int, error)
	GetCommentCount(ctx context.Context, postID uuid.UUID) (int, error)
	GetPost(ctx context.Context, postID uuid.UUID) (*postaggre.Post, error)

	// Love endpoints
	GetLove(ctx context.Context, req *models.GetLoveReq) (*loveaggre.Love, error)
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

func WithMongoLoveRepository(client *mongo.Client) PostConfiguration {
	return func(ps *postService) error {
		repo := mongo_love.New(client)
		ps.loveRepo = repo
		return nil
	}
}

func WithKafkaPostEventPublisher(producer *producer.KafkaProducer, repo outbox_repo.EventRepository) PostConfiguration {
	return func(ps *postService) error {
		publisher := kafka.NewPostEventPublisher(producer, repo)
		ps.postEventPublisher = publisher
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
	err = ps.postEventPublisher.Publish(ctx, kafkamodels.PostEvent{
		ID:        createdPost.GetPostID(),
		AuthorID:  createdPost.GetAuthorID(),
		CreatedAt: createdPost.GetCreatedAt(),
		Status:    kafkamodels.POST_CREATED,
	})
	if err != nil {
		_, dbErr := ps.postRepo.DeleteByUUID(ctx, createdPost.GetPostID())
		if dbErr != nil {
			logger.ErrorData("Finished CreatePost: FAILED", logging.Data{"error": err.Error()})
			return nil, err
		}
		logger.ErrorData("Finished CreatePost: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish CreatePost: Successful")
	return createdPost, nil
}

func (ps *postService) LoveReactPost(ctx context.Context, req *models.LoveReactReq) (*loveaggre.Love, error) {
	logger.Info("Start LoveReactPost")

	post, err := ps.postRepo.GetByUUID(ctx, req.TargetID)
	if err != nil {
		logger.ErrorData("Finish LoveReactPost: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}
	err = post.AddLoveByAuthorIDAndSave(req.AuthorID, ps.loveRepo)
	if err != nil {
		logger.ErrorData("Finish LoveReactPost: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	love, err := ps.loveRepo.GetByTargetIDAndAuthorID(ctx, req.AuthorID, req.TargetID)
	if err != nil {
		logger.ErrorData("Finish LoveReactPost: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish LoveReactPost: Successful")
	return love, nil
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
			if err == mongo_post.ErrPostNotExist {
				return
			}
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

func (ps *postService) GetLoveCount(ctx context.Context, postID uuid.UUID) (int, error) {
	logger.Info("Start GetLoveCount")

	count, err := ps.loveRepo.CountLoveByTargetID(ctx, postID)
	if err != nil {
		logger.ErrorData("Finish GetLoveCount: Failed", logging.Data{"error": err.Error()})
		return 0, err
	}

	logger.Info("Finish LoveReactComment: Successful")
	return count, nil
}

func (ps *postService) GetCommentCount(ctx context.Context, postID uuid.UUID) (int, error) {
	logger.Info("Start GetCommentCount")

	count, err := ps.commentRepo.CountCommentByParentID(ctx, postID)
	if err != nil {
		logger.ErrorData("Finish GetCommentCount: Failed", logging.Data{"error": err.Error()})
		return 0, err
	}

	logger.Info("Finish GetCommentCount: Successful")
	return count, nil
}

func (ps *postService) GetPost(ctx context.Context, postID uuid.UUID) (*postaggre.Post, error) {
	logger.Info("Start GetPost")

	post, err := ps.postRepo.GetByUUID(ctx, postID)
	if err != nil {
		logger.ErrorData("Finish GetPost: Failed", logging.Data{"error": err.Error()})
	}

	logger.Info("Finish GetPost: Successful")
	return post, nil
}

func (ps *postService) RemoveLoveReactPost(ctx context.Context, req *models.RemoveLoveReactReq) error {
	logger.Info("Start RemoveLoveReactPost")

	post, err := ps.postRepo.GetByUUID(ctx, req.TargetID)
	if err != nil {
		logger.ErrorData("Finish RemoveLoveReactPost: Failed", logging.Data{"error": err.Error()})
		return err
	}
	err = post.RemoveLoveByAuthorIDAndDelete(req.AuthorID, ps.loveRepo)
	if err != nil {
		logger.ErrorData("Finish RemoveLoveReactPost: Failed", logging.Data{"error": err.Error()})
		return err
	}

	if exists, err := ps.loveRepo.ExistsLoveByTargetIDAndAuthorID(ctx, req.AuthorID, req.TargetID); err != nil {
		logger.ErrorData("Finish RemoveLoveReactPost: Failed", logging.Data{"error": err.Error()})
		return err
	} else if exists {
		logger.ErrorData("Finish RemoveLoveReactComment: Failed", logging.Data{"error": "failed to remove react"})
		return status.Errorf(codes.Internal, "failed to remove react")
	}

	logger.Info("Finish RemoveLoveReactPost: Successful")
	return nil
}

func (ps *postService) GetLove(ctx context.Context, req *models.GetLoveReq) (*loveaggre.Love, error) {
	logger.Info("Start GetLove")

	love, err := ps.loveRepo.GetByTargetIDAndAuthorID(ctx, req.AuthorID, req.TargetID)
	if err != nil {
		logger.ErrorData("Finish GetLove: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetLove: Successful")
	return love, nil
}
