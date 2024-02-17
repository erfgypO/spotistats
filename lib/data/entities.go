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

type Datapoint struct {
	Id        interface{} `bson:"_id,omitempty"`
	Owner     interface{}
	Data      spotify.CurrentlyPlaying
	CreatedAt int64
}
