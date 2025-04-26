package ws

import (
	r "socialAPI/internal/storage/repository"

	"go.uber.org/zap"
)

type Message struct {
	IncomingMessage
	SenderID uint `json:"sender_id"`
}

type Hub struct {
	clients     map[*Client]bool
	broadcast   chan Message
	Register    chan *Client
	unregister  chan *Client
	messageRepo r.MessageRepository
	chatRepo    r.ChatRepository
	logger      *zap.SugaredLogger
}

func NewHub(messageRepo r.MessageRepository, chatRepo r.ChatRepository, logger *zap.SugaredLogger) *Hub {
	return &Hub{
		broadcast:   make(chan Message),
		Register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		messageRepo: messageRepo,
		chatRepo:    chatRepo,
		logger:      logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.logger.Infow("Registering new client",
				"clientID", client.userID)

			h.clients[client] = true
			go client.ReadPump()
			go client.WritePump()

		case client := <-h.unregister:
			h.logger.Infow("Unregistering client",
				"clientID", client.userID)

			delete(h.clients, client)
			close(client.send)

		case message := <-h.broadcast:
			h.logger.Infow("Broadcasting message",
				"senderID", message.SenderID,
				"chatID", message.IncomingMessage.ChatID,
				"content", message.IncomingMessage.Content)

			exists, err := h.chatRepo.ExistsID(message.ChatID)
			if !exists {
				h.logger.Warnw("Chat doesnt exists", "senderID", message.SenderID,
					"chatID", message.IncomingMessage.ChatID,
					"content", message.IncomingMessage.Content)
			}

			if err != nil {
				h.logger.Errorw("Error checking existence chatID", "senderID", message.SenderID,
					"chatID", message.IncomingMessage.ChatID,
					"content", message.IncomingMessage.Content, "error", err)
			}

			err = h.messageRepo.Create(message.ChatID, message.SenderID, message.Content)
			if err != nil {
				h.logger.Errorw("Error creating message", "senderID", message.SenderID,
					"chatID", message.IncomingMessage.ChatID,
					"content", message.IncomingMessage.Content, "error", err)
				continue
			}
			for client := range h.clients {
				client.send <- message
			}
		}
	}
}
