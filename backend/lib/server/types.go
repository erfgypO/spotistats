package server

import "github.com/erfgypO/spotistats/lib/data"

type ErrorResponse struct {
	Error string `json:"error"`
}

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresAt   int64  `json:"expiresAt"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DataPercentage struct {
	Name           string  `json:"name"`
	SpotifyUrl     string  `json:"spotifyUrl"`
	Percentage     float64 `json:"percentage"`
	DatapointCount int     `json:"datapointCount"`
}

type UserResponse struct {
	Id                 string `json:"id"`
	Username           string `json:"username"`
	DisplayName        string `json:"displayName"`
	ConnectedToSpotify bool   `json:"connectedToSpotify"`
	DatapointCount     int64  `json:"datapointCount"`
}

type HourlyStats struct {
	Hour     int      `json:"hour"`
	Seconds  int      `json:"seconds"`
	SongName string   `json:"songName"`
	Color    data.RGB `json:"color"`
}
