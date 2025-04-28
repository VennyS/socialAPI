package ws

import (
	r "socialAPI/internal/storage/repository"

	"go.uber.org/zap"
)

// Типы команд для хаба
type HubCommandType int

const (
	CommandRegister HubCommandType = iota
	CommandUnregister
	CommandBroadcast
)

// Структура команды
type HubCommand struct {
	Type    HubCommandType
	Client  *Client
	Message Message
}

// Интерфейс Hub
type Hub interface {
	Run()
	RegisterClient(client *Client)
	UnregisterClient(client *Client)
	BroadcastMessage(msg Message)
}

// Структура сообщения
type Message struct {
	IncomingMessage
	SenderID uint `json:"sender_id"`
}

// Реализация Hub
type hub struct {
	clients     map[*Client]bool
	commands    chan HubCommand
	messageRepo r.MessageRepository
	chatRepo    r.ChatRepository
	logger      *zap.SugaredLogger
}

// Конструктор
func NewHub(messageRepo r.MessageRepository, chatRepo r.ChatRepository, logger *zap.SugaredLogger) Hub {
	return &hub{
		clients:     make(map[*Client]bool),
		commands:    make(chan HubCommand),
		messageRepo: messageRepo,
		chatRepo:    chatRepo,
		logger:      logger,
	}
}

// Запуск хаба
func (h *hub) Run() {
	for cmd := range h.commands {
		switch cmd.Type {
		case CommandRegister:
			h.handleRegister(cmd.Client)
		case CommandUnregister:
			h.handleUnregister(cmd.Client)
		case CommandBroadcast:
			h.handleBroadcast(cmd.Message)
		}
	}
}

// Регистрация клиента
func (h *hub) RegisterClient(client *Client) {
	h.commands <- HubCommand{Type: CommandRegister, Client: client}
}

// Удаление клиента
func (h *hub) UnregisterClient(client *Client) {
	h.commands <- HubCommand{Type: CommandUnregister, Client: client}
}

// Отправка сообщения
func (h *hub) BroadcastMessage(msg Message) {
	h.commands <- HubCommand{Type: CommandBroadcast, Message: msg}
}

// Внутренние обработчики:

func (h *hub) handleRegister(client *Client) {
	h.logger.Infow("Registering new client", "clientID", client.userID)

	h.clients[client] = true

	go client.ReadPump()
	go client.WritePump()
}

func (h *hub) handleUnregister(client *Client) {
	h.logger.Infow("Unregistering client", "clientID", client.userID)

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
	}
}

func (h *hub) handleBroadcast(msg Message) {
	h.logger.Infow("Broadcasting message",
		"senderID", msg.SenderID,
		"chatID", msg.ChatID,
		"content", msg.Content)

	exists, err := h.chatRepo.ExistsID(msg.ChatID)
	if err != nil {
		h.logger.Errorw("Error checking chat existence",
			"chatID", msg.ChatID,
			"error", err)
		return
	}

	if !exists {
		h.logger.Warnw("Chat does not exist",
			"chatID", msg.ChatID)
		return
	}

	err = h.messageRepo.Create(msg.ChatID, msg.SenderID, msg.Content)
	if err != nil {
		h.logger.Errorw("Error creating message",
			"chatID", msg.ChatID,
			"error", err)
		return
	}

	for client := range h.clients {
		select {
		case client.send <- msg:
		default:
			h.handleUnregister(client)
		}
	}
}
