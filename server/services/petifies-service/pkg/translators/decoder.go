package translators

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	petifiesProtoV1 "github.com/vantoan19/Petifies/proto/petifies-service/v1"
	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
)

var (
	MustBeProtoReqErr     = status.Error(codes.InvalidArgument, "must be proto request")
	MustBeProtoRespErr    = status.Error(codes.InvalidArgument, "must be proto response")
	MustBeEndpointReqErr  = status.Error(codes.InvalidArgument, "must be endpoints' request")
	MustBeEndpointRespErr = status.Error(codes.InvalidArgument, "must be endpoints' response")
)

func DecodeCreatePetifiesRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.CreatePetifiesRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	ownerId, err := uuid.Parse(req.GetOwnerId())
	if err != nil {
		return nil, err
	}

	return &models.CreatePetifiesReq{
		OwnerID:     ownerId,
		Type:        req.Type.String(),
		Title:       req.Title,
		Description: req.Description,
		PetName:     req.PetName,
		Images:      commonutils.Map2(req.Images, func(i *commonProto.Image) models.Image { return models.Image{URI: i.Uri, Description: i.Description} }),
		Address:     *decodeAddressProtoModel(req.Address),
	}, nil
}

func DecodePetifiesResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*petifiesProtoV1.Petifies)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return decodePetifiesProtoModel(resp)
}

func DecodeCreatePetifiesSessionRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.CreatePetifiesSessionRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	creatorId, err := uuid.Parse(req.GetCreatorId())
	if err != nil {
		return nil, err
	}
	petifiesId, err := uuid.Parse(req.GetPetifiesId())
	if err != nil {
		return nil, err
	}

	return &models.CreatePetifiesSessionReq{
		CreatorID:  creatorId,
		PetifiesID: petifiesId,
		FromTime:   req.FromTime.AsTime(),
		ToTime:     req.ToTime.AsTime(),
	}, nil
}

func DecodePetifiesSessionResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*petifiesProtoV1.PetifiesSession)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return decodePetifiesSessionProtoModel(resp)
}

func DecodeCreatePetifiesProposalRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.CreatePetifiesProposalRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, err
	}
	petifiesSessionID, err := uuid.Parse(req.GetPetifiesSessionId())
	if err != nil {
		return nil, err
	}

	return &models.CreatePetifiesProposalReq{
		UserID:            userID,
		PetifiesSessionID: petifiesSessionID,
		Proposal:          req.Proposal,
	}, nil
}

func DecodePetifiesProposalResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*petifiesProtoV1.PetifiesProposal)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return decodePetifiesProposalProtoModel(resp)
}

func DecodeCreateReviewRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.CreateReviewRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	authorID, err := uuid.Parse(req.GetAuthorId())
	if err != nil {
		return nil, err
	}
	petifiesID, err := uuid.Parse(req.GetPetifiesId())
	if err != nil {
		return nil, err
	}

	return &models.CreateReviewReq{
		AuthorID:   authorID,
		PetifiesID: petifiesID,
		Review:     req.Review,
		Image: models.Image{
			URI:         req.Image.Uri,
			Description: req.Image.Description,
		},
	}, nil
}

func DecodeReviewResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*petifiesProtoV1.Review)
	if !ok {
		return nil, MustBeEndpointRespErr
	}

	return decodeReviewProtoModel(resp)
}

func DecodeEditPetifiesRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.EditPetifiesRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	return &models.EditPetifiesReq{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		PetName:     req.PetName,
		Images:      commonutils.Map2(req.Images, func(i *commonProto.Image) models.Image { return models.Image{URI: i.Uri, Description: i.Description} }),
		Address:     *decodeAddressProtoModel(req.Address),
	}, nil
}

func DecodeEditPetifiesSessionRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.EditPetifiesSessionRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	return &models.EditPetifiesSessionReq{
		ID:       id,
		FromTime: req.FromTime.AsTime(),
		ToTime:   req.ToTime.AsTime(),
	}, nil
}

func DecodeEditPetifiesProposalRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.EditPetifiesProposalRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	return &models.EditPetifiesProposalReq{
		ID:       id,
		Proposal: req.Proposal,
	}, nil
}

func DecodeEditReviewRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.EditReviewRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	return &models.EditReviewReq{
		ID:     id,
		Review: req.Review,
		Image:  models.Image{URI: req.Image.Uri, Description: req.Image.Description},
	}, nil
}

func DecodeGetPetifiesByIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.GetPetifiesByIdRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	return &models.GetPetifiesByIdReq{
		ID: id,
	}, nil
}

func DecodeGetSessionByIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.GetSessionByIdRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	return &models.GetSessionByIdReq{
		ID: id,
	}, nil
}

func DecodeGetProposalByIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.GetProposalByIdRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	return &models.GetProposalByIdReq{
		ID: id,
	}, nil
}

func DecodeGetReviewByIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.GetReviewByIdRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	return &models.GetReviewByIdReq{
		ID: id,
	}, nil
}

func DecodeListPetifiesByIdsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.ListPetifiesByIdsRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	var ids []uuid.UUID
	for _, i := range req.PetifiesIds {
		id, err := uuid.Parse(i)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return &models.ListPetifiesByIdsReq{
		PetifiesIDs: ids,
	}, nil
}

func DecodeListSessionsByIdsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.ListSessionsByIdsRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	var ids []uuid.UUID
	for _, i := range req.PetifiesSessionIds {
		id, err := uuid.Parse(i)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return &models.ListSessionsByIdsReq{
		PetifiesSessionIDs: ids,
	}, nil
}

func DecodeListProposalsByIdsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.ListProposalsByIdsRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	var ids []uuid.UUID
	for _, i := range req.PetifiesProposalIds {
		id, err := uuid.Parse(i)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return &models.ListProposalsByIdsReq{
		PetifiesProposalIDs: ids,
	}, nil
}

func DecodeListReviewsByIdsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.ListReviewsByIdsRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	var ids []uuid.UUID
	for _, i := range req.ReviewIds {
		id, err := uuid.Parse(i)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return &models.ListReviewsByIdsReq{
		ReviewIDs: ids,
	}, nil
}

func DecodeListPetifiesByOwnerIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.ListPetifiesByOwnerIdRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	ownerID, err := uuid.Parse(req.OwnerId)
	if err != nil {
		return nil, err
	}
	afterID, err := uuid.Parse(req.AfterId)
	if err != nil {
		afterID = uuid.Nil
	}

	return &models.ListPetifiesByOwnerIdReq{
		OwnerID:  ownerID,
		PageSize: int(req.PageSize),
		AfterID:  afterID,
	}, nil
}

func DecodeListSessionsByPetifiesIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.ListSessionsByPetifiesIdRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	petifiesID, err := uuid.Parse(req.PetifiesId)
	if err != nil {
		return nil, err
	}
	afterID, err := uuid.Parse(req.AfterId)
	if err != nil {
		afterID = uuid.Nil
	}

	return &models.ListSessionsByPetifiesIdReq{
		PetifiesID: petifiesID,
		PageSize:   int(req.PageSize),
		AfterID:    afterID,
	}, nil
}

func DecodeListProposalsBySessionIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.ListProposalsBySessionIdRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	sessionID, err := uuid.Parse(req.PetifiesSessionId)
	if err != nil {
		return nil, err
	}
	afterID, err := uuid.Parse(req.AfterId)
	if err != nil {
		afterID = uuid.Nil
	}

	return &models.ListProposalsBySessionIdReq{
		PetifiesSessionID: sessionID,
		PageSize:          int(req.PageSize),
		AfterID:           afterID,
	}, nil
}

func DecodeListReviewsByPetifiesIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.ListReviewsByPetifiesIdRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	petifiesID, err := uuid.Parse(req.PetifiesId)
	if err != nil {
		return nil, err
	}
	afterID, err := uuid.Parse(req.AfterId)
	if err != nil {
		afterID = uuid.Nil
	}

	return &models.ListReviewsByPetifiesIdReq{
		PetifiesID: petifiesID,
		PageSize:   int(req.PageSize),
		AfterID:    afterID,
	}, nil
}

func DecodeListProposalsByUserIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.ListProposalsByUserIdRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	afterID, err := uuid.Parse(req.AfterId)
	if err != nil {
		afterID = uuid.Nil
	}

	return &models.ListProposalsByUserIdReq{
		UserId:   userId,
		PageSize: int(req.PageSize),
		AfterID:  afterID,
	}, nil
}

func DecodeListReviewsByUserIdRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.ListReviewsByUserIdRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	afterID, err := uuid.Parse(req.AfterId)
	if err != nil {
		afterID = uuid.Nil
	}

	return &models.ListReviewsByUserIdReq{
		UserId:   userId,
		PageSize: int(req.PageSize),
		AfterID:  afterID,
	}, nil
}

func DecodeManyPetifesResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*petifiesProtoV1.ManyPetifies)
	if !ok {
		return nil, MustBeProtoRespErr
	}

	var petifies []*models.Petifies
	for _, p := range resp.Petifies {
		p_, err := decodePetifiesProtoModel(p)
		if err != nil {
			return nil, err
		}
		petifies = append(petifies, p_)
	}

	return &models.ManyPetifies{
		Petifies: petifies,
	}, nil
}

func DecodeManyPetifesSessionsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*petifiesProtoV1.ManyPetifiesSessions)
	if !ok {
		return nil, MustBeProtoRespErr
	}

	var petifiesSessions []*models.PetifiesSession
	for _, p := range resp.PetifiesSessions {
		p_, err := decodePetifiesSessionProtoModel(p)
		if err != nil {
			return nil, err
		}
		petifiesSessions = append(petifiesSessions, p_)
	}

	return &models.ManyPetifiesSessions{
		PetifiesSessions: petifiesSessions,
	}, nil
}

func DecodeManyPetifesProposalsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*petifiesProtoV1.ManyPetifiesProposals)
	if !ok {
		return nil, MustBeProtoRespErr
	}

	var petifiesProposals []*models.PetifiesProposal
	for _, p := range resp.PetifiesProposals {
		p_, err := decodePetifiesProposalProtoModel(p)
		if err != nil {
			return nil, err
		}
		petifiesProposals = append(petifiesProposals, p_)
	}

	return &models.ManyPetifiesProposals{
		PetifiesProposals: petifiesProposals,
	}, nil
}

func DecodeManyReviewsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*petifiesProtoV1.ManyReviews)
	if !ok {
		return nil, MustBeProtoRespErr
	}

	var reviews []*models.Review
	for _, r := range resp.Reviews {
		review, err := decodeReviewProtoModel(r)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return &models.ManyReviews{
		Reviews: reviews,
	}, nil
}

func DecodeAcceptProposalRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.AcceptProposalRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	sessionId, err := uuid.Parse(req.SessionId)
	if err != nil {
		return nil, err
	}
	proposalId, err := uuid.Parse(req.ProposalId)
	if err != nil {
		return nil, err
	}

	return &models.AcceptProposalReq{
		UserId:     userId,
		SessionId:  sessionId,
		ProposalId: proposalId,
	}, nil
}

func DecodeAcceptProposalResponse(_ context.Context, response interface{}) (interface{}, error) {
	_, ok := response.(*petifiesProtoV1.AcceptProposalResponse)
	if !ok {
		return nil, MustBeProtoRespErr
	}

	return &models.AcceptProposalResp{}, nil
}

func DecodeRejectProposalRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.RejectProposalRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	sessionId, err := uuid.Parse(req.SessionId)
	if err != nil {
		return nil, err
	}
	proposalId, err := uuid.Parse(req.ProposalId)
	if err != nil {
		return nil, err
	}

	return &models.RejectProposalReq{
		UserId:     userId,
		SessionId:  sessionId,
		ProposalId: proposalId,
	}, nil
}

func DecodeRejectProposalResponse(_ context.Context, response interface{}) (interface{}, error) {
	_, ok := response.(*petifiesProtoV1.RejectProposalResponse)
	if !ok {
		return nil, MustBeProtoRespErr
	}

	return &models.RejectProposalResp{}, nil
}

func DecodeCancelProposalRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*petifiesProtoV1.CancelProposalRequest)
	if !ok {
		return nil, MustBeProtoReqErr
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}
	proposalId, err := uuid.Parse(req.ProposalId)
	if err != nil {
		return nil, err
	}

	return &models.CancelProposalReq{
		UserId:     userId,
		ProposalId: proposalId,
	}, nil
}

func DecodeCancelProposalResponse(_ context.Context, response interface{}) (interface{}, error) {
	_, ok := response.(*petifiesProtoV1.CancelProposalResponse)
	if !ok {
		return nil, MustBeProtoRespErr
	}

	return &models.RejectProposalResp{}, nil
}

func decodePetifiesProtoModel(petifies *petifiesProtoV1.Petifies) (*models.Petifies, error) {
	id, err := uuid.Parse(petifies.GetId())
	if err != nil {
		return nil, err
	}
	ownerID, err := uuid.Parse(petifies.GetOwnerId())
	if err != nil {
		return nil, err
	}

	return &models.Petifies{
		ID:          id,
		OwnerID:     ownerID,
		Type:        petifies.Type.String(),
		Title:       petifies.Title,
		Description: petifies.Description,
		PetName:     petifies.PetName,
		Images: commonutils.Map2(petifies.Images, func(i *commonProto.Image) models.Image {
			return models.Image{URI: i.Uri, Description: i.Description}
		}),
		Status:    petifies.Status.String(),
		Address:   *decodeAddressProtoModel(petifies.Address),
		CreatedAt: petifies.CreatedAt.AsTime(),
		UpdatedAt: petifies.UpdatedAt.AsTime(),
	}, nil
}

func decodePetifiesSessionProtoModel(session *petifiesProtoV1.PetifiesSession) (*models.PetifiesSession, error) {
	id, err := uuid.Parse(session.GetId())
	if err != nil {
		return nil, err
	}
	petifiesID, err := uuid.Parse(session.GetPetifiesId())
	if err != nil {
		return nil, err
	}

	return &models.PetifiesSession{
		ID:         id,
		PetifiesID: petifiesID,
		FromTime:   session.FromTime.AsTime(),
		ToTime:     session.ToTime.AsTime(),
		Status:     session.Status.String(),
		CreatedAt:  session.CreatedAt.AsTime(),
		UpdatedAt:  session.UpdatedAt.AsTime(),
	}, nil
}

func decodePetifiesProposalProtoModel(proposal *petifiesProtoV1.PetifiesProposal) (*models.PetifiesProposal, error) {
	id, err := uuid.Parse(proposal.GetId())
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(proposal.GetUserId())
	if err != nil {
		return nil, err
	}
	petifiesSessionID, err := uuid.Parse(proposal.GetPetifiesSessionId())
	if err != nil {
		return nil, err
	}

	return &models.PetifiesProposal{
		ID:                id,
		UserID:            userID,
		PetifiesSessionID: petifiesSessionID,
		Proposal:          proposal.Proposal,
		Status:            proposal.Status.String(),
		CreatedAt:         proposal.CreatedAt.AsTime(),
		UpdatedAt:         proposal.UpdatedAt.AsTime(),
	}, nil
}

func decodeReviewProtoModel(review *petifiesProtoV1.Review) (*models.Review, error) {
	id, err := uuid.Parse(review.GetId())
	if err != nil {
		return nil, err
	}
	petifiesID, err := uuid.Parse(review.GetPetifiesId())
	if err != nil {
		return nil, err
	}
	authorID, err := uuid.Parse(review.GetAuthorId())
	if err != nil {
		return nil, err
	}

	return &models.Review{
		ID:         id,
		PetifiesID: petifiesID,
		AuthorID:   authorID,
		Review:     review.Review,
		Image: models.Image{
			URI:         review.Image.Uri,
			Description: review.Image.Description,
		},
		CreatedAt: review.CreatedAt.AsTime(),
		UpdatedAt: review.UpdatedAt.AsTime(),
	}, nil
}

func decodeAddressProtoModel(address *commonProto.Address) *models.Address {
	return &models.Address{
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
