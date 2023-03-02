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
}

func NewPostServer(endpoints endpointsV1.PostEndpoints) postProtoV1.PostServiceServer {
	return &gRPCPostServer{
		createPost: grpctransport.NewServer(
			endpoints.CreatePost,
			translators.DecodeCreatePostRequest,
			translators.EncodeCreatePostResponse,
		),
		createComment: grpctransport.NewServer(
			endpoints.CreateComment,
			translators.DecodeCreateCommentRequest,
			translators.EncodeCreateCommentResponse,
		),
	}
}

func (s *gRPCPostServer) CreatePost(ctx context.Context, req *commonProto.CreatePostRequest) (*commonProto.Post, error) {
	_, resp, err := s.createPost.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*commonProto.Post), nil
}

func (s *gRPCPostServer) CreateComment(ctx context.Context, req *commonProto.CreateCommentRequest) (*commonProto.Comment, error) {
	_, resp, err := s.createComment.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*commonProto.Comment), nil
}
