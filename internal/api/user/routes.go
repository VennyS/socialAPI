package user

import (
	"socialAPI/internal/api"
	"socialAPI/internal/api/service/user"
	"socialAPI/internal/shared"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UserController struct {
	userService  user.UserService
	tokenService shared.TokenService
	logger       *zap.SugaredLogger
}

func NewAuthController(userService user.UserService, tokenService shared.TokenService, logger *zap.SugaredLogger) *UserController {
	return &UserController{userService: userService, tokenService: tokenService, logger: logger}
}

func (u UserController) RegisterRoutes(r *chi.Mux) {
	r.Route("/v1/user", func(r chi.Router) {
		r.With(api.AuthMiddleware(&u.tokenService, u.logger)).Get("/", u.GetAllHandler())
	})
}
