package utils

import (
	"errors"
	"fmt"
	"github.com/FeelDat/urlshort/internal/app/models"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserIDFromToken(t string, key models.JWTKey) (string, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(key), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("unexpected claims type")
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		return "", errors.New("userID is not a string")
	}

	return userID, nil
}
