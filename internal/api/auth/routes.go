package auth

import (
	"socialAPI/internal/api"
	"socialAPI/internal/api/service/auth"
	"socialAPI/internal/shared"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type AuthController struct {
	authService  auth.AuthService
	tokenService shared.TokenService
	logger       *zap.SugaredLogger
}

func NewAuthController(authService auth.AuthService, tokenService shared.TokenService, logger *zap.SugaredLogger) *AuthController {
	return &AuthController{authService: authService, tokenService: tokenService, logger: logger}
}

func (a AuthController) RegisterRoutes(r *chi.Mux) {
	r.Route("/v1/auth", func(r chi.Router) {
		r.With(api.JsonBodyMiddleware[auth.UserRequest](a.logger)).Post("/login", a.LoginHandler())
		r.With(api.JsonBodyMiddleware[auth.UserRequest](a.logger)).Post("/register", a.RegisterHandler())
		r.With(api.JsonBodyMiddleware[auth.RefreshRequest](a.logger)).Post("/refresh", a.RefreshHandler())
		r.With(api.JsonBodyMiddleware[auth.RefreshRequest](a.logger)).Post("/logout", a.LogoutHandler())
	})
}
