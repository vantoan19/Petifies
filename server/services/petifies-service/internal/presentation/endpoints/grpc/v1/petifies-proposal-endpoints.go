package v1

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	petifiesproposalservice "github.com/vantoan19/Petifies/server/services/petifies-service/internal/application/services/petifies-proposal-service"
	petifiesproposalaggre "github.com/vantoan19/Petifies/server/services/petifies-service/internal/domain/aggregates/petifies_proposal"
	"github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
)

type PetifiesProposalEndpoints struct {
	CreatePetifiesProposal   endpoint.Endpoint
	EditPetifiesProposal     endpoint.Endpoint
	GetProposalById          endpoint.Endpoint
	ListProposalsByIds       endpoint.Endpoint
	ListProposalsBySessionId endpoint.Endpoint
	CancelProposal           endpoint.Endpoint
	ListProposalsByUserId    endpoint.Endpoint
}

func NewPetifiesProposalEndpoints(ps petifiesproposalservice.PetifesProposalService) PetifiesProposalEndpoints {
	return PetifiesProposalEndpoints{
		CreatePetifiesProposal:   makeCreatePetifiesProposalEndpoint(ps),
		EditPetifiesProposal:     makeEditPetifiesProposalEndpoint(ps),
		GetProposalById:          makeGetProposalByIdEndpoint(ps),
		ListProposalsByIds:       makeListProposalsByIdsEndpoint(ps),
		ListProposalsBySessionId: makeListProposalsBySessionIdEndpoint(ps),
		CancelProposal:           makeCancelProposalEndpoint(ps),
		ListProposalsByUserId:    makeListProposalsByUserIdEndpoint(ps),
	}
}

// ================ Endpoint Makers ==================

func makeCreatePetifiesProposalEndpoint(ps petifiesproposalservice.PetifesProposalService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.CreatePetifiesProposalReq)
		result, err := ps.CreatePetifiesProposal(ctx, req)
		if err != nil {
			return nil, err
		}

		return mapPetifiesProposalAggregateToPetifiesProposalModel(result), nil
	}
}

func makeEditPetifiesProposalEndpoint(ps petifiesproposalservice.PetifesProposalService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.EditPetifiesProposalReq)
		result, err := ps.EditPetifiesProposal(ctx, req)
		if err != nil {
			return nil, err
		}

		return mapPetifiesProposalAggregateToPetifiesProposalModel(result), nil
	}
}

func makeGetProposalByIdEndpoint(ps petifiesproposalservice.PetifesProposalService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.GetProposalByIdReq)
		result, err := ps.GetPetifiesProposalById(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		return mapPetifiesProposalAggregateToPetifiesProposalModel(result), nil
	}
}

func makeListProposalsByIdsEndpoint(ps petifiesproposalservice.PetifesProposalService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListProposalsByIdsReq)
		results, err := ps.ListPetifiesProposalsByIds(ctx, req.PetifiesProposalIDs)
		if err != nil {
			return nil, err
		}

		return &models.ManyPetifiesProposals{
			PetifiesProposals: commonutils.Map2(results, func(p *petifiesproposalaggre.PetifiesProposalAggre) *models.PetifiesProposal {
				return mapPetifiesProposalAggregateToPetifiesProposalModel(p)
			}),
		}, nil
	}
}

func makeListProposalsBySessionIdEndpoint(ps petifiesproposalservice.PetifesProposalService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListProposalsBySessionIdReq)
		results, err := ps.ListPetifiesProposalsBySessionId(ctx, req.PetifiesSessionID, req.PageSize, req.AfterID)
		if err != nil {
			return nil, err
		}

		return &models.ManyPetifiesProposals{
			PetifiesProposals: commonutils.Map2(results, func(p *petifiesproposalaggre.PetifiesProposalAggre) *models.PetifiesProposal {
				return mapPetifiesProposalAggregateToPetifiesProposalModel(p)
			}),
		}, nil
	}
}

func makeListProposalsByUserIdEndpoint(ps petifiesproposalservice.PetifesProposalService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.ListProposalsByUserIdReq)
		results, err := ps.ListPetifiesProposalsByUserId(ctx, req.UserId, req.PageSize, req.AfterID)
		if err != nil {
			return nil, err
		}

		return &models.ManyPetifiesProposals{
			PetifiesProposals: commonutils.Map2(results, func(p *petifiesproposalaggre.PetifiesProposalAggre) *models.PetifiesProposal {
				return mapPetifiesProposalAggregateToPetifiesProposalModel(p)
			}),
		}, nil
	}
}

func makeCancelProposalEndpoint(ps petifiesproposalservice.PetifesProposalService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*models.CancelProposalReq)
		err = ps.CancelPetifiesProposal(ctx, req)
		if err != nil {
			return nil, err
		}

		return &models.CancelProposalResp{}, nil
	}
}

// ================ Mappers ======================

func mapPetifiesProposalAggregateToPetifiesProposalModel(petifyProposal *petifiesproposalaggre.PetifiesProposalAggre) *models.PetifiesProposal {
	return &models.PetifiesProposal{
		ID:                petifyProposal.GetID(),
		UserID:            petifyProposal.GetUserID(),
		PetifiesSessionID: petifyProposal.GetPetifiesSessionID(),
		Proposal:          petifyProposal.GetProposal(),
		Status:            string(petifyProposal.GetStatus()),
		CreatedAt:         petifyProposal.GetCreatedAt(),
		UpdatedAt:         petifyProposal.GetUpdatedAt(),
	}
}
