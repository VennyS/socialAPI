package chat

import (
	"socialAPI/internal/api"
	"socialAPI/internal/api/service/chat"
	"socialAPI/internal/shared"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ChatController struct {
	chatService  chat.ChatService
	tokenService shared.TokenService
	logger       *zap.SugaredLogger
}

func NewChatController(chatService chat.ChatService, tokenService shared.TokenService, logger *zap.SugaredLogger) *ChatController {
	return &ChatController{chatService: chatService, tokenService: tokenService, logger: logger}
}

func (c ChatController) RegisterRoutes(r *chi.Mux) {
	r.Route("/v1/chat", func(r chi.Router) {
		r.With(api.AuthMiddleware(&c.tokenService, c.logger), api.JsonBodyMiddleware[chat.CreateRequest](c.logger)).Post("/", c.CreateHandler())
		r.With(api.AuthMiddleware(&c.tokenService, c.logger), api.JsonBodyMiddleware[chat.CreateRequest](c.logger)).Patch("/{id}", c.UpdateHandler())
	})
}
