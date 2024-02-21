package server

import (
	"errors"
	"github.com/erfgypO/spotistats/lib/data"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"os"
	"time"
)

func createJWT(user data.UserEntity) (TokenResponse, error) {
	expiresAt := time.Now().Add(time.Hour * 2).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":      user.Id,
		"username": user.Username,
		"exp":      expiresAt,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{
		AccessToken: tokenString,
		ExpiresAt:   expiresAt,
	}, nil
}

func getClaimsFromContext(c echo.Context) (jwt.MapClaims, error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil, errors.New("JWT token missing or invalid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to cast claims as jwt.MapClaims")
	}
	return claims, nil
}
