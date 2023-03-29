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
	redisLoveCache "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/infra/repositories/love/redis"
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
	loveCacheRepo    repositories.LoveCacheRepository
}

type PostService interface {
	CreatePost(ctx context.Context, req *models.UserCreatePostReq) (*models.PostWithUserInfo, error)
	CreateComment(ctx context.Context, req *models.UserCreateCommentReq) (*models.CommentWithUserInfo, error)
	ToggleLoveReactPost(ctx context.Context, req *models.UserToggleLoveReq) (*models.UserToggleLoveResp, error)

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
	ListCommentsWithUserInfosByParentID(ctx context.Context, parentID uuid.UUID, pageSize int, afterCommentID uuid.UUID) ([]*models.CommentWithUserInfo, error)

	// Love Service
	GetLove(ctx context.Context, authorID, targetID uuid.UUID) (*postModels.Love, error)
	ExistsLove(ctx context.Context, authorID, targetID uuid.UUID) (bool, error)
	GetLoveWithUserInfo(ctx context.Context, authorID, targetID uuid.UUID) (*models.LoveWithUserInfo, error)
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

func WithRedisLoveCacheRepository(client *redis.Client) PostConfiguration {
	return func(ps *postService) error {
		repo := redisLoveCache.NewRedisLoveCacheRepository(client)
		ps.loveCacheRepo = repo
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
		err := ps.postCacheRepo.SetPostContent(context.Background(), postResp.ID, *postResp)
		if err != nil {
			logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
		}
	}()

	logger.Info("Finished CreatePost: SUCCESSFUL")
	return aggregatePostWithUserInfo(postResp, userResp, 0, 0, false), nil
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
		err := ps.commentCacheRepo.SetCommentContent(context.Background(), commentResp.ID, *commentResp)
		if err != nil {
			logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
		}
	}()

	logger.Info("Finished CreateComment: SUCCESSFUL")
	return aggregateCommentWithUserInfo(commentResp, userResp, 0, 0, false), nil
}

func (ps *postService) GetPostContent(ctx context.Context, postID uuid.UUID) (*postModels.Post, error) {
	logger.Info("Start GetPostContent")

	// Get from cache
	if exist, _ := ps.postCacheRepo.ExistsPostContent(ctx, postID); exist {
		logger.Info("Executing GetPostContent: getting post content info from cache")
		post, err := ps.postCacheRepo.GetPostContent(ctx, postID)
		if err != nil {
			logger.WarningData("Executing GetPostContent: failed to get post content from cache", logging.Data{"error": err.Error()})
			return nil, err
		}
		return post, nil
	}

	// Get from user service
	logger.Info("Executing GetPostContent: forwarding the request to Post Service")
	resp, err := ps.postClient.GetPost(ctx, &postModels.GetPostReq{PostID: postID})
	if err != nil {
		logger.ErrorData("Finished GetPostContent: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}
	// save to cache
	go func() {
		err := ps.postCacheRepo.SetPostContent(context.Background(), postID, *resp)
		if err != nil {
			logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
		}
	}()

	logger.Info("Finished GetPostContent: SUCCESSFUL")
	return resp, nil
}

func (ps *postService) GetPostLoveCount(ctx context.Context, postID uuid.UUID) (int, error) {
	logger.Info("Start GetPostLoveCount")

	// Get from post service
	logger.Info("Executing GetPostLoveCount: forwarding the request to Post Service")
	resp, err := ps.postClient.GetLoveCount(ctx, &postModels.GetLoveCountReq{TargetID: postID, IsPostParent: true})
	if err != nil {
		logger.ErrorData("Finished GetPostLoveCount: FAILED", logging.Data{"error": err.Error()})
		return 0, err
	}

	logger.Info("Finished GetPostLoveCount: SUCCESSFUL")
	return resp.Count, nil
}

func (ps *postService) GetPostCommentCount(ctx context.Context, postID uuid.UUID) (int, error) {
	logger.Info("Start GetPostCommentCount")

	// Get from post service
	logger.Info("Executing GetPostCommentCount: forwarding the request to Post Service")
	resp, err := ps.postClient.GetCommentCount(ctx, &postModels.GetCommentCountReq{ParentID: postID, IsPostParent: true})
	if err != nil {
		logger.ErrorData("Finished GetPostCommentCount: FAILED", logging.Data{"error": err.Error()})
		return 0, err
	}

	logger.Info("Finished GetPostCommentCount: SUCCESSFUL")
	return resp.Count, nil
}

func (ps *postService) GetCommentContent(ctx context.Context, commentID uuid.UUID) (*postModels.Comment, error) {
	logger.Info("Start GetCommentContent")

	// Get from post service
	logger.Info("Executing GetCommentContent: forwarding the request to Post Service")
	resp, err := ps.postClient.GetComment(ctx, &postModels.GetCommentReq{CommentID: commentID})
	if err != nil {
		logger.ErrorData("Finished GetCommentContent: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished GetCommentContent: SUCCESSFUL")
	return resp, nil
}

func (ps *postService) GetCommentLoveCount(ctx context.Context, commentID uuid.UUID) (int, error) {
	logger.Info("Start GetCommentLoveCount")

	// Get from post service
	logger.Info("Executing GetCommentLoveCount: forwarding the request to Post Service")
	resp, err := ps.postClient.GetLoveCount(ctx, &postModels.GetLoveCountReq{TargetID: commentID, IsPostParent: false})
	if err != nil {
		logger.ErrorData("Finished GetCommentLoveCount: FAILED", logging.Data{"error": err.Error()})
		return 0, err
	}

	logger.Info("Finished GetCommentLoveCount: SUCCESSFUL")
	return resp.Count, nil
}

func (ps *postService) GetCommentSubCommentCount(ctx context.Context, commendID uuid.UUID) (int, error) {
	logger.Info("Start GetCommentSubCommentCount")

	// Get from post service
	logger.Info("Executing GetCommentSubCommentCount: forwarding the request to Post Service")
	resp, err := ps.postClient.GetCommentCount(ctx, &postModels.GetCommentCountReq{ParentID: commendID, IsPostParent: false})
	if err != nil {
		logger.ErrorData("Finished GetCommentSubCommentCount: FAILED", logging.Data{"error": err.Error()})
		return 0, err
	}

	logger.Info("Finished GetCommentSubCommentCount: SUCCESSFUL")
	return resp.Count, nil
}

func (ps *postService) GetLove(ctx context.Context, authorID, targetID uuid.UUID) (*postModels.Love, error) {
	logger.Info("Start GetLove")

	// Get from cache
	if exist, _ := ps.loveCacheRepo.ExistsLove(ctx, authorID, targetID); exist {
		logger.Info("Executing GetLove: getting love info from cache")
		love, err := ps.loveCacheRepo.GetLove(ctx, authorID, targetID)
		if err != nil {
			logger.WarningData("Executing GetLove: failed to get love", logging.Data{"error": err.Error()})
		}
		return love, nil
	}

	// Get from post service
	logger.Info("Executing GetLove: forwarding the request to Post Service")
	resp, err := ps.postClient.GetLove(ctx, &postModels.GetLoveReq{AuthorID: authorID, TargetID: targetID})
	if err != nil {
		logger.ErrorData("Finished GetLove: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}
	// save to cache
	go func() {
		err := ps.loveCacheRepo.SetLove(context.Background(), authorID, targetID, *resp)
		if err != nil {
			logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
		}
	}()

	logger.Info("Finished GetLove: SUCCESSFUL")
	return resp, nil
}

func (ps *postService) ExistsLove(ctx context.Context, authorID, targetID uuid.UUID) (bool, error) {
	logger.Info("Start ExistsLove")

	love, err := ps.GetLove(ctx, authorID, targetID)
	if e, ok := status.FromError(err); ok {
		if e.Code() == codes.NotFound {
			return false, nil
		}
	}
	if err != nil {
		logger.ErrorData("Finished ExistsLove: FAILED", logging.Data{"error": err.Error()})
		return false, err
	}

	logger.Info("Finished ExistsLove: SUCCESSFUL")
	return (love != nil), nil
}

func (ps *postService) GetLoveWithUserInfo(ctx context.Context, authorID, targetID uuid.UUID) (*models.LoveWithUserInfo, error) {
	logger.Info("Start GetLoveWithUserInfo")

	userResp, err := ps.userService.GetMyInfo(ctx)
	if err != nil {
		logger.ErrorData("Finished GetLoveWithUserInfo: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	love, err := ps.GetLove(ctx, authorID, targetID)
	if err != nil {
		logger.ErrorData("Finished GetLoveWithUserInfo: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished GetLoveWithUserInfo: SUCCESSFUL")
	return aggregateLoveWithUserInfo(love, userResp), nil
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
	hasReacted, err := ps.ExistsLove(ctx, userResp.ID, postResp.ID)
	if err != nil {
		logger.ErrorData("Finished EditPost: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	// save to cache
	go func() {
		err := ps.postCacheRepo.SetPostContent(context.Background(), req.PostID, *postResp)
		if err != nil {
			logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
		}
	}()

	logger.Info("Finished EditPost: SUCCESSFUL")
	return aggregatePostWithUserInfo(postResp, userResp, loveCount, commentCount, hasReacted), nil
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
	hasReacted, err := ps.ExistsLove(ctx, userResp.ID, commentResp.ID)
	if err != nil {
		logger.ErrorData("Finished EditPost: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	// save to cache
	go func() {
		err := ps.commentCacheRepo.SetCommentContent(context.Background(), req.CommentID, *commentResp)
		if err != nil {
			logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
		}
	}()

	logger.Info("Finished EditComment: SUCCESSFUL")
	return aggregateCommentWithUserInfo(commentResp, userResp, loveCount, commentCount, hasReacted), nil
}

func (ps *postService) GetPostWithUserInfo(ctx context.Context, postID uuid.UUID) (*models.PostWithUserInfo, error) {
	logger.Info("Start GetPostWithUserInfo")

	var (
		post         *postModels.Post
		loveCount    int
		commentCount int
		hasReacted   bool
		wg           sync.WaitGroup
		errsChan     = make(chan error, 4)
	)

	wg.Add(4)
	go func() {
		defer wg.Done()
		post_, err := ps.GetPostContent(ctx, postID)
		if err != nil {
			errsChan <- err
		}
		post = post_
	}()

	go func() {
		defer wg.Done()
		loveCount_, err := ps.GetPostLoveCount(ctx, postID)
		if err != nil {
			errsChan <- err
		}
		loveCount = loveCount_
	}()

	go func() {
		defer wg.Done()
		commentCount_, err := ps.GetPostCommentCount(ctx, postID)
		if err != nil {
			errsChan <- err
		}
		commentCount = commentCount_
	}()

	go func() {
		defer wg.Done()
		userID, err := commonutils.GetUserID(ctx)
		if err != nil {
			errsChan <- err
			return
		}
		hasReacted_, err := ps.ExistsLove(ctx, userID, postID)
		if err != nil {
			errsChan <- err
			return
		}
		hasReacted = hasReacted_
	}()

	wg.Wait()
	close(errsChan)
	errs := commonutils.ToSlice(errsChan)
	if len(errs) > 0 {
		logger.ErrorData("Finish GetPostWithUserInfo: Failed", logging.Data{"error": errs[0].Error()})
		return nil, errs[0]
	}

	user, err := ps.userService.GetUser(ctx, post.AuthorID)
	if err != nil {
		logger.ErrorData("Finished GetPostWithUserInfo: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished GetPostWithUserInfo: SUCCESSFUL")
	return aggregatePostWithUserInfo(post, user, loveCount, commentCount, hasReacted), nil
}

func (ps *postService) GetCommentWithUserInfo(ctx context.Context, commentID uuid.UUID) (*models.CommentWithUserInfo, error) {
	logger.Info("Start GetCommentWithUserInfo")

	var (
		comment         *postModels.Comment
		loveCount       int
		subcommentCount int
		hasReacted      bool
		wg              sync.WaitGroup
		errsChan        = make(chan error, 4)
	)

	wg.Add(4)
	go func() {
		defer wg.Done()
		comment_, err := ps.GetCommentContent(ctx, commentID)
		if err != nil {
			errsChan <- err
		}
		comment = comment_
	}()

	go func() {
		defer wg.Done()
		loveCount_, err := ps.GetCommentLoveCount(ctx, commentID)
		if err != nil {
			errsChan <- err
		}
		loveCount = loveCount_
	}()

	go func() {
		defer wg.Done()
		subcommentCount_, err := ps.GetCommentSubCommentCount(ctx, commentID)
		if err != nil {
			errsChan <- err
		}
		subcommentCount = subcommentCount_
	}()

	go func() {
		defer wg.Done()
		userID, err := commonutils.GetUserID(ctx)
		if err != nil {
			errsChan <- err
			return
		}
		hasReacted_, err := ps.ExistsLove(ctx, userID, commentID)
		if err != nil {
			errsChan <- err
			return
		}
		hasReacted = hasReacted_
	}()

	wg.Wait()
	close(errsChan)
	errs := commonutils.ToSlice(errsChan)
	if len(errs) > 0 {
		logger.ErrorData("Finish GetCommentWithUserInfo: Failed", logging.Data{"error": errs[0].Error()})
		return nil, errs[0]
	}

	user, err := ps.userService.GetUser(ctx, comment.AuthorID)
	if err != nil {
		logger.ErrorData("Finished GetCommentWithUserInfo: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished GetCommentWithUserInfo: SUCCESSFUL")
	return aggregateCommentWithUserInfo(comment, user, loveCount, subcommentCount, hasReacted), nil
}

func (ps *postService) ListPostsWithUserInfos(ctx context.Context, postIDs []uuid.UUID) ([]*models.PostWithUserInfo, error) {
	logger.Info("Start ListPostsWithUserInfos")

	var wg sync.WaitGroup
	resultsChan := make(chan *models.PostWithUserInfo, len(postIDs))
	errsChan := make(chan error, len(postIDs))
	postMap := make(map[uuid.UUID]*models.PostWithUserInfo)

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
	for _, p := range results {
		postMap[p.ID] = p
	}

	logger.Info("Finish ListPostsWithUserInfos: Successful")
	return commonutils.Map2(postIDs, func(pId uuid.UUID) *models.PostWithUserInfo { return postMap[pId] }), nil
}

func (ps *postService) ListCommentsWithUserInfos(ctx context.Context, commentIDs []uuid.UUID) ([]*models.CommentWithUserInfo, error) {
	logger.Info("Start ListCommentsWithUserInfos")

	var wg sync.WaitGroup
	resultsChan := make(chan *models.CommentWithUserInfo, len(commentIDs))
	errsChan := make(chan error, len(commentIDs))
	commentMap := make(map[uuid.UUID]*models.CommentWithUserInfo)

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
	for _, c := range results {
		commentMap[c.ID] = c
	}

	logger.Info("Finish ListPostsWithUserInfos: Successful")
	return commonutils.Map2(commentIDs, func(cId uuid.UUID) *models.CommentWithUserInfo { return commentMap[cId] }), nil
}

func (ps *postService) ToggleLoveReactPost(ctx context.Context, req *models.UserToggleLoveReq) (*models.UserToggleLoveResp, error) {
	logger.Info("Start ToggleLoveReactPost")

	userID, err := commonutils.GetUserID(ctx)
	if err != nil {
		logger.ErrorData("Finished ToggleLoveReactPost: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	exists, err := ps.ExistsLove(ctx, userID, req.TargetID)
	if err != nil {
		logger.ErrorData("Finished ToggleLoveReactPost: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	if !exists {
		logger.Info("Executing ToggleLoveReactPost: user hasn't reacted to the post, toggle the love react on")
		loveResp, err := ps.postClient.LoveReact(ctx, &postModels.LoveReactReq{
			TargetID:     req.TargetID,
			IsTargetPost: req.IsPostTarget,
			AuthorID:     userID,
		})
		if err != nil {
			logger.ErrorData("Finished ToggleLoveReactPost: FAILED", logging.Data{"error": err.Error()})
			return nil, err
		}

		go func() {
			err := ps.loveCacheRepo.SetLove(context.Background(), loveResp.AuthorID, loveResp.TargetID, *loveResp)
			if err != nil {
				logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
			}
		}()

		logger.Info("Finished ToggleLoveReactPost: Successful")
		return &models.UserToggleLoveResp{
			HasReacted: true,
		}, nil
	} else {
		logger.Info("Executing ToggleLoveReactPost: user already reacted to the post, toggle the love react off")
		_, err := ps.postClient.RemoveLoveReact(ctx, &postModels.RemoveLoveReactReq{
			TargetID:     req.TargetID,
			IsTargetPost: req.IsPostTarget,
			AuthorID:     userID,
		})
		if err != nil {
			logger.ErrorData("Finished ToggleLoveReactPost: FAILED", logging.Data{"error": err.Error()})
			return nil, err
		}

		go func() {
			err := ps.loveCacheRepo.RemoveLove(context.Background(), userID, req.TargetID)
			if err != nil {
				logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
			}
		}()

		logger.Info("Finished ToggleLoveReactPost: Successful")
		return &models.UserToggleLoveResp{
			HasReacted: false,
		}, nil
	}
}

func (ps *postService) ListCommentsWithUserInfosByParentID(ctx context.Context, parentID uuid.UUID, pageSize int, afterCommentID uuid.UUID) ([]*models.CommentWithUserInfo, error) {
	logger.Info("Start ListCommentsWithUserInfosByParentID")

	listResp, err := ps.postClient.ListCommentIDsByParentID(ctx, &postModels.ListCommentIDsByParentIDReq{
		ParentID:       parentID,
		PageSize:       pageSize,
		AfterCommentID: afterCommentID,
	})
	if err != nil {
		logger.ErrorData("Finished ListCommentsWithUserInfosByParentID: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	comments, err := ps.ListCommentsWithUserInfos(ctx, listResp.CommentIDs)
	if err != nil {
		logger.ErrorData("Finished ListCommentsWithUserInfosByParentID: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished ListCommentsWithUserInfosByParentID: Successful")
	return comments, nil
}

func aggregatePostWithUserInfo(post *postModels.Post, user *userModels.User, loveCount, commentCount int, hasReacted bool) *models.PostWithUserInfo {
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
		HasReacted:   hasReacted,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
	}
}

func aggregateCommentWithUserInfo(comment *postModels.Comment, user *userModels.User, loveCount, subCommentCount int, hasReacted bool) *models.CommentWithUserInfo {
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
		HasReacted:      hasReacted,
		LoveCount:       loveCount,
		SubcommentCount: subCommentCount,
		CreatedAt:       comment.CreatedAt,
		UpdatedAt:       comment.UpdatedAt,
	}
}

func aggregateLoveWithUserInfo(love *postModels.Love, user *userModels.User) *models.LoveWithUserInfo {
	return &models.LoveWithUserInfo{
		ID:           love.ID,
		TargetID:     love.TargetID,
		IsPostTarget: love.IsPostTarget,
		Author: models.BasicUserInfo{
			ID:         user.ID,
			Email:      user.Email,
			UserAvatar: "",
			FirstName:  user.FirstName,
			LastName:   user.LastName,
		},
		CreatedAt: love.CreatedAt,
	}
}
