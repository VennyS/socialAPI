package api

import (
	"context"
	"encoding/json"
	"net/http"
	l "socialAPI/internal/lib"
	"socialAPI/internal/shared"
	"strings"
)

type key string

var (
	DataKey   key = "data"
	UserIDKey key = "userID" // New context key for storing user ID
)

// AuthMiddleware checks for valid JWT token and extracts user ID
func AuthMiddleware(tokenService *shared.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				l.SendMessage(w, r, http.StatusUnauthorized, "Authorization header is required")
				return
			}

			// Check if the header has the "Bearer " prefix
			if !strings.HasPrefix(authHeader, "Bearer ") {
				l.SendMessage(w, r, http.StatusUnauthorized, "Invalid authorization format")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := tokenService.ValidateToken(tokenString)
			if err != nil {
				l.SendMessage(w, r, http.StatusUnauthorized, "Invalid token: "+err.Error())
				return
			}

			// Add user ID to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func JsonBodyMiddleware[T any]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body T

			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				l.SendMessage(w, r, http.StatusBadRequest, "invalid request body")
				return
			}

			if err := l.Validate.Struct(body); err != nil {
				l.SendMessage(w, r, http.StatusBadRequest, err.Error())
				return
			}

			ctx := context.WithValue(r.Context(), DataKey, body)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
