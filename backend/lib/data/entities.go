package data

import spotify "github.com/erfgypO/spotistats/lib/spotifyClient"

type UserEntity struct {
	Id                 interface{} `bson:"_id,omitempty"`
	Uid                string
	Username           string
	Password           string
	DisplayName        string
	Token              TokenEntity
	State              string
	ConnectedToSpotify bool
}

type TokenEntity struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
}

type RGB struct {
	R uint32 `json:"r"`
	G uint32 `json:"g"`
	B uint32 `json:"b"`
}

type Datapoint struct {
	Id           interface{}              `json:"id" bson:"_id,omitempty"`
	Owner        interface{}              `json:"owner"`
	Data         spotify.CurrentlyPlaying `json:"data"`
	CreatedAt    int64                    `json:"createdat"`
	PrimaryColor RGB                      `json:"primarycolor"`
}
