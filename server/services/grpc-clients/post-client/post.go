package postclient

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	postProtoV1 "github.com/vantoan19/Petifies/proto/post-service/v1"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/models"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/translators"
)

var logger = logging.New("PostClient")

type postClient struct {
	createPost    endpoint.Endpoint
	createComment endpoint.Endpoint
	loveReact     endpoint.Endpoint
	editPost      endpoint.Endpoint
	editComment   endpoint.Endpoint
	listComments  endpoint.Endpoint
	listPosts     endpoint.Endpoint
}

type PostClient interface {
	CreatePost(ctx context.Context, req *models.CreatePostReq) (*models.Post, error)
	CreateComment(ctx context.Context, req *models.CreateCommentReq) (*models.Comment, error)
	LoveReact(ctx context.Context, req *models.LoveReactReq) (*models.Love, error)
	EditPost(ctx context.Context, req *models.EditPostReq) (*models.Post, error)
	EditComment(ctx context.Context, req *models.EditCommentReq) (*models.Comment, error)
	ListPosts(ctx context.Context, req *models.ListPostsReq) (*models.ListPostsResp, error)
	ListComments(ctx context.Context, req *models.ListCommentsReq) (*models.ListCommentsResp, error)
}

func New(conn *grpc.ClientConn) PostClient {
	return &postClient{
		createPost: grpctransport.NewClient(
			conn,
			"PostService",
			"CreatePost",
			translators.EncodeCreatePostRequest,
			translators.DecodePostResponse,
			commonProto.Post{},
		).Endpoint(),
		createComment: grpctransport.NewClient(
			conn,
			"PostService",
			"CreateComment",
			translators.EncodeCreateCommentRequest,
			translators.DecodeCommentResponse,
			commonProto.Comment{},
		).Endpoint(),
		loveReact: grpctransport.NewClient(
			conn,
			"PostService",
			"LoveReact",
			translators.EncodeLoveReactRequest,
			translators.DecodeLoveResponse,
			commonProto.Love{},
		).Endpoint(),
		editPost: grpctransport.NewClient(
			conn,
			"PostService",
			"EditPost",
			translators.EncodeEditPostRequest,
			translators.DecodePostResponse,
			commonProto.Post{},
		).Endpoint(),
		editComment: grpctransport.NewClient(
			conn,
			"PostService",
			"EditComment",
			translators.EncodeEditCommentRequest,
			translators.DecodeCommentResponse,
			commonProto.Comment{},
		).Endpoint(),
		listPosts: grpctransport.NewClient(
			conn,
			"PostService",
			"ListPosts",
			translators.EncodeListPostsRequest,
			translators.DecodeListPostsResponse,
			postProtoV1.ListPostsResponse{},
		).Endpoint(),
		listComments: grpctransport.NewClient(
			conn,
			"PostService",
			"ListComments",
			translators.EncodeListCommentsRequest,
			translators.DecodeListCommentsResponse,
			postProtoV1.ListCommentsResponse{},
		).Endpoint(),
	}
}

func (pc *postClient) CreatePost(ctx context.Context, req *models.CreatePostReq) (*models.Post, error) {
	logger.Info("Start CreatePost")

	resp, err := pc.createPost(ctx, req)
	if err != nil {
		logger.ErrorData("Finish CreatePost: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish CreatePost: Successful")
	return resp.(*models.Post), nil
}

func (pc *postClient) CreateComment(ctx context.Context, req *models.CreateCommentReq) (*models.Comment, error) {
	logger.Info("Start CreateComment")

	resp, err := pc.createComment(ctx, req)
	if err != nil {
		logger.ErrorData("Finish CreateComment: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish CreateComment: Successful")
	return resp.(*models.Comment), nil
}

func (pc *postClient) LoveReact(ctx context.Context, req *models.LoveReactReq) (*models.Love, error) {
	logger.Info("Start LoveReact")

	resp, err := pc.loveReact(ctx, req)
	if err != nil {
		logger.ErrorData("Finish LoveReact: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish LoveReact: Successful")
	return resp.(*models.Love), nil
}

func (pc *postClient) EditPost(ctx context.Context, req *models.EditPostReq) (*models.Post, error) {
	logger.Info("Start EditPost")

	resp, err := pc.editPost(ctx, req)
	if err != nil {
		logger.ErrorData("Finish EditPost: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish EditPost: Successful")
	return resp.(*models.Post), nil
}

func (pc *postClient) EditComment(ctx context.Context, req *models.EditCommentReq) (*models.Comment, error) {
	logger.Info("Start EditComment")

	resp, err := pc.editComment(ctx, req)
	if err != nil {
		logger.ErrorData("Finish EditComment: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish EditComment: Successful")
	return resp.(*models.Comment), nil
}

func (pc *postClient) ListPosts(ctx context.Context, req *models.ListPostsReq) (*models.ListPostsResp, error) {
	logger.Info("Start ListPosts")

	resp, err := pc.listPosts(ctx, req)
	if err != nil {
		logger.ErrorData("Finish ListPosts: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListPosts: Successful")
	return resp.(*models.ListPostsResp), nil
}

func (pc *postClient) ListComments(ctx context.Context, req *models.ListCommentsReq) (*models.ListCommentsResp, error) {
	logger.Info("Start ListComments")

	resp, err := pc.listComments(ctx, req)
	if err != nil {
		logger.ErrorData("Finish ListComments: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListComments: Successful")
	return resp.(*models.ListCommentsResp), nil
}
