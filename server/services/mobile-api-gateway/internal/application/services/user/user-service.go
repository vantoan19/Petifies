package userservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	commonProto "github.com/vantoan19/Petifies/proto/common"
	commonutils "github.com/vantoan19/Petifies/server/libs/common-utils"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	userclient "github.com/vantoan19/Petifies/server/services/grpc-clients/user-client"
	"github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/domain/repositories"
	redisUserCache "github.com/vantoan19/Petifies/server/services/mobile-api-gateway/internal/infra/repositories/user/redis"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/models"
	"github.com/vantoan19/Petifies/server/services/user-service/pkg/translator"
)

var logger = logging.New("MobileGateway.UserService")

type UserConfiguration func(us *userService) error

type userService struct {
	userClient userclient.UserClient
	userCacheRepo repositories.UserCacheRepository
}

type UserService interface {
	CreateUser(ctx context.Context, req *commonProto.CreateUserRequest) (*commonProto.User, error)
	Login(ctx context.Context, req *commonProto.LoginRequest) (*commonProto.LoginResponse, error)
	RefreshToken(ctx context.Context, req *commonProto.RefreshTokenRequest) (*commonProto.RefreshTokenResponse, error)
	GetMyInfo(ctx context.Context) (*models.User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error)
}

func NewUserService(userClientConn *grpc.ClientConn, cfgs ...UserConfiguration) (UserService, error) {
	us := &userService{
		userClient: userclient.New(userClientConn),
	}
	for _, cfg := range cfgs {
		err := cfg(us)
		if err != nil {
			return nil, err
		}
	}
	return us, nil
}

func WithRedisUserCacheRepository(client *redis.Client) UserConfiguration {
	return func(us *userService) error {
		repo := redisUserCache.NewRedisUserCacheRepository(client)
		us.userCacheRepo = repo
		return nil
	}
}

func (s *userService) CreateUser(ctx context.Context, req *commonProto.CreateUserRequest) (*commonProto.User, error) {
	logger.Info("Start UserService.CreateUser")

	logger.Info("Executing UserService.CreateUser: forwarding the request to UserService")
	resp, err := s.userClient.CreateUserForward(ctx, req)
	if err != nil {
		logger.ErrorData("Finished UserService.CreateUser: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	// save to cache
	go func ()  {
		user, _ := translator.DecodeCreateUserResponse(ctx, resp)
		userModel, ok := user.(models.User)
		if ok {
			err := s.userCacheRepo.SetUser(ctx, userModel.ID, userModel)
			if err != nil {
				logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
			}
		}
	}()

	logger.Info("Finished UserService.CreateUser: SUCCESSFUL")
	return resp, nil
}

func (s *userService) Login(ctx context.Context, req *commonProto.LoginRequest) (*commonProto.LoginResponse, error) {
	logger.Info("Start UserService.Login")

	logger.Info("Executing UserService.Login: forwarding the request to UserService")
	resp, err := s.userClient.LoginForward(ctx, req)
	if err != nil {
		logger.ErrorData("Finished UserService.Login: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	// save to cache
	go func ()  {
		loginResp, _ := translator.DecodeLoginResponse(ctx, resp)
		loginRespModel, ok := loginResp.(models.LoginResp)
		if ok {
			if exist, err := s.userCacheRepo.ExistsUser(ctx, loginRespModel.User.ID); !exist && err == nil {
				err := s.userCacheRepo.SetUser(ctx, loginRespModel.User.ID, loginRespModel.User)
				if err != nil {
					logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
				}
			}
		}
	}()

	logger.Info("Finished UserService.Login: SUCCESSFUL")
	return resp, nil
}

func (s *userService) RefreshToken(ctx context.Context, req *commonProto.RefreshTokenRequest) (*commonProto.RefreshTokenResponse, error) {
	logger.Info("Start UserService.RefreshToken")

	logger.Info("Executing UserService.RefreshToken: forwarding the request to UserService")
	resp, err := s.userClient.RefreshTokenForward(ctx, req)
	if err != nil {
		logger.ErrorData("Finished UserService.RefreshToken: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished UserService.RefreshToken: SUCCESSFUL")
	return resp, nil
}

func (s *userService) GetMyInfo(ctx context.Context) (*models.User, error) {
	logger.Info("Start UserService.GetMyInfo")

	userID, err := commonutils.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	var user *models.User
	// Get from cache
	if exist, err := s.userCacheRepo.ExistsUser(ctx, userID); exist {
		logger.Info("Executing UserService.GetMyInfo: getting user info from cache")
		user_, err := s.userCacheRepo.GetUser(ctx, userID)
		if err != nil {
			logger.ErrorData("Finished UserService.GetMyInfo: FAILED", logging.Data{"error": err.Error()})
			return nil, err
		}
		user = user_
	} else if err != nil {
		logger.ErrorData("Finished UserService.GetMyInfo: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	} else { // Get from user service
		logger.Info("Executing UserService.GetMyInfo: forwarding the request to UserService")
		resp, err := s.userClient.GetUser(ctx, userID)
		if err != nil {
			logger.ErrorData("Finished UserService.GetMyInfo: FAILED", logging.Data{"error": err.Error()})
			return nil, err
		}
		// save to cache
		go func ()  {
			err := s.userCacheRepo.SetUser(ctx, userID, *resp)
			if err != nil {
				logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
			}
		}()
		user = resp
	}

	logger.Info("Finished UserService.GetMyInfo: SUCCESSFUL")
	return user, nil
}

func (s *userService) GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	logger.Info("Start GetUser")

	var user *models.User
	// Get from cache
	if exist, err := s.userCacheRepo.ExistsUser(ctx, userID); exist {
		logger.Info("Executing GetUser: getting user info from cache")
		user_, err := s.userCacheRepo.GetUser(ctx, userID)
		if err != nil {
			logger.ErrorData("Finished GetUser: FAILED", logging.Data{"error": err.Error()})
			return nil, err
		}
		user = user_
	} else if err != nil {
		logger.ErrorData("Finished GetUser: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	} else { // Get from user service
		logger.Info("Executing GetUser: forwarding the request to UserService")
		resp, err := s.userClient.GetUser(ctx, userID)
		if err != nil {
			logger.ErrorData("Finished GetUser: FAILED", logging.Data{"error": err.Error()})
			return nil, err
		}
		// save to cache
		go func ()  {
			err := s.userCacheRepo.SetUser(ctx, userID, *resp)
			if err != nil {
				logger.WarningData("Error at setting cache", logging.Data{"error": err.Error()})
			}
		}()
		user = resp
	}

	logger.Info("Finished GetUser: SUCCESSFUL")
	return user, nil
}
