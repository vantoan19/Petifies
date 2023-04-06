package v1

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	petifiesProtoV1 "github.com/vantoan19/Petifies/proto/petifies-service/v1"
	endpointsV1 "github.com/vantoan19/Petifies/server/services/petifies-service/internal/presentation/endpoints/grpc/v1"
	"github.com/vantoan19/Petifies/server/services/petifies-service/pkg/translators"
)

type gRPCPetifiesServer struct {
	createPetifies         grpctransport.Handler
	createPetifiesSession  grpctransport.Handler
	createPetifiesProposal grpctransport.Handler
	createReview           grpctransport.Handler
	editPetifies           grpctransport.Handler
	editPetifiesSession    grpctransport.Handler
	editPetifiesProposal   grpctransport.Handler
	editReview             grpctransport.Handler

	getPetifiesById          grpctransport.Handler
	getSessionById           grpctransport.Handler
	getProposalById          grpctransport.Handler
	getReviewById            grpctransport.Handler
	listPetifiesByIds        grpctransport.Handler
	listSessionsByIds        grpctransport.Handler
	listProposalsByIds       grpctransport.Handler
	listReviewsByIds         grpctransport.Handler
	listPetifiesByOwnerId    grpctransport.Handler
	listSessionsByPetifiesId grpctransport.Handler
	listProposalsBySessionId grpctransport.Handler
	listReviewsByPetifiesId  grpctransport.Handler
}

func NewPetifiesServer(
	petifiesEndpoints endpointsV1.PetifiesEndpoints,
	petifiesSessionEndpoints endpointsV1.PetifiesSessionEndpoints,
	petifiesProposalEndpoints endpointsV1.PetifiesProposalEndpoints,
	reviewEndpoints endpointsV1.ReviewEndpoints,
) petifiesProtoV1.PetifiesServiceServer {
	return &gRPCPetifiesServer{
		createPetifies: grpctransport.NewServer(
			petifiesEndpoints.CreatePetify,
			translators.DecodeCreatePetifiesRequest,
			translators.EncodePetifiesResponse,
		),
		createPetifiesSession: grpctransport.NewServer(
			petifiesSessionEndpoints.CreatePetifiesSession,
			translators.DecodeCreatePetifiesSessionRequest,
			translators.EncodePetifiesSessionResponse,
		),
		createPetifiesProposal: grpctransport.NewServer(
			petifiesProposalEndpoints.CreatePetifiesProposal,
			translators.DecodeCreatePetifiesProposalRequest,
			translators.EncodePetifiesProposalResponse,
		),
		createReview: grpctransport.NewServer(
			reviewEndpoints.CreateReview,
			translators.DecodeCreateReviewRequest,
			translators.EncodeReviewResponse,
		),
		editPetifies: grpctransport.NewServer(
			petifiesEndpoints.EditPetify,
			translators.DecodeEditPetifiesRequest,
			translators.EncodePetifiesResponse,
		),
		editPetifiesSession: grpctransport.NewServer(
			petifiesSessionEndpoints.EditPetifiesSession,
			translators.DecodeEditPetifiesSessionRequest,
			translators.EncodePetifiesSessionResponse,
		),
		editPetifiesProposal: grpctransport.NewServer(
			petifiesProposalEndpoints.EditPetifiesProposal,
			translators.DecodeEditPetifiesProposalRequest,
			translators.EncodePetifiesProposalResponse,
		),
		editReview: grpctransport.NewServer(
			reviewEndpoints.EditReview,
			translators.DecodeEditReviewRequest,
			translators.EncodeReviewResponse,
		),
		getPetifiesById: grpctransport.NewServer(
			petifiesEndpoints.GetPetifyById,
			translators.DecodeGetPetifiesByIdRequest,
			translators.EncodePetifiesResponse,
		),
		getSessionById: grpctransport.NewServer(
			petifiesSessionEndpoints.GetSessionById,
			translators.DecodeGetSessionByIdRequest,
			translators.EncodePetifiesSessionResponse,
		),
		getProposalById: grpctransport.NewServer(
			petifiesProposalEndpoints.GetProposalById,
			translators.DecodeGetProposalByIdRequest,
			translators.EncodePetifiesProposalResponse,
		),
		getReviewById: grpctransport.NewServer(
			reviewEndpoints.GetReviewById,
			translators.DecodeGetReviewByIdRequest,
			translators.EncodeReviewResponse,
		),
		listPetifiesByIds: grpctransport.NewServer(
			petifiesEndpoints.ListPetifiesByIds,
			translators.DecodeListPetifiesByIdsRequest,
			translators.EncodeManyPetifiesResponse,
		),
		listSessionsByIds: grpctransport.NewServer(
			petifiesSessionEndpoints.ListSessionsByIds,
			translators.DecodeListSessionsByIdsRequest,
			translators.EncodeManyPetifiesSessionsResponse,
		),
		listProposalsByIds: grpctransport.NewServer(
			petifiesProposalEndpoints.ListProposalsByIds,
			translators.DecodeListProposalsByIdsRequest,
			translators.EncodeManyPetifiesProposalsResponse,
		),
		listReviewsByIds: grpctransport.NewServer(
			reviewEndpoints.ListReviewsByIds,
			translators.DecodeListReviewsByIdsRequest,
			translators.EncodeManyReviewsResponse,
		),
		listPetifiesByOwnerId: grpctransport.NewServer(
			petifiesEndpoints.ListPetifiesByOwnerId,
			translators.DecodeListPetifiesByOwnerIdRequest,
			translators.EncodeManyPetifiesResponse,
		),
		listSessionsByPetifiesId: grpctransport.NewServer(
			petifiesSessionEndpoints.ListSessionsByPetifiesId,
			translators.DecodeListSessionsByPetifiesIdRequest,
			translators.EncodeManyPetifiesSessionsResponse,
		),
		listProposalsBySessionId: grpctransport.NewServer(
			petifiesProposalEndpoints.ListProposalsBySessionId,
			translators.DecodeListProposalsBySessionIdRequest,
			translators.EncodeManyPetifiesProposalsResponse,
		),
		listReviewsByPetifiesId: grpctransport.NewServer(
			reviewEndpoints.ListReviewsByPetifiesId,
			translators.DecodeListReviewsByPetifiesIdRequest,
			translators.EncodeManyReviewsResponse,
		),
	}
}

func (s *gRPCPetifiesServer) CreatePetifies(ctx context.Context, req *petifiesProtoV1.CreatePetifiesRequest) (*petifiesProtoV1.Petifies, error) {
	_, resp, err := s.createPetifies.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.Petifies), nil
}

func (s *gRPCPetifiesServer) CreatePetifiesSession(ctx context.Context, req *petifiesProtoV1.CreatePetifiesSessionRequest) (*petifiesProtoV1.PetifiesSession, error) {
	_, resp, err := s.createPetifiesSession.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.PetifiesSession), nil
}

func (s *gRPCPetifiesServer) CreatePetifiesProposal(ctx context.Context, req *petifiesProtoV1.CreatePetifiesProposalRequest) (*petifiesProtoV1.PetifiesProposal, error) {
	_, resp, err := s.createPetifiesProposal.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.PetifiesProposal), nil
}

func (s *gRPCPetifiesServer) CreateReview(ctx context.Context, req *petifiesProtoV1.CreateReviewRequest) (*petifiesProtoV1.Review, error) {
	_, resp, err := s.createReview.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.Review), nil
}

func (s *gRPCPetifiesServer) EditPetifies(ctx context.Context, req *petifiesProtoV1.EditPetifiesRequest) (*petifiesProtoV1.Petifies, error) {
	_, resp, err := s.editPetifies.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.Petifies), nil
}

func (s *gRPCPetifiesServer) EditPetifiesSession(ctx context.Context, req *petifiesProtoV1.EditPetifiesSessionRequest) (*petifiesProtoV1.PetifiesSession, error) {
	_, resp, err := s.editPetifiesSession.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.PetifiesSession), nil
}

func (s *gRPCPetifiesServer) EditPetifiesProposal(ctx context.Context, req *petifiesProtoV1.EditPetifiesProposalRequest) (*petifiesProtoV1.PetifiesProposal, error) {
	_, resp, err := s.editPetifiesProposal.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.PetifiesProposal), nil
}

func (s *gRPCPetifiesServer) EditReview(ctx context.Context, req *petifiesProtoV1.EditReviewRequest) (*petifiesProtoV1.Review, error) {
	_, resp, err := s.editReview.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.Review), nil
}

func (s *gRPCPetifiesServer) GetPetifiesById(ctx context.Context, req *petifiesProtoV1.GetPetifiesByIdRequest) (*petifiesProtoV1.Petifies, error) {
	_, resp, err := s.getPetifiesById.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.Petifies), nil
}

func (s *gRPCPetifiesServer) GetSessionById(ctx context.Context, req *petifiesProtoV1.GetSessionByIdRequest) (*petifiesProtoV1.PetifiesSession, error) {
	_, resp, err := s.getSessionById.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.PetifiesSession), nil
}

func (s *gRPCPetifiesServer) GetProposalById(ctx context.Context, req *petifiesProtoV1.GetProposalByIdRequest) (*petifiesProtoV1.PetifiesProposal, error) {
	_, resp, err := s.getProposalById.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.PetifiesProposal), nil
}

func (s *gRPCPetifiesServer) GetReviewById(ctx context.Context, req *petifiesProtoV1.GetReviewByIdRequest) (*petifiesProtoV1.Review, error) {
	_, resp, err := s.getReviewById.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.Review), nil
}

func (s *gRPCPetifiesServer) ListPetifiesByIds(ctx context.Context, req *petifiesProtoV1.ListPetifiesByIdsRequest) (*petifiesProtoV1.ManyPetifies, error) {
	_, resp, err := s.listPetifiesByIds.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.ManyPetifies), nil
}

func (s *gRPCPetifiesServer) ListSessionsByIds(ctx context.Context, req *petifiesProtoV1.ListSessionsByIdsRequest) (*petifiesProtoV1.ManyPetifiesSessions, error) {
	_, resp, err := s.listPetifiesByIds.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.ManyPetifiesSessions), nil
}

func (s *gRPCPetifiesServer) ListProposalsByIds(ctx context.Context, req *petifiesProtoV1.ListProposalsByIdsRequest) (*petifiesProtoV1.ManyPetifiesProposals, error) {
	_, resp, err := s.listProposalsByIds.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.ManyPetifiesProposals), nil
}

func (s *gRPCPetifiesServer) ListReviewsByIds(ctx context.Context, req *petifiesProtoV1.ListReviewsByIdsRequest) (*petifiesProtoV1.ManyReviews, error) {
	_, resp, err := s.listReviewsByIds.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.ManyReviews), nil
}

func (s *gRPCPetifiesServer) ListPetifiesByOwnerId(
	ctx context.Context,
	req *petifiesProtoV1.ListPetifiesByOwnerIdRequest,
) (*petifiesProtoV1.ManyPetifies, error) {
	_, resp, err := s.listPetifiesByOwnerId.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.ManyPetifies), nil
}

func (s *gRPCPetifiesServer) ListSessionsByPetifiesId(
	ctx context.Context,
	req *petifiesProtoV1.ListSessionsByPetifiesIdRequest,
) (*petifiesProtoV1.ManyPetifiesSessions, error) {
	_, resp, err := s.listPetifiesByOwnerId.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.ManyPetifiesSessions), nil
}

func (s *gRPCPetifiesServer) ListProposalsBySessionId(
	ctx context.Context,
	req *petifiesProtoV1.ListProposalsBySessionIdRequest,
) (*petifiesProtoV1.ManyPetifiesProposals, error) {
	_, resp, err := s.listProposalsBySessionId.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.ManyPetifiesProposals), nil
}

func (s *gRPCPetifiesServer) ListReviewsByPetifiesId(
	ctx context.Context,
	req *petifiesProtoV1.ListReviewsByPetifiesIdRequest,
) (*petifiesProtoV1.ManyReviews, error) {
	_, resp, err := s.listReviewsByPetifiesId.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*petifiesProtoV1.ManyReviews), nil
}
