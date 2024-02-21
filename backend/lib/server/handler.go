package server

import (
	"context"
	"errors"
	"github.com/erfgypO/spotistats/lib/data"
	spotify "github.com/erfgypO/spotistats/lib/spotifyClient"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	uuid "github.com/nu7hatch/gouuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

func HandleGetAuthUrl(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return errors.New("JWT token missing or invalid")
	}
	claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
	if !ok {
		return errors.New("failed to cast claims as jwt.MapClaims")
	}

	responseType := "code"
	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	scope := "user-read-currently-playing"

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

	redirectUrl := os.Getenv("REDIRECT_URL")

	response := struct {
		Url string `json:"url"`
	}{
		Url: "https://accounts.spotify.com/authorize?client_id=" + clientId + "&response_type=" + responseType + "&redirect_uri=" + redirectUrl + "&scope=" + scope + "&state=" + user.State,
	}

	return c.JSON(http.StatusOK, response)
}

func HandleAuthRedirect(c echo.Context) error {
	code := c.QueryParams().Get("code")
	state := c.QueryParams().Get("state")
	client := spotify.CreateSpotifyClient()
	token, err := client.GetAccessToken(code)

	if err != nil {
		return err
	}

	user, err := client.GetUser(token.AccessToken)
	if err != nil {
		return err
	}

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		return err
	}

	defer func() {
		_ = mongoClient.Disconnect(context.TODO())
	}()

	collection := mongoClient.Database("spotistats").Collection(data.UserCollectionName)

	filter := bson.D{{"state", state}}

	tokenEntity := data.TokenEntity{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.ExpiresAt,
	}

	update := bson.D{{"$set", bson.D{{"token", tokenEntity}, {"displayname", user.DisplayName}, {"connectedtospotify", true}}}}

	_, err = collection.UpdateOne(context.TODO(), filter, update, options.Update())

	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}

func HandleGetMe(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return errors.New("JWT token missing or invalid")
	}
	claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
	if !ok {
		return errors.New("failed to cast claims as jwt.MapClaims")
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

func HandleSignUp(c echo.Context) error {
	signUpRequest := new(LoginRequest)
	if err := c.Bind(&signUpRequest); err != nil {
		c.Logger().Error(err)
		return err
	}

	mongoClient, err := data.CreateClient()
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	defer func() {
		_ = mongoClient.Disconnect(context.TODO())
	}()

	collection := mongoClient.Database("spotistats").Collection(data.UserCollectionName)

	filter := bson.D{{"username", strings.ToLower(signUpRequest.Username)}}

	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	if count > 0 {
		errorResponse := ErrorResponse{Error: "User already exists"}
		return c.JSON(http.StatusBadRequest, errorResponse)
	}

	stateUuid, err := uuid.NewV4()
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(signUpRequest.Password), 14)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	user := data.UserEntity{
		Username:           signUpRequest.Username,
		Password:           string(passwordHash),
		DisplayName:        "",
		Token:              data.TokenEntity{},
		Uid:                "",
		State:              signUpRequest.Username + stateUuid.String(),
		ConnectedToSpotify: false,
	}

	res, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	user.Id = res.InsertedID

	tokenResponse, err := createJWT(user)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return c.JSON(http.StatusOK, tokenResponse)
}

func HandleSignIn(c echo.Context) error {
	signInRequest := new(LoginRequest)
	if err := c.Bind(&signInRequest); err != nil {
		c.Logger().Error(err)
		return err
	}

	mongoClient, err := data.CreateClient()
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	defer func() {
		_ = mongoClient.Disconnect(context.TODO())
	}()

	var user data.UserEntity
	collection := mongoClient.Database("spotistats").Collection(data.UserCollectionName)
	err = collection.FindOne(context.TODO(), bson.D{{"username", signInRequest.Username}}).Decode(&user)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signInRequest.Password)) != nil {
		response := ErrorResponse{Error: "Invalid username or password"}
		return c.JSON(http.StatusBadRequest, response)
	}

	tokenResponse, err := createJWT(user)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return c.JSON(http.StatusOK, tokenResponse)
}

func HandleGetStats(c echo.Context) error {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return errors.New("JWT token missing or invalid")
	}
	claims, ok := token.Claims.(jwt.MapClaims) // by default claims is of type `jwt.MapClaims`
	if !ok {
		return errors.New("failed to cast claims as jwt.MapClaims")
	}

	afterQuery := c.QueryParams().Get("after")
	after, err := strconv.ParseInt(afterQuery, 10, 64)
	if err != nil {
		after = 0
	}

	mongoClient, err := data.CreateClient()
	if err != nil {
		return err
	}

	defer func() {
		_ = mongoClient.Disconnect(context.TODO())
	}()

	collection := mongoClient.Database("spotistats").Collection(data.DatapointCollectionName)

	id, err := primitive.ObjectIDFromHex(claims["uid"].(string))
	if err != nil {
		return err
	}

	findOptions := options.Find().SetSort(bson.D{{"createdat", -1}})
	filter := bson.D{
		{"owner", id},
		{"createdat", bson.D{{"$gte", after}}},
	}
	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return err
	}

	var datapoints []data.Datapoint
	err = cursor.All(context.TODO(), &datapoints)
	if err != nil {
		return err
	}

	type NameCount struct {
		Name  string
		Count int
	}

	artistData := make(map[string]NameCount)
	artistCount := 0

	trackData := make(map[string]NameCount)
	trackCount := 0

	for _, datapoint := range datapoints {
		if datapoint.Data.Item.Type == "track" {
			//artistCount[datapoint.Data.Item.Artists[0].Name]++
			artistData[datapoint.Data.Item.Artists[0].ID] = NameCount{
				Name:  datapoint.Data.Item.Artists[0].Name,
				Count: artistData[datapoint.Data.Item.Artists[0].ID].Count + 1,
			}
			artistCount++

			trackData[datapoint.Data.Item.ID] = NameCount{
				Name:  datapoint.Data.Item.Name,
				Count: trackData[datapoint.Data.Item.ID].Count + 1,
			}
			trackCount++
		}
	}

	response := struct {
		Artists []DataPercentage `json:"artists"`
		Tracks  []DataPercentage `json:"tracks"`
	}{
		Artists: make([]DataPercentage, 0),
		Tracks:  make([]DataPercentage, 0),
	}

	for _, item := range artistData {
		percentage := float64(item.Count) / float64(trackCount)
		response.Artists = append(response.Artists, DataPercentage{
			Name:           item.Name,
			Percentage:     percentage,
			DatapointCount: item.Count,
		})
	}

	sort.Slice(response.Artists, func(i, j int) bool {
		return response.Artists[i].DatapointCount > response.Artists[j].DatapointCount
	})

	for _, item := range trackData {
		percentage := float64(item.Count) / float64(trackCount)
		response.Tracks = append(response.Tracks, DataPercentage{
			Name:           item.Name,
			Percentage:     percentage,
			DatapointCount: item.Count,
		})
	}

	sort.Slice(response.Tracks, func(i, j int) bool {
		return response.Tracks[i].DatapointCount > response.Tracks[j].DatapointCount
	})

	return c.JSON(http.StatusOK, response)
}
