// Package shared provides utility functions for various common tasks.
package shared

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

// GetUserIDFromToken parses a JWT token to extract the userID. The function checks the token's
// validity by validating its signing method and then tries to extract the "userID" claim.
//
// Parameters:
//   - t:   The JWT token as a string.
//   - key: The secret key used to sign the token.
//
// Returns:
//   - A string representing the userID if found and valid, or an error if the token is invalid,
//     the signing method is unexpected, or the userID claim is not present or is of an unexpected type.
func GetUserIDFromToken(t string, key string) (string, error) {
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
