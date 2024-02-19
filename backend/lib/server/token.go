package server

import (
	"github.com/erfgypO/spotistats/lib/data"
	"github.com/golang-jwt/jwt/v5"
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

func verifyJWT(tokenString string) (jwt.MapClaims, bool) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ok
	}

	return claims, ok
}
