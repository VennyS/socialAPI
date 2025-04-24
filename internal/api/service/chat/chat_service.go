package chat

import (
	"net/http"
	"socialAPI/internal/shared"
	r "socialAPI/internal/storage/repository"

	"go.uber.org/zap"
)

type ChatService interface {
	Create(CreateRequest) *shared.HttpError
}

type chatService struct {
	userRepo r.UserRepository
	chatRepo r.ChatRepository
	logger   *zap.SugaredLogger
}

func NewChatService(chatRepo r.ChatRepository, userRepo r.UserRepository, logger *zap.SugaredLogger) ChatService {
	return &chatService{chatRepo: chatRepo, userRepo: userRepo, logger: logger}
}

func (c chatService) Create(req CreateRequest) *shared.HttpError {
	c.logger.Infow("Attemting to create chat", "userIDs", req.UserIDs, "name", req.Name)
	exists, err := c.userRepo.IDsExists(req.UserIDs)
	if err != nil {
		c.logger.Errorw("Error checking user IDs existence", "userIDs", req.UserIDs, "name", req.Name, "error", err)
		return shared.InternalError
	}

	if !exists {
		c.logger.Warnw("Some user IDs do not exist", "userIDs", req.UserIDs, "name", req.Name)
		return shared.NewHttpError("some user IDs do not exist", http.StatusBadRequest)
	}

	exists, err = c.chatRepo.Exists(req.UserIDs)

	if exists {
		c.logger.Warnw("Chat with the same users already exists", "userIDs", req.UserIDs, "name", req.Name)
		return shared.NewHttpError("chat with the same users already exists", http.StatusBadRequest)
	}

	if err != nil {
		c.logger.Errorw("Error checking chat with this userIDs", "userIDs", req.UserIDs, "name", req.Name, "error", err)
		return shared.InternalError
	}

	err = c.chatRepo.Create(req.Name, req.UserIDs)
	if err != nil {
		c.logger.Errorw("Error creating chat", "userIDs", req.UserIDs, "name", req.Name, "error", err)
		return shared.InternalError
	}

	c.logger.Infow("Chat successfully created", "userIDs", req.UserIDs, "name", req.Name)
	return nil
}
