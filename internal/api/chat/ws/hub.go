package ws

import "go.uber.org/zap"

type Message struct {
	IncomingMessage
	SenderID uint `json:"sender_id"`
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	Register   chan *Client
	unregister chan *Client
	logger     *zap.SugaredLogger
}

func NewHub(logger *zap.SugaredLogger) *Hub {
	return &Hub{
		broadcast:  make(chan Message),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		logger:     logger,
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
				"message", message.IncomingMessage.Text)

			for client := range h.clients {
				client.send <- message
			}
		}
	}
}
