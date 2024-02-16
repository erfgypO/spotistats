package spotistats

import (
	"context"
	"github.com/erfgypO/spotistats/bin/data"
	spotify "github.com/erfgypO/spotistats/bin/spotifyClient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"sync"
	"time"
)

func scrapDataForUser(user data.UserEntity, results chan<- data.Datapoint, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Get data for user: %s", user.Id)

	client := spotify.CreateSpotifyClient()
	currentlyPlaying, err := client.GetCurrentlyPlaying(user.Token.AccessToken)
	if err != nil {
		log.Printf("Error getting data for user: %s", user.Uid)
		return
	}

	if !currentlyPlaying.IsPlaying {
		log.Printf("User: %s is not playing anything", user.Uid)
		return
	}

	results <- data.Datapoint{
		Owner:     user.Id,
		Data:      currentlyPlaying,
		CreatedAt: time.Now().Unix(),
	}
}

func runScrapper() {
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {

		log.Printf("Error connecting to mongo: %s", err)
		return
	}

	defer func() {
		_ = mongoClient.Disconnect(context.TODO())
	}()

	collection := mongoClient.Database("spotistats").Collection(data.UserCollectionName)

	var userEntities []data.UserEntity
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Printf("Error getting data from mongo: %s", err)
		return
	}

	err = cursor.All(context.TODO(), &userEntities)
	if err != nil {
		log.Printf("Error getting data from mongo: %s", err)
	}

	var wg sync.WaitGroup

	results := make(chan data.Datapoint)

	for _, userEntity := range userEntities {
		wg.Add(1)
		go func(user data.UserEntity) {
			if user.Token.ExpiresAt-5 <= time.Now().Unix() {
				client := spotify.CreateSpotifyClient()
				token, err := client.RefreshAccessToken(user.Token.RefreshToken)
				if err != nil {
					log.Printf("Error refreshing token for user: %s", user.Uid)
					return
				}

				user.Token.AccessToken = token.AccessToken
				user.Token.RefreshToken = token.RefreshToken
				user.Token.ExpiresAt = token.ExpiresAt

				filter := bson.D{{"uid", user.Uid}}
				update := bson.D{{"$set", bson.D{{"token", user.Token}}}}

				_, err = collection.UpdateOne(context.TODO(), filter, update)
				if err != nil {
					log.Printf("Error updating token for user: %s", user.Uid)
					return
				}
			}

			scrapDataForUser(user, results, &wg)
		}(userEntity)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var currentlyPlaying []interface{}
	for result := range results {
		currentlyPlaying = append(currentlyPlaying, result)
	}

	if len(currentlyPlaying) == 0 {
		log.Printf("No data to insert")
		return
	}

	_, err = mongoClient.Database("spotistats").Collection(data.DatapointCollectionName).InsertMany(context.TODO(), currentlyPlaying)
	if err != nil {
		log.Printf("Error inserting data to mongo: %s", err)
	}
}
