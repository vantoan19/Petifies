package relationshipservice

import (
	"context"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	"github.com/vantoan19/Petifies/server/libs/logging-config"
	relationshipclient "github.com/vantoan19/Petifies/server/services/grpc-clients/relationship-client"
	userservice "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/application/services/user"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/domain/repositories"
	redisRelationshipCache "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/infra/repositories/relationship/redis"
	"github.com/vantoan19/Petifies/server/services/relationship-service/pkg/models"
)

var logger = logging.New("MobileGateway.Relationship")

type RelationshipConfiguration func(rs *relationshipService) error

type relationshipService struct {
	relationshipClient    relationshipclient.RelationshipClient
	relationshipCacheRepo repositories.RelationshipCacheRepository
	userService           userservice.UserService
}

type RelationshipService interface {
	AddRelationship(ctx context.Context, req *models.AddRelationshipReq) (*models.AddRelationshipResp, error)
	RemoveRelationship(ctx context.Context, req *models.RemoveRelationshipReq) (*models.RemoveRelationshipResp, error)
	ListFollowers(ctx context.Context, req *models.ListFollowersReq) (*models.ListFollowersResp, error)
	ListFollowings(ctx context.Context, req *models.ListFollowingsReq) (*models.ListFollowingsResp, error)
}

func NewRelationshipService(relationshipClientConn *grpc.ClientConn, userService userservice.UserService, cfgs ...RelationshipConfiguration) (RelationshipService, error) {
	rs := &relationshipService{
		relationshipClient: relationshipclient.New(relationshipClientConn),
		userService:        userService,
	}
	for _, cfg := range cfgs {
		err := cfg(rs)
		if err != nil {
			return nil, err
		}
	}
	return rs, nil
}

func WithRedisRelationshipCacheRepository(client *redis.Client, relationshipClient relationshipclient.RelationshipClient) RelationshipConfiguration {
	return func(rs *relationshipService) error {
		repo := redisRelationshipCache.NewRedisRelationshipCacheRepository(client, relationshipClient)
		rs.relationshipCacheRepo = repo
		return nil
	}
}

func (rs *relationshipService) AddRelationship(ctx context.Context, req *models.AddRelationshipReq) (*models.AddRelationshipResp, error) {
	logger.Info("Start AddRelationship")

	addResp, err := rs.relationshipClient.AddRelationship(ctx, req)
	if err != nil {
		logger.ErrorData("Finished AddRelationship: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	// caching
	go func() {
		if req.RelationshipType == "FOLLOW" {
			fromUserFollowings, err := rs.relationshipClient.ListFollowings(context.Background(), req.FromUserID)
			if err == nil {
				err = rs.relationshipCacheRepo.SetFollowingsInfo(context.Background(), req.FromUserID, fromUserFollowings)
				if err != nil {
					logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
				}
			}
			toUserFollowers, err := rs.relationshipClient.ListFollowers(context.Background(), req.FromUserID)
			if err == nil {
				err = rs.relationshipCacheRepo.SetFollowersInfo(context.Background(), req.ToUserID, toUserFollowers)
				if err != nil {
					logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
				}
			}
		}
	}()

	logger.Info("Finish AddRelationship: Successful")
	return addResp, nil
}

func (rs *relationshipService) RemoveRelationship(ctx context.Context, req *models.RemoveRelationshipReq) (*models.RemoveRelationshipResp, error) {
	logger.Info("Start RemoveRelationship")

	removeResp, err := rs.relationshipClient.RemoveRelationship(ctx, req)
	if err != nil {
		logger.ErrorData("Finished RemoveRelationship: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	// caching
	go func() {
		if req.RelationshipType == "FOLLOW" {
			fromUserFollowings, err := rs.relationshipClient.ListFollowings(context.Background(), req.FromUserID)
			if err == nil {
				err = rs.relationshipCacheRepo.SetFollowingsInfo(context.Background(), req.FromUserID, fromUserFollowings)
				if err != nil {
					logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
				}
			}
			toUserFollowers, err := rs.relationshipClient.ListFollowers(context.Background(), req.FromUserID)
			if err == nil {
				err = rs.relationshipCacheRepo.SetFollowersInfo(context.Background(), req.ToUserID, toUserFollowers)
				if err != nil {
					logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
				}
			}
		}
	}()

	logger.Info("Finish RemoveRelationship: Successful")
	return removeResp, nil
}

func (rs *relationshipService) ListFollowers(ctx context.Context, req *models.ListFollowersReq) (*models.ListFollowersResp, error) {
	logger.Info("Start ListFollowers")

	followers, err := rs.relationshipCacheRepo.GetFollowersInfo(ctx, req.UserID)
	if err != nil {
		logger.ErrorData("Finished ListFollowers: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListFollowers: Successful")
	return followers, nil
}

func (rs *relationshipService) ListFollowings(ctx context.Context, req *models.ListFollowingsReq) (*models.ListFollowingsResp, error) {
	logger.Info("Start ListFollowings")

	followings, err := rs.relationshipCacheRepo.GetFollowingsInfo(ctx, req.UserID)
	if err != nil {
		logger.ErrorData("Finished ListFollowings: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finish ListFollowings: Successful")
	return followings, nil
}
