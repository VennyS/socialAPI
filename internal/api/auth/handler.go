package auth

import (
	"net/http"
	"socialAPI/internal/api"
	"socialAPI/internal/api/service/auth"
	"socialAPI/internal/lib"

	"github.com/go-chi/render"
)

func (c AuthController) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := r.Context().Value(api.DataKey).(auth.UserRequest)

		c.logger.Infow("Login attempt", "email", req.Email)

		tokenPair, hErr := c.authService.Authenticate(req)
		if hErr != nil {
			c.logger.Warnw("Login failed", "error", hErr.Error(), "email", req.Email)
			lib.SendMessage(w, r, hErr.StatusCode, hErr.Error())
			return
		}

		c.logger.Infow("Login success", "email", req.Email)

		response := auth.LoginResponse{AccessToken: tokenPair.AccessToken, RefreshToken: tokenPair.RefreshToken}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response)
	}
}

func (c AuthController) RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := r.Context().Value(api.DataKey).(auth.UserRequest)

		c.logger.Infow("Register attempt", "email", req.Email)

		hErr := c.authService.Register(req)
		if hErr != nil {
			c.logger.Warnw("Registration failed", "error", hErr.Error(), "email", req.Email)
			lib.SendMessage(w, r, hErr.StatusCode, hErr.Error())
			return
		}

		c.logger.Infow("Registration success", "email", req.Email)

		lib.SendMessage(w, r, http.StatusCreated, "user created successfully")
	}
}

func (c AuthController) RefreshHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := r.Context().Value(api.DataKey).(auth.RefreshRequest)

		c.logger.Infow("Refresh token attempt", "refresh token", req.Refresh)

		tokenPair, err := c.authService.Refresh(req)
		if err != nil {
			c.logger.Warnw("Token refresh failed", "error", err.Error(), "refresh token", req.Refresh)
			lib.SendMessage(w, r, err.StatusCode, err.Error())
			return
		}

		c.logger.Infow("Token refresh success", "refresh token", req.Refresh)

		response := auth.LoginResponse{AccessToken: tokenPair.AccessToken, RefreshToken: tokenPair.RefreshToken}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response)
	}
}

func (c AuthController) LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := r.Context().Value(api.DataKey).(auth.RefreshRequest)

		c.logger.Infow("Logout attempt", "refresh token", req.Refresh)

		err := c.authService.Revoke(req)
		if err != nil {
			c.logger.Warnw("Logout failed", "error", err.Error(), "refresh token", req.Refresh)
			lib.SendMessage(w, r, err.StatusCode, err.Error())
			return
		}

		c.logger.Infow("Logout success", "refresh token", req.Refresh)

		lib.SendMessage(w, r, http.StatusOK, "revoked successfully")
	}
}
