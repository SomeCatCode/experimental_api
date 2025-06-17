package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MeterValue struct {
	ID      bson.ObjectID `json:"id" bson:"_id"`
	MeterId string        `json:"meter_id" bson:"meter_id"`

	Reading float64 `json:"reading" bson:"reading"`
	RawData string  `json:"raw_data" bson:"raw_data"`

	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
}
