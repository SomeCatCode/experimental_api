package organisation

import (
	"context"

	"github.com/SomeCatCode/experimental_api/model"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoRepository struct {
	Database *mongo.Database
}

func (repo *MongoRepository) Insert(ctx context.Context, organisation model.Organisation) error {
	collection := repo.Database.Collection("organisations")
	_, err := collection.InsertOne(ctx, organisation)
	if err != nil {
		return err
	}

	return nil
}

func (repo *MongoRepository) FindByID(ctx context.Context, id string) (*model.Organisation, error) {
	collection := repo.Database.Collection("organisations")
	var organisation model.Organisation
	err := collection.FindOne(ctx, map[string]string{"_id": id}).Decode(&organisation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No document found
		}
		return nil, err // Other error
	}

	return &organisation, nil
}
