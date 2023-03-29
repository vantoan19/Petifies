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
	createPost               endpoint.Endpoint
	createComment            endpoint.Endpoint
	loveReact                endpoint.Endpoint
	editPost                 endpoint.Endpoint
	editComment              endpoint.Endpoint
	listComments             endpoint.Endpoint
	listPosts                endpoint.Endpoint
	getLoveCount             endpoint.Endpoint
	getCommentCount          endpoint.Endpoint
	getPost                  endpoint.Endpoint
	getComment               endpoint.Endpoint
	removeLoveReact          endpoint.Endpoint
	getLove                  endpoint.Endpoint
	listCommentIDsByParentID endpoint.Endpoint
}

type PostClient interface {
	CreatePost(ctx context.Context, req *models.CreatePostReq) (*models.Post, error)
	CreateComment(ctx context.Context, req *models.CreateCommentReq) (*models.Comment, error)
	LoveReact(ctx context.Context, req *models.LoveReactReq) (*models.Love, error)
	EditPost(ctx context.Context, req *models.EditPostReq) (*models.Post, error)
	EditComment(ctx context.Context, req *models.EditCommentReq) (*models.Comment, error)
	ListPosts(ctx context.Context, req *models.ListPostsReq) (*models.ListPostsResp, error)
	ListComments(ctx context.Context, req *models.ListCommentsReq) (*models.ListCommentsResp, error)
	GetLoveCount(ctx context.Context, req *models.GetLoveCountReq) (*models.GetLoveCountResp, error)
	GetCommentCount(ctx context.Context, req *models.GetCommentCountReq) (*models.GetCommentCountResp, error)
	GetPost(ctx context.Context, req *models.GetPostReq) (*models.Post, error)
	GetComment(ctx context.Context, req *models.GetCommentReq) (*models.Comment, error)
	RemoveLoveReact(ctx context.Context, req *models.RemoveLoveReactReq) (*models.RemoveLoveReactResp, error)
	GetLove(ctx context.Context, req *models.GetLoveReq) (*models.Love, error)
	ListCommentIDsByParentID(ctx context.Context, req *models.ListCommentIDsByParentIDReq) (*models.ListCommentIDsByParentIDResp, error)
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
		getLoveCount: grpctransport.NewClient(
			conn,
			"PostService",
			"GetLoveCount",
			translators.EncodeGetLoveCountRequest,
			translators.DecodeGetLoveCountResponse,
			postProtoV1.GetLoveCountReponse{},
		).Endpoint(),
		getCommentCount: grpctransport.NewClient(
			conn,
			"PostService",
			"GetCommentCount",
			translators.EncodeGetCommentCountRequest,
			translators.DecodeGetCommentCountResponse,
			postProtoV1.GetCommentCountReponse{},
		).Endpoint(),
		getPost: grpctransport.NewClient(
			conn,
			"PostService",
			"GetPost",
			translators.EncodeGetPostRequest,
			translators.DecodePostResponse,
			commonProto.Post{},
		).Endpoint(),
		getComment: grpctransport.NewClient(
			conn,
			"PostService",
			"GetComment",
			translators.EncodeGetCommentRequest,
			translators.DecodeCommentResponse,
			commonProto.Comment{},
		).Endpoint(),
		removeLoveReact: grpctransport.NewClient(
			conn,
			"PostService",
			"RemoveLoveReact",
			translators.EncodeRemoveLoveReactRequest,
			translators.DecodeRemoveLoveReactResponse,
			postProtoV1.RemoveLoveReactResponse{},
		).Endpoint(),
		getLove: grpctransport.NewClient(
			conn,
			"PostService",
			"GetLove",
			translators.EncodeGetLoveRequest,
			translators.DecodeLoveResponse,
			commonProto.Love{},
		).Endpoint(),
		listCommentIDsByParentID: grpctransport.NewClient(
			conn,
			"PostService",
			"ListCommentIDsByParentID",
			translators.EncodeListCommentIDsByParentIDRequest,
			translators.DecodeListCommentIDsByParentIDResponse,
			postProtoV1.ListCommentIDsByParentIDResponse{},
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

func (pc *postClient) GetLoveCount(ctx context.Context, req *models.GetLoveCountReq) (*models.GetLoveCountResp, error) {
	logger.Info("Start GetLoveCount")

	resp, err := pc.getLoveCount(ctx, req)
	if err != nil {
		logger.ErrorData("Finish GetLoveCount: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetLoveCount: Successful")
	return resp.(*models.GetLoveCountResp), nil
}

func (pc *postClient) GetCommentCount(ctx context.Context, req *models.GetCommentCountReq) (*models.GetCommentCountResp, error) {
	logger.Info("Start GetCommentCount")

	resp, err := pc.getCommentCount(ctx, req)
	if err != nil {
		logger.ErrorData("Finish GetCommentCount: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetCommentCount: Successful")
	return resp.(*models.GetCommentCountResp), nil
}

func (pc *postClient) GetPost(ctx context.Context, req *models.GetPostReq) (*models.Post, error) {
	logger.Info("Start GetPost")

	resp, err := pc.getPost(ctx, req)
	if err != nil {
		logger.ErrorData("Finish GetPost: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetPost: Successful")
	return resp.(*models.Post), nil
}

func (pc *postClient) GetComment(ctx context.Context, req *models.GetCommentReq) (*models.Comment, error) {
	logger.Info("Start GetComment")

	resp, err := pc.getComment(ctx, req)
	if err != nil {
		logger.ErrorData("Finish GetComment: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetComment: Successful")
	return resp.(*models.Comment), nil
}

func (pc *postClient) RemoveLoveReact(ctx context.Context, req *models.RemoveLoveReactReq) (*models.RemoveLoveReactResp, error) {
	logger.Info("Start RemoveLoveReact")

	resp, err := pc.removeLoveReact(ctx, req)
	if err != nil {
		logger.ErrorData("Finish RemoveLoveReact: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish RemoveLoveReact: Successful")
	return resp.(*models.RemoveLoveReactResp), nil
}

func (pc *postClient) GetLove(ctx context.Context, req *models.GetLoveReq) (*models.Love, error) {
	logger.Info("Start GetLove")

	resp, err := pc.getLove(ctx, req)
	if err != nil {
		logger.ErrorData("Finish GetLove: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetLove: Successful")
	return resp.(*models.Love), nil
}

func (pc *postClient) ListCommentIDsByParentID(ctx context.Context, req *models.ListCommentIDsByParentIDReq) (*models.ListCommentIDsByParentIDResp, error) {
	logger.Info("Start ListCommentsByParentID")

	resp, err := pc.listCommentIDsByParentID(ctx, req)
	if err != nil {
		logger.ErrorData("Finish ListCommentsByParentID: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListCommentsByParentID: Successful")
	return resp.(*models.ListCommentIDsByParentIDResp), nil
}
