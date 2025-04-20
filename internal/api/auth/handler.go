package auth

import (
	"net/http"
	"socialAPI/internal/api"
	"socialAPI/internal/api/service/auth"
	l "socialAPI/internal/lib"

	"github.com/go-chi/render"
)

func (c AuthController) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := l.GetFromContext[auth.UserRequest](r.Context(), api.DataKey)
		if err != nil {
			http.Error(w, "invalid login data", http.StatusBadRequest)
			return
		}

		access, refresh, err := c.authService.Authenticate(req)

		response := auth.LoginResponse{AccessToken: access, RefreshToken: refresh}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response)
	}
}

func (c AuthController) RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
func (c AuthController) RefreshHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
func (c AuthController) LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
