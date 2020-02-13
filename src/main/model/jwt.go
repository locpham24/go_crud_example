package model

import (
	"github.com/dgrijalva/jwt-go"
	"main/conf"
	"time"
)

type JwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func CreateJwtToken(name string, id string) (string, error) {
	jwtClaims := JwtClaims{
		name,
		jwt.StandardClaims{
			Id:        id,
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtClaims)
	token, err := rawToken.SignedString(conf.JwtKey)
	if err != nil {
		return "", err
	}
	return token, nil
}
