package user

import (
	"socialAPI/internal/api/middleware"
	"socialAPI/internal/shared"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UserController struct {
	userService  UserService
	tokenService shared.TokenService
	logger       *zap.SugaredLogger
}

func NewAuthController(userService UserService, tokenService shared.TokenService, logger *zap.SugaredLogger) *UserController {
	return &UserController{userService: userService, tokenService: tokenService, logger: logger}
}

func (u UserController) RegisterRoutes(r *chi.Mux) {
	r.Route("/v1/user", func(r chi.Router) {
		r.With(middleware.AuthMiddleware(u.tokenService, u.logger)).Get("/", u.GetAllHandler())
	})
}
