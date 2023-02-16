package userservice

import (
	"context"
	"database/sql"
	"errors"

	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/user-service/cmd"
	userAggre "github.com/vantoan19/Petifies/server/services/user-service/internal/domain/aggregates/user"
	userRepo "github.com/vantoan19/Petifies/server/services/user-service/internal/domain/repositories/user"
	postgreRepo "github.com/vantoan19/Petifies/server/services/user-service/internal/domain/repositories/user/postgres"
	"github.com/vantoan19/Petifies/server/services/user-service/internal/jwt"
	"github.com/vantoan19/Petifies/server/services/user-service/internal/utils"
)

var logger = logging.New("UserService.Service")

type UserConfiguration func(us *userService) error

type userService struct {
	userRepository userRepo.UserRepository
	tokenMaker     jwt.TokenMaker
}

type UserService interface {
	CreateUser(ctx context.Context, email, password, firstName, lastName string) (*userAggre.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

func New(cfgs ...UserConfiguration) (UserService, error) {
	us := &userService{
		tokenMaker: jwt.NewJWTMaker(cmd.Conf.TokenSecretKey),
	}
	for _, cfg := range cfgs {
		err := cfg(us)
		if err != nil {
			return nil, err
		}
	}
	return us, nil
}

func WithUserRepository(ur userRepo.UserRepository) UserConfiguration {
	return func(us *userService) error {
		us.userRepository = ur
		return nil
	}
}

func WithPostgreUserRepository(db *sql.DB) UserConfiguration {
	return func(us *userService) error {
		pgRepo, _ := postgreRepo.New(db)
		us.userRepository = pgRepo
		return nil
	}
}

func (s *userService) CreateUser(ctx context.Context, email, password, firstName, lastName string) (*userAggre.User, error) {
	logger.Info("Start UserService.CreateUser")

	userAg, errs := userAggre.New(email, password, firstName, lastName)
	if errs.Exist() {
		logger.ErrorData("Finished UserService.CreateUser: FAILED", logging.Data{"error": errs.Error()})
		return nil, errs[0]
	}
	createdUser, err := s.userRepository.Add(userAg)
	if err != nil {
		logger.ErrorData("Finished UserService.CreateUser: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished UserService.CreateUser: SUCCESSFUL")
	return &createdUser, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	logger.Info("Start UserService.Login")

	userAg, err := s.userRepository.GetByEmail(email)
	if err != nil {
		logger.ErrorData("Finished UserService.Login: FAILED", logging.Data{"error": err.Error()})
		return "", err
	}

	if !utils.ComparePassword(password, userAg.GetPassword()) {
		logger.Error("Finished UserService.Login: FAILED: incorrect password")
		return "", errors.New("incorrect password")
	}

	logger.Info("Finished UserService.Login: SUCCESSFUL")
	return s.tokenMaker.CreateToken(userAg.GetID(), cmd.Conf.AccessTokenDuration)
}
