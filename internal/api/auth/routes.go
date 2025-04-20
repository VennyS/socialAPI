package auth

import (
	"socialAPI/internal/api"
	"socialAPI/internal/api/service/auth"
	"socialAPI/internal/shared"

	"github.com/go-chi/chi/v5"
)

type AuthController struct {
	authService  auth.AuthService
	tokenService shared.TokenService
}

func NewAuthController(authService auth.AuthService, tokenService shared.TokenService) *AuthController {
	return &AuthController{authService: authService, tokenService: tokenService}
}

func (a AuthController) RegisterRoutes(r *chi.Mux) {
	r.Route("/auth", func(r chi.Router) {
		r.With(api.JsonBodyMiddleware[auth.UserRequest]()).Post("/login", a.LoginHandler())
		r.With(api.JsonBodyMiddleware[auth.UserRequest]()).Post("/register", a.RegisterHandler())
		r.With(api.JsonBodyMiddleware[auth.RefreshRequest]()).Post("/refresh", a.RefreshHandler())
		r.With(api.JsonBodyMiddleware[auth.RefreshRequest]()).Post("/logout", a.LogoutHandler())
	})
}
