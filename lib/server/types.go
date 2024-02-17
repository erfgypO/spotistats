package server

type ErrorResponse struct {
	Error string `json:"error"`
}

type TokenResponse struct {
	AccessToken string `json:"accessToken"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
