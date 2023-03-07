package models

import "github.com/google/uuid"

type AddRelationshipResp struct {
	Message string
}

type RemoveRelationshipResp struct {
	Message string
}

type ListFollowersResp struct {
	FollowerIDs []uuid.UUID
}

type ListFollowingsResp struct {
	FollowingIDs []uuid.UUID
}
