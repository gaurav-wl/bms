package server

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gauravcoco/bms/dbModels"
	"time"
)

func getJWTToken(userAuthID string, jwtKey string) (string, error) {
	claims := &dbModels.Claims{
		UUID: userAuthID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
