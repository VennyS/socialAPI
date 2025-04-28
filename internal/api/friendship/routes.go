package friendship

import (
	"socialAPI/internal/api/middleware"
	"socialAPI/internal/shared"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type FriendshipController struct {
	friendshipService FriendshipService
	tokenService      shared.TokenService
	logger            *zap.SugaredLogger
}

func NewFriendshipController(friendshipService FriendshipService, tokenService shared.TokenService, logger *zap.SugaredLogger) *FriendshipController {
	return &FriendshipController{friendshipService: friendshipService, tokenService: tokenService, logger: logger}
}

func (f FriendshipController) RegisterRoutes(r *chi.Mux) {
	r.Route("/v1/friendship", func(r chi.Router) {
		r.With(middleware.AuthMiddleware(f.tokenService, f.logger), middleware.JsonBodyMiddleware[FriendshipPostRequest](f.logger)).Post("/", f.SendRequestHandler())
		r.With(middleware.AuthMiddleware(f.tokenService, f.logger)).Get("/", f.GetFriendsHandler())
		r.With(middleware.AuthMiddleware(f.tokenService, f.logger), middleware.JsonBodyMiddleware[ChangeStatusRequest](f.logger)).Patch("/{id}", f.PutStatusHandler())
	})
}
