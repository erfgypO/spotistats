package server

import (
	"github.com/erfgypO/spotistats/bin/data"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

func createJWT(user data.UserEntity) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":      user.Id,
		"username": user.Username,
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func verifyJWT(tokenString string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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
