package v1

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	petifiesservice "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/application/services/petifies"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/pkg/models"
)

type PetifiesEndpoint struct {
	CreatePetifies           endpoint.Endpoint
	CreatePetifiesSession    endpoint.Endpoint
	CreatePetifiesProposal   endpoint.Endpoint
	CreateReview             endpoint.Endpoint
	ListNearByPetifies       endpoint.Endpoint
	ListPetifiesByUserId     endpoint.Endpoint
	ListSessionsByPetifiesId endpoint.Endpoint
	ListProposalsBySessionId endpoint.Endpoint
	ListProposalsByUserId    endpoint.Endpoint
	ListReviewsByPetifiesId  endpoint.Endpoint
	ListReviewsByUserId      endpoint.Endpoint
}

func NewPetifiesEndpoints(s petifiesservice.PetifiesService) PetifiesEndpoint {
	return PetifiesEndpoint{
		CreatePetifies:           makeCreatePetifiesEndpoint(s),
		CreatePetifiesSession:    makeCreatePetifiesSessionEndpoint(s),
		CreatePetifiesProposal:   makeCreatePetifiesProposalEndpoint(s),
		CreateReview:             makeCreateReviewEndpoint(s),
		ListNearByPetifies:       makeListNearByPetifiesEndpoint(s),
		ListPetifiesByUserId:     makeListPetifiesByUserIdEndpoint(s),
		ListSessionsByPetifiesId: makeListSessionsByPetifiesIdEndpoint(s),
		ListProposalsBySessionId: makeListProposalsBySessionIdEndpoint(s),
		ListProposalsByUserId:    makeListProposalsByUserIdEndpoint(s),
		ListReviewsByPetifiesId:  makeListReviewsByPetifiesIdEndpoint(s),
		ListReviewsByUserId:      makeListReviewsByUserIdEndpoint(s),
	}
}

func makeCreatePetifiesEndpoint(s petifiesservice.PetifiesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.UserCreatePetifiesReq)
		resp, err := s.UserCreatePetifies(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

func makeCreatePetifiesSessionEndpoint(s petifiesservice.PetifiesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.UserCreatePetifiesSessionReq)
		resp, err := s.UserCreateSession(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

func makeCreatePetifiesProposalEndpoint(s petifiesservice.PetifiesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.UserCreatePetifiesProposalReq)
		resp, err := s.UserCreateProposal(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

func makeCreateReviewEndpoint(s petifiesservice.PetifiesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.UserCreateReviewReq)
		resp, err := s.UserCreateReview(ctx, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

func makeListNearByPetifiesEndpoint(s petifiesservice.PetifiesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListNearByPetifiesReq)
		resp, err := s.ListNearByPetifies(ctx, req)
		if err != nil {
			return nil, err
		}
		return &models.ListNearByPetifiesResp{
			Petifies: resp,
		}, nil
	}
}

func makeListPetifiesByUserIdEndpoint(s petifiesservice.PetifiesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListPetifiesByUserIdReq)
		resp, err := s.ListPetifiesByUserId(ctx, req.UserId, req.PageSize, req.AfterId)
		if err != nil {
			return nil, err
		}
		return &models.ListPetifiesByUserIdResp{
			Petifies: resp,
		}, nil
	}
}

// ListSessionsByPetifiesId endpoint.Endpoint
// ListProposalsBySessionId endpoint.Endpoint
// ListProposalsByUserId    endpoint.Endpoint
// ListReviewsByPetifiesId  endpoint.Endpoint
// ListReviewsByUserId      endpoint.Endpoint

func makeListSessionsByPetifiesIdEndpoint(s petifiesservice.PetifiesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListSessionsByPetifiesIdReq)
		resp, err := s.ListSessionsByPetifiesId(ctx, req.PetifiesId, req.PageSize, req.AfterId)
		if err != nil {
			return nil, err
		}
		return &models.ListSessionsByPetifiesIdResp{
			Sessions: resp,
		}, nil
	}
}

func makeListProposalsBySessionIdEndpoint(s petifiesservice.PetifiesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListProposalsBySessionIdReq)
		resp, err := s.ListProposalsBySessionId(ctx, req.SessionId, req.PageSize, req.AfterId)
		if err != nil {
			return nil, err
		}
		return &models.ListProposalsBySessionIdResp{
			Proposals: resp,
		}, nil
	}
}

func makeListReviewsByPetifiesIdEndpoint(s petifiesservice.PetifiesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListReviewsByPetifiesIdReq)
		resp, err := s.ListReviewsByPetifiesId(ctx, req.PetifiesId, req.PageSize, req.AfterId)
		if err != nil {
			return nil, err
		}
		return &models.ListReviewsByPetifiesIdResp{
			Reviews: resp,
		}, nil
	}
}

func makeListReviewsByUserIdEndpoint(s petifiesservice.PetifiesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListReviewsByUserIdReq)
		resp, err := s.ListReviewsByUserId(ctx, req.UserId, req.PageSize, req.AfterId)
		if err != nil {
			return nil, err
		}
		return &models.ListReviewsByUserIdResp{
			Reviews: resp,
		}, nil
	}
}

func makeListProposalsByUserIdEndpoint(s petifiesservice.PetifiesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListProposalsByUserIdReq)
		resp, err := s.ListProposalsByUserId(ctx, req.UserId, req.PageSize, req.AfterId)
		if err != nil {
			return nil, err
		}
		return &models.ListProposalsByUserIdResp{
			Proposals: resp,
		}, nil
	}
}
