package server

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
	Percentage     float64 `json:"percentage"`
	DatapointCount int     `json:"datapointCount"`
}
