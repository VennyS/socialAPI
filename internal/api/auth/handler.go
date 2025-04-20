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
		req := r.Context().Value(api.IdKey).(auth.UserRequest)

		access, refresh, hErr := c.authService.Authenticate(req)
		if hErr != nil {
			l.SendMessage(w, r, hErr.StatusCode, hErr.Error())
			return
		}

		response := auth.LoginResponse{AccessToken: access, RefreshToken: refresh}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response)
	}
}

func (c AuthController) RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := r.Context().Value(api.IdKey).(auth.UserRequest)

		hErr := c.authService.Register(req)
		if hErr != nil {
			l.SendMessage(w, r, hErr.StatusCode, hErr.Error())
			return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, map[string]string{
			"message": "user created successfully",
		})
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
