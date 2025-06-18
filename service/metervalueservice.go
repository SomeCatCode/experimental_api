package service

import (
	"context"
	"fmt"
	"time"

	"github.com/lan1hotspotgmbh/ms_meter/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MeterValueService struct {
	Collection *mongo.Collection
}

func NewMeterValueService(db *mongo.Database) *MeterValueService {
	return &MeterValueService{
		Collection: db.Collection("meter_values"),
	}
}

func (s *MeterValueService) Create(ctx context.Context, meter *model.MeterValue) error {
	now := time.Now()
	meter.ID = bson.NewObjectID()
	meter.CreatedAt = &now
	meter.UpdatedAt = &now

	_, err := s.Collection.InsertOne(ctx, meter)
	return err
}

func (s *MeterValueService) GetAll(ctx context.Context) ([]model.MeterValue, error) {
	cursor, err := s.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var values []model.MeterValue
	err = cursor.All(ctx, &values)
	return values, err
}

func (s *MeterValueService) GetAllWithFilters(ctx context.Context, filters map[string]string, search string, limit, start int) ([]model.MeterValue, error) {
	query := bson.M{}

	for key, val := range filters {
		switch key {
		case "MeterId":
			query["meter_id"] = val
		case "RawData":
			query["raw_data"] = val
		}
	}

	if search != "" {
		orConditions := bson.A{
			bson.M{"meter_id": bson.M{"$regex": search, "$options": "i"}},
			bson.M{"raw_data": bson.M{"$regex": search, "$options": "i"}},
		}
		query["$or"] = orConditions
	}

	findOptions := options.Find().SetLimit(int64(limit)).SetSkip(int64(start))
	cursor, err := s.Collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var values []model.MeterValue
	err = cursor.All(ctx, &values)
	return values, err
}

func (s *MeterValueService) GetByID(ctx context.Context, id string) (*model.MeterValue, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	var obj model.MeterValue
	err = s.Collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&obj)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &obj, nil
}

func (s *MeterValueService) UpdateByID(ctx context.Context, id string, meter *model.MeterValue) (*model.MeterValue, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	now := time.Now()
	meter.UpdatedAt = &now

	update := bson.M{
		"$set": bson.M{
			"meter_id":   meter.MeterId,
			"reading":    meter.Reading,
			"raw_data":   meter.RawData,
			"updated_at": meter.UpdatedAt,
		},
	}

	_, err = s.Collection.UpdateByID(ctx, oid, update)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

func (s *MeterValueService) DeleteByID(ctx context.Context, id string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	_, err = s.Collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: oid}})
	return err
}
