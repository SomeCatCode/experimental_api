package model

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MeterValueIndex struct {
	Data          []MeterValue `json:"data"`
	CountTotal    int64        `json:"count_total"`
	CountFiltered int64        `json:"count_filtered"`
	Start         uint64       `json:"start"`
	Limit         uint64       `json:"limit"`
}

type MeterValue struct {
	ID        bson.ObjectID `json:"id" bson:"_id"`
	MeterId   string        `json:"meter_id" bson:"meter_id"`
	Reading   float64       `json:"reading" bson:"reading"`
	RawData   string        `json:"raw_data" bson:"raw_data"`
	CreatedAt *time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time    `json:"updated_at" bson:"updated_at"`
}

func (mv *MeterValue) Validate() error {
	if strings.TrimSpace(mv.MeterId) == "" {
		return errors.New("meter_id is required")
	}
	if mv.Reading < 0 {
		return errors.New("reading must be a positive number")
	}
	if mv.RawData == "" {
		return errors.New("raw_data is required")
	}
	if mv.CreatedAt == nil {
		return errors.New("created_at is required")
	}
	if mv.UpdatedAt == nil {
		return errors.New("updated_at is required")
	}
	if mv.CreatedAt.After(time.Now()) || mv.UpdatedAt.After(time.Now()) {
		return errors.New("created_at and updated_at cannot be in the future")
	}
	return nil
}
