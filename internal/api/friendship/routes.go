package friendship

import (
	"socialAPI/internal/api"
	"socialAPI/internal/api/service/friendship"
	"socialAPI/internal/shared"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type FriendshipController struct {
	friendshipService friendship.FriendshipService
	tokenService      shared.TokenService
	logger            *zap.SugaredLogger
}

func NewFriendshipController(friendshipService friendship.FriendshipService, tokenService shared.TokenService, logger *zap.SugaredLogger) *FriendshipController {
	return &FriendshipController{friendshipService: friendshipService, tokenService: tokenService, logger: logger}
}

func (f FriendshipController) RegisterRoutes(r *chi.Mux) {
	r.Route("/v1/friendship", func(r chi.Router) {
		r.With(api.AuthMiddleware(&f.tokenService, f.logger), api.JsonBodyMiddleware[friendship.FriendshipPostRequest](f.logger)).Post("/", f.SendRequestHandler())
		r.With(api.AuthMiddleware(&f.tokenService, f.logger)).Get("/", f.GetFriendsHandler())
		r.With(api.AuthMiddleware(&f.tokenService, f.logger), api.JsonBodyMiddleware[friendship.ChangeStatusRequest](f.logger)).Patch("/{id}", f.PutStatusHandler())
	})
}
