package translators

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	petifiesProtoV1 "github.com/vantoan19/Petifies/proto/petifies-service/v1"
	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
)

func EncodeCreatePetifiesRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.CreatePetifiesReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.CreatePetifiesRequest{
		OwnerId:     req.OwnerID.String(),
		Type:        getPetifiesType(req.Type),
		Title:       req.Title,
		Description: req.Description,
		PetName:     req.PetName,
		Images: commonutils.Map2(req.Images, func(i models.Image) *commonProto.Image {
			return &commonProto.Image{Uri: i.URI, Description: i.Description}
		}),
		Address: encodeAddressModel(&req.Address),
	}, nil
}

func EncodePetifiesResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.Petifies)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return encodePetifiesModel(resp), nil
}

func EncodeCreatePetifiesSessionRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.CreatePetifiesSessionReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.CreatePetifiesSessionRequest{
		PetifiesId: req.PetifiesID.String(),
		FromTime:   timestamppb.New(req.FromTime),
		ToTime:     timestamppb.New(req.ToTime),
	}, nil
}

func EncodePetifiesSessionResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.PetifiesSession)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return encodePetifiesSessionModel(resp), nil
}

func EncodeCreatePetifiesProposalRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.CreatePetifiesProposalReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.CreatePetifiesProposalRequest{
		UserId:            req.UserID.String(),
		PetifiesSessionId: req.PetifiesSessionID.String(),
		Proposal:          req.Proposal,
	}, nil
}

func EncodePetifiesProposalResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.PetifiesProposal)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return encodePetifiesProposalModel(resp), nil
}

func EncodeCreateReviewRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.CreateReviewReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.CreateReviewRequest{
		PetifiesId: req.PetifiesID.String(),
		AuthorId:   req.AuthorID.String(),
		Review:     req.Review,
		Image:      &commonProto.Image{Uri: req.Image.URI, Description: req.Image.Description},
	}, nil
}

func EncodeReviewResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.Review)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return encodeReviewModel(resp), nil
}

func EncodeEditPetifiesRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.EditPetifiesReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.EditPetifiesRequest{
		Id:          req.ID.String(),
		Title:       req.Title,
		Description: req.Description,
		PetName:     req.PetName,
		Images: commonutils.Map2(req.Images, func(i models.Image) *commonProto.Image {
			return &commonProto.Image{Uri: i.URI, Description: i.Description}
		}),
		Address: encodeAddressModel(&req.Address),
	}, nil
}

func EncodeEditPetifiesSessionRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.EditPetifiesSessionReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.EditPetifiesSessionRequest{
		Id:       req.ID.String(),
		FromTime: timestamppb.New(req.FromTime),
		ToTime:   timestamppb.New(req.ToTime),
	}, nil
}

func EncodeEditPetifiesProposalRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.EditPetifiesProposalReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.EditPetifiesProposalRequest{
		Id:       req.ID.String(),
		Proposal: req.Proposal,
	}, nil
}

func EncodeEditReviewRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.EditReviewReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.EditReviewRequest{
		Id:     req.ID.String(),
		Review: req.Review,
		Image:  &commonProto.Image{Uri: req.Image.URI, Description: req.Image.Description},
	}, nil
}

func EncodeGetPetifiesByIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.GetPetifiesByIdReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.GetPetifiesByIdRequest{
		Id: req.ID.String(),
	}, nil
}

func EncodeGetSessionByIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.GetSessionByIdReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.GetSessionByIdRequest{
		Id: req.ID.String(),
	}, nil
}

func EncodeGetProposalByIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.GetProposalByIdReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.GetProposalByIdRequest{
		Id: req.ID.String(),
	}, nil
}

func EncodeGetReviewByIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.GetReviewByIdReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.GetReviewByIdRequest{
		Id: req.ID.String(),
	}, nil
}

func EncodeListPetifiesByIdsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.ListPetifiesByIdsReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.ListPetifiesByIdsRequest{
		PetifiesIds: commonutils.Map2(req.PetifiesIDs, func(id uuid.UUID) string { return id.String() }),
	}, nil
}

func EncodeListSessionsByIdsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.ListSessionsByIdsReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.ListSessionsByIdsRequest{
		PetifiesSessionIds: commonutils.Map2(req.PetifiesSessionIDs, func(id uuid.UUID) string { return id.String() }),
	}, nil
}

func EncodeListProposalsByIdsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.ListProposalsByIdsReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.ListProposalsByIdsRequest{
		PetifiesProposalIds: commonutils.Map2(req.PetifiesProposalIDs, func(id uuid.UUID) string { return id.String() }),
	}, nil
}

func EncodeListReviewsByIdsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.ListReviewsByIdsReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.ListReviewsByIdsRequest{
		ReviewIds: commonutils.Map2(req.ReviewIDs, func(id uuid.UUID) string { return id.String() }),
	}, nil
}

func EncodeListPetifiesByOwnerIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.ListPetifiesByOwnerIdReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.ListPetifiesByOwnerIdRequest{
		OwnerId:  req.OwnerID.String(),
		PageSize: int32(req.PageSize),
		AfterId:  req.AfterID.String(),
	}, nil
}

func EncodeListSessionsByPetifiesIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.ListSessionsByPetifiesIdReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.ListSessionsByPetifiesIdRequest{
		PetifiesId: req.PetifiesID.String(),
		PageSize:   int32(req.PageSize),
		AfterId:    req.AfterID.String(),
	}, nil
}

func EncodeListProposalsBySessionIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.ListProposalsBySessionIdReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.ListProposalsBySessionIdRequest{
		PetifiesSessionId: req.PetifiesSessionID.String(),
		PageSize:          int32(req.PageSize),
		AfterId:           req.AfterID.String(),
	}, nil
}

func EncodeListReviewsByPetifiesIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*models.ListReviewsByPetifiesIdReq)
	if !ok {
		return nil, MustBeEndpointReqErr
	}

	return &petifiesProtoV1.ListReviewsByPetifiesIdRequest{
		PetifiesId: req.PetifiesID.String(),
		PageSize:   int32(req.PageSize),
		AfterId:    req.AfterID.String(),
	}, nil
}

func EncodeManyPetifiesResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ManyPetifies)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &petifiesProtoV1.ManyPetifies{
		Petifies: commonutils.Map2(resp.Petifies,
			func(p *models.Petifies) *petifiesProtoV1.Petifies { return encodePetifiesModel(p) }),
	}, nil
}

func EncodeManyPetifiesSessionsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ManyPetifiesSessions)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &petifiesProtoV1.ManyPetifiesSessions{
		PetifiesSessions: commonutils.Map2(resp.PetifiesSessions,
			func(p *models.PetifiesSession) *petifiesProtoV1.PetifiesSession { return encodePetifiesSessionModel(p) }),
	}, nil
}

func EncodeManyPetifiesProposalsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ManyPetifiesProposals)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &petifiesProtoV1.ManyPetifiesProposals{
		PetifiesProposals: commonutils.Map2(resp.PetifiesProposals,
			func(p *models.PetifiesProposal) *petifiesProtoV1.PetifiesProposal {
				return encodePetifiesProposalModel(p)
			}),
	}, nil
}

func EncodeManyReviewsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*models.ManyReviews)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return &petifiesProtoV1.ManyReviews{
		Reviews: commonutils.Map2(resp.Reviews, func(r *models.Review) *petifiesProtoV1.Review { return encodeReviewModel(r) }),
	}, nil
}

func encodePetifiesModel(petifies *models.Petifies) *petifiesProtoV1.Petifies {
	return &petifiesProtoV1.Petifies{
		Id:          petifies.ID.String(),
		OwnerId:     petifies.OwnerID.String(),
		Type:        getPetifiesType(petifies.Type),
		Title:       petifies.Title,
		Description: petifies.Description,
		PetName:     petifies.PetName,
		Images: commonutils.Map2(petifies.Images, func(i models.Image) *commonProto.Image {
			return &commonProto.Image{Uri: i.URI, Description: i.Description}
		}),
		Address:   encodeAddressModel(&petifies.Address),
		Status:    getPetifiesStatus(petifies.Status),
		CreatedAt: timestamppb.New(petifies.CreatedAt),
		UpdatedAt: timestamppb.New(petifies.UpdatedAt),
	}
}

func encodePetifiesSessionModel(session *models.PetifiesSession) *petifiesProtoV1.PetifiesSession {
	return &petifiesProtoV1.PetifiesSession{
		Id:         session.ID.String(),
		PetifiesId: session.PetifiesID.String(),
		FromTime:   timestamppb.New(session.FromTime),
		ToTime:     timestamppb.New(session.ToTime),
		Status:     getPetifiesSessionStatus(session.Status),
		CreatedAt:  timestamppb.New(session.CreatedAt),
		UpdatedAt:  timestamppb.New(session.UpdatedAt),
	}
}

func encodePetifiesProposalModel(proposal *models.PetifiesProposal) *petifiesProtoV1.PetifiesProposal {
	return &petifiesProtoV1.PetifiesProposal{
		Id:                proposal.ID.String(),
		UserId:            proposal.UserID.String(),
		PetifiesSessionId: proposal.PetifiesSessionID.String(),
		Proposal:          proposal.Proposal,
		Status:            getPetifiesProposalStatus(proposal.Status),
		CreatedAt:         timestamppb.New(proposal.CreatedAt),
		UpdatedAt:         timestamppb.New(proposal.UpdatedAt),
	}
}

func encodeReviewModel(review *models.Review) *petifiesProtoV1.Review {
	return &petifiesProtoV1.Review{
		Id:         review.ID.String(),
		PetifiesId: review.PetifiesID.String(),
		AuthorId:   review.AuthorID.String(),
		Review:     review.Review,
		Image:      &commonProto.Image{Uri: review.Image.URI, Description: review.Image.Description},
		CreatedAt:  timestamppb.New(review.CreatedAt),
		UpdatedAt:  timestamppb.New(review.UpdatedAt),
	}
}

func encodeAddressModel(address *models.Address) *commonProto.Address {
	return &commonProto.Address{
		AddressLineOne: address.AddressLineOne,
		AddressLineTwo: address.AddressLineTwo,
		Street:         address.Street,
		District:       address.District,
		City:           address.City,
		Region:         address.Region,
		PostalCode:     address.PostalCode,
		Country:        address.Country,
		Longitude:      address.Longitude,
		Latitude:       address.Latitude,
	}
}

func getPetifiesType(petifiesType string) petifiesProtoV1.PetifiesType {
	switch petifiesType {
	case "DOG_WALKING":
		return petifiesProtoV1.PetifiesType_PETIFIES_TYPE_DOG_WALKING
	case "CAT_PLAYING":
		return petifiesProtoV1.PetifiesType_PETIFIES_TYPE_CAT_PLAYING
	case "DOG_SITTING":
		return petifiesProtoV1.PetifiesType_PETIFIES_TYPE_DOG_SITTING
	case "CAT_SITTING":
		return petifiesProtoV1.PetifiesType_PETIFIES_TYPE_CAT_SITTING
	case "DOG_ADOPTION":
		return petifiesProtoV1.PetifiesType_PETIFIES_TYPE_DOG_ADOPTION
	case "CAT_ADOPTION":
		return petifiesProtoV1.PetifiesType_PETIFIES_TYPE_CAT_ADOPTION
	default:
		return petifiesProtoV1.PetifiesType_PETIFIES_TYPE_UNKNOWN
	}
}

func getPetifiesStatus(petifiesStatus string) petifiesProtoV1.PetifiesStatus {
	switch petifiesStatus {
	case "UNAVAILABLE":
		return petifiesProtoV1.PetifiesStatus_PETIFIES_STATUS_UNAVAILABLE
	case "ON_A_SESSION":
		return petifiesProtoV1.PetifiesStatus_PETIFIES_STATUS_ON_A_SESSION
	case "DELETED":
		return petifiesProtoV1.PetifiesStatus_PETIFIES_STATUS_DELETED
	default:
		return petifiesProtoV1.PetifiesStatus_PETIFIES_STATUS_UNKNOWN
	}
}

func getPetifiesSessionStatus(petifiesSessionStatus string) petifiesProtoV1.PetifiesSessionStatus {
	switch petifiesSessionStatus {
	case "WAITING_FOR_PROPOSAL":
		return petifiesProtoV1.PetifiesSessionStatus_PETIFIES_SESSION_STATUS_WAITING_FOR_PROPOSAL
	case "PROPOSAL_ACCEPTED":
		return petifiesProtoV1.PetifiesSessionStatus_PETIFIES_SESSION_STATUS_PROPOSAL_ACCEPTED
	case "ON_GOING":
		return petifiesProtoV1.PetifiesSessionStatus_PETIFIES_SESSION_STATUS_ON_GOING
	case "ENDED":
		return petifiesProtoV1.PetifiesSessionStatus_PETIFIES_SESSION_STATUS_ENDED
	default:
		return petifiesProtoV1.PetifiesSessionStatus_PETIFIES_SESSION_STATUS_UNKNOWN
	}
}

func getPetifiesProposalStatus(petifiesProposal string) petifiesProtoV1.PetifiesProposalStatus {
	switch petifiesProposal {
	case "WAITING_FOR_ACCEPTANCE":
		return petifiesProtoV1.PetifiesProposalStatus_PETIFIES_PROPOSAL_STATUS_WAITING_FOR_ACCEPTANCE
	case "ACCEPTED":
		return petifiesProtoV1.PetifiesProposalStatus_PETIFIES_PROPOSAL_STATUS_ACCEPTED
	case "CANCELLED":
		return petifiesProtoV1.PetifiesProposalStatus_PETIFIES_PROPOSAL_STATUS_CANCELLED
	case "REJECTED":
		return petifiesProtoV1.PetifiesProposalStatus_PETIFIES_PROPOSAL_STATUS_REJECTED
	case "SESSION_CLOSED":
		return petifiesProtoV1.PetifiesProposalStatus_PETIFIES_PROPOSAL_STATUS_SESSION_CLOSED
	default:
		return petifiesProtoV1.PetifiesProposalStatus_PETIFIES_PROPOSAL_STATUS_UNKNOWN
	}
}
