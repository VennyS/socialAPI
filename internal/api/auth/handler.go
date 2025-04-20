package auth

import (
	"net/http"
	"socialAPI/internal/api"
	"socialAPI/internal/api/service/auth"
	"socialAPI/internal/lib"
	l "socialAPI/internal/lib"

	"github.com/go-chi/render"
)

func (c AuthController) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := r.Context().Value(api.DataKey).(auth.UserRequest)

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
		req := r.Context().Value(api.DataKey).(auth.UserRequest)

		hErr := c.authService.Register(req)
		if hErr != nil {
			l.SendMessage(w, r, hErr.StatusCode, hErr.Error())
			return
		}

		lib.SendMessage(w, r, http.StatusCreated, "user created successfully")
	}
}

func (c AuthController) RefreshHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := r.Context().Value(api.DataKey).(auth.RefreshRequest)

		access, refresh, err := c.authService.Refresh(req)
		if err != nil {
			lib.SendMessage(w, r, err.StatusCode, err.Error())
			return
		}

		response := auth.LoginResponse{AccessToken: access, RefreshToken: refresh}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response)
	}
}

func (c AuthController) LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := r.Context().Value(api.DataKey).(auth.RefreshRequest)

		err := c.authService.Revoke(req)
		if err != nil {
			l.SendMessage(w, r, err.StatusCode, err.Error())
			return
		}

		lib.SendMessage(w, r, http.StatusOK, "revoked succefully")
	}
}
