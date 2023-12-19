// Package custommiddleware provides custom middlewares for use in an HTTP server.
package custommiddleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"os"
	"time"
)

// AuthMiddleware is a struct responsible for handling authentication via JWT.
type AuthMiddleware struct {
	key string
}

// NewAuthMiddleware initializes and returns an instance of AuthMiddleware
// with the JWT key retrieved from the environment variable JWT_KEY.
func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{
		key: os.Getenv("JWT_KEY"),
	}
}

// AuthMiddleware is a middleware function that checks for the presence and validity
// of a JWT token in the request's cookies. If no valid token is found, a new one is created
// and set as a cookie in the response.
func (m *AuthMiddleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("jwt")
		if err != nil || !m.validToken(cookie.Value) {
			token, err := m.createToken()
			if err != nil {
				http.Error(w, "Issue with creating JWT token", http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "jwt",
				Value:    token,
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
			})

			r.AddCookie(&http.Cookie{
				Name:     "jwt",
				Value:    token,
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
			})
		}
		next.ServeHTTP(w, r)
	})
}

// validToken checks the validity of a given JWT token.
func (m *AuthMiddleware) validToken(t string) bool {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(m.key), nil
	})

	if err != nil || !token.Valid {
		return false
	}

	return true
}

// createToken generates a new JWT token with claims that include a new userID,
// an authorization flag set to true, and an expiration set 24 hours from the creation time.
func (m *AuthMiddleware) createToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	userID := uuid.NewString()

	claims["authorized"] = true
	claims["userID"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	m.key = "8PNHgjK2kPunGpzMgL0ZmMdJCRKy2EnL/Cg0GbnELLI="
	tokenString, err := token.SignedString([]byte(m.key))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
