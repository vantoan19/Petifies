package neo4j

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	useraggre "github.com/vantoan19/Petifies/server/services/relationship-service/internal/domain/aggregates/user"
	"github.com/vantoan19/Petifies/server/services/relationship-service/internal/domain/aggregates/user/entities"
	"github.com/vantoan19/Petifies/server/services/relationship-service/internal/domain/aggregates/user/valueobjects"
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
		result, err := tx.Run("MATCH (u:User {id: $id}) RETURN u", map[string]interface{}{
			"id": id.String(),
		})
		if err != nil {
			return nil, err
		}

		if !result.Next() {
			return nil, fmt.Errorf("user not found")
		}
		email, ok := result.Record().Get("u.email")
		if !ok {
			return nil, fmt.Errorf("cannot retrieve email")
		}

		user, err := useraggre.NewUserAggregate(entities.User{
			ID:    id,
			Email: email.(string),
		})

		result, err = tx.Run("MATCH (u:User {id: $id})-[r]->(v:User) RETURN r.id, r.from_user_id, r.to_user_id, r.type, v.id, v.name", map[string]interface{}{
			"id": id.String(),
		})
		if err != nil {
			return nil, err
		}

		for result.Next() {
			relationship, err := parseRelationship(result.Record())
			if err != nil {
				return nil, err
			}

			_, err = user.AddRelationship(*relationship)
			if err != nil {
				return nil, err
			}
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
func (ur *UserRepository) SaveUser(user *useraggre.UserAggregate) error {
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

	// Create relationship edges
	for _, relationships := range user.GetRelationships() {
		for _, relationship := range relationships {
			_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
				_, err := tx.Run("MATCH (u1:User {id: $from_user_id}), (u2:User {id: $to_user_id}) CREATE (u1)-[r:"+string(relationship.Type)+" {id: $id, type: $type, from_user_id: $from_user_id, to_user_id: $to_user_id}]->(u2)", map[string]interface{}{
					"from_user_id": relationship.FromUserID,
					"to_user_id":   relationship.ToUserID,
					"id":           relationship.ID,
					"type":         string(relationship.Type),
				})
				return nil, err
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Update updates an existing UserAggregate in Neo4j
func (ur *UserRepository) UpdateUser(user *useraggre.UserAggregate) error {
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

	// Delete existing relationships
	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run("MATCH (u:User {id: $id})-[r]->() DELETE r", map[string]interface{}{
			"id": user.GetID().String(),
		})
		return nil, err
	})
	if err != nil {
		return err
	}

	// Create relationship edges
	for _, relationships := range user.GetRelationships() {
		for _, relationship := range relationships {
			_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
				_, err := tx.Run("MATCH (u1:User {id: $from_user_id}), (u2:User {id: $to_user_id}) CREATE (u1)-[r:"+string(relationship.Type)+" {id: $id, type: $type, from_user_id: $from_user_id, to_user_id: $to_user_id}]->(u2)", map[string]interface{}{
					"from_user_id": relationship.FromUserID,
					"to_user_id":   relationship.ToUserID,
					"id":           relationship.ID,
					"type":         string(relationship.Type),
				})
				return nil, err
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Delete deletes an existing UserAggregate from Neo4j by ID
func (ur *UserRepository) Delete(id uuid.UUID) error {
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

func parseRelationship(record *neo4j.Record) (*entities.Relationship, error) {
	id, ok := record.Get("r.id")
	if !ok {
		return nil, fmt.Errorf("invalid relationship ID")
	}
	fromUserID, ok := record.Get("r.from_user_id")
	if !ok {
		return nil, fmt.Errorf("invalid relationship from user ID")
	}
	toUserID, ok := record.Get("r.to_user_id")
	if !ok {
		return nil, fmt.Errorf("invalid relationship to user ID")
	}
	relationshipType, ok := record.Get("r.type")
	if !ok {
		return nil, fmt.Errorf("invalid relationship type")
	}

	id_, err := uuid.Parse(id.(string))
	if err != nil {
		return nil, err
	}
	fromUserID_, err := uuid.Parse(fromUserID.(string))
	if err != nil {
		return nil, err
	}
	toUserID_, err := uuid.Parse(toUserID.(string))
	if err != nil {
		return nil, err
	}

	return &entities.Relationship{
		ID:         id_,
		FromUserID: fromUserID_,
		ToUserID:   toUserID_,
		Type:       valueobjects.RelationshipType(relationshipType.(string)),
	}, nil
}
