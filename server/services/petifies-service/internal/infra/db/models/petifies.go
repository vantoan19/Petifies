package models

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	URI         string `bson:"url"`
	Description string `bson:"description"`
}

type Coordinates struct {
	Longitude float64 `bson:"longitude"`
	Latitude  float64 `bson:"latitude"`
}

type Address struct {
	AddressLineOne string      `bson:"address_line_one"`
	AddressLineTwo string      `bson:"address_line_two"`
	Street         string      `bson:"street"`
	District       string      `bson:"district"`
	City           string      `bson:"city"`
	Region         string      `bson:"region"`
	PostalCode     string      `bson:"postal_code"`
	Country        string      `bson:"country"`
	Coordinates    Coordinates `bson:"coordinates"`
}

type Petifies struct {
	ID          uuid.UUID `bson:"id"`
	OwnerID     uuid.UUID `bson:"owner_id"`
	Type        string    `bson:"type"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	Address     Address   `bson:"address"`
	PetName     string    `bson:"pet_name"`
	Images      []Image   `bson:"images"`
	Status      string    `bson:"status"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}
