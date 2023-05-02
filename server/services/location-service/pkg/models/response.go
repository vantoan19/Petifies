package models

import "github.com/google/uuid"

type Location struct {
	ID           uuid.UUID
	EntityID     uuid.UUID
	LocationType string
}

type ListNearByLocationsByTypeResp struct {
	Locations []*Location
}
