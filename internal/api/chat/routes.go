package chat

import (
	"socialAPI/internal/api/middleware"
	"socialAPI/internal/shared"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ChatController struct {
	chatService  ChatService
	tokenService shared.TokenService
	logger       *zap.SugaredLogger
}

func NewChatController(chatService ChatService, tokenService shared.TokenService, logger *zap.SugaredLogger) *ChatController {
	return &ChatController{chatService: chatService, tokenService: tokenService, logger: logger}
}

func (c ChatController) RegisterRoutes(r *chi.Mux) {
	r.Route("/v1/chat", func(r chi.Router) {
		r.With(middleware.AuthMiddleware(&c.tokenService, c.logger)).Get("/", c.GetAllHandler())
		r.With(middleware.AuthMiddleware(&c.tokenService, c.logger)).Get("/ws", c.BroadcastHandler())
		r.With(middleware.AuthMiddleware(&c.tokenService, c.logger)).Get("/{id}", c.GetOneHandler())
		r.With(middleware.AuthMiddleware(&c.tokenService, c.logger), middleware.JsonBodyMiddleware[CreateRequest](c.logger)).Post("/", c.CreateHandler())
		r.With(middleware.AuthMiddleware(&c.tokenService, c.logger), middleware.JsonBodyMiddleware[CreateRequest](c.logger)).Patch("/{id}", c.UpdateHandler())
	})
}
