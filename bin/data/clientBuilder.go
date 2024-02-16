package data

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

func CreateClient() (*mongo.Client, error) {
	mongoClient, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(os.Getenv("MONGODB_URI")),
	)
	if err != nil {
		return nil, err
	}

	return mongoClient, nil
}
