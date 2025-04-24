package chat

import (
	"net/http"
	"socialAPI/internal/shared"
	r "socialAPI/internal/storage/repository"

	"go.uber.org/zap"
)

type ChatService interface {
	Create(req CreateRequest) *shared.HttpError
	Update(id uint, req CreateRequest) *shared.HttpError
}

type chatService struct {
	userRepo r.UserRepository
	chatRepo r.ChatRepository
	logger   *zap.SugaredLogger
}

func NewChatService(chatRepo r.ChatRepository, userRepo r.UserRepository, logger *zap.SugaredLogger) ChatService {
	return &chatService{chatRepo: chatRepo, userRepo: userRepo, logger: logger}
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

	if exists {
		c.logger.Warnw("Chat with the same users already exists", "userIDs", req.UserIDs, "name", req.Name)
		return shared.NewHttpError("chat with the same users already exists", http.StatusBadRequest)
	}

	if err != nil {
		c.logger.Errorw("Error checking chat with this userIDs", "userIDs", req.UserIDs, "name", req.Name, "error", err)
		return shared.InternalError
	}

	return nil
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
