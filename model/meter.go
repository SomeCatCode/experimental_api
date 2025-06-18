package model

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MeterIndex struct {
	Data          []Meter `json:"data"`
	CountTotal    int64   `json:"count_total"`
	CountFiltered int64   `json:"count_filtered"`
	Start         uint64  `json:"start"`
	Limit         uint64  `json:"limit"`
}

type Meter struct {
	ID             bson.ObjectID `json:"id" bson:"_id"`
	UUID           uuid.UUID     `json:"uuid" bson:"uuid"`
	OrganisationId string        `json:"organisation_id" bson:"organisation_id"`
	LocationId     string        `json:"location_id" bson:"location_id"`
	Name           string        `json:"name" bson:"name"`
	Description    string        `json:"description" bson:"description"`
	SerialNumber   string        `json:"serial_number" bson:"serial_number"`
	Manufacturer   string        `json:"manufacturer" bson:"manufacturer"`
	Model          string        `json:"model" bson:"model"`
	PublicKey      string        `json:"public_key" bson:"public_key"`
	CreatedAt      *time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt      *time.Time    `json:"updated_at" bson:"updated_at"`
}

func (m *Meter) Validate() error {
	if strings.TrimSpace(m.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(m.SerialNumber) == "" {
		return errors.New("serial_number is required")
	}
	return nil
}
