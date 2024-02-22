package server

import (
	"context"
	"github.com/erfgypO/spotistats/lib/data"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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

	datapointCountFilter := bson.D{{"owner", id}, {"data.item.type", "track"}}
	datapointCount, err := collection.Database().Collection(data.DatapointCollectionName).CountDocuments(context.TODO(), datapointCountFilter)

	userResponse := UserResponse{
		Id:                 user.Id.(primitive.ObjectID).Hex(),
		Username:           user.Username,
		DisplayName:        user.DisplayName,
		ConnectedToSpotify: user.ConnectedToSpotify,
		DatapointCount:     datapointCount,
	}

	return c.JSON(http.StatusOK, userResponse)
}

func HandleUpdatePassword(c echo.Context) error {
	claims, err := getClaimsFromContext(c)
	if err != nil {
		return err
	}

	updatePasswordRequest := new(LoginRequest)
	if err := c.Bind(&updatePasswordRequest); err != nil {
		c.Logger().Error(err)
		return err
	}

	if len(updatePasswordRequest.Password) < 8 {
		errorResponse := ErrorResponse{Error: "Password must be at least 8 characters long"}
		return c.JSON(http.StatusBadRequest, errorResponse)
	}

	mongoClient, err := data.CreateClient()
	if err != nil {
		return err
	}

	defer func() {
		_ = mongoClient.Disconnect(context.TODO())
	}()

	collection := mongoClient.Database("spotistats").Collection(data.UserCollectionName)

	id, err := primitive.ObjectIDFromHex(claims["uid"].(string))
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatePasswordRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(context.TODO(), bson.D{{"_id", id}}, bson.D{{"$set", bson.D{{"password", string(hashedPassword)}}}})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
