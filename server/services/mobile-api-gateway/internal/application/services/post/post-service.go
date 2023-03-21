package postservice

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	postclient "github.com/vantoan19/Petifies/server/services/grpc-clients/post-client"
	userclient "github.com/vantoan19/Petifies/server/services/grpc-clients/user-client"
	userservice "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/application/services/user"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/domain/repositories"
	redisCommentCache "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/infra/repositories/comment/redis"
	redisPostCache "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/infra/repositories/post/redis"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/pkg/models"
	postModels "github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
	userModels "github.com/vantoan19/Petifies/server/services/user-service/pkg/models"
)

var logger = logging.New("MobileGateway.PostService")

type PostConfiguration func(ps *postService) error

type postService struct {
	postClient       postclient.PostClient
	userClient       userclient.UserClient
	userService      userservice.UserService
	postCacheRepo    repositories.PostCacheRepository
	commentCacheRepo repositories.CommentCacheRepository
}

type PostService interface {
	CreatePost(ctx context.Context, req *models.UserCreatePostReq) (*models.PostWithUserInfo, error)
	CreateComment(ctx context.Context, req *models.UserCreateCommentReq) (*models.CommentWithUserInfo, error)
	EditPost(ctx context.Context, req *models.UserEditPostReq) (*models.PostWithUserInfo, error)
	EditComment(ctx context.Context, req *models.UserEditCommentReq) (*models.CommentWithUserInfo, error)
	GetPostContent(ctx context.Context, postID uuid.UUID) (*postModels.Post, error)
	GetPostLoveCount(ctx context.Context, postID uuid.UUID) (int, error)
	GetPostCommentCount(ctx context.Context, postID uuid.UUID) (int, error)
	GetCommentContent(ctx context.Context, commendID uuid.UUID) (*postModels.Comment, error)
	GetCommentLoveCount(ctx context.Context, commendID uuid.UUID) (int, error)
	GetCommentSubCommentCount(ctx context.Context, commendID uuid.UUID) (int, error)
	GetPostWithUserInfo(ctx context.Context, postID uuid.UUID) (*models.PostWithUserInfo, error)
	GetCommentWithUserInfo(ctx context.Context, commentID uuid.UUID) (*models.CommentWithUserInfo, error)
	ListPostsWithUserInfos(ctx context.Context, postIDs []uuid.UUID) ([]*models.PostWithUserInfo, error)
	ListCommentsWithUserInfos(ctx context.Context, commentIDs []uuid.UUID) ([]*models.CommentWithUserInfo, error)
}

func NewPostService(postClientConn *grpc.ClientConn, userClientConn *grpc.ClientConn, userService userservice.UserService, cfgs ...PostConfiguration) (PostService, error) {
	ps := &postService{
		postClient:  postclient.New(postClientConn),
		userClient:  userclient.New(userClientConn),
		userService: userService,
	}
	for _, cfg := range cfgs {
		err := cfg(ps)
		if err != nil {
			return nil, err
		}
	}
	return ps, nil
}

func WithRedisPostCacheRepository(client *redis.Client) PostConfiguration {
	return func(ps *postService) error {
		repo := redisPostCache.NewRedisPostCacheRepository(client)
		ps.postCacheRepo = repo
		return nil
	}
}

func WithRedisCommentCacheRepository(client *redis.Client) PostConfiguration {
	return func(ps *postService) error {
		repo := redisCommentCache.NewRedisCommentCacheRepository(client)
		ps.commentCacheRepo = repo
		return nil
	}
}

func (ps *postService) CreatePost(ctx context.Context, req *models.UserCreatePostReq) (*models.PostWithUserInfo, error) {
	logger.Info("Start CreatePost")

	userResp, err := ps.userService.GetMyInfo(ctx)
	if err != nil {
		logger.ErrorData("Finished CreatePost: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Executing CreatePost: delegating the request to PostService")
	postResp, err := ps.postClient.CreatePost(ctx, &postModels.CreatePostReq{
		AuthorID:    userResp.ID,
		Visibility:  req.Visibility,
		Activity:    req.Activity,
		TextContent: req.TextContent,
		Images:      req.Images,
		Videos:      req.Videos,
	})
	if err != nil {
		logger.ErrorData("Finished CreatePost: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	// Set cache
	go func() {
		err := ps.postCacheRepo.SetPostContent(ctx, postResp.ID, *postResp)
		if err != nil {
			logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
		}
		err = ps.postCacheRepo.SetPostCommentCount(ctx, postResp.ID, 0)
		if err != nil {
			logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
		}
		err = ps.postCacheRepo.SetPostLoveCount(ctx, postResp.ID, 0)
		if err != nil {
			logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
		}
	}()

	logger.Info("Finished CreatePost: SUCCESSFUL")
	return aggregatePostWithUserInfo(postResp, userResp, 0, 0), nil
}

func (ps *postService) CreateComment(ctx context.Context, req *models.UserCreateCommentReq) (*models.CommentWithUserInfo, error) {
	logger.Info("Start CreateComment")

	userResp, err := ps.userService.GetMyInfo(ctx)
	if err != nil {
		logger.ErrorData("Finished CreatePost: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Executing CreateComment: delegating the request to PostService")
	commentResp, err := ps.postClient.CreateComment(ctx, &postModels.CreateCommentReq{
		AuthorID:     userResp.ID,
		PostID:       req.PostID,
		ParentID:     req.ParentID,
		IsPostParent: req.IsParentPost,
		Content:      req.Content,
		ImageContent: req.Image,
		VideoContent: req.Video,
	})
	if err != nil {
		logger.ErrorData("Finished CreateComment: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	// Set cache
	go func() {
		err := ps.commentCacheRepo.SetCommentContent(ctx, commentResp.ID, *commentResp)
		if err != nil {
			logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
		}
		err = ps.commentCacheRepo.SetCommentLoveCount(ctx, commentResp.ID, 0)
		if err != nil {
			logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
		}
		err = ps.commentCacheRepo.SetCommentSubCommentCount(ctx, commentResp.ID, 0)
		if err != nil {
			logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
		}
	}()

	logger.Info("Finished CreateComment: SUCCESSFUL")
	return aggregateCommentWithUserInfo(commentResp, userResp, 0, 0), nil
}

func (ps *postService) GetPostContent(ctx context.Context, postID uuid.UUID) (*postModels.Post, error) {
	logger.Info("Start GetPostContent")

	var post *postModels.Post
	// Get from cache
	if exist, err := ps.postCacheRepo.ExistsPostContent(ctx, postID); exist {
		logger.Info("Executing GetPostContent: getting post content info from cache")
		post_, err := ps.postCacheRepo.GetPostContent(ctx, postID)
		if err != nil {
			logger.ErrorData("Finished GetPostContent: FAILED", logging.Data{"error": err.Error()})
			return nil, err
		}
		post = post_
	} else if err != nil {
		logger.ErrorData("Finished GetPostContent: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	} else { // Get from user service
		logger.Info("Executing GetPostContent: forwarding the request to Post Service")
		resp, err := ps.postClient.GetPost(ctx, &postModels.GetPostReq{PostID: postID})
		if err != nil {
			logger.ErrorData("Finished GetPostContent: FAILED", logging.Data{"error": err.Error()})
			return nil, err
		}
		// save to cache
		go func() {
			err := ps.postCacheRepo.SetPostContent(ctx, postID, *resp)
			if err != nil {
				logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
			}
		}()
		post = resp
	}

	logger.Info("Finished GetPostContent: SUCCESSFUL")
	return post, nil
}

func (ps *postService) GetPostLoveCount(ctx context.Context, postID uuid.UUID) (int, error) {
	logger.Info("Start GetPostLoveCount")

	var count int
	// Get from cache
	if exist, err := ps.postCacheRepo.ExistsPostLoveCount(ctx, postID); exist {
		logger.Info("Executing GetPostLoveCount: getting post love count info from cache")
		count_, err := ps.postCacheRepo.GetPostLoveCount(ctx, postID)
		if err != nil {
			logger.ErrorData("Finished GetPostLoveCount: FAILED", logging.Data{"error": err.Error()})
			return 0, err
		}
		count = count_
	} else if err != nil {
		logger.ErrorData("Finished GetPostLoveCount: FAILED", logging.Data{"error": err.Error()})
		return 0, err
	} else { // Get from post service
		logger.Info("Executing GetPostLoveCount: forwarding the request to Post Service")
		resp, err := ps.postClient.GetLoveCount(ctx, &postModels.GetLoveCountReq{TargetID: postID, IsPostParent: true})
		if err != nil {
			logger.ErrorData("Finished GetPostLoveCount: FAILED", logging.Data{"error": err.Error()})
			return 0, err
		}
		// save to cache
		go func() {
			err := ps.postCacheRepo.SetPostLoveCount(ctx, postID, resp.Count)
			if err != nil {
				logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
			}
		}()
		count = resp.Count
	}

	logger.Info("Finished GetPostLoveCount: SUCCESSFUL")
	return count, nil
}

func (ps *postService) GetPostCommentCount(ctx context.Context, postID uuid.UUID) (int, error) {
	logger.Info("Start GetPostCommentCount")

	var count int
	// Get from cache
	if exist, err := ps.postCacheRepo.ExistsPostCommentCount(ctx, postID); exist {
		logger.Info("Executing GetPostCommentCount: getting post love count info from cache")
		count_, err := ps.postCacheRepo.GetPostCommentCount(ctx, postID)
		if err != nil {
			logger.ErrorData("Finished GetPostCommentCount: FAILED", logging.Data{"error": err.Error()})
			return 0, err
		}
		count = count_
	} else if err != nil {
		logger.ErrorData("Finished GetPostCommentCount: FAILED", logging.Data{"error": err.Error()})
		return 0, err
	} else { // Get from post service
		logger.Info("Executing GetPostCommentCount: forwarding the request to Post Service")
		resp, err := ps.postClient.GetCommentCount(ctx, &postModels.GetCommentCountReq{ParentID: postID, IsPostParent: true})
		if err != nil {
			logger.ErrorData("Finished GetPostCommentCount: FAILED", logging.Data{"error": err.Error()})
			return 0, err
		}
		// save to cache
		go func() {
			err := ps.postCacheRepo.SetPostCommentCount(ctx, postID, resp.Count)
			if err != nil {
				logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
			}
		}()
		count = resp.Count
	}

	logger.Info("Finished GetPostCommentCount: SUCCESSFUL")
	return count, nil
}

func (ps *postService) GetCommentContent(ctx context.Context, commentID uuid.UUID) (*postModels.Comment, error) {
	logger.Info("Start GetCommentContent")

	var comment *postModels.Comment
	// Get from cache
	if exist, err := ps.commentCacheRepo.ExistsCommentContent(ctx, commentID); exist {
		logger.Info("Executing GetCommentContent: getting comment content info from cache")
		comment_, err := ps.commentCacheRepo.GetCommentContent(ctx, commentID)
		if err != nil {
			logger.ErrorData("Finished GetCommentContent: FAILED", logging.Data{"error": err.Error()})
			return nil, err
		}
		comment = comment_
	} else if err != nil {
		logger.ErrorData("Finished GetCommentContent: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	} else { // Get from post service
		logger.Info("Executing GetCommentContent: forwarding the request to Post Service")
		resp, err := ps.postClient.GetComment(ctx, &postModels.GetCommentReq{CommentID: commentID})
		if err != nil {
			logger.ErrorData("Finished GetCommentContent: FAILED", logging.Data{"error": err.Error()})
			return nil, err
		}
		// save to cache
		go func() {
			err := ps.commentCacheRepo.SetCommentContent(ctx, commentID, *resp)
			if err != nil {
				logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
			}
		}()
		comment = resp
	}

	logger.Info("Finished GetCommentContent: SUCCESSFUL")
	return comment, nil
}

func (ps *postService) GetCommentLoveCount(ctx context.Context, commentID uuid.UUID) (int, error) {
	logger.Info("Start GetCommentLoveCount")

	var count int
	// Get from cache
	if exist, err := ps.commentCacheRepo.ExistsCommentLoveCount(ctx, commentID); exist {
		logger.Info("Executing GetCommentLoveCount: getting post love count info from cache")
		count_, err := ps.commentCacheRepo.GetCommentLoveCount(ctx, commentID)
		if err != nil {
			logger.ErrorData("Finished GetCommentLoveCount: FAILED", logging.Data{"error": err.Error()})
			return 0, err
		}
		count = count_
	} else if err != nil {
		logger.ErrorData("Finished GetCommentLoveCount: FAILED", logging.Data{"error": err.Error()})
		return 0, err
	} else { // Get from post service
		logger.Info("Executing GetCommentLoveCount: forwarding the request to Post Service")
		resp, err := ps.postClient.GetLoveCount(ctx, &postModels.GetLoveCountReq{TargetID: commentID, IsPostParent: false})
		if err != nil {
			logger.ErrorData("Finished GetCommentLoveCount: FAILED", logging.Data{"error": err.Error()})
			return 0, err
		}
		// save to cache
		go func() {
			err := ps.commentCacheRepo.SetCommentLoveCount(ctx, commentID, resp.Count)
			if err != nil {
				logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
			}
		}()
		count = resp.Count
	}

	logger.Info("Finished GetCommentLoveCount: SUCCESSFUL")
	return count, nil
}

func (ps *postService) GetCommentSubCommentCount(ctx context.Context, commendID uuid.UUID) (int, error) {
	logger.Info("Start GetCommentSubCommentCount")

	var count int
	// Get from cache
	if exist, err := ps.commentCacheRepo.ExistsCommentSubCommentCount(ctx, commendID); exist {
		logger.Info("Executing GetCommentSubCommentCount: getting post love count info from cache")
		count_, err := ps.commentCacheRepo.GetCommentSubCommentCount(ctx, commendID)
		if err != nil {
			logger.ErrorData("Finished GetCommentSubCommentCount: FAILED", logging.Data{"error": err.Error()})
			return 0, err
		}
		count = count_
	} else if err != nil {
		logger.ErrorData("Finished GetCommentSubCommentCount: FAILED", logging.Data{"error": err.Error()})
		return 0, err
	} else { // Get from post service
		logger.Info("Executing GetCommentSubCommentCount: forwarding the request to Post Service")
		resp, err := ps.postClient.GetCommentCount(ctx, &postModels.GetCommentCountReq{ParentID: commendID, IsPostParent: false})
		if err != nil {
			logger.ErrorData("Finished GetCommentSubCommentCount: FAILED", logging.Data{"error": err.Error()})
			return 0, err
		}
		// save to cache
		go func() {
			err := ps.commentCacheRepo.SetCommentSubCommentCount(ctx, commendID, resp.Count)
			if err != nil {
				logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
			}
		}()
		count = resp.Count
	}

	logger.Info("Finished GetCommentSubCommentCount: SUCCESSFUL")
	return count, nil
}

func (ps *postService) EditPost(ctx context.Context, req *models.UserEditPostReq) (*models.PostWithUserInfo, error) {
	logger.Info("Start EditPost")

	userResp, err := ps.userService.GetMyInfo(ctx)
	if err != nil {
		logger.ErrorData("Finished CreatePost: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	post, err := ps.GetPostContent(ctx, req.PostID)
	if err != nil {
		logger.ErrorData("Finished EditPost: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	if post.AuthorID != userResp.ID {
		logger.ErrorData("Finished EditPost: FAILED", logging.Data{"error": err.Error()})
		return nil, status.Errorf(codes.PermissionDenied, "editor doesn't have permission to edit this post")
	}

	logger.Info("Executing EditPost: delegating the request to PostService")
	postResp, err := ps.postClient.EditPost(ctx, &postModels.EditPostReq{
		ID:         req.PostID,
		Visibility: req.Visibility,
		Activity:   req.Activity,
		Content:    req.Content,
		Images:     req.Images,
		Videos:     req.Videos,
	})
	if err != nil {
		logger.ErrorData("Finished EditPost: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	loveCount, err := ps.GetPostLoveCount(ctx, req.PostID)
	if err != nil {
		logger.ErrorData("Finished EditPost: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}
	commentCount, err := ps.GetPostCommentCount(ctx, req.PostID)
	if err != nil {
		logger.ErrorData("Finished EditPost: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	// save to cache
	go func() {
		err := ps.postCacheRepo.SetPostContent(ctx, req.PostID, *postResp)
		if err != nil {
			logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
		}
	}()

	logger.Info("Finished EditPost: SUCCESSFUL")
	return aggregatePostWithUserInfo(postResp, userResp, loveCount, commentCount), nil
}

func (ps *postService) EditComment(ctx context.Context, req *models.UserEditCommentReq) (*models.CommentWithUserInfo, error) {
	logger.Info("Start EditComment")

	userResp, err := ps.userService.GetMyInfo(ctx)
	if err != nil {
		logger.ErrorData("Finished CreatePost: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	comment, err := ps.GetCommentContent(ctx, req.CommentID)
	if err != nil {
		logger.ErrorData("Finished EditComment: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	if comment.AuthorID != userResp.ID {
		logger.ErrorData("Finished EditComment: FAILED", logging.Data{"error": err.Error()})
		return nil, status.Errorf(codes.PermissionDenied, "editor doesn't have permission to edit this comment")
	}

	logger.Info("Executing EditComment: delegating the request to PostService")
	commentResp, err := ps.postClient.EditComment(ctx, &postModels.EditCommentReq{
		ID:      req.CommentID,
		Content: req.Content,
		Image:   req.Image,
		Video:   req.Video,
	})
	if err != nil {
		logger.ErrorData("Finished EditComment: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	loveCount, err := ps.GetCommentLoveCount(ctx, req.CommentID)
	if err != nil {
		logger.ErrorData("Finished EditPost: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}
	commentCount, err := ps.GetCommentSubCommentCount(ctx, req.CommentID)
	if err != nil {
		logger.ErrorData("Finished EditPost: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	// save to cache
	go func() {
		err := ps.commentCacheRepo.SetCommentContent(ctx, req.CommentID, *commentResp)
		if err != nil {
			logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
		}
	}()

	logger.Info("Finished EditComment: SUCCESSFUL")
	return aggregateCommentWithUserInfo(commentResp, userResp, loveCount, commentCount), nil
}

func (ps *postService) GetPostWithUserInfo(ctx context.Context, postID uuid.UUID) (*models.PostWithUserInfo, error) {
	logger.Info("Start GetPostWithUserInfo")

	postResp, err := ps.GetPostContent(ctx, postID)
	if err != nil {
		logger.ErrorData("Finished GetPostWithUserInfo: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	userResp, err := ps.userService.GetUser(ctx, postResp.AuthorID)
	if err != nil {
		logger.ErrorData("Finished GetPostWithUserInfo: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	loveCount, err := ps.GetPostLoveCount(ctx, postID)
	if err != nil {
		logger.ErrorData("Finished GetPostWithUserInfo: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	commentCount, err := ps.GetPostCommentCount(ctx, postID)
	if err != nil {
		logger.ErrorData("Finished GetPostWithUserInfo: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished GetPostWithUserInfo: SUCCESSFUL")
	return aggregatePostWithUserInfo(postResp, userResp, loveCount, commentCount), nil
}

func (ps *postService) GetCommentWithUserInfo(ctx context.Context, commentID uuid.UUID) (*models.CommentWithUserInfo, error) {
	logger.Info("Start GetCommentWithUserInfo")

	commentResp, err := ps.GetCommentContent(ctx, commentID)
	if err != nil {
		logger.ErrorData("Finished GetCommentWithUserInfo: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	userResp, err := ps.userService.GetUser(ctx, commentResp.AuthorID)
	if err != nil {
		logger.ErrorData("Finished GetCommentWithUserInfo: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	loveCount, err := ps.GetCommentLoveCount(ctx, commentID)
	if err != nil {
		logger.ErrorData("Finished GetCommentWithUserInfo: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	subcommentCount, err := ps.GetCommentSubCommentCount(ctx, commentID)
	if err != nil {
		logger.ErrorData("Finished GetCommentWithUserInfo: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished GetCommentWithUserInfo: SUCCESSFUL")
	return aggregateCommentWithUserInfo(commentResp, userResp, loveCount, subcommentCount), nil
}

func (ps *postService) ListPostsWithUserInfos(ctx context.Context, postIDs []uuid.UUID) ([]*models.PostWithUserInfo, error) {
	logger.Info("Start ListPostsWithUserInfos")

	var wg sync.WaitGroup
	resultsChan := make(chan *models.PostWithUserInfo, len(postIDs))
	errsChan := make(chan error, len(postIDs))

	for _, id := range postIDs {
		wg.Add(1)
		go func(id uuid.UUID) {
			defer wg.Done()
			post, err := ps.GetPostWithUserInfo(ctx, id)
			if e, ok := status.FromError(err); ok {
				if e.Code() == codes.NotFound {
					return
				}
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
	errs := commonutils.ToSlice(errsChan)
	results := commonutils.ToSlice(resultsChan)
	if len(errs) > 0 {
		logger.ErrorData("Finish ListPostsWithUserInfos: Failed", logging.Data{"error": errs[0].Error()})
		return nil, errs[0]
	}

	logger.Info("Finish ListPostsWithUserInfos: Successful")
	return results, nil
}

func (ps *postService) ListCommentsWithUserInfos(ctx context.Context, commentIDs []uuid.UUID) ([]*models.CommentWithUserInfo, error) {
	logger.Info("Start ListCommentsWithUserInfos")

	var wg sync.WaitGroup
	resultsChan := make(chan *models.CommentWithUserInfo, len(commentIDs))
	errsChan := make(chan error, len(commentIDs))

	for _, id := range commentIDs {
		wg.Add(1)
		go func(id uuid.UUID) {
			defer wg.Done()
			comment, err := ps.GetCommentWithUserInfo(ctx, id)
			if e, ok := status.FromError(err); ok {
				if e.Code() == codes.NotFound {
					return
				}
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
	errs := commonutils.ToSlice(errsChan)
	results := commonutils.ToSlice(resultsChan)
	if len(errs) > 0 {
		logger.ErrorData("Finish ListCommentsWithUserInfos: Failed", logging.Data{"error": errs[0].Error()})
		return nil, errs[0]
	}

	logger.Info("Finish ListPostsWithUserInfos: Successful")
	return results, nil
}

func aggregatePostWithUserInfo(post *postModels.Post, user *userModels.User, loveCount, commentCount int) *models.PostWithUserInfo {
	return &models.PostWithUserInfo{
		ID: post.ID,
		Author: models.BasicUserInfo{
			ID:         user.ID,
			Email:      user.Email,
			UserAvatar: "",
			FirstName:  user.FirstName,
			LastName:   user.LastName,
		},
		Content:      post.Content,
		Images:       post.Images,
		Videos:       post.Videos,
		LoveCount:    loveCount,
		CommentCount: commentCount,
		Visibility:   post.Visibility,
		Activity:     post.Activity,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
	}
}

func aggregateCommentWithUserInfo(comment *postModels.Comment, user *userModels.User, loveCount, subCommentCount int) *models.CommentWithUserInfo {
	return &models.CommentWithUserInfo{
		ID: comment.ID,
		Author: models.BasicUserInfo{
			ID:         user.ID,
			Email:      user.Email,
			UserAvatar: "",
			FirstName:  user.FirstName,
			LastName:   user.LastName,
		},
		PostID:          comment.PostID,
		ParentID:        comment.ParentID,
		IsPostParent:    comment.IsPostParent,
		Content:         comment.Content,
		Image:           comment.Image,
		Video:           comment.Video,
		LoveCount:       loveCount,
		SubcommentCount: subCommentCount,
		CreatedAt:       comment.CreatedAt,
		UpdatedAt:       comment.UpdatedAt,
	}
}
