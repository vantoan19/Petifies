package v1

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	relationshipProtoV1 "github.com/vantoan19/Petifies/proto/relationship-service/v1"
	endpointsV1 "github.com/vantoan19/Petifies/server/services/relationship-service/internal/presentation/endpoints/grpc/v1"
	"github.com/vantoan19/Petifies/server/services/relationship-service/pkg/translator"
)

type gRPCRelationshipServer struct {
	addRelationship    grpctransport.Handler
	removeRelationship grpctransport.Handler
	listFollowers      grpctransport.Handler
	listFollowings     grpctransport.Handler
}

func NewRelationshipServer(endpoints endpointsV1.RelationshipEndpoints) relationshipProtoV1.RelationshipServiceServer {
	return &gRPCRelationshipServer{
		addRelationship: grpctransport.NewServer(
			endpoints.AddRelationship,
			translator.DecodeAddRelationshipRequest,
			translator.EncodeAddRelationshipResponse,
		),
		removeRelationship: grpctransport.NewServer(
			endpoints.RemoveRelationship,
			translator.DecodeRemoveRelationshipRequest,
			translator.EncodeRemoveRelationshipResponse,
		),
		listFollowers: grpctransport.NewServer(
			endpoints.ListFollowers,
			translator.DecodeListFollowersRequest,
			translator.EncodeListFollowersResponse,
		),
		listFollowings: grpctransport.NewServer(
			endpoints.ListFollowings,
			translator.DecodeListFollowingsRequest,
			translator.EncodeListFollowingsResponse,
		),
	}
}

func (s *gRPCRelationshipServer) AddRelationship(ctx context.Context, req *relationshipProtoV1.AddRelationshipRequest) (*relationshipProtoV1.AddRelationshipResponse, error) {
	_, resp, err := s.addRelationship.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*relationshipProtoV1.AddRelationshipResponse), nil
}

func (s *gRPCRelationshipServer) RemoveRelationship(ctx context.Context, req *relationshipProtoV1.RemoveRelationshipRequest) (*relationshipProtoV1.RemoveRelationshipResponse, error) {
	_, resp, err := s.removeRelationship.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*relationshipProtoV1.RemoveRelationshipResponse), nil
}

func (s *gRPCRelationshipServer) ListFollowers(ctx context.Context, req *relationshipProtoV1.ListFollowersRequest) (*relationshipProtoV1.ListFollowersResponse, error) {
	_, resp, err := s.listFollowers.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*relationshipProtoV1.ListFollowersResponse), nil
}

func (s *gRPCRelationshipServer) ListFollowings(ctx context.Context, req *relationshipProtoV1.ListFollowingsRequest) (*relationshipProtoV1.ListFollowingsResponse, error) {
	_, resp, err := s.listFollowings.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*relationshipProtoV1.ListFollowingsResponse), nil
}
