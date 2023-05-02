package cassandra

import (
	"context"
	"errors"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vantoan19/Petifies/server/libs/logging-config"
	useraggre "github.com/vantoan19/Petifies/server/services/newfeed-service/internal/domain/aggregates/user"
	"github.com/vantoan19/Petifies/server/services/newfeed-service/internal/infra/repositories/db/mapper"
	"github.com/vantoan19/Petifies/server/services/newfeed-service/internal/infra/repositories/db/models"
)

var logger = logging.New("NewFeedService.UserRepository")

type UserRepository struct {
	session *gocql.Session
}

func NewCassandraUserRepository(session *gocql.Session) (*UserRepository, error) {
	return &UserRepository{
		session: session,
	}, nil
}

func (ur *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*useraggre.UserAggre, error) {
	logger.Info("Start GetByID")

	var user models.User
	query := `SELECT id, email FROM users WHERE id=?`
	if err := ur.session.Query(query, id.String()).Scan(&user.ID, &user.Email); err != nil {
		logger.ErrorData("Finish GetByID: FAILED", logging.Data{"error": err.Error()})
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	logger.Info("Finish GetByID: SUCCESSFUL")
	return mapper.DbUserToUserAggregate(&user)
}

func (ur *UserRepository) Save(ctx context.Context, user useraggre.UserAggre) (*useraggre.UserAggre, error) {
	logger.Info("Start Save")

	userModel, err := mapper.UserAggregateToUserDb(&user)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO users(id, email) VALUES(?,?)`
	if err := ur.session.Query(query, userModel.ID, userModel.Email).Exec(); err != nil {
		logger.ErrorData("Finish Save: FAILED", logging.Data{"error": err.Error()})
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	logger.Info("Finish Save: SUCCESSFUL")
	return mapper.DbUserToUserAggregate(userModel)
}

func (ur *UserRepository) Update(ctx context.Context, user useraggre.UserAggre) (*useraggre.UserAggre, error) {
	return nil, errors.New("Not implemented")
}

func (ur *UserRepository) DeleteByID(ctx context.Context, id uuid.UUID) (*useraggre.UserAggre, error) {
	logger.Info("Start DeleteByID")

	var user models.User
	query := `DELETE FROM users WHERE id=? IF EXISTS`
	if err := ur.session.Query(query, id.String()).Scan(&user.ID, &user.Email); err != nil {
		logger.ErrorData("Finish DeleteByID: FAILED", logging.Data{"error": err.Error()})
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	logger.Info("Finish DeleteByID: SUCCESSFUL")
	return mapper.DbUserToUserAggregate(&user)
}
