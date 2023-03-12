package neo4j

import (
	"context"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vantoan19/Petifies/server/libs/logging-config"
	useraggre "github.com/vantoan19/Petifies/server/services/relationship-service/internal/domain/aggregates/user"
	"github.com/vantoan19/Petifies/server/services/relationship-service/internal/domain/aggregates/user/entities"
)

var logger = logging.New("RelationshipService.UserRepository")

type UserRepository struct {
	db neo4j.Driver
}

func NewNeo4jUserRepository(db neo4j.Driver) (*UserRepository, error) {
	return &UserRepository{
		db: db,
	}, nil
}

// Get retrieves a UserAggregate from Neo4j by ID
func (ur *UserRepository) GetByUUID(_ context.Context, id uuid.UUID) (*useraggre.UserAggregate, error) {
	logger.Info("Start GetByUUID")

	session := ur.db.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Get user
		result, err := tx.Run("MATCH (u:User {id: $id}) RETURN u.email", map[string]interface{}{
			"id": id.String(),
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		if !result.Next() {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		email, _ := result.Record().Get("u.email")
		user, err := useraggre.NewUserAggregate(entities.User{
			ID:    id,
			Email: email.(string),
		})
		if err != nil {
			return nil, err
		}

		// Get Following
		result, err = tx.Run("MATCH (u:User {id: $id})-[r:FOLLOW]->(v:User) RETURN v.id", map[string]interface{}{
			"id": id.String(),
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		for result.Next() {
			userID_, _ := result.Record().Get("v.id")
			userID, err := uuid.Parse(userID_.(string))
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, err.Error())
			}

			err = user.Follow(userID)
			if err != nil {
				return nil, err
			}
		}
		if result.Err() != nil {
			return nil, status.Errorf(codes.Internal, result.Err().Error())
		}

		// Get Followers
		result, err = tx.Run("MATCH (u:User {id: $id})<-[r:FOLLOW]-(v:User) RETURN v.id", map[string]interface{}{
			"id": id.String(),
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		for result.Next() {
			userID_, _ := result.Record().Get("v.id")
			userID, err := uuid.Parse(userID_.(string))
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, err.Error())
			}
			user.AddFollower(userID)
		}
		if result.Err() != nil {
			return nil, status.Errorf(codes.Internal, result.Err().Error())
		}

		return user, nil
	})
	if err != nil {
		logger.ErrorData("Finished GetByUUID: FAILED", logging.Data{"error": err.Error()})
		return nil, err
	}

	logger.Info("Finished GetByUUID: SUCCESSFUL")
	return result.(*useraggre.UserAggregate), nil
}

// Save inserts a new UserAggregate into Neo4j
func (ur *UserRepository) SaveUser(_ context.Context, user *useraggre.UserAggregate) error {
	logger.Info("Start SaveUser")

	session := ur.db.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	// Create the user
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("CREATE (u:User {id: $id, email: $email})", map[string]interface{}{
			"id":    user.GetID().String(),
			"email": user.GetEmail(),
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		return result, nil
	})
	if err != nil {
		logger.ErrorData("Finished SaveUser: FAILED", logging.Data{"error": err.Error()})
		return err
	}

	// Create following relationships
	for _, followingID := range user.GetFollowings() {
		_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			result, err := tx.Run("MATCH (u1:User {id: $user_id}), (u2:User {id: $following_id}) CREATE (u1)-[r:FOLLOW {from_user_id: $user_id, to_user_id: $following_id}]->(u2)", map[string]interface{}{
				"user_id":      user.GetID().String(),
				"following_id": followingID.String(),
			})
			if err != nil {
				return nil, status.Errorf(codes.Internal, err.Error())
			}
			return result, nil
		})
		if err != nil {
			logger.ErrorData("Finished SaveUser: FAILED", logging.Data{"error": err.Error()})
			return err
		}
	}

	logger.Info("Finish SaveUser: SUCCESSFUL")
	return nil
}

// Update updates an existing UserAggregate in Neo4j
func (ur *UserRepository) UpdateUser(_ context.Context, user *useraggre.UserAggregate) error {
	logger.Info("Start UpdateUser")
	session := ur.db.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (u:User {id: $id}) SET u.email = $email", map[string]interface{}{
			"id":    user.GetID().String(),
			"email": user.GetEmail(),
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		return result, nil
	})
	if err != nil {
		logger.ErrorData("Finished UpdateUser: FAILED", logging.Data{"error": err.Error()})
		return err
	}

	// Delete existing following relationships
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (u:User {id: $id})-[r]->() DELETE r", map[string]interface{}{
			"id": user.GetID().String(),
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		return result, nil
	})
	if err != nil {
		logger.ErrorData("Finished UpdateUser: FAILED", logging.Data{"error": err.Error()})
		return err
	}

	// Create following relationships
	for _, followingID := range user.GetFollowings() {
		_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			result, err := tx.Run("MATCH (u1:User {id: $user_id}), (u2:User {id: $following_id}) CREATE (u1)-[r:FOLLOW {from_user_id: $user_id, to_user_id: $following_id}]->(u2)", map[string]interface{}{
				"user_id":      user.GetID().String(),
				"following_id": followingID.String(),
			})
			if err != nil {
				return nil, status.Errorf(codes.Internal, err.Error())
			}
			return result, nil
		})
		if err != nil {
			logger.ErrorData("Finished UpdateUser: FAILED", logging.Data{"error": err.Error()})
			return err
		}
	}

	logger.Info("Finish UpdateUser: SUCCESSFUL")
	return nil
}

// Delete deletes an existing UserAggregate from Neo4j by ID
func (ur *UserRepository) Delete(_ context.Context, id uuid.UUID) error {
	logger.Info("Start Delete")

	session := ur.db.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run("MATCH (u:User {id: $id}) DETACH DELETE u", map[string]interface{}{
			"id": id.String(),
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		return result, nil
	})
	if err != nil {
		logger.ErrorData("Finished Delete: FAILED", logging.Data{"error": err.Error()})
		return err
	}

	logger.Info("Finish Delete: SUCCESSFUL")
	return err
}
