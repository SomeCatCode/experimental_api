package model

import (
	"time"

	"github.com/google/uuid"
)

type Organisation struct {
	ID   string    `json:"id" bson:"_id"`
	UUID uuid.UUID `json:"uid" bson:"uid"`

	ParentId string `json:"parent_id" bson:"parent_id"`

	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`

	Homepage string `json:"website_url" bson:"website_url"`
	Email    string `json:"contact_email" bson:"contact_email"`
	Phone    string `json:"contact_phone" bson:"contact_phone"`
	Address  string `json:"address" bson:"address"`

	Latetude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`

	Style string `json:"style" bson:"style"`

	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
}
