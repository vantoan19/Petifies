package relationshipservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"github.com/vantoan19/Petifies/server/libs/logging-config"
	useraggre "github.com/vantoan19/Petifies/server/services/relationship-service/internal/domain/aggregates/user"
	neo4jRepo "github.com/vantoan19/Petifies/server/services/relationship-service/internal/infra/repositories/user/neo4j"
	"github.com/vantoan19/Petifies/server/services/relationship-service/pkg/models"
)

var logger = logging.New("RelationshipService.RelationshipSvc")

type RelationshipServiceConfiguration func(rs *relationshipService) error

type relationshipService struct {
	userRepo *neo4jRepo.UserRepository
}

type RelationshipService interface {
	AddRelationship(ctx context.Context, req *models.AddRelationshipReq) (*useraggre.UserAggregate, error)
	RemoveRelationship(ctx context.Context, req *models.RemoveRelationshipReq) (*useraggre.UserAggregate, error)
	ListFollowers(ctx context.Context, req *models.ListFollowersReq) ([]uuid.UUID, error)
	ListFollowings(ctx context.Context, req *models.ListFollowingsReq) ([]uuid.UUID, error)
}

func NewRelationshipService(cfgs ...RelationshipServiceConfiguration) (RelationshipService, error) {
	rs := &relationshipService{}
	for _, cfg := range cfgs {
		if err := cfg(rs); err != nil {
			return nil, err
		}
	}
	return rs, nil
}

func WithNeo4jUserRepository(db neo4j.Driver) RelationshipServiceConfiguration {
	return func(rs *relationshipService) error {
		repo, err := neo4jRepo.NewNeo4jUserRepository(db)
		if err != nil {
			return err
		}
		rs.userRepo = repo
		return nil
	}
}

func (rs *relationshipService) AddRelationship(ctx context.Context, req *models.AddRelationshipReq) (*useraggre.UserAggregate, error) {
	logger.InfoData("Start AddRelationship", logging.Data{
		"from": req.FromUserID.String(),
		"to":   req.ToUserID.String(),
		"type": req.RelationshipType,
	})

	fromUser, err := rs.userRepo.GetByUUID(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	toUser, err := rs.userRepo.GetByUUID(ctx, req.ToUserID)
	if err != nil {
		return nil, err
	}

	if req.RelationshipType == "FOLLOW" {
		err = fromUser.Follow(toUser.GetID())
		if err != nil {
			return nil, err
		}
	}

	err = rs.userRepo.UpdateUser(ctx, fromUser)
	if err != nil {
		return nil, err
	}

	logger.InfoData("Finish AddRelationship: Success", logging.Data{
		"from": req.FromUserID.String(),
		"to":   req.ToUserID.String(),
		"type": req.RelationshipType,
	})
	return fromUser, nil
}

func (rs *relationshipService) RemoveRelationship(ctx context.Context, req *models.RemoveRelationshipReq) (*useraggre.UserAggregate, error) {
	logger.InfoData("Start RemoveRelationship", logging.Data{
		"from": req.FromUserID.String(),
		"to":   req.ToUserID.String(),
		"type": req.RelationshipType,
	})

	fromUser, err := rs.userRepo.GetByUUID(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	toUser, err := rs.userRepo.GetByUUID(ctx, req.ToUserID)
	if err != nil {
		return nil, err
	}
	if req.RelationshipType == "FOLLOW" {
		err = fromUser.Unfollow(toUser.GetID())
		if err != nil {
			return nil, err
		}
	}

	err = rs.userRepo.UpdateUser(ctx, fromUser)
	if err != nil {
		return nil, err
	}

	logger.InfoData("Finish RemoveRelationship: Success", logging.Data{
		"from": req.FromUserID.String(),
		"to":   req.ToUserID.String(),
		"type": req.RelationshipType,
	})
	return fromUser, nil
}

func (rs *relationshipService) ListFollowers(ctx context.Context, req *models.ListFollowersReq) ([]uuid.UUID, error) {
	logger.InfoData("Start ListFollowers", logging.Data{"user": req.UserID.String()})

	user, err := rs.userRepo.GetByUUID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	logger.InfoData("Finish ListFollowers: Failed", logging.Data{"user": req.UserID.String()})
	return append(user.GetFollowers(), user.GetID()), nil
}

func (rs *relationshipService) ListFollowings(ctx context.Context, req *models.ListFollowingsReq) ([]uuid.UUID, error) {
	logger.InfoData("Start ListFollowings", logging.Data{"user": req.UserID.String()})

	user, err := rs.userRepo.GetByUUID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	logger.InfoData("Finish ListFollowings: Failed", logging.Data{"user": req.UserID.String()})
	return append(user.GetFollowings(), user.GetID()), nil
}
