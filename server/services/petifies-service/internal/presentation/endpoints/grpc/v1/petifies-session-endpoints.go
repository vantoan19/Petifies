package v1

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	petifiessessionservice "github.com/vantoan19/Petifies/server/services/petifies-service/internal/application/services/petifies-session-service"
	petifiesssionaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_session"
	"github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
)

type PetifiesSessionEndpoints struct {
	CreatePetifiesSession    endpoint.Endpoint
	EditPetifiesSession      endpoint.Endpoint
	GetSessionById           endpoint.Endpoint
	ListSessionsByIds        endpoint.Endpoint
	ListSessionsByPetifiesId endpoint.Endpoint
	AcceptProposal           endpoint.Endpoint
	RejectProposal           endpoint.Endpoint
}

func NewPetifiesSessionEndpoints(ps petifiessessionservice.PetifesSessionService) PetifiesSessionEndpoints {
	return PetifiesSessionEndpoints{
		CreatePetifiesSession:    makeCreatePetifiesSessionEndpoint(ps),
		EditPetifiesSession:      makeEditPetifiesSessionEndpoint(ps),
		GetSessionById:           makeGetSessionByIdEndpoint(ps),
		ListSessionsByIds:        makeListSessionsByIdsEndpoint(ps),
		ListSessionsByPetifiesId: makeListSessionsByPetifiesIdEndpoint(ps),
		AcceptProposal:           makeAcceptProposalEndpoint(ps),
		RejectProposal:           makeRejectProposalEndpoint(ps),
	}
}

// ================ Endpoint Makers ==================

func makeCreatePetifiesSessionEndpoint(ps petifiessessionservice.PetifesSessionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.CreatePetifiesSessionReq)
		result, err := ps.CreatePetifiesSession(ctx, req)
		if err != nil {
			return nil, err
		}

		return mapPetifiesSessionAggregateToPetifiesSessionModel(result), nil
	}
}

func makeEditPetifiesSessionEndpoint(ps petifiessessionservice.PetifesSessionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.EditPetifiesSessionReq)
		result, err := ps.EditPetifiesSession(ctx, req)
		if err != nil {
			return nil, err
		}

		return mapPetifiesSessionAggregateToPetifiesSessionModel(result), nil
	}
}

func makeGetSessionByIdEndpoint(ps petifiessessionservice.PetifesSessionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.GetSessionByIdReq)
		result, err := ps.GetById(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		return mapPetifiesSessionAggregateToPetifiesSessionModel(result), nil
	}
}

func makeListSessionsByIdsEndpoint(ps petifiessessionservice.PetifesSessionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListSessionsByIdsReq)
		results, err := ps.ListByIds(ctx, req.PetifiesSessionIDs)
		if err != nil {
			return nil, err
		}

		return &models.ManyPetifiesSessions{
			PetifiesSessions: commonutils.Map2(results, func(p *petifiesssionaggre.PetifiesSessionAggre) *models.PetifiesSession {
				return mapPetifiesSessionAggregateToPetifiesSessionModel(p)
			}),
		}, nil
	}
}

func makeListSessionsByPetifiesIdEndpoint(ps petifiessessionservice.PetifesSessionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListSessionsByPetifiesIdReq)
		results, err := ps.ListByPetifiesId(ctx, req.PetifiesID, req.PageSize, req.AfterID)
		if err != nil {
			return nil, err
		}

		return &models.ManyPetifiesSessions{
			PetifiesSessions: commonutils.Map2(results, func(p *petifiesssionaggre.PetifiesSessionAggre) *models.PetifiesSession {
				return mapPetifiesSessionAggregateToPetifiesSessionModel(p)
			}),
		}, nil
	}
}

func makeAcceptProposalEndpoint(ps petifiessessionservice.PetifesSessionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.AcceptProposalReq)
		err = ps.AcceptProposal(ctx, req)
		if err != nil {
			return nil, err
		}

		return &models.AcceptProposalResp{}, nil
	}
}

func makeRejectProposalEndpoint(ps petifiessessionservice.PetifesSessionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.RejectProposalReq)
		err = ps.RejectProposal(ctx, req)
		if err != nil {
			return nil, err
		}

		return &models.RejectProposalResp{}, nil
	}
}

// ================ Mappers ======================

func mapPetifiesSessionAggregateToPetifiesSessionModel(petifyProposal *petifiesssionaggre.PetifiesSessionAggre) *models.PetifiesSession {
	return &models.PetifiesSession{
		ID:         petifyProposal.GetID(),
		PetifiesID: petifyProposal.GetPetifiesID(),
		FromTime:   petifyProposal.GetTime().GetFromTime(),
		ToTime:     petifyProposal.GetTime().GetToTime(),
		Status:     string(petifyProposal.GetStatus()),
		CreatedAt:  petifyProposal.GetCreatedAt(),
		UpdatedAt:  petifyProposal.GetUpdatedAt(),
	}
}
