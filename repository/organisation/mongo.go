package organisation

import (
	"context"
	"fmt"

	"github.com/SomeCatCode/experimental_api/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoRepository struct {
	Collection string
	Database   *mongo.Database
}

func (repo *MongoRepository) ToObjectId(id string) (bson.ObjectID, error) {
	return bson.ObjectIDFromHex(id)
}

func (repo *MongoRepository) Insert(ctx context.Context, organisation model.Organisation) error {
	collection := repo.Database.Collection(repo.Collection)
	_, err := collection.InsertOne(ctx, organisation)
	if err != nil {
		return err
	}

	return nil
}

func (repo *MongoRepository) FindByID(ctx context.Context, id string) (*model.Organisation, error) {
	// String to ObjectID conversion
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	// Create Collection
	collection := repo.Database.Collection(repo.Collection)

	// Create filter for the query
	var obj model.Organisation
	filter := bson.D{{Key: "_id", Value: oid}}
	err = collection.FindOne(ctx, filter).Decode(&obj)

	// Handle errors
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No document found
		}
		return nil, err // Other error
	}
	return &obj, nil
}

func (repo *MongoRepository) UpdateByID(ctx context.Context, id string, organisation model.Organisation) error {
	collection := repo.Database.Collection(repo.Collection)
	_, err := collection.UpdateOne(ctx, map[string]string{"_id": id}, map[string]interface{}{
		"$set": organisation,
	})
	if err != nil {
		return err
	}

	return nil
}

func (repo *MongoRepository) DeleteByID(ctx context.Context, id string) error {
	collection := repo.Database.Collection(repo.Collection)
	_, err := collection.DeleteOne(ctx, map[string]string{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

func (repo *MongoRepository) List(ctx context.Context) ([]model.Organisation, error) {
	collection := repo.Database.Collection(repo.Collection)
	cursor, err := collection.Find(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var organisations []model.Organisation
	for cursor.Next(ctx) {
		var organisation model.Organisation
		if err := cursor.Decode(&organisation); err != nil {
			return nil, err
		}
		organisations = append(organisations, organisation)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return organisations, nil
}
