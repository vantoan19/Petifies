package v1

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	userProtoV1 "github.com/vantoan19/Petifies/proto/user-service/v1"
	userEndpointsV1 "github.com/vantoan19/Petifies/server/services/user-services/user-service/internal/endpoints/grpc/v1"
)

type gRPCUserServer struct {
	createUser grpctransport.Handler
}

func New(endpoints userEndpointsV1.UserEndpoints) userProtoV1.UserServiceServer {
	return &gRPCUserServer{
		createUser: grpctransport.NewServer(
			endpoints.CreateUser,
			decodeCreateUserRequest,
			encodeCreateUserResponse,
		),
	}
}

func (s *gRPCUserServer) CreateUser(ctx context.Context, req *userProtoV1.CreateUserRequest) (*userProtoV1.User, error) {
	_, resp, err := s.createUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*userProtoV1.User), nil
}

func decodeCreateUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*userProtoV1.CreateUserRequest)
	return userEndpointsV1.CreateUserReq{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}, nil
}

func encodeCreateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(userEndpointsV1.CreateUserResp)
	return &userProtoV1.User{
		Id:        resp.ID.String(),
		Email:     resp.Email,
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		CreatedAt: timestamppb.New(resp.CreatedAt),
		UpdatedAt: timestamppb.New(resp.UpdatedAt),
	}, nil
}
