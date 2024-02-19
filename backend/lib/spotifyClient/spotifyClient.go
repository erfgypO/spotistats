package spotifyClient

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type SpotifyClient struct {
	clientId     string
	clientSecret string
}

func CreateSpotifyClient() SpotifyClient {
	return SpotifyClient{
		clientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
		clientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
	}
}

func (client *SpotifyClient) GetAccessToken(code string) (Token, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", os.Getenv("REDIRECT_URL"))

	r, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return Token{}, err
	}
	authToken := base64.StdEncoding.EncodeToString([]byte(os.Getenv("SPOTIFY_CLIENT_ID") + ":" + os.Getenv("SPOTIFY_CLIENT_SECRET")))
	r.Header.Add("Authorization", "Basic "+authToken)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return Token{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Token{}, err
	}

	var token Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return Token{}, err
	}
	token.ExpiresAt = time.Now().Unix() + int64(token.ExpiresIn)

	return token, nil
}

func (client *SpotifyClient) RefreshAccessToken(refreshToken string) (Token, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	r, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return Token{}, err
	}
	authToken := base64.StdEncoding.EncodeToString([]byte(os.Getenv("SPOTIFY_CLIENT_ID") + ":" + os.Getenv("SPOTIFY_CLIENT_SECRET")))
	r.Header.Add("Authorization", "Basic "+authToken)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return Token{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Token{}, err
	}

	var token Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return Token{}, err
	}
	token.ExpiresAt = time.Now().Unix() + int64(token.ExpiresIn)
	token.RefreshToken = refreshToken
	return token, nil
}

func (client *SpotifyClient) GetUser(accessToken string) (User, error) {
	r, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return User{}, err
	}
	r.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return User{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return User{}, err
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (client *SpotifyClient) GetCurrentlyPlaying(accessToken string) (CurrentlyPlaying, error) {
	r, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/currently-playing", nil)
	if err != nil {
		return CurrentlyPlaying{}, err
	}
	r.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return CurrentlyPlaying{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return CurrentlyPlaying{}, err
	}

	var currentlyPlaying CurrentlyPlaying
	err = json.Unmarshal(body, &currentlyPlaying)
	if err != nil {
		return CurrentlyPlaying{}, err
	}

	return currentlyPlaying, nil
}
