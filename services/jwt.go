package services

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

type JWTClaims struct {
	Authorized bool   `json:"authorized,omitempty"`
	Email      string `json:"email,omitempty"`
	Exp        int64  `json:"exp,omitempty"`
}

func CreateJWT(email string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix() //Expire time is one week
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_JWT")))
}

func ParseJWTToken(token string) (*JWTClaims, error) {
	tokenString := token
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_JWT")), nil
	})
	if err != nil {
		return nil, err
	}
	var result JWTClaims

	if err := mapstructure.Decode(claims, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
