package neo4j

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	useraggre "github.com/vantoan19/Petifies/server/services/relationship-service/internal/domain/aggregates/user"
	"github.com/vantoan19/Petifies/server/services/relationship-service/internal/domain/aggregates/user/entities"
)

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
	session := ur.db.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Get user
		result, err := tx.Run("MATCH (u:User {id: $id}) RETURN u.email", map[string]interface{}{
			"id": id.String(),
		})
		if err != nil {
			return nil, err
		}
		if !result.Next() {
			return nil, fmt.Errorf("user not found")
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
			return nil, err
		}
		for result.Next() {
			userID_, _ := result.Record().Get("v.id")
			userID, err := uuid.Parse(userID_.(string))
			if err != nil {
				return nil, err
			}

			err = user.Follow(userID)
			if err != nil {
				return nil, err
			}
		}
		if result.Err() != nil {
			return nil, result.Err()
		}

		// Get Followers
		result, err = tx.Run("MATCH (u:User {id: $id})<-[r:FOLLOW]-(v:User) RETURN v.id", map[string]interface{}{
			"id": id.String(),
		})
		if err != nil {
			return nil, err
		}
		for result.Next() {
			userID_, _ := result.Record().Get("v.id")
			userID, err := uuid.Parse(userID_.(string))
			if err != nil {
				return nil, err
			}
			user.AddFollower(userID)
		}
		if result.Err() != nil {
			return nil, result.Err()
		}

		return user, nil
	})
	if err != nil {
		return nil, err
	}

	return result.(*useraggre.UserAggregate), nil
}

// Save inserts a new UserAggregate into Neo4j
func (ur *UserRepository) SaveUser(_ context.Context, user *useraggre.UserAggregate) error {
	session := ur.db.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	// Create the user
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run("CREATE (u:User {id: $id, email: $email})", map[string]interface{}{
			"id":    user.GetID().String(),
			"email": user.GetEmail(),
		})
		return nil, err
	})
	if err != nil {
		return err
	}

	// Create following relationships
	for _, followingID := range user.GetFollowings() {
		_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			_, err := tx.Run("MATCH (u1:User {id: $user_id}), (u2:User {id: $following_id}) CREATE (u1)-[r:FOLLOW {from_user_id: $user_id, to_user_id: $following_id}]->(u2)", map[string]interface{}{
				"user_id":      user.GetID().String(),
				"following_id": followingID.String(),
			})
			return nil, err
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Update updates an existing UserAggregate in Neo4j
func (ur *UserRepository) UpdateUser(_ context.Context, user *useraggre.UserAggregate) error {
	session := ur.db.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run("MATCH (u:User {id: $id}) SET u.email = $email", map[string]interface{}{
			"id":    user.GetID().String(),
			"email": user.GetEmail(),
		})
		return nil, err
	})
	if err != nil {
		return err
	}

	// Delete existing following relationships
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run("MATCH (u:User {id: $id})-[r]->() DELETE r", map[string]interface{}{
			"id": user.GetID().String(),
		})
		return nil, err
	})
	if err != nil {
		return err
	}

	// Create following relationships
	for _, followingID := range user.GetFollowings() {
		_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			_, err := tx.Run("MATCH (u1:User {id: $user_id}), (u2:User {id: $following_id}) CREATE (u1)-[r:FOLLOW {from_user_id: $user_id, to_user_id: $following_id}]->(u2)", map[string]interface{}{
				"user_id":      user.GetID().String(),
				"following_id": followingID.String(),
			})
			return nil, err
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete deletes an existing UserAggregate from Neo4j by ID
func (ur *UserRepository) Delete(_ context.Context, id uuid.UUID) error {
	session := ur.db.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run("MATCH (u:User {id: $id}) DETACH DELETE u", map[string]interface{}{
			"id": id.String(),
		})
		return nil, err
	})

	return err
}
