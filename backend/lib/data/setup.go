package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func SetupDB() error {
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		return err
	}

	defer func() {
		_ = mongoClient.Disconnect(context.TODO())
	}()

	collection := mongoClient.Database("spotistats").Collection("users")

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{"username", 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{"state", 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	name, err := collection.Indexes().CreateMany(context.TODO(), indexes)
	if err != nil {
		return err
	} else {
		log.Printf("Created index: %s", name)
	}

	return nil
}
