package chat

import (
	"net/http"
	"socialAPI/internal/api/middleware"
	"socialAPI/internal/lib"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (c ChatController) GetOneHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatIDParam := chi.URLParam(r, "id")
		chatIDUint64, err := strconv.ParseUint(chatIDParam, 10, 32)
		if err != nil {
			c.logger.Warnw("Invalid chat ID parameter", "chatID", chatIDParam, "error", err.Error())
			lib.SendMessage(w, r, http.StatusBadRequest, "Invalid id parameter")
			return
		}
		chatID := uint(chatIDUint64)

		c.logger.Infow("Handling GetOne request", "chatID", chatID)

		chat, hErr := c.chatService.GetOne(chatID)
		if hErr != nil {
			c.logger.Errorw("Failed to get chat", "chatID", chatID, "error", err)
			lib.SendMessage(w, r, hErr.StatusCode, hErr.Error())
			return
		}

		c.logger.Infow("Chat successfully retrieved", "chatID", chatID)

		render.Status(r, http.StatusOK)
		render.JSON(w, r, chat)
	}
}

func (c ChatController) GetAllHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.logger.Infow("Handling GetAll request")

		chats, err := c.chatService.GetAll()
		if err != nil {
			c.logger.Errorw("Failed to get chats", "error", err)
			lib.SendMessage(w, r, err.StatusCode, err.Error())
			return
		}

		c.logger.Infow("Chats successfully retrieved")

		render.Status(r, http.StatusOK)
		render.JSON(w, r, chats)
	}
}

func (c ChatController) CreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := r.Context().Value(middleware.DataKey).(CreateRequest)

		c.logger.Infow("Create chat", "userIDs", req.UserIDs, "name", req.Name)
		err := c.chatService.Create(req)
		if err != nil {
			c.logger.Errorw("Error while creating chat", "userIDs", req.UserIDs, "name", req.Name, "error", err)
			lib.SendMessage(w, r, err.StatusCode, err.Error())
			return
		}

		c.logger.Infow("response chat successfully created sent", "userIDs", req.UserIDs, "name", req.Name)
		lib.SendMessage(w, r, http.StatusCreated, "Chat successfully created")
	}
}

func (c ChatController) UpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := r.Context().Value(middleware.DataKey).(CreateRequest)
		chatID := r.Context().Value(middleware.UserIDKey).(uint)

		c.logger.Infow("Update chat", "chatID", chatID, "userIDs", req.UserIDs, "name", req.Name)
		err := c.chatService.Update(chatID, req)
		if err != nil {
			c.logger.Errorw("Error while updating chat", "chatID", chatID, "userIDs", req.UserIDs, "name", req.Name, "error", err)
			lib.SendMessage(w, r, err.StatusCode, err.Error())
			return
		}

		c.logger.Infow("response chat changed successfully sent", "chatID", chatID, "userIDs", req.UserIDs, "name", req.Name)
		lib.SendMessage(w, r, http.StatusOK, "Chat changed successfully")
	}
}

func (c ChatController) BroadcastHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		senderID := r.Context().Value(middleware.UserIDKey).(uint)
		err := c.chatService.HandleWebSocket(senderID, w, r)

		if err != nil {
			lib.SendMessage(w, r, err.StatusCode, err.Error())
			return
		}
	}
}
