package chat

import (
	"net/http"
	chatWS "socialAPI/internal/api/chat/ws"
	"socialAPI/internal/setting/cfg"
	"socialAPI/internal/shared"
	r "socialAPI/internal/storage/repository"

	"go.uber.org/zap"
)

type ChatService interface {
	GetOne(chatID uint) (*r.ChatDTO, *shared.HttpError)
	GetAll() (*[]r.ChatDTO, *shared.HttpError)
	Create(req CreateRequest) *shared.HttpError
	Update(id uint, req CreateRequest) *shared.HttpError
	HandleWebSocket(userID uint, w http.ResponseWriter, r *http.Request) *shared.HttpError
}

type chatService struct {
	userRepo   r.UserRepository
	chatRepo   r.ChatRepository
	hub        chatWS.Hub
	wsUpgrader cfg.Upgrader
	logger     *zap.SugaredLogger
}

func NewChatService(chatRepo r.ChatRepository, userRepo r.UserRepository, hub chatWS.Hub, wsUpgrader cfg.Upgrader, logger *zap.SugaredLogger) ChatService {
	return &chatService{chatRepo: chatRepo, userRepo: userRepo, hub: hub, wsUpgrader: wsUpgrader, logger: logger}
}

func (c chatService) checksUsersAndChatExistense(req CreateRequest) *shared.HttpError {
	exists, err := c.userRepo.IDsExists(req.UserIDs)
	if err != nil {
		c.logger.Errorw("Error checking user IDs existence", "userIDs", req.UserIDs, "name", req.Name, "error", err)
		return shared.InternalError
	}

	if !exists {
		c.logger.Warnw("Some user IDs do not exist", "userIDs", req.UserIDs, "name", req.Name)
		return shared.NewHttpError("some user IDs do not exist", http.StatusBadRequest)
	}

	exists, err = c.chatRepo.ExistsSetUserIDs(req.UserIDs)

	if err != nil {
		c.logger.Errorw("Error checking chat with this userIDs", "userIDs", req.UserIDs, "name", req.Name, "error", err)
		return shared.InternalError
	}

	if exists {
		c.logger.Warnw("Chat with the same users already exists", "userIDs", req.UserIDs, "name", req.Name)
		return shared.NewHttpError("chat with the same users already exists", http.StatusBadRequest)
	}

	return nil
}

func (c chatService) GetOne(id uint) (*r.ChatDTO, *shared.HttpError) {
	c.logger.Infow("Fetching chat", "chatID", id)

	exists, err := c.chatRepo.ExistsID(id)
	if err != nil {
		c.logger.Errorw("Failed to check chat existence", "chatID", id, "error", err)
		return nil, shared.InternalError
	}

	if !exists {
		c.logger.Warnw("Chat not found", "chatID", id)
		return nil, shared.NewHttpError("chat not found", http.StatusNotFound)
	}

	chat, err := c.chatRepo.GetOne(id)
	if err != nil {
		c.logger.Errorw("Failed to fetch chat", "chatID", id, "error", err)
		return nil, shared.InternalError
	}

	c.logger.Infow("Chat successfully fetched", "chatID", id)

	chatDTO := chat.ConvertToDTO()
	c.logger.Infow("Chat successfully converted to DTO", "chatID", id)

	return chatDTO, nil
}

func (c chatService) GetAll() (*[]r.ChatDTO, *shared.HttpError) {
	c.logger.Infow("Fetching chats")

	chats, err := c.chatRepo.GetAll() // Получаем список чатов
	if err != nil {
		c.logger.Errorw("Failed to fetch chats", "error", err)
		return nil, shared.InternalError
	}

	c.logger.Infow("Chats successfully fetched")

	chatDTOs := []r.ChatDTO{}
	for _, chat := range chats {
		chatDTOs = append(chatDTOs, *chat.ConvertToDTO()) // Конвертируем каждый чат
	}

	c.logger.Infow("Chats successfully converted to DTOs")

	return &chatDTOs, nil
}

func (c chatService) Create(req CreateRequest) *shared.HttpError {
	c.logger.Infow("Attemting to create chat", "userIDs", req.UserIDs, "name", req.Name)
	hErr := c.checksUsersAndChatExistense(req)
	if hErr != nil {
		return hErr
	}

	err := c.chatRepo.Create(req.Name, req.UserIDs)
	if err != nil {
		c.logger.Errorw("Error creating chat", "userIDs", req.UserIDs, "name", req.Name, "error", err)
		return shared.InternalError
	}

	c.logger.Infow("Chat created successfully", "userIDs", req.UserIDs, "name", req.Name)
	return nil
}

func (c chatService) Update(id uint, req CreateRequest) *shared.HttpError {
	c.logger.Infow("Attempting to update chat", "chatID", id, "userIDs", req.UserIDs, "name", req.Name)

	hErr := c.checksUsersAndChatExistense(req)
	if hErr != nil {
		c.logger.Warnw("User validation failed during chat update", "chatID", id, "userIDs", req.UserIDs, "name", req.Name, "error", hErr)
		return hErr
	}

	exists, err := c.chatRepo.ExistsID(id)
	if err != nil {
		c.logger.Errorw("Failed to check chat existence", "chatID", id, "error", err)
		return shared.InternalError
	}

	if !exists {
		c.logger.Warnw("Chat not found", "chatID", id)
		return shared.NewHttpError("chat not found", http.StatusNotFound)
	}

	err = c.chatRepo.Update(id, req.Name, req.UserIDs)
	if err != nil {
		c.logger.Errorw("Error while updating chat", "chatID", id, "userIDs", req.UserIDs, "name", req.Name, "error", err)
		return shared.InternalError
	}

	c.logger.Infow("Chat updated successfully", "chatID", id, "userIDs", req.UserIDs, "name", req.Name)
	return nil
}

func (c chatService) HandleWebSocket(userID uint, w http.ResponseWriter, r *http.Request) *shared.HttpError {
	c.logger.Infow("Handling WebSocket connection", "userID", userID)

	chatIDs, err := c.chatRepo.GetChatIDsByUserID(userID)
	if err != nil {
		c.logger.Errorw("Failed to get chat IDs", "userID", userID, "error", err)
		return shared.InternalError
	}

	conn, err := c.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		c.logger.Errorw("WebSocket upgrade failed", "userID", userID, "error", err)
		return shared.NewHttpError("WebSocket upgrade failed", http.StatusBadRequest)
	}

	c.logger.Infow("WebSocket connection established", "userID", userID)

	chatIDMap := make(map[uint]bool)
	for _, id := range chatIDs {
		chatIDMap[id] = true
	}

	client := chatWS.NewClient(conn, make(chan chatWS.Message, 256), c.hub, userID, chatIDMap, c.logger)

	c.hub.RegisterClient(client)

	c.logger.Infow("Client registered in hub", "userID", userID)

	return nil
}
