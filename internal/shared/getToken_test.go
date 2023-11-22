package shared

import (
	"github.com/golang-jwt/jwt/v5"
	"testing"
	"time"
)

func TestGetUserIDFromToken(t *testing.T) {
	// Setup
	key := "secret"
	validTokenString := createTestJWT("123456", key)
	wrongMethodTokenString := createTestJWTWithWrongMethod("123456", key)

	tests := []struct {
		name    string
		token   string
		key     string
		want    string
		wantErr bool
	}{
		{
			name:    "Valid Token",
			token:   validTokenString,
			key:     key,
			want:    "123456",
			wantErr: false,
		},
		{
			name:    "Invalid Token",
			token:   "invalidToken",
			key:     key,
			want:    "",
			wantErr: true,
		},
		{
			name:    "Wrong Method Token",
			token:   wrongMethodTokenString,
			key:     key,
			want:    "",
			wantErr: true,
		},
		// Additional test cases can be added here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserIDFromToken(tt.token, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserIDFromToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserIDFromToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to create a JWT for testing
func createTestJWT(userID string, key string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(key))
	return tokenString
}

// Helper function to create a JWT with a wrong signing method for testing
func createTestJWTWithWrongMethod(userID string, key string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	return tokenString
}
