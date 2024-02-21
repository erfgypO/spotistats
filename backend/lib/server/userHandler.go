package server

import (
	"context"
	"github.com/erfgypO/spotistats/lib/data"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func HandleGetMe(c echo.Context) error {
	claims, err := getClaimsFromContext(c)
	if err != nil {
		return err
	}

	mongoClient, err := data.CreateClient()
	if err != nil {
		return err
	}

	defer func() {
		_ = mongoClient.Disconnect(context.TODO())
	}()

	var user data.UserEntity
	collection := mongoClient.Database("spotistats").Collection(data.UserCollectionName)
	id, err := primitive.ObjectIDFromHex(claims["uid"].(string))
	if err != nil {
		return err
	}
	err = collection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&user)
	if err != nil {
		return err
	}

	datapointCount, err := collection.Database().Collection(data.DatapointCollectionName).CountDocuments(context.TODO(), bson.D{{"owner", id}})

	userResponse := UserResponse{
		Id:                 user.Id.(primitive.ObjectID).Hex(),
		Username:           user.Username,
		DisplayName:        user.DisplayName,
		ConnectedToSpotify: user.ConnectedToSpotify,
		DatapointCount:     datapointCount,
	}

	return c.JSON(http.StatusOK, userResponse)
}
