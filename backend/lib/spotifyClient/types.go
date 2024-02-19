package spotifyClient

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64
}

type User struct {
	DisplayName string `json:"display_name"`
	Id          string `json:"id"`
}

type CurrentlyPlaying struct {
	Timestamp            int64   `json:"timestamp"`
	Context              Context `json:"context"`
	ProgressMs           int     `json:"progress_ms"`
	Item                 Item    `json:"item"`
	CurrentlyPlayingType string  `json:"currently_playing_type"`
	Actions              Actions `json:"actions"`
	IsPlaying            bool    `json:"is_playing"`
}
type ExternalUrls struct {
	Spotify string `json:"spotify"`
}
type Context struct {
	ExternalUrls ExternalUrls `json:"external_urls"`
	Href         string       `json:"href"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}
type Artists struct {
	ExternalUrls ExternalUrls `json:"external_urls"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}
type Images struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}
type Album struct {
	AlbumType            string       `json:"album_type"`
	Artists              []Artists    `json:"artists"`
	AvailableMarkets     []string     `json:"available_markets"`
	ExternalUrls         ExternalUrls `json:"external_urls"`
	Href                 string       `json:"href"`
	ID                   string       `json:"id"`
	Images               []Images     `json:"images"`
	Name                 string       `json:"name"`
	ReleaseDate          string       `json:"release_date"`
	ReleaseDatePrecision string       `json:"release_date_precision"`
	TotalTracks          int          `json:"total_tracks"`
	Type                 string       `json:"type"`
	URI                  string       `json:"uri"`
}
type ExternalIds struct {
	Isrc string `json:"isrc"`
}
type Item struct {
	Album            Album        `json:"album"`
	Artists          []Artists    `json:"artists"`
	AvailableMarkets []string     `json:"available_markets"`
	DiscNumber       int          `json:"disc_number"`
	DurationMs       int          `json:"duration_ms"`
	Explicit         bool         `json:"explicit"`
	ExternalIds      ExternalIds  `json:"external_ids"`
	ExternalUrls     ExternalUrls `json:"external_urls"`
	Href             string       `json:"href"`
	ID               string       `json:"id"`
	IsLocal          bool         `json:"is_local"`
	Name             string       `json:"name"`
	Popularity       int          `json:"popularity"`
	PreviewURL       string       `json:"preview_url"`
	TrackNumber      int          `json:"track_number"`
	Type             string       `json:"type"`
	URI              string       `json:"uri"`
}
type Disallows struct {
	Pausing bool `json:"pausing"`
}
type Actions struct {
	Disallows Disallows `json:"disallows"`
}
