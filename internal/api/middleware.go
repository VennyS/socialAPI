package api

import (
	"context"
	"encoding/json"
	"net/http"
	l "socialAPI/internal/lib"
	"socialAPI/internal/shared"
	"strings"

	"go.uber.org/zap"
)

type key string

var (
	DataKey   key = "data"
	UserIDKey key = "userID"
)

func AuthMiddleware(tokenService *shared.TokenService, logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warnw("Authorization header is missing", "error", "Authorization header required")
				l.SendMessage(w, r, http.StatusUnauthorized, "Authorization header is required")
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				logger.Warnw("Invalid authorization format", "header", authHeader)
				l.SendMessage(w, r, http.StatusUnauthorized, "Invalid authorization format")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := tokenService.ValidateToken(tokenString)
			if err != nil {
				logger.Warnw("Invalid token", "error", err.Error())
				l.SendMessage(w, r, http.StatusUnauthorized, "Invalid token: "+err.Error())
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func JsonBodyMiddleware[T any](logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body T

			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				logger.Warnw("Failed to decode JSON body", "error", err.Error())
				l.SendMessage(w, r, http.StatusBadRequest, "invalid request body")
				return
			}

			if err := shared.Validate.Struct(body); err != nil {
				logger.Warnw("Validation failed", "error", err.Error())
				l.SendMessage(w, r, http.StatusBadRequest, err.Error())
				return
			}

			ctx := context.WithValue(r.Context(), DataKey, body)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
