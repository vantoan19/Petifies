package v1

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	postProtoV1 "github.com/vantoan19/Petifies/proto/post-service/v1"
	endpointsV1 "github.com/vantoan19/Petifies/server/services/post-service/internal/presentation/endpoints/grpc/v1"
	"github.com/vantoan19/Petifies/server/services/post-service/pkg/translators"
)

type gRPCPostServer struct {
	createPost    grpctransport.Handler
	createComment grpctransport.Handler
	loveReact     grpctransport.Handler
	editPost      grpctransport.Handler
	editComment   grpctransport.Handler
	listComments  grpctransport.Handler
	listPosts     grpctransport.Handler
}

func NewPostServer(endpoints endpointsV1.PostEndpoints) postProtoV1.PostServiceServer {
	return &gRPCPostServer{
		createPost: grpctransport.NewServer(
			endpoints.CreatePost,
			translators.DecodeCreatePostRequest,
			translators.EncodePostResponse,
		),
		createComment: grpctransport.NewServer(
			endpoints.CreateComment,
			translators.DecodeCreateCommentRequest,
			translators.EncodeCommentResponse,
		),
		loveReact: grpctransport.NewServer(
			endpoints.LoveReact,
			translators.DecodeLoveReactRequest,
			translators.EncodeLoveResponse,
		),
		editPost: grpctransport.NewServer(
			endpoints.EditPost,
			translators.DecodeEditPostRequest,
			translators.EncodePostResponse,
		),
		editComment: grpctransport.NewServer(
			endpoints.EditComment,
			translators.DecodeEditCommentRequest,
			translators.EncodeCommentResponse,
		),
		listComments: grpctransport.NewServer(
			endpoints.ListComments,
			translators.DecodeListCommentsRequest,
			translators.EncodeListCommentsResponse,
		),
		listPosts: grpctransport.NewServer(
			endpoints.ListPosts,
			translators.DecodeListPostsRequest,
			translators.EncodeListPostsResponse,
		),
	}
}

func (s *gRPCPostServer) CreatePost(ctx context.Context, req *postProtoV1.CreatePostRequest) (*commonProto.Post, error) {
	_, resp, err := s.createPost.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*commonProto.Post), nil
}

func (s *gRPCPostServer) CreateComment(ctx context.Context, req *postProtoV1.CreateCommentRequest) (*commonProto.Comment, error) {
	_, resp, err := s.createComment.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*commonProto.Comment), nil
}

func (s *gRPCPostServer) LoveReact(ctx context.Context, req *commonProto.LoveReactRequest) (*commonProto.Love, error) {
	_, resp, err := s.loveReact.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*commonProto.Love), nil
}

func (s *gRPCPostServer) EditPost(ctx context.Context, req *postProtoV1.EditPostRequest) (*commonProto.Post, error) {
	_, resp, err := s.editPost.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*commonProto.Post), nil
}

func (s *gRPCPostServer) EditComment(ctx context.Context, req *postProtoV1.EditCommentRequest) (*commonProto.Comment, error) {
	_, resp, err := s.editComment.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*commonProto.Comment), nil
}

func (s *gRPCPostServer) ListComments(ctx context.Context, req *postProtoV1.ListCommentsRequest) (*postProtoV1.ListCommentsResponse, error) {
	_, resp, err := s.listComments.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*postProtoV1.ListCommentsResponse), nil
}

func (s *gRPCPostServer) ListPosts(ctx context.Context, req *postProtoV1.ListPostsRequest) (*postProtoV1.ListPostsResponse, error) {
	_, resp, err := s.listPosts.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*postProtoV1.ListPostsResponse), nil
}