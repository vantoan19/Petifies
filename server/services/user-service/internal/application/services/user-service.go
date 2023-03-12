package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/vantoan19/Petifies/server/infrastructure/kafka/models"
	"github.com/vantoan19/Petifies/server/infrastructure/kafka/producer"
	outbox_repo "github.com/vantoan19/Petifies/server/infrastructure/outbox/repository"
	"github.com/vantoan19/Petifies/server/libs/logging-config"
	"github.com/vantoan19/Petifies/server/services/user-service/cmd"
	"github.com/vantoan19/Petifies/server/services/user-service/internal/application/handlers/jwt"
	userAggre "github.com/vantoan19/Petifies/server/services/user-service/internal/domain/aggregates/user"
	"github.com/vantoan19/Petifies/server/services/user-service/internal/domain/aggregates/user/entities"
	userRepo "github.com/vantoan19/Petifies/server/services/user-service/internal/domain/aggregates/user/repository"
	"github.com/vantoan19/Petifies/server/services/user-service/internal/domain/publisher"
	"github.com/vantoan19/Petifies/server/services/user-service/internal/infra/publishers/kafka"
	"github.com/vantoan19/Petifies/server/services/user-service/internal/infra/repositories/user/postgres"
	"github.com/vantoan19/Petifies/server/services/user-service/internal/utils"
)

var logger = logging.New("UserService.Service")

type UserConfiguration func(us *userService) error

type userService struct {
	userRepository userRepo.UserRepository
	userPublisher  publisher.UserRequestMessagePublisher
	tokenMaker     jwt.TokenMaker
}

type UserService interface {
	CreateUser(ctx context.Context, email, password, firstName, lastName string) (*userAggre.User, error)
	Login(ctx context.Context, email, password string) (uuid.UUID, string, time.Time, string, time.Time, *userAggre.User, error)
	VerifyToken(ctx context.Context, token string) (string, error)
	RefreshToken(ctx context.Context, token string) (string, time.Time, error)
	GetUser(ctx context.Context, id uuid.UUID) (*userAggre.User, error)
}

func NewUserService(cfgs ...UserConfiguration) (UserService, error) {
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
		pgRepo, _ := postgres.New(db)
		us.userRepository = pgRepo
		return nil
	}
}

func WithKafkaUserEventPublisher(producer *producer.KafkaProducer, repo outbox_repo.EventRepository) UserConfiguration {
	return func(us *userService) error {
		publisher := kafka.NewUserEventPublisher(producer, repo)
		us.userPublisher = publisher
		return nil
	}
}

// CreateUser
func (s *userService) CreateUser(ctx context.Context, email, password, firstName, lastName string) (*userAggre.User, error) {
	logger.Info("Start UserService.CreateUser")

	newUser, errs := userAggre.New(email, password, firstName, lastName, false)
	if errs.Exist() {
		logger.ErrorData("Finished UserService.CreateUser: FAILED", logging.Data{"error": errs.Error()})
		return nil, status.Errorf(codes.InvalidArgument, "%s", errs)
	}
	createdUser, err := s.userRepository.SaveUser(ctx, newUser)
	if err != nil {
		logger.ErrorData("Finished UserService.CreateUser: FAILED", logging.Data{"error": err.Error()})
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = s.userPublisher.Publish(ctx, models.UserEvent{
		ID:        createdUser.GetID(),
		Email:     createdUser.GetEmail(),
		CreatedAt: createdUser.GetCreatedAt(),
		Status:    models.USER_CREATED,
	})
	if err != nil {
		_, dbErr := s.userRepository.DeleteByUUID(ctx, createdUser.GetID())
		if dbErr != nil {
			logger.ErrorData("Finished UserService.CreateUser: FAILED", logging.Data{"error": err.Error()})
			return nil, status.Error(codes.Internal, err.Error())
		}
		logger.ErrorData("Finished UserService.CreateUser: FAILED", logging.Data{"error": err.Error()})
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Info("Finished UserService.CreateUser: SUCCESSFUL")
	return createdUser, nil
}

// Login
func (s *userService) Login(ctx context.Context, email, password string) (uuid.UUID, string, time.Time, string, time.Time, *userAggre.User, error) {
	logger.Info("Start UserService.Login")

	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		logger.ErrorData("Finished UserService.Login: FAILED", logging.Data{
			"error": err.Error(),
		})
		return uuid.UUID{}, "", time.Time{}, "", time.Time{}, nil, status.Error(codes.NotFound, err.Error())
	}

	if !utils.ComparePassword(password, user.GetPassword()) {
		logger.Error("Finished UserService.Login: FAILED: incorrect password")
		return uuid.UUID{}, "", time.Time{}, "", time.Time{}, nil, status.Error(codes.PermissionDenied, "incorrect password")
	}

	refreshToken, refreshClaim, err := s.tokenMaker.CreateToken(user.GetID(), uuid.New(), cmd.Conf.RefreshTokenDuration)
	if err != nil {
		logger.ErrorData("Finished UserService.Login: FAILED", logging.Data{
			"error": err.Error(),
		})
		return uuid.UUID{}, "", time.Time{}, "", time.Time{}, nil, status.Error(codes.Internal, err.Error())
	}
	p, _ := peer.FromContext(ctx)
	session := entities.Session{
		ID:           refreshClaim.ID,
		UserID:       refreshClaim.UserID,
		RefreshToken: refreshToken,
		ClientIP:     p.Addr.String(),
		IsDisabled:   false,
		ExpiresAt:    refreshClaim.ExpriresAt,
		CreatedAt:    refreshClaim.IssuedAt,
	}
	err = user.AddSession(session)
	if err != nil {
		logger.ErrorData("Finished UserService.Login: FAILED", logging.Data{
			"error": err.Error(),
		})
		return uuid.UUID{}, "", time.Time{}, "", time.Time{}, nil, status.Error(codes.Internal, err.Error())
	}
	updatedUser, err := s.userRepository.UpdateUser(ctx, user)
	if err != nil {
		logger.ErrorData("Finished UserService.Login: FAILED", logging.Data{
			"error": err.Error(),
		})
		return uuid.UUID{}, "", time.Time{}, "", time.Time{}, nil, status.Error(codes.Internal, err.Error())
	}

	addedSession, err := updatedUser.GetSessionById(session.ID)
	if err != nil || addedSession.RefreshToken != refreshToken {
		logger.ErrorData("Finished UserService.Login: FAILED", logging.Data{
			"error": "failed to add new session to the db",
		})
		return uuid.UUID{}, "", time.Time{}, "", time.Time{}, nil, status.Error(codes.Internal, "failed to add new session to the db")
	}

	accessToken, accessClaim, err := s.tokenMaker.CreateToken(user.GetID(), refreshClaim.ID, cmd.Conf.AccessTokenDuration)
	if err != nil {
		logger.ErrorData("Finished UserService.Login: FAILED", logging.Data{
			"error": err.Error(),
		})
		return uuid.UUID{}, "", time.Time{}, "", time.Time{}, nil, status.Error(codes.Internal, err.Error())
	}

	logger.Info("Finished UserService.Login: SUCCESSFUL")
	return addedSession.ID, addedSession.RefreshToken, addedSession.ExpiresAt, accessToken, accessClaim.ExpriresAt, updatedUser, nil
}

// VerifyToken
func (s *userService) VerifyToken(ctx context.Context, token string) (string, error) {
	logger.Info("Start UserService.VerifyToken")

	claims, err := s.tokenMaker.VerifyToken(token)
	if err != nil {
		logger.ErrorData("Finished UserService.VerifyToken: FAILED", logging.Data{
			"error": err.Error(),
		})
		if errors.Is(err, jwt.ExpriedTokenError) {
			return "", status.Error(codes.Unauthenticated, err.Error())
		}
		return "", status.Error(codes.Internal, err.Error())
	}

	userAg, err := s.userRepository.GetByUUID(ctx, claims.UserID)
	if err != nil {
		logger.ErrorData("Finished UserService.VerifyToken: FAILED", logging.Data{
			"error": err.Error(),
		})
		return "", status.Error(codes.NotFound, err.Error())
	}
	session, err := userAg.GetSessionById(claims.SessionID)
	if err != nil {
		logger.ErrorData("Finished UserService.VerifyToken: FAILED", logging.Data{
			"error": "session not found",
		})
		return "", status.Error(codes.NotFound, err.Error())
	}
	if session.IsDisabled {
		logger.ErrorData("Finished UserService.VerifyToken: FAILED", logging.Data{
			"error": "session is blocked",
		})
		return "", status.Error(codes.NotFound, "session is blocked")
	}

	logger.Info("Finished UserService.VerifyToken: SUCCESSFUL")
	return userAg.GetID().String(), nil
}

// RefreshToken
func (s *userService) RefreshToken(ctx context.Context, token string) (string, time.Time, error) {
	logger.Info("Start UserService.RefreshToken")

	claims, err := s.tokenMaker.VerifyToken(token)
	if err != nil {
		logger.ErrorData("Finished UserService.RefreshToken: FAILED", logging.Data{
			"error": err.Error(),
		})
		if errors.Is(err, jwt.ExpriedTokenError) {
			return "", time.Time{}, status.Error(codes.Unauthenticated, err.Error())
		}
		return "", time.Time{}, status.Error(codes.Internal, err.Error())
	}

	user, err := s.userRepository.GetByUUID(ctx, claims.UserID)
	if err != nil {
		logger.ErrorData("Finished UserService.RefreshToken: FAILED", logging.Data{
			"error": err.Error(),
		})
		return "", time.Time{}, status.Error(codes.Internal, err.Error())
	}

	session, err := user.GetSessionById(claims.ID)
	if err != nil {
		logger.ErrorData("Finished UserService.RefreshToken: FAILED", logging.Data{
			"error": "session doesn't exist",
		})
		return "", time.Time{}, status.Error(codes.NotFound, err.Error())
	}
	if session.IsDisabled {
		logger.ErrorData("Finished UserService.RefreshToken: FAILED", logging.Data{
			"error": "session is blocked",
		})
		return "", time.Time{}, status.Error(codes.PermissionDenied, "session is blocked")
	}
	if claims.UserID != session.UserID {
		logger.ErrorData("Finished UserService.RefreshToken: FAILED", logging.Data{
			"error": "incorrect session user",
		})
		return "", time.Time{}, status.Error(codes.PermissionDenied, "incorrect session user")
	}
	if session.HasExpired() {
		logger.ErrorData("Finished UserService.RefreshToken: FAILED", logging.Data{
			"error": "refresh token has expired",
		})
		return "", time.Time{}, status.Error(codes.Unauthenticated, "refresh token has expired")
	}

	accessToken, accessClaim, err := s.tokenMaker.CreateToken(user.GetID(), session.ID, cmd.Conf.AccessTokenDuration)
	if err != nil {
		logger.ErrorData("Finished UserService.RefreshToken: FAILED", logging.Data{
			"error": err.Error(),
		})
		return "", time.Time{}, status.Error(codes.Internal, err.Error())
	}

	return accessToken, accessClaim.ExpriresAt, nil
}

// RefreshToken
func (s *userService) GetUser(ctx context.Context, id uuid.UUID) (*userAggre.User, error) {
	logger.Info("Start UserService.GetMyInfo")

	user, err := s.userRepository.GetByUUID(ctx, id)
	if err != nil {
		logger.ErrorData("Finished UserService.GetMyInfo: FAILED", logging.Data{
			"error": err.Error(),
		})
		return nil, status.Error(codes.Internal, err.Error())
	}

	return user, nil
}
