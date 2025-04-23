package friendship

import (
	"socialAPI/internal/api"
	"socialAPI/internal/api/service/friendship"
	"socialAPI/internal/shared"

	"github.com/go-chi/chi/v5"
)

type FriendshipController struct {
	friendshipService friendship.FriendshipService
	tokenService      shared.TokenService
}

func NewFriendshipController(friendshipService friendship.FriendshipService, tokenService shared.TokenService) *FriendshipController {
	return &FriendshipController{friendshipService: friendshipService, tokenService: tokenService}
}

func (f FriendshipController) RegisterRoutes(r *chi.Mux) {
	r.Route("/v1/friendship", func(r chi.Router) {
		r.With(api.AuthMiddleware(&f.tokenService), api.JsonBodyMiddleware[friendship.FriendshipPostRequest]()).Post("/", f.SendRequestHandler())
		r.With(api.AuthMiddleware(&f.tokenService)).Get("/", f.GetFriendsHandler())
		r.With(api.AuthMiddleware(&f.tokenService), api.JsonBodyMiddleware[friendship.ChangeStatusRequest]()).Patch("/{id}", f.PutStatusHandler())
	})
}
