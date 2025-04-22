package user

import (
	"socialAPI/internal/api"
	"socialAPI/internal/api/service/user"
	"socialAPI/internal/shared"

	"github.com/go-chi/chi/v5"
)

type UserController struct {
	userService  user.UserService
	tokenService shared.TokenService
}

func NewAuthController(userService user.UserService, tokenService shared.TokenService) *UserController {
	return &UserController{userService: userService, tokenService: tokenService}
}

func (a UserController) RegisterRoutes(r *chi.Mux) {
	r.Route("/v1/user", func(r chi.Router) {
		r.With(api.AuthMiddleware(&a.tokenService)).Get("/", a.GetAllHandler())
	})
}
