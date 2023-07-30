package custommiddleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"os"
	"time"
)

type AuthMiddleware struct {
	key string
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{
		key: os.Getenv("JWT_KEY"),
	}
}

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
		}
		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) validToken(t string) bool {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return m.key, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	return true
}

func (m *AuthMiddleware) createToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	userID := uuid.NewString()

	claims["authorized"] = true
	claims["userID"] = userID // Установите имя пользователя или идентификатор здесь
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(m.key))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
