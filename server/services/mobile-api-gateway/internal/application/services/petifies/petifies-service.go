package petifiesservice

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	locationclient "github.com/vantoan19/Petifies/server/services/grpc-clients/location-client"
	petifiesclient "github.com/vantoan19/Petifies/server/services/grpc-clients/petifies-client"
	userclient "github.com/vantoan19/Petifies/server/services/grpc-clients/user-client"
	locationModels "github.com/vantoan19/Petifies/server/services/location-service/pkg/models"
	userservice "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/application/services/user"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/domain/repositories"
	redisPetifiesCache "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/infra/repositories/petifies/redis"
	redisProposalCache "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/infra/repositories/proposal/redis"
	redisReviewCache "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/infra/repositories/review/redis"
	redisSessionCache "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/infra/repositories/session/redis"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/pkg/models"
	petifiesModels "github.com/vantoan19/Petifies/server/services/petifies-service/pkg/models"
	userModels "github.com/vantoan19/Petifies/server/services/user-service/pkg/models"
)

var logger = logging.New("MobileGateway.PetifiesService")

type PetifiesConfiguration func(ps *petifiesService) error

type petifiesService struct {
	petifiesClient    petifiesclient.PetifiesClient
	userClient        userclient.UserClient
	locationClient    locationclient.LocationClient
	userService       userservice.UserService
	petifiesCacheRepo repositories.PetifiesCacheRepository
	proposalCacheRepo repositories.PetifiesProposalCacheRepository
	sessionCacheRepo  repositories.PetifiesSessionCacheRepository
	reviewCacheRepo   repositories.ReviewCacheRepository
}

type PetifiesService interface {
	UserCreatePetifies(ctx context.Context, req *models.UserCreatePetifiesReq) (*models.PetifiesWithUserInfo, error)
	UserCreateSession(ctx context.Context, req *models.UserCreatePetifiesSessionReq) (*models.PetifiesSession, error)
	UserCreateProposal(ctx context.Context, req *models.UserCreatePetifiesProposalReq) (*models.PetifiesProposalWithUserInfo, error)
	UserCreateReview(ctx context.Context, req *models.UserCreateReviewReq) (*models.ReviewWithUserInfo, error)

	GetPetifies(ctx context.Context, id uuid.UUID) (*models.PetifiesWithUserInfo, error)
	GetPetifiesSession(ctx context.Context, id uuid.UUID) (*models.PetifiesSession, error)
	GetPetifiesProposal(ctx context.Context, id uuid.UUID) (*models.PetifiesProposalWithUserInfo, error)
	GetReview(ctx context.Context, id uuid.UUID) (*models.ReviewWithUserInfo, error)

	ListNearByPetifies(ctx context.Context, req *models.ListNearByPetifiesReq) ([]*models.PetifiesWithUserInfo, error)
	ListPetifiesByUserId(ctx context.Context, userId uuid.UUID, pageSize int, afterId uuid.UUID) ([]*models.PetifiesWithUserInfo, error)
	ListSessionsByPetifiesId(ctx context.Context, petifiesId uuid.UUID, pageSize int, afterId uuid.UUID) ([]*models.PetifiesSession, error)
	ListProposalsByUserId(ctx context.Context, userId uuid.UUID, pageSize int, afterId uuid.UUID) ([]*models.PetifiesProposalWithUserInfo, error)
	ListProposalsBySessionId(ctx context.Context, sessionId uuid.UUID, pageSize int, afterId uuid.UUID) ([]*models.PetifiesProposalWithUserInfo, error)
	ListReviewsByUserId(ctx context.Context, userId uuid.UUID, pageSize int, afterId uuid.UUID) ([]*models.ReviewWithUserInfo, error)
	ListReviewsByPetifiesId(ctx context.Context, petifiesId uuid.UUID, pageSize int, afterId uuid.UUID) ([]*models.ReviewWithUserInfo, error)
}

func NewPetifiesService(
	petifiesClientConn *grpc.ClientConn,
	userClientConn *grpc.ClientConn,
	locationClientConn *grpc.ClientConn,
	userService userservice.UserService,
	cfgs ...PetifiesConfiguration,
) (PetifiesService, error) {
	ps := &petifiesService{
		petifiesClient: petifiesclient.New(petifiesClientConn),
		userClient:     userclient.New(userClientConn),
		userService:    userService,
		locationClient: locationclient.New(locationClientConn),
	}
	for _, cfg := range cfgs {
		err := cfg(ps)
		if err != nil {
			return nil, err
		}
	}
	return ps, nil
}

func WithRedisPetifiesCacheRepository(client *redis.Client, petifiesClient petifiesclient.PetifiesClient) PetifiesConfiguration {
	return func(ps *petifiesService) error {
		repo := redisPetifiesCache.NewRedisPetifiesCacheRepository(client, petifiesClient)
		ps.petifiesCacheRepo = repo
		return nil
	}
}

func WithRedisPetifiesSessionCacheRepository(client *redis.Client, petifiesClient petifiesclient.PetifiesClient) PetifiesConfiguration {
	return func(ps *petifiesService) error {
		repo := redisSessionCache.NewRedisPetifiesSessionCacheRepository(client, petifiesClient)
		ps.sessionCacheRepo = repo
		return nil
	}
}

func WithRedisPetifiesProposalCacheRepository(client *redis.Client, petifiesClient petifiesclient.PetifiesClient) PetifiesConfiguration {
	return func(ps *petifiesService) error {
		repo := redisProposalCache.NewRedisPetifiesProposalCacheRepository(client, petifiesClient)
		ps.proposalCacheRepo = repo
		return nil
	}
}

func WithRedisReviewCacheRepository(client *redis.Client, petifiesClient petifiesclient.PetifiesClient) PetifiesConfiguration {
	return func(ps *petifiesService) error {
		repo := redisReviewCache.NewRedisReviewCacheRepository(client, petifiesClient)
		ps.reviewCacheRepo = repo
		return nil
	}
}

func (ps *petifiesService) UserCreatePetifies(ctx context.Context, req *models.UserCreatePetifiesReq) (*models.PetifiesWithUserInfo, error) {
	logger.Info("Start UserCreatePetifies")

	userResp, err := ps.userService.GetMyInfo(ctx)
	if err != nil {
		logger.ErrorData("Finished UserCreatePetifies: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Executing UserCreatePetifies: delegating the request to PetifiesService")
	petifiesResp, err := ps.petifiesClient.CreatePetifies(ctx, &petifiesModels.CreatePetifiesReq{
		OwnerID:     userResp.ID,
		Type:        req.Type,
		Title:       req.Title,
		Description: req.Description,
		PetName:     req.PetName,
		Images:      req.Images,
		Address:     req.Address,
	})
	if err != nil {
		logger.ErrorData("Finished UserCreatePetifies: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	// Set cache
	go func() {
		err := ps.petifiesCacheRepo.SetPetifies(context.Background(), petifiesResp.ID, petifiesResp)
		if err != nil {
			logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
		}
	}()

	logger.Info("Finished UserCreatePetifies: SUCCESSFUL")
	return aggregatePetifiesWithUserInfo(petifiesResp, userResp), nil
}

func (ps *petifiesService) UserCreateSession(ctx context.Context, req *models.UserCreatePetifiesSessionReq) (*models.PetifiesSession, error) {
	logger.Info("Start UserCreateSession")

	userResp, err := ps.userService.GetMyInfo(ctx)
	if err != nil {
		logger.ErrorData("Finished UserCreateSession: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Executing UserCreateSession: delegating the request to PetifiesService")
	sessionResp, err := ps.petifiesClient.CreatePetifiesSession(ctx, &petifiesModels.CreatePetifiesSessionReq{
		CreatorID:  userResp.ID,
		PetifiesID: req.PetifiesId,
		FromTime:   req.FromTime,
		ToTime:     req.ToTime,
	})
	if err != nil {
		logger.ErrorData("Finished UserCreateSession: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	// Set cache
	go func() {
		err := ps.sessionCacheRepo.SetPetifiesSession(context.Background(), sessionResp.ID, sessionResp)
		if err != nil {
			logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
		}
	}()

	logger.Info("Finished UserCreateSession: SUCCESSFUL")
	return convertPetifiesSessionToGatewaySession(sessionResp), nil
}

func (ps *petifiesService) UserCreateProposal(ctx context.Context, req *models.UserCreatePetifiesProposalReq) (*models.PetifiesProposalWithUserInfo, error) {
	logger.Info("Start UserCreateProposal")

	userResp, err := ps.userService.GetMyInfo(ctx)
	if err != nil {
		logger.ErrorData("Finished UserCreateProposal: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Executing UserCreateProposal: delegating the request to PetifiesService")
	proposalResp, err := ps.petifiesClient.CreatePetifiesProposal(ctx, &petifiesModels.CreatePetifiesProposalReq{
		UserID:            userResp.ID,
		PetifiesSessionID: req.PetifiesSessionId,
		Proposal:          req.Proposal,
	})
	if err != nil {
		logger.ErrorData("Finished UserCreateProposal: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	// Set cache
	go func() {
		err := ps.proposalCacheRepo.SetPetifiesProposal(context.Background(), proposalResp.ID, proposalResp)
		if err != nil {
			logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
		}
	}()

	logger.Info("Finished UserCreateProposal: SUCCESSFUL")
	return aggregatePetifiesProposalWithUserInfo(proposalResp, userResp), nil
}

func (ps *petifiesService) UserCreateReview(ctx context.Context, req *models.UserCreateReviewReq) (*models.ReviewWithUserInfo, error) {
	logger.Info("Start UserCreateReview")

	userResp, err := ps.userService.GetMyInfo(ctx)
	if err != nil {
		logger.ErrorData("Finished UserCreateReview: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Executing UserCreateReview: delegating the request to PetifiesService")
	reviewResp, err := ps.petifiesClient.CreateReview(ctx, &petifiesModels.CreateReviewReq{
		PetifiesID: req.PetifiesId,
		AuthorID:   userResp.ID,
		Review:     req.Review,
		Image:      req.Image,
	})
	if err != nil {
		logger.ErrorData("Finished UserCreateReview: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	// Set cache
	go func() {
		err := ps.reviewCacheRepo.SetReview(context.Background(), reviewResp.ID, reviewResp)
		if err != nil {
			logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
		}
	}()

	logger.Info("Finished UserCreateReview: SUCCESSFUL")
	return aggregateReviewWithUserInfo(reviewResp, userResp), nil
}

func (ps *petifiesService) GetPetifies(ctx context.Context, id uuid.UUID) (*models.PetifiesWithUserInfo, error) {
	logger.Info("Start GetPetifies")

	petifies, err := ps.petifiesCacheRepo.GetPetifies(ctx, id)
	if err != nil {
		logger.ErrorData("Finished GetPetifies: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	userResp, err := ps.userService.GetUser(ctx, petifies.OwnerID)
	if err != nil {
		logger.ErrorData("Finished GetPetifies: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished GetPetifies: SUCCESSFUL")
	return aggregatePetifiesWithUserInfo(petifies, userResp), nil
}

func (ps *petifiesService) GetPetifiesSession(ctx context.Context, id uuid.UUID) (*models.PetifiesSession, error) {
	logger.Info("Start GetPetifiesSession")

	session, err := ps.sessionCacheRepo.GetPetifiesSession(ctx, id)
	if err != nil {
		logger.ErrorData("Finished GetPetifiesSession: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished GetPetifiesSession: SUCCESSFUL")
	return convertPetifiesSessionToGatewaySession(session), nil
}

func (ps *petifiesService) GetPetifiesProposal(ctx context.Context, id uuid.UUID) (*models.PetifiesProposalWithUserInfo, error) {
	logger.Info("Start GetPetifiesProposal")

	proposal, err := ps.proposalCacheRepo.GetPetifiesProposal(ctx, id)
	if err != nil {
		logger.ErrorData("Finished GetPetifiesProposal: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	userResp, err := ps.userService.GetUser(ctx, proposal.UserID)
	if err != nil {
		logger.ErrorData("Finished GetPetifiesProposal: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished GetPetifiesProposal: SUCCESSFUL")
	return aggregatePetifiesProposalWithUserInfo(proposal, userResp), nil
}

func (ps *petifiesService) GetReview(ctx context.Context, id uuid.UUID) (*models.ReviewWithUserInfo, error) {
	logger.Info("Start GetReview")

	review, err := ps.reviewCacheRepo.GetReview(ctx, id)
	if err != nil {
		logger.ErrorData("Finished GetReview: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	userResp, err := ps.userService.GetUser(ctx, review.AuthorID)
	if err != nil {
		logger.ErrorData("Finished GetReview: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished GetReview: SUCCESSFUL")
	return aggregateReviewWithUserInfo(review, userResp), nil
}

func (ps *petifiesService) ListNearByPetifies(ctx context.Context, req *models.ListNearByPetifiesReq) ([]*models.PetifiesWithUserInfo, error) {
	logger.Info("Start ListNearByPetifies")

	locationType, err := convertPetifiesTypeToLocationType(req.Type)
	if err != nil {
		logger.ErrorData("Finished ListNearByPetifies: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}
	petifiesLocations, err := ps.locationClient.ListNearByLocationsByType(ctx, &locationModels.ListNearByLocationsByTypeReq{
		LocationType: locationType,
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
		Radius:       req.Radius,
		PageSize:     int(req.PageSize),
		Offset:       req.Offset,
	})
	if err != nil {
		logger.ErrorData("Finished ListNearByPetifies: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	petifiesIds := commonutils.Map2(petifiesLocations.Locations, func(l *locationModels.Location) uuid.UUID { return l.EntityID })

	petifies, err := ps.petifiesCacheRepo.ListPetifies(ctx, petifiesIds)
	if err != nil {
		logger.ErrorData("Finished ListNearByPetifies: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	userIds := commonutils.Map2(petifies, func(p *petifiesModels.Petifies) uuid.UUID { return p.OwnerID })
	users, err := ps.userService.ListUsersByIds(ctx, userIds)
	if err != nil {
		logger.ErrorData("Finished ListNearByPetifies: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	result := make([]*models.PetifiesWithUserInfo, 0)
	for idx, petify := range petifies {
		if petify != nil && users[idx] != nil {
			p := aggregatePetifiesWithUserInfo(petify, users[idx])
			result = append(result, p)
		}
	}

	logger.Info("Finished ListNearByPetifies: SUCCESSFUL")
	return result, nil
}

func (ps *petifiesService) ListPetifiesByUserId(ctx context.Context, userId uuid.UUID, pageSize int, afterId uuid.UUID) ([]*models.PetifiesWithUserInfo, error) {
	petifies, err := ps.petifiesClient.ListPetifiesByUserId(ctx, &petifiesModels.ListPetifiesByOwnerIdReq{
		OwnerID:  userId,
		PageSize: pageSize,
		AfterID:  afterId,
	})
	if err != nil {
		logger.ErrorData("Finished ListPetifiesByUserId: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	userIds := commonutils.Map2(petifies.Petifies, func(p *petifiesModels.Petifies) uuid.UUID { return p.OwnerID })
	users, err := ps.userService.ListUsersByIds(ctx, userIds)
	if err != nil {
		logger.ErrorData("Finished ListPetifiesByUserId: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	result := make([]*models.PetifiesWithUserInfo, 0)
	for idx, petify := range petifies.Petifies {
		if petify != nil && users[idx] != nil {
			p := aggregatePetifiesWithUserInfo(petify, users[idx])
			result = append(result, p)
		}
	}

	logger.Info("Finished ListPetifiesByUserId: SUCCESSFUL")
	return result, nil
}

func (ps *petifiesService) ListSessionsByPetifiesId(ctx context.Context, petifiesId uuid.UUID, pageSize int, afterId uuid.UUID) ([]*models.PetifiesSession, error) {
	logger.Info("Start ListSessionsByPetifiesId")
	sessions, err := ps.petifiesClient.ListSessionsByPetifiesId(ctx, &petifiesModels.ListSessionsByPetifiesIdReq{
		PetifiesID: petifiesId,
		PageSize:   pageSize,
		AfterID:    afterId,
	})
	if err != nil {
		logger.ErrorData("Finished ListPetifiesByUserId: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished ListPetifiesByUserId: SUCCESSFUL")
	return commonutils.Map2(sessions.PetifiesSessions, func(s *petifiesModels.PetifiesSession) *models.PetifiesSession {
		return convertPetifiesSessionToGatewaySession(s)
	}), nil
}

func (ps *petifiesService) ListProposalsByUserId(ctx context.Context, userId uuid.UUID, pageSize int, afterId uuid.UUID) ([]*models.PetifiesProposalWithUserInfo, error) {
	logger.Info("Start ListProposalsByUserId")
	proposals, err := ps.petifiesClient.ListProposalsByUserId(ctx, &petifiesModels.ListProposalsByUserIdReq{
		UserId:   userId,
		PageSize: pageSize,
		AfterID:  afterId,
	})
	if err != nil {
		logger.ErrorData("Finished ListProposalsByUserId: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	userIds := commonutils.Map2(proposals.PetifiesProposals, func(p *petifiesModels.PetifiesProposal) uuid.UUID { return p.UserID })
	users, err := ps.userService.ListUsersByIds(ctx, userIds)
	if err != nil {
		logger.ErrorData("Finished ListProposalsByUserId: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	result := make([]*models.PetifiesProposalWithUserInfo, 0)
	for idx, proposal := range proposals.PetifiesProposals {
		if proposal != nil && users[idx] != nil {
			p := aggregatePetifiesProposalWithUserInfo(proposal, users[idx])
			result = append(result, p)
		}
	}

	logger.Info("Finished ListProposalsByUserId: SUCCESSFUL")
	return result, nil
}

func (ps *petifiesService) ListProposalsBySessionId(ctx context.Context, sessionId uuid.UUID, pageSize int, afterId uuid.UUID) ([]*models.PetifiesProposalWithUserInfo, error) {
	logger.Info("Start ListProposalsBySessionId")
	proposals, err := ps.petifiesClient.ListProposalsBySessionId(ctx, &petifiesModels.ListProposalsBySessionIdReq{
		PetifiesSessionID: sessionId,
		PageSize:          pageSize,
		AfterID:           afterId,
	})
	if err != nil {
		logger.ErrorData("Finished ListProposalsBySessionId: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	userIds := commonutils.Map2(proposals.PetifiesProposals, func(p *petifiesModels.PetifiesProposal) uuid.UUID { return p.UserID })
	users, err := ps.userService.ListUsersByIds(ctx, userIds)
	if err != nil {
		logger.ErrorData("Finished ListProposalsBySessionId: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	result := make([]*models.PetifiesProposalWithUserInfo, 0)
	for idx, proposal := range proposals.PetifiesProposals {
		if proposal != nil && users[idx] != nil {
			p := aggregatePetifiesProposalWithUserInfo(proposal, users[idx])
			result = append(result, p)
		}
	}

	logger.Info("Finished ListProposalsByUserId: SUCCESSFUL")
	return result, nil
}

func (ps *petifiesService) ListReviewsByUserId(ctx context.Context, userId uuid.UUID, pageSize int, afterId uuid.UUID) ([]*models.ReviewWithUserInfo, error) {
	logger.Info("Start ListReviewsByUserId")
	reviews, err := ps.petifiesClient.ListReviewsByUserId(ctx, &petifiesModels.ListReviewsByUserIdReq{
		UserId:   userId,
		PageSize: pageSize,
		AfterID:  afterId,
	})
	if err != nil {
		logger.ErrorData("Finished ListReviewsByUserId: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	userIds := commonutils.Map2(reviews.Reviews, func(r *petifiesModels.Review) uuid.UUID { return r.AuthorID })
	users, err := ps.userService.ListUsersByIds(ctx, userIds)
	if err != nil {
		logger.ErrorData("Finished ListReviewsByUserId: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	result := make([]*models.ReviewWithUserInfo, 0)
	for idx, review := range reviews.Reviews {
		if review != nil && users[idx] != nil {
			r := aggregateReviewWithUserInfo(review, users[idx])
			result = append(result, r)
		}
	}

	logger.Info("Finished ListReviewsByUserId: SUCCESSFUL")
	return result, nil
}

func (ps *petifiesService) ListReviewsByPetifiesId(ctx context.Context, petifiesId uuid.UUID, pageSize int, afterId uuid.UUID) ([]*models.ReviewWithUserInfo, error) {
	logger.Info("Start ListReviewsByUserId")
	reviews, err := ps.petifiesClient.ListReviewsByPetifiesId(ctx, &petifiesModels.ListReviewsByPetifiesIdReq{
		PetifiesID: petifiesId,
		PageSize:   pageSize,
		AfterID:    afterId,
	})
	if err != nil {
		logger.ErrorData("Finished ListReviewsByUserId: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	userIds := commonutils.Map2(reviews.Reviews, func(r *petifiesModels.Review) uuid.UUID { return r.AuthorID })
	users, err := ps.userService.ListUsersByIds(ctx, userIds)
	if err != nil {
		logger.ErrorData("Finished ListReviewsByUserId: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	result := make([]*models.ReviewWithUserInfo, 0)
	for idx, review := range reviews.Reviews {
		if review != nil && users[idx] != nil {
			r := aggregateReviewWithUserInfo(review, users[idx])
			result = append(result, r)
		}
	}

	logger.Info("Finished ListReviewsByUserId: SUCCESSFUL")
	return result, nil
}

func convertPetifiesTypeToLocationType(t string) (string, error) {
	switch t {
	case "PETIFIES_TYPE_DOG_WALKING":
		return "LOCATION_TYPE_PETIFIES_DOG_WALKING", nil
	case "PETIFIES_TYPE_CAT_PLAYING":
		return "LOCATION_TYPE_PETIFIES_CAT_PLAYING", nil
	case "PETIFIES_TYPE_DOG_SITTING":
		return "LOCATION_TYPE_PETIFIES_DOG_SITTING", nil
	case "PETIFIES_TYPE_CAT_SITTING":
		return "LOCATION_TYPE_PETIFIES_CAT_SITTING", nil
	case "PETIFIES_TYPE_DOG_ADOPTION":
		return "LOCATION_TYPE_PETIFIES_DOG_ADOPTION", nil
	case "PETIFIES_TYPE_CAT_ADOPTION":
		return "LOCATION_TYPE_PETIFIES_CAT_ADOPTION", nil
	default:
		return "", status.Errorf(codes.Unknown, "Unknown Petifies Type")
	}
}

func aggregatePetifiesWithUserInfo(petifies *petifiesModels.Petifies, user *userModels.User) *models.PetifiesWithUserInfo {
	return &models.PetifiesWithUserInfo{
		Id: petifies.ID,
		Owner: models.BasicUserInfo{
			ID:         user.ID,
			Email:      user.Email,
			UserAvatar: "",
			FirstName:  user.FirstName,
			LastName:   user.LastName,
		},
		Type:        petifies.Type,
		Title:       petifies.Title,
		Description: petifies.Description,
		PetName:     petifies.PetName,
		Images:      petifies.Images,
		Status:      petifies.Status,
		Address:     petifies.Address,
		CreatedAt:   petifies.CreatedAt,
		UpdatedAt:   petifies.UpdatedAt,
	}
}

func convertPetifiesSessionToGatewaySession(session *petifiesModels.PetifiesSession) *models.PetifiesSession {
	return &models.PetifiesSession{
		Id:         session.ID,
		PetifiesId: session.PetifiesID,
		FromTime:   session.FromTime,
		ToTime:     session.ToTime,
		Status:     session.Status,
		CreatedAt:  session.CreatedAt,
		UpdatedAt:  session.UpdatedAt,
	}
}

func aggregatePetifiesProposalWithUserInfo(proposal *petifiesModels.PetifiesProposal, user *userModels.User) *models.PetifiesProposalWithUserInfo {
	return &models.PetifiesProposalWithUserInfo{
		Id: proposal.ID,
		User: models.BasicUserInfo{
			ID:         user.ID,
			Email:      user.Email,
			UserAvatar: "",
			FirstName:  user.FirstName,
			LastName:   user.LastName,
		},
		PetifiesSessionId: proposal.PetifiesSessionID,
		Proposal:          proposal.Proposal,
		Status:            proposal.Status,
		CreatedAt:         proposal.CreatedAt,
		UpdatedAt:         proposal.UpdatedAt,
	}
}

func aggregateReviewWithUserInfo(review *petifiesModels.Review, user *userModels.User) *models.ReviewWithUserInfo {
	return &models.ReviewWithUserInfo{
		Id:         review.ID,
		PetifiesId: review.PetifiesID,
		Author: models.BasicUserInfo{
			ID:         user.ID,
			Email:      user.Email,
			UserAvatar: "",
			FirstName:  user.FirstName,
			LastName:   user.LastName,
		},
		Review:    review.Review,
		Image:     review.Image,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}
}
