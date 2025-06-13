package model

import (
	"time"
)

type Organisation struct {
	ID             string `json:"id" bson:"_id"`
	OrganisationId string `json:"organisation_id" bson:"organisation_id"`

	Name        string  `json:"name" bson:"name"`
	Description string  `json:"description" bson:"description"`
	Address     Address `json:"address" bson:"address"`

	Phone    string `json:"phone" bson:"phone"`
	Email    string `json:"email" bson:"email"`
	Homepage string `json:"homepage" bson:"homepage"`

	Latetude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`

	Style string `json:"style" bson:"style"`

	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
}

type Address struct {
	Street  string `json:"street" bson:"street"`
	City    string `json:"city" bson:"city"`
	State   string `json:"state" bson:"state"`
	Country string `json:"country" bson:"country"`
	Postal  string `json:"postal" bson:"postal"`
}
