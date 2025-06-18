package service

import (
	"context"
	"fmt"
	"time"

	"github.com/lan1hotspotgmbh/ms_meter/model"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MeterService struct {
	Collection *mongo.Collection
}

func NewMeterService(db *mongo.Database) *MeterService {
	return &MeterService{
		Collection: db.Collection("meters"),
	}
}

func (s *MeterService) Create(ctx context.Context, meter *model.Meter) error {
	now := time.Now()
	meter.ID = bson.NewObjectID()
	meter.UUID = uuid.New()
	meter.CreatedAt = &now
	meter.UpdatedAt = &now

	_, err := s.Collection.InsertOne(ctx, meter)
	return err
}

func (s *MeterService) GetAll(ctx context.Context) ([]model.Meter, error) {
	cursor, err := s.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var meters []model.Meter
	err = cursor.All(ctx, &meters)
	return meters, err
}

func (s *MeterService) GetAllWithFilters(ctx context.Context, filters map[string]string, search string, limit, start int) ([]model.Meter, error) {
	query := bson.M{}

	// Direktfilter
	for key, val := range filters {
		switch key {
		case "ID":
			if oid, err := bson.ObjectIDFromHex(val); err == nil {
				query["_id"] = oid
			}
		case "UUID":
			query["uuid"] = val
		case "Name":
			query["name"] = val
		case "Description":
			query["description"] = val
		case "SerialNumber":
			query["serial_number"] = val
		case "Model":
			query["model"] = val
		case "Manufacturer":
			query["manufacturer"] = val
		}
	}

	// Volltextsuche
	if search != "" {
		orConditions := bson.A{
			bson.M{"name": bson.M{"$regex": search, "$options": "i"}},
			bson.M{"description": bson.M{"$regex": search, "$options": "i"}},
			bson.M{"serial_number": bson.M{"$regex": search, "$options": "i"}},
			bson.M{"model": bson.M{"$regex": search, "$options": "i"}},
			bson.M{"manufacturer": bson.M{"$regex": search, "$options": "i"}},
		}
		query["$or"] = orConditions
	}

	findOptions := options.Find().SetLimit(int64(limit)).SetSkip(int64(start))
	cursor, err := s.Collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var meters []model.Meter
	err = cursor.All(ctx, &meters)
	return meters, err
}

func (s *MeterService) GetByID(ctx context.Context, id string) (*model.Meter, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	var obj model.Meter
	err = s.Collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&obj)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &obj, nil
}

func (s *MeterService) UpdateByID(ctx context.Context, id string, meter *model.Meter) (*model.Meter, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	now := time.Now()
	meter.UpdatedAt = &now

	update := bson.M{
		"$set": bson.M{
			"name":            meter.Name,
			"description":     meter.Description,
			"serial_number":   meter.SerialNumber,
			"manufacturer":    meter.Manufacturer,
			"model":           meter.Model,
			"public_key":      meter.PublicKey,
			"organisation_id": meter.OrganisationId,
			"location_id":     meter.LocationId,
			"updated_at":      meter.UpdatedAt,
		},
	}

	_, err = s.Collection.UpdateByID(ctx, oid, update)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

func (s *MeterService) DeleteByID(ctx context.Context, id string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	_, err = s.Collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: oid}})
	return err
}
