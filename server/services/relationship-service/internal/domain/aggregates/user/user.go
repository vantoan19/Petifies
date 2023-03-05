package useraggre

import (
	"errors"

	"github.com/google/uuid"
	"github.com/vantoan19/Petifies/server/services/relationship-service/internal/domain/aggregates/user/entities"
	"github.com/vantoan19/Petifies/server/services/relationship-service/internal/domain/aggregates/user/valueobjects"
)

var (
	ErrRelationshipAlreadyExists = errors.New("relationship already exist")
	ErrRelationshipNotExist      = errors.New("relationship does not exist")
	ErrExceedRelationshipLimit   = errors.New("exceeds relationship limit")
)

// UserAggregate represents the aggregate for a user and their relationships
type UserAggregate struct {
	user          *entities.User
	relationships map[valueobjects.RelationshipType][]*entities.Relationship
}

func NewUserAggregate(user entities.User) (*UserAggregate, error) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	return &UserAggregate{
		user:          &user,
		relationships: make(map[valueobjects.RelationshipType][]*entities.Relationship),
	}, nil
}

func (u *UserAggregate) AddRelationship(relationship entities.Relationship) (entities.Relationship, error) {
	if len(u.relationships[relationship.Type]) > 2000 {
		return entities.Relationship{}, ErrExceedRelationshipLimit
	}

	for _, r := range u.relationships[relationship.Type] {
		if r.ToUserID == relationship.ToUserID {
			return entities.Relationship{}, ErrRelationshipAlreadyExists
		}
	}

	u.relationships[relationship.Type] = append(u.relationships[relationship.Type], &relationship)
	return relationship, nil
}

func (u *UserAggregate) GetRelationshipsByType(relationshipType valueobjects.RelationshipType) []entities.Relationship {
	res := make([]entities.Relationship, 0)
	for _, r := range u.relationships[relationshipType] {
		res = append(res, *r)
	}
	return res
}

func (u *UserAggregate) DeleteRelationship(toUserID uuid.UUID, relationshipType valueobjects.RelationshipType) error {
	for i, r := range u.relationships[relationshipType] {
		if r.ToUserID == toUserID {
			u.relationships[relationshipType] = append(u.relationships[relationshipType][:i], u.relationships[relationshipType][i+1:]...)
			return nil
		}
	}
	return ErrRelationshipNotExist
}

func (u *UserAggregate) GetRelationships() map[valueobjects.RelationshipType][]entities.Relationship {
	res := make(map[valueobjects.RelationshipType][]entities.Relationship)
	for t, rs := range u.relationships {
		rs_ := make([]entities.Relationship, 0)
		for _, r := range rs {
			rs_ = append(rs_, *r)
		}
		res[t] = rs_
	}
	return res
}

// ========= Aggregate Root Getter =========

func (u *UserAggregate) GetID() uuid.UUID {
	return u.user.ID
}

func (u *UserAggregate) GetEmail() string {
	return u.user.Email
}
