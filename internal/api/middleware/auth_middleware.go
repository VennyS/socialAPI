package middleware

import (
	"context"
	"net/http"
	"socialAPI/internal/lib"
	"socialAPI/internal/shared"
	"strings"

	"go.uber.org/zap"
)

var UserIDKey key = "userID"

func AuthMiddleware(tokenService shared.TokenService, logger *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warnw("Authorization header is missing", "error", "Authorization header required")
				lib.SendMessage(w, r, http.StatusUnauthorized, "Authorization header is required")
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				logger.Warnw("Invalid authorization format", "header", authHeader)
				lib.SendMessage(w, r, http.StatusUnauthorized, "Invalid authorization format")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := tokenService.ValidateToken(tokenString)
			if err != nil {
				logger.Warnw("Invalid token", "error", err.Error())
				lib.SendMessage(w, r, http.StatusUnauthorized, "Invalid token: "+err.Error())
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
