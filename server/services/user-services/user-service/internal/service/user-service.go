package userservice

import (
	"context"
	"database/sql"

	userAggre "github.com/vantoan19/Petifies/server/services/user-services/user-service/internal/domain/aggregates/user"
	userRepo "github.com/vantoan19/Petifies/server/services/user-services/user-service/internal/domain/repositories/user"
	postgreRepo "github.com/vantoan19/Petifies/server/services/user-services/user-service/internal/domain/repositories/user/postgres"
)

// var logger = logging.New("UserService.Service")

type UserConfiguration func(us *userService) error

type userService struct {
	userRepository userRepo.UserRepository
}

type UserService interface {
	CreateUser(ctx context.Context, email, password, lastName, firstName string) (userAggre.User, error)
}

func New(cfgs ...UserConfiguration) (UserService, error) {
	us := &userService{}
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

func (s *userService) CreateUser(ctx context.Context, email, password, lastName, firstName string) (userAggre.User, error) {
	userAg, err := userAggre.New(email, password, lastName, firstName)
	if err != nil {
		return userAggre.User{}, err
	}

	createdUser, err := s.userRepository.Add(userAg)
	if err != nil {
		return userAggre.User{}, err
	}

	return createdUser, nil
}
