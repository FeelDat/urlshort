package custommiddleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware_AuthMiddleware(t *testing.T) {
	type fields struct {
		key string
	}
	type args struct {
		next http.Handler
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		setupRequest   func(r *http.Request)
		expectedStatus int
		expectJWTToken bool
	}{
		{
			name: "No Token in Request",
			fields: fields{
				key: "testKey",
			},
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
			setupRequest:   func(r *http.Request) {},
			expectedStatus: http.StatusOK,
			expectJWTToken: true,
		},
		{
			name: "Invalid Token in Request",
			fields: fields{
				key: "testKey",
			},
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
			setupRequest: func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:  "jwt",
					Value: "invalidToken",
				})
			},
			expectedStatus: http.StatusOK,
			expectJWTToken: true,
		},
		// Additional test cases can be added here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AuthMiddleware{
				key: tt.fields.key,
			}
			handler := m.AuthMiddleware(tt.args.next)
			rr := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Set up the request as per the test case
			tt.setupRequest(req)

			handler.ServeHTTP(rr, req)

			// Check the status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("%s: handler returned wrong status code: got %v, want %v", tt.name, rr.Code, tt.expectedStatus)
			}

			// Check for JWT token presence in the response
			foundJWTToken := false
			for _, cookie := range rr.Result().Cookies() {
				if cookie.Name == "jwt" {
					foundJWTToken = true
					rr.Result().Body.Close() // Close the response body here
					break
				}
			}
			if foundJWTToken != tt.expectJWTToken {
				t.Errorf("%s: expected JWT token presence: %v, found: %v", tt.name, tt.expectJWTToken, foundJWTToken)
			}
		})
	}
}

//TODO implement tests for other middlewares
