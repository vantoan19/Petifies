package models

import "github.com/google/uuid"

type AddRelationshipReq struct {
	FromUserID       uuid.UUID
	ToUserID         uuid.UUID
	RelationshipType string
}

type RemoveRelationshipReq struct {
	FromUserID       uuid.UUID
	ToUserID         uuid.UUID
	RelationshipType string
}

type ListFollowersReq struct {
	UserID uuid.UUID
}

type ListFollowingsReq struct {
	UserID uuid.UUID
}
