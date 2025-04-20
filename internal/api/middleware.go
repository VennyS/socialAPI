package api

import (
	"context"
	"encoding/json"
	"net/http"
	l "socialAPI/internal/lib"

	"github.com/go-playground/validator/v10"
)

type key string

var DataKey key = "data"

func JsonBodyMiddleware[T any]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body T
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				l.SendError(w, r, http.StatusBadRequest, "invalid request body")
				return
			}

			validate := validator.New()

			if err := validate.Struct(body); err != nil {
				l.SendError(w, r, http.StatusBadRequest, err.Error())
				return
			}

			ctx := context.WithValue(r.Context(), DataKey, body)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
