package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("sdmn340dcnsdlkasldfj8e3wr4qwejrnwemrnwejkrnkwe")

func GenerateToken(userId uint, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   userId,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
		"username": username,
	})

	return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorClaimsInvalid)
		}
		return jwtSecret, nil
	})
	return token, err
}
