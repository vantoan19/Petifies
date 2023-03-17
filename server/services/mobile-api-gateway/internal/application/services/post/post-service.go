package services

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	postclient "github.com/vantoan19/Petifies/server/services/grpc-clients/post-client"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/pkg/models"
	postModels "github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
)

var logger = logging.New("MobileGateway.PostService")

type PostConfiguration func(us *postService) error

type postService struct {
	postClient postclient.PostClient
}

type PostService interface {
	CreatePost(ctx context.Context, req *models.UserCreatePostReq) (*postModels.Post, error)
	CreateComment(ctx context.Context, req *models.UserCreateCommentReq) (*postModels.Comment, error)
	EditPost(ctx context.Context, req *models.UserEditPostReq) (*postModels.Post, error)
	EditComment(ctx context.Context, req *models.UserEditCommentReq) (*postModels.Comment, error)
}

func NewPostService(conn *grpc.ClientConn, cfgs ...PostConfiguration) (PostService, error) {
	ps := &postService{
		postClient: postclient.New(conn),
	}
	for _, cfg := range cfgs {
		err := cfg(ps)
		if err != nil {
			return nil, err
		}
	}
	return ps, nil
}

func (ps *postService) CreatePost(ctx context.Context, req *models.UserCreatePostReq) (*postModels.Post, error) {
	logger.Info("Start CreatePost")

	userID, err := commonutils.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	logger.Info("Executing CreatePost: delegating the request to PostService")
	resp, err := ps.postClient.CreatePost(ctx, &postModels.CreatePostReq{
		AuthorID:    userID,
		TextContent: req.TextContent,
		Images:      req.Images,
		Videos:      req.Videos,
	})
	if err != nil {
		logger.ErrorData("Finished CreatePost: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished CreatePost: SUCCESSFUL")
	return resp, nil
}

func (ps *postService) CreateComment(ctx context.Context, req *models.UserCreateCommentReq) (*postModels.Comment, error) {
	logger.Info("Start CreateComment")

	userID, err := commonutils.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	logger.Info("Executing CreateComment: delegating the request to PostService")
	resp, err := ps.postClient.CreateComment(ctx, &postModels.CreateCommentReq{
		AuthorID:     userID,
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

	logger.Info("Finished CreateComment: SUCCESSFUL")
	return resp, nil
}

func (ps *postService) EditPost(ctx context.Context, req *models.UserEditPostReq) (*postModels.Post, error) {
	logger.Info("Start EditPost")

	userID, err := commonutils.GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	listResp, err := ps.postClient.ListPosts(ctx, &postModels.ListPostsReq{
		PostIDs: []uuid.UUID{req.PostID},
	})
	if err != nil {
		logger.ErrorData("Finished EditPost: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	if listResp.Posts[0].AuthorID != userID {
		logger.ErrorData("Finished EditPost: FAILED", logging.Data{"error": err.Error()})
		return nil, status.Errorf(codes.PermissionDenied, "editor doesn't have permission to edit this post")
	}

	logger.Info("Executing EditPost: delegating the request to PostService")
	resp, err := ps.postClient.EditPost(ctx, &postModels.EditPostReq{
		ID:      req.PostID,
		Content: req.Content,
		Images:  req.Images,
		Videos:  req.Videos,
	})
	if err != nil {
		logger.ErrorData("Finished EditPost: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished EditPost: SUCCESSFUL")
	return resp, nil
}

func (ps *postService) EditComment(ctx context.Context, req *models.UserEditCommentReq) (*postModels.Comment, error) {
	logger.Info("Start EditComment")

	userID, err := commonutils.GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	listResp, err := ps.postClient.ListComments(ctx, &postModels.ListCommentsReq{
		CommentIDs: []uuid.UUID{req.CommentID},
	})
	if err != nil {
		logger.ErrorData("Finished EditComment: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	if listResp.Comments[0].AuthorID != userID {
		logger.ErrorData("Finished EditComment: FAILED", logging.Data{"error": err.Error()})
		return nil, status.Errorf(codes.PermissionDenied, "editor doesn't have permission to edit this comment")
	}

	logger.Info("Executing EditComment: delegating the request to PostService")
	resp, err := ps.postClient.EditComment(ctx, &postModels.EditCommentReq{
		ID:      req.CommentID,
		Content: req.Content,
		Image:   req.Image,
		Video:   req.Video,
	})
	if err != nil {
		logger.ErrorData("Finished EditComment: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished EditComment: SUCCESSFUL")
	return resp, nil
}
