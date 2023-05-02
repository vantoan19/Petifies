package models

import "github.com/google/uuid"

type AddRelationshipResp struct {
	Message string
}

type RemoveRelationshipResp struct {
	Message string
}

type ListFollowersResp struct {
	FollowerIDs []uuid.UUID `json:"follower_ids"`
}

type ListFollowingsResp struct {
	FollowingIDs []uuid.UUID `json:"following_ids"`
}
