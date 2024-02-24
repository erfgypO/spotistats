package server

import (
	"context"
	"errors"
	"github.com/erfgypO/spotistats/lib/data"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"sort"
	"strconv"
)

func HandleGetStats(c echo.Context) error {
	claims, err := getClaimsFromContext(c)
	if err != nil {
		return err
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
		Url   string
	}

	artistData := make(map[string]NameCount)
	artistCount := 0

	trackData := make(map[string]NameCount)
	trackCount := 0

	albumData := make(map[string]NameCount)
	albumCount := 0

	for _, datapoint := range datapoints {
		if datapoint.Data.Item.Type != "track" {
			continue
		}

		artistData[datapoint.Data.Item.Artists[0].ID] = NameCount{
			Name:  datapoint.Data.Item.Artists[0].Name,
			Count: artistData[datapoint.Data.Item.Artists[0].ID].Count + 1,
			Url:   datapoint.Data.Item.Artists[0].ExternalUrls.Spotify,
		}
		artistCount++

		trackData[datapoint.Data.Item.ID] = NameCount{
			Name:  datapoint.Data.Item.Name,
			Count: trackData[datapoint.Data.Item.ID].Count + 1,
			Url:   datapoint.Data.Item.ExternalUrls.Spotify,
		}
		trackCount++

		albumData[datapoint.Data.Item.Album.ID] = NameCount{
			Name:  datapoint.Data.Item.Album.Name,
			Count: albumData[datapoint.Data.Item.Album.ID].Count + 1,
			Url:   datapoint.Data.Item.Album.ExternalUrls.Spotify,
		}
		albumCount++
	}

	response := struct {
		Artists []DataPercentage `json:"artists"`
		Tracks  []DataPercentage `json:"tracks"`
		Albums  []DataPercentage `json:"albums"`
	}{
		Artists: make([]DataPercentage, 0),
		Tracks:  make([]DataPercentage, 0),
		Albums:  make([]DataPercentage, 0),
	}

	for _, item := range artistData {
		percentage := float64(item.Count) / float64(trackCount)
		response.Artists = append(response.Artists, DataPercentage{
			Name:           item.Name,
			Percentage:     percentage,
			DatapointCount: item.Count,
			SpotifyUrl:     item.Url,
		})
	}

	sort.Slice(response.Artists, func(i, j int) bool {
		return response.Artists[i].DatapointCount >= response.Artists[j].DatapointCount
	})

	for _, item := range trackData {
		percentage := float64(item.Count) / float64(trackCount)
		response.Tracks = append(response.Tracks, DataPercentage{
			Name:           item.Name,
			Percentage:     percentage,
			DatapointCount: item.Count,
			SpotifyUrl:     item.Url,
		})
	}

	sort.Slice(response.Tracks, func(i, j int) bool {
		return response.Tracks[i].DatapointCount >= response.Tracks[j].DatapointCount
	})

	for _, item := range albumData {
		percentage := float64(item.Count) / float64(trackCount)
		response.Albums = append(response.Albums, DataPercentage{
			Name:           item.Name,
			Percentage:     percentage,
			DatapointCount: item.Count,
			SpotifyUrl:     item.Url,
		})
	}

	sort.Slice(response.Albums, func(i, j int) bool {
		return response.Albums[i].DatapointCount >= response.Albums[j].DatapointCount
	})

	return c.JSON(http.StatusOK, response)
}

func HandleGetStatsGroupedByHour(c echo.Context) error {
	claims, err := getClaimsFromContext(c)
	if err != nil {
		return err
	}

	afterQuery := c.QueryParams().Get("after")
	after, err := strconv.ParseInt(afterQuery, 10, 64)
	if err != nil {
		return errors.New("failed to parse after query param")
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
	/*
		addFieldsStage := bson.D{{"$addFields", bson.D{
			{"date", bson.D{
				{"$toDate", bson.D{
					{"$multiply", bson.A{"$createdat", 1000}},
				}}}}, {"songName", "$data.item.name"}}}}

		matchStage := bson.D{{"$match", bson.D{
			{"owner", id},
			{"date", bson.D{{"$gte", bson.D{{"$toDate", after}}}}},
			{"songName", bson.D{{"$ne", ""}}},
		}}}

		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{
				{"hour", bson.D{{"$hour", "$date"}}},
				{"songName", "$songName"},
			}},

			{"seconds", bson.D{{"$sum", 10}}},
		}}}

		projectStage := bson.D{{"$project", bson.D{
			{"_id", 0},
			{"hour", bson.D{{"$add", bson.A{bson.D{{"$hour", "$date"}}, 1}}}},
			{"seconds", 1},
			{"songName", "$_id.songName"},
		}}}
	*/

	addFieldsStage := bson.D{
		{"$addFields", bson.D{
			{"date", bson.D{
				{"$toDate", bson.D{
					{"$multiply", bson.A{"$createdat", 1000}},
				}},
			}},
			{"songName", "$data.item.name"},
		}},
	}

	matchStage := bson.D{
		{"$match", bson.D{
			{"owner", id},
			{"date", bson.D{{"$gte", after}}},
			{"songName", bson.D{{"$ne", ""}}},
		}},
	}

	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", bson.D{
				{"hour", bson.D{{"$hour", "$date"}}},
				{"songName", "$songName"},
			}},
			{"seconds", bson.D{{"$sum", 10}}},
		}},
	}

	projectStage := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"hour", bson.D{{"$add", bson.A{bson.D{{"$hour", "$date"}}, 1}}}},
			{"seconds", 1},
			{"songName", "$_id.songName"},
		}},
	}

	cursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{addFieldsStage, matchStage, groupStage, projectStage})
	if err != nil {
		return err
	}

	var hourlyStats []HourlyStats
	err = cursor.All(context.TODO(), &hourlyStats)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, hourlyStats)
}
