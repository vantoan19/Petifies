package petifiesclient

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

	petifiesProtoV1 "github.com/vantoan19/Petifies/proto/petifies-service/v1"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
	"github.com/vantoan19/Petifies/server/services/petifies-service/pkg/translators"
)

var logger = logging.New("PetifiesClient")

const petifiesService = "petifies_service.v1.PetifiesService"

type petifiesClient struct {
	createPetifies           endpoint.Endpoint
	createPetifiesSession    endpoint.Endpoint
	createPetifiesProposal   endpoint.Endpoint
	createReview             endpoint.Endpoint
	getPetifiesById          endpoint.Endpoint
	getPetifiesSessionById   endpoint.Endpoint
	getPetifiesProposalById  endpoint.Endpoint
	getReviewById            endpoint.Endpoint
	listPetifiesByIds        endpoint.Endpoint
	listPetifiesByUserId     endpoint.Endpoint
	listSessionsByIds        endpoint.Endpoint
	listProposalsByIds       endpoint.Endpoint
	listReviewsByIds         endpoint.Endpoint
	listReviewsByUserId      endpoint.Endpoint
	listReviewsByPetifiesId  endpoint.Endpoint
	listProposalsByUserId    endpoint.Endpoint
	listProposalsBySessionId endpoint.Endpoint
	listSessionsByPetifiesId endpoint.Endpoint
}

type PetifiesClient interface {
	CreatePetifies(ctx context.Context, req *models.CreatePetifiesReq) (*models.Petifies, error)
	CreatePetifiesSession(ctx context.Context, req *models.CreatePetifiesSessionReq) (*models.PetifiesSession, error)
	CreatePetifiesProposal(ctx context.Context, req *models.CreatePetifiesProposalReq) (*models.PetifiesProposal, error)
	CreateReview(ctx context.Context, req *models.CreateReviewReq) (*models.Review, error)
	GetPetifiesById(ctx context.Context, req *models.GetPetifiesByIdReq) (*models.Petifies, error)
	GetPetifiesSessionById(ctx context.Context, req *models.GetSessionByIdReq) (*models.PetifiesSession, error)
	GetPetifiesProposalById(ctx context.Context, req *models.GetProposalByIdReq) (*models.PetifiesProposal, error)
	GetReviewById(ctx context.Context, req *models.GetReviewByIdReq) (*models.Review, error)
	ListPetifiesByIds(ctx context.Context, req *models.ListPetifiesByIdsReq) (*models.ManyPetifies, error)
	ListPetifiesByUserId(ctx context.Context, req *models.ListPetifiesByOwnerIdReq) (*models.ManyPetifies, error)
	ListSessionsByIds(ctx context.Context, req *models.ListSessionsByIdsReq) (*models.ManyPetifiesSessions, error)
	ListSessionsByPetifiesId(ctx context.Context, req *models.ListSessionsByPetifiesIdReq) (*models.ManyPetifiesSessions, error)
	ListProposalsByIds(ctx context.Context, req *models.ListProposalsByIdsReq) (*models.ManyPetifiesProposals, error)
	ListProposalsByUserId(ctx context.Context, req *models.ListProposalsByUserIdReq) (*models.ManyPetifiesProposals, error)
	ListReviewsByIds(ctx context.Context, req *models.ListReviewsByIdsReq) (*models.ManyReviews, error)
	ListReviewsByUserId(ctx context.Context, req *models.ListReviewsByUserIdReq) (*models.ManyReviews, error)
	ListReviewsByPetifiesId(ctx context.Context, req *models.ListReviewsByPetifiesIdReq) (*models.ManyReviews, error)
	ListProposalsBySessionId(ctx context.Context, req *models.ListProposalsBySessionIdReq) (*models.ManyPetifiesProposals, error)
}

func New(conn *grpc.ClientConn) PetifiesClient {
	return &petifiesClient{
		createPetifies: grpctransport.NewClient(
			conn,
			petifiesService,
			"CreatePetifies",
			translators.EncodeCreatePetifiesRequest,
			translators.DecodePetifiesResponse,
			petifiesProtoV1.Petifies{},
		).Endpoint(),
		createPetifiesSession: grpctransport.NewClient(
			conn,
			petifiesService,
			"CreatePetifiesSession",
			translators.EncodeCreatePetifiesSessionRequest,
			translators.DecodePetifiesSessionResponse,
			petifiesProtoV1.PetifiesSession{},
		).Endpoint(),
		createPetifiesProposal: grpctransport.NewClient(
			conn,
			petifiesService,
			"CreatePetifiesProposal",
			translators.EncodeCreatePetifiesProposalRequest,
			translators.DecodePetifiesProposalResponse,
			petifiesProtoV1.PetifiesProposal{},
		).Endpoint(),
		createReview: grpctransport.NewClient(
			conn,
			petifiesService,
			"CreateReview",
			translators.EncodeCreateReviewRequest,
			translators.DecodeReviewResponse,
			petifiesProtoV1.Review{},
		).Endpoint(),
		getPetifiesById: grpctransport.NewClient(
			conn,
			petifiesService,
			"GetPetifiesById",
			translators.EncodeGetPetifiesByIdRequest,
			translators.DecodePetifiesResponse,
			petifiesProtoV1.Petifies{},
		).Endpoint(),
		getPetifiesSessionById: grpctransport.NewClient(
			conn,
			petifiesService,
			"GetSessionById",
			translators.EncodeGetSessionByIdRequest,
			translators.DecodePetifiesProposalResponse,
			petifiesProtoV1.PetifiesSession{},
		).Endpoint(),
		getPetifiesProposalById: grpctransport.NewClient(
			conn,
			petifiesService,
			"GetProposalById",
			translators.EncodeGetProposalByIdRequest,
			translators.DecodePetifiesProposalResponse,
			petifiesProtoV1.PetifiesProposal{},
		).Endpoint(),
		getReviewById: grpctransport.NewClient(
			conn,
			petifiesService,
			"GetReviewById",
			translators.EncodeGetReviewByIdRequest,
			translators.DecodeReviewResponse,
			petifiesProtoV1.Review{},
		).Endpoint(),
		listPetifiesByIds: grpctransport.NewClient(
			conn,
			petifiesService,
			"ListPetifiesByIds",
			translators.EncodeListPetifiesByIdsRequest,
			translators.DecodeManyPetifesResponse,
			petifiesProtoV1.ManyPetifies{},
		).Endpoint(),
		listSessionsByIds: grpctransport.NewClient(
			conn,
			petifiesService,
			"ListSessionsByIds",
			translators.EncodeListSessionsByIdsRequest,
			translators.DecodeManyPetifesSessionsResponse,
			petifiesProtoV1.ManyPetifiesSessions{},
		).Endpoint(),
		listProposalsByIds: grpctransport.NewClient(
			conn,
			petifiesService,
			"ListProposalsByIds",
			translators.EncodeListProposalsByIdsRequest,
			translators.DecodeManyPetifesProposalsResponse,
			petifiesProtoV1.ManyPetifiesProposals{},
		).Endpoint(),
		listReviewsByIds: grpctransport.NewClient(
			conn,
			petifiesService,
			"ListReviewsByIds",
			translators.EncodeListReviewsByIdsRequest,
			translators.DecodeManyReviewsResponse,
			petifiesProtoV1.ManyReviews{},
		).Endpoint(),
		listPetifiesByUserId: grpctransport.NewClient(
			conn,
			petifiesService,
			"ListPetifiesByOwnerId",
			translators.EncodeListPetifiesByOwnerIdRequest,
			translators.DecodeManyPetifesResponse,
			petifiesProtoV1.ManyPetifies{},
		).Endpoint(),
		listReviewsByUserId: grpctransport.NewClient(
			conn,
			petifiesService,
			"ListReviewsByUserId",
			translators.EncodeListReviewsByUserIdRequest,
			translators.DecodeManyReviewsResponse,
			petifiesProtoV1.ManyReviews{},
		).Endpoint(),
		listProposalsByUserId: grpctransport.NewClient(
			conn,
			petifiesService,
			"ListProposalsByUserId",
			translators.EncodeListProposalsByUserIdRequest,
			translators.DecodeManyPetifesProposalsResponse,
			petifiesProtoV1.ManyPetifiesProposals{},
		).Endpoint(),
		listSessionsByPetifiesId: grpctransport.NewClient(
			conn,
			petifiesService,
			"ListSessionsByPetifiesId",
			translators.EncodeListSessionsByPetifiesIdRequest,
			translators.DecodeManyPetifesSessionsResponse,
			petifiesProtoV1.ManyPetifiesSessions{},
		).Endpoint(),
		listProposalsBySessionId: grpctransport.NewClient(
			conn,
			petifiesService,
			"ListProposalsBySessionId",
			translators.EncodeListProposalsBySessionIdRequest,
			translators.DecodeManyPetifesProposalsResponse,
			petifiesProtoV1.ManyPetifiesProposals{},
		).Endpoint(),
		listReviewsByPetifiesId: grpctransport.NewClient(
			conn,
			petifiesService,
			"ListReviewsByPetifiesId",
			translators.EncodeListReviewsByPetifiesIdRequest,
			translators.DecodeManyReviewsResponse,
			petifiesProtoV1.ManyReviews{},
		).Endpoint(),
	}
}

func (p *petifiesClient) CreatePetifies(ctx context.Context, req *models.CreatePetifiesReq) (*models.Petifies, error) {
	logger.Info("Start CreatePetifies")

	resp, err := p.createPetifies(ctx, req)
	if err != nil {
		logger.ErrorData("Finish CreatePetifies: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish CreatePetifies: Successful")
	return resp.(*models.Petifies), nil
}

func (p *petifiesClient) CreatePetifiesSession(ctx context.Context, req *models.CreatePetifiesSessionReq) (*models.PetifiesSession, error) {
	logger.Info("Start CreatePetifiesSession")

	resp, err := p.createPetifiesSession(ctx, req)
	if err != nil {
		logger.ErrorData("Finish CreatePetifiesSession: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish CreatePetifiesSession: Successful")
	return resp.(*models.PetifiesSession), nil
}

func (p *petifiesClient) CreatePetifiesProposal(ctx context.Context, req *models.CreatePetifiesProposalReq) (*models.PetifiesProposal, error) {
	logger.Info("Start CreatePetifiesProposal")

	resp, err := p.createPetifiesProposal(ctx, req)
	if err != nil {
		logger.ErrorData("Finish CreatePetifiesProposal: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish CreatePetifiesProposal: Successful")
	return resp.(*models.PetifiesProposal), nil
}

func (p *petifiesClient) CreateReview(ctx context.Context, req *models.CreateReviewReq) (*models.Review, error) {
	logger.Info("Start CreateReview")

	resp, err := p.createReview(ctx, req)
	if err != nil {
		logger.ErrorData("Finish CreateReview: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish CreateReview: Successful")
	return resp.(*models.Review), nil
}

func (p *petifiesClient) GetPetifiesById(ctx context.Context, req *models.GetPetifiesByIdReq) (*models.Petifies, error) {
	logger.Info("Start GetPetifiesById")

	resp, err := p.getPetifiesById(ctx, req)
	if err != nil {
		logger.ErrorData("Finish GetPetifiesById: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetPetifiesById: Successful")
	return resp.(*models.Petifies), nil
}

func (p *petifiesClient) GetPetifiesSessionById(ctx context.Context, req *models.GetSessionByIdReq) (*models.PetifiesSession, error) {
	logger.Info("Start GetPetifiesSessionById")

	resp, err := p.getPetifiesSessionById(ctx, req)
	if err != nil {
		logger.ErrorData("Finish GetPetifiesSessionById: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetPetifiesSessionById: Successful")
	return resp.(*models.PetifiesSession), nil
}

func (p *petifiesClient) GetPetifiesProposalById(ctx context.Context, req *models.GetProposalByIdReq) (*models.PetifiesProposal, error) {
	logger.Info("Start GetPetifiesProposalById")

	resp, err := p.getPetifiesProposalById(ctx, req)
	if err != nil {
		logger.ErrorData("Finish GetPetifiesProposalById: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetPetifiesProposalById: Successful")
	return resp.(*models.PetifiesProposal), nil
}

func (p *petifiesClient) GetReviewById(ctx context.Context, req *models.GetReviewByIdReq) (*models.Review, error) {
	logger.Info("Start GetReviewById")

	resp, err := p.getReviewById(ctx, req)
	if err != nil {
		logger.ErrorData("Finish GetReviewById: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish GetReviewById: Successful")
	return resp.(*models.Review), nil
}

func (p *petifiesClient) ListPetifiesByIds(ctx context.Context, req *models.ListPetifiesByIdsReq) (*models.ManyPetifies, error) {
	logger.Info("Start ListPetifiesByIds")

	resp, err := p.listPetifiesByIds(ctx, req)
	if err != nil {
		logger.ErrorData("Finish ListPetifiesByIds: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListPetifiesByIds: Successful")
	return resp.(*models.ManyPetifies), nil
}

func (p *petifiesClient) ListSessionsByIds(ctx context.Context, req *models.ListSessionsByIdsReq) (*models.ManyPetifiesSessions, error) {
	logger.Info("Start ListSessionsByIds")

	resp, err := p.listSessionsByIds(ctx, req)
	if err != nil {
		logger.ErrorData("Finish ListSessionsByIds: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListSessionsByIds: Successful")
	return resp.(*models.ManyPetifiesSessions), nil
}

func (p *petifiesClient) ListProposalsByIds(ctx context.Context, req *models.ListProposalsByIdsReq) (*models.ManyPetifiesProposals, error) {
	logger.Info("Start ListProposalsByIds")

	resp, err := p.listProposalsByIds(ctx, req)
	if err != nil {
		logger.ErrorData("Finish ListProposalsByIds: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListProposalsByIds: Successful")
	return resp.(*models.ManyPetifiesProposals), nil
}

func (p *petifiesClient) ListReviewsByIds(ctx context.Context, req *models.ListReviewsByIdsReq) (*models.ManyReviews, error) {
	logger.Info("Start ListReviewsByIds")

	resp, err := p.listReviewsByIds(ctx, req)
	if err != nil {
		logger.ErrorData("Finish ListReviewsByIds: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListReviewsByIds: Successful")
	return resp.(*models.ManyReviews), nil
}

func (p *petifiesClient) ListReviewsByUserId(ctx context.Context, req *models.ListReviewsByUserIdReq) (*models.ManyReviews, error) {
	logger.Info("Start ListReviewsByUserId")

	resp, err := p.listReviewsByUserId(ctx, req)
	if err != nil {
		logger.ErrorData("Finish ListReviewsByUserId: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListReviewsByUserId: Successful")
	return resp.(*models.ManyReviews), nil
}

func (p *petifiesClient) ListProposalsByUserId(ctx context.Context, req *models.ListProposalsByUserIdReq) (*models.ManyPetifiesProposals, error) {
	logger.Info("Start ListProposalsByUserId")

	resp, err := p.listProposalsByUserId(ctx, req)
	if err != nil {
		logger.ErrorData("Finish ListProposalsByUserId: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListProposalsByUserId: Successful")
	return resp.(*models.ManyPetifiesProposals), nil
}

func (p *petifiesClient) ListPetifiesByUserId(ctx context.Context, req *models.ListPetifiesByOwnerIdReq) (*models.ManyPetifies, error) {
	logger.Info("Start ListPetifiesByUserId")

	resp, err := p.listPetifiesByUserId(ctx, req)
	if err != nil {
		logger.ErrorData("Finish ListPetifiesByUserId: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListPetifiesByUserId: Successful")
	return resp.(*models.ManyPetifies), nil
}

func (p *petifiesClient) ListSessionsByPetifiesId(ctx context.Context, req *models.ListSessionsByPetifiesIdReq) (*models.ManyPetifiesSessions, error) {
	logger.Info("Start ListSessionsByPetifiesId")

	resp, err := p.listSessionsByPetifiesId(ctx, req)
	if err != nil {
		logger.ErrorData("Finish ListSessionsByPetifiesId: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListSessionsByPetifiesId: Successful")
	return resp.(*models.ManyPetifiesSessions), nil
}

func (p *petifiesClient) ListProposalsBySessionId(ctx context.Context, req *models.ListProposalsBySessionIdReq) (*models.ManyPetifiesProposals, error) {
	logger.Info("Start ListProposalsBySessionId")

	resp, err := p.listProposalsBySessionId(ctx, req)
	if err != nil {
		logger.ErrorData("Finish ListProposalsBySessionId: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListProposalsBySessionId: Successful")
	return resp.(*models.ManyPetifiesProposals), nil
}

func (p *petifiesClient) ListReviewsByPetifiesId(ctx context.Context, req *models.ListReviewsByPetifiesIdReq) (*models.ManyReviews, error) {
	logger.Info("Start ListReviewsByPetifiesId")

	resp, err := p.listReviewsByPetifiesId(ctx, req)
	if err != nil {
		logger.ErrorData("Finish ListReviewsByPetifiesId: Failed", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListReviewsByPetifiesId: Successful")
	return resp.(*models.ManyReviews), nil
}
