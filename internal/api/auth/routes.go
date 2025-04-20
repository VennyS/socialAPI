package auth

import (
	"socialAPI/internal/api"
	"socialAPI/internal/api/service/auth"

	"github.com/go-chi/chi/v5"
)

type AuthController struct {
	authService auth.AuthService
}

func NewAuthController(authService auth.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (a AuthController) RegisterRoutes(r *chi.Mux) {
	r.Route("/auth", func(r chi.Router) {
		r.With(api.JsonBodyMiddleware[auth.UserRequest]()).Post("/login", a.LoginHandler())
		r.With(api.JsonBodyMiddleware[auth.UserRequest]()).Post("/register", a.RegisterHandler())
		r.With(api.JsonBodyMiddleware[auth.RefreshRequest]()).Post("/refresh", a.RefreshHandler())
	})
}
