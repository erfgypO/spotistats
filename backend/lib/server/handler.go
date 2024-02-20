package server

import (
	"context"
	"encoding/json"
	"github.com/erfgypO/spotistats/lib/data"
	spotify "github.com/erfgypO/spotistats/lib/spotifyClient"
	uuid "github.com/nu7hatch/gouuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

func writeJsonResponse(writer http.ResponseWriter, response interface{}, statusCode int) {
	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	_, _ = writer.Write(responseBytes)
}

func UseAuth(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(401)
		} else {
			token := strings.Replace(auth, "Bearer ", "", 1)
			claims, ok := verifyJWT(token)

			if !ok {
				w.WriteHeader(401)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "uid", claims["uid"]))
			next(w, r)
		}
	}
}

func HandleGetAuthUrl(writer http.ResponseWriter, request *http.Request) {
	responseType := "code"
	clientId := os.Getenv("SPOTIFY_CLIENT_ID")
	scope := "user-read-currently-playing"

	mongoClient, err := data.CreateClient()
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	defer func() {
		_ = mongoClient.Disconnect(context.TODO())
	}()

	var user data.UserEntity
	collection := mongoClient.Database("spotistats").Collection(data.UserCollectionName)
	id, err := primitive.ObjectIDFromHex(request.Context().Value("uid").(string))
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}
	err = collection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&user)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	redirectUrl := os.Getenv("REDIRECT_URL")

	response := struct {
		Url string `json:"url"`
	}{
		Url: "https://accounts.spotify.com/authorize?client_id=" + clientId + "&response_type=" + responseType + "&redirect_uri=" + redirectUrl + "&scope=" + scope + "&state=" + user.State,
	}

	writeJsonResponse(writer, response, 200)
}

func HandleAuthRedirect(writer http.ResponseWriter, request *http.Request) {
	code := request.URL.Query().Get("code")
	state := request.URL.Query().Get("state")
	client := spotify.CreateSpotifyClient()
	token, err := client.GetAccessToken(code)

	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	user, err := client.GetUser(token.AccessToken)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
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
		log.Println(err)
		writer.WriteHeader(500)
		return
	}
}

func HandlePing(writer http.ResponseWriter, request *http.Request) {
	uid := request.Context().Value("uid").(interface{})
	log.Printf("User %s is pinging", uid)
	writer.WriteHeader(200)
}

func HandleGetMe(writer http.ResponseWriter, request *http.Request) {
	mongoClient, err := data.CreateClient()
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	defer func() {
		_ = mongoClient.Disconnect(context.TODO())
	}()

	var user data.UserEntity
	collection := mongoClient.Database("spotistats").Collection(data.UserCollectionName)
	id, err := primitive.ObjectIDFromHex(request.Context().Value("uid").(string))
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}
	err = collection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&user)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	datapointCount, err := collection.Database().Collection(data.DatapointCollectionName).CountDocuments(context.TODO(), bson.D{{"owner", id}})

	userResponse := struct {
		Id                 string `json:"id"`
		Username           string `json:"username"`
		DisplayName        string `json:"displayName"`
		ConnectedToSpotify bool   `json:"connectedToSpotify"`
		DatapointCount     int64  `json:"datapointCount"`
	}{
		Id:                 user.Id.(primitive.ObjectID).Hex(),
		Username:           user.Username,
		DisplayName:        user.DisplayName,
		ConnectedToSpotify: user.ConnectedToSpotify,
		DatapointCount:     datapointCount,
	}

	writeJsonResponse(writer, userResponse, 200)
}

func HandleSignUp(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	var signUpRequest LoginRequest
	err = json.Unmarshal(body, &signUpRequest)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	mongoClient, err := data.CreateClient()
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	defer func() {
		_ = mongoClient.Disconnect(context.TODO())
	}()

	collection := mongoClient.Database("spotistats").Collection(data.UserCollectionName)

	filter := bson.D{{"username", strings.ToLower(signUpRequest.Username)}}

	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	if count > 0 {
		log.Printf("User %s already exists", signUpRequest.Username)
		errorResponse := ErrorResponse{Error: "User already exists"}

		writeJsonResponse(writer, errorResponse, 400)
		return
	}

	stateUuid, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(signUpRequest.Password), 14)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
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
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	user.Id = res.InsertedID

	tokenResponse, err := createJWT(user)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	writeJsonResponse(writer, tokenResponse, 200)
}

func HandleSignIn(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	var signInRequest LoginRequest
	err = json.Unmarshal(body, &signInRequest)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	mongoClient, err := data.CreateClient()
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}
	defer func() {
		_ = mongoClient.Disconnect(context.TODO())
	}()

	var user data.UserEntity
	collection := mongoClient.Database("spotistats").Collection(data.UserCollectionName)
	err = collection.FindOne(context.TODO(), bson.D{{"username", signInRequest.Username}}).Decode(&user)

	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signInRequest.Password)) != nil {
		log.Println(err)
		response := ErrorResponse{Error: "Invalid username or password"}
		writeJsonResponse(writer, response, 401)
		return
	}

	tokenResponse, err := createJWT(user)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	writeJsonResponse(writer, tokenResponse, 200)
}

func HandleGetStats(writer http.ResponseWriter, request *http.Request) {
	afterQuery := request.URL.Query().Get("after")
	after, err := strconv.ParseInt(afterQuery, 10, 64)
	if err != nil {
		after = 0
	}

	mongoClient, err := data.CreateClient()
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	defer func() {
		_ = mongoClient.Disconnect(context.TODO())
	}()

	collection := mongoClient.Database("spotistats").Collection(data.DatapointCollectionName)

	id, err := primitive.ObjectIDFromHex(request.Context().Value("uid").(string))
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	findOptions := options.Find().SetSort(bson.D{{"createdat", -1}})
	filter := bson.D{
		{"owner", id},
		{"createdat", bson.D{{"$gte", after}}},
	}
	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
	}

	var datapoints []data.Datapoint
	err = cursor.All(context.TODO(), &datapoints)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		return
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
			for _, artist := range datapoint.Data.Item.Artists {
				artistData[artist.ID] = NameCount{
					Name:  artist.Name,
					Count: artistData[artist.Name].Count + 1,
				}
				artistCount++
			}

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

	writeJsonResponse(writer, response, 200)
}
