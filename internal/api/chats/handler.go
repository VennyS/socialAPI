package chat

import (
	"net/http"
	"socialAPI/internal/api"
	"socialAPI/internal/api/service/chat"
	"socialAPI/internal/lib"
)

func (c ChatController) CreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := r.Context().Value(api.DataKey).(chat.CreateRequest)

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
		req := r.Context().Value(api.DataKey).(chat.CreateRequest)
		chatID := r.Context().Value(api.UserIDKey).(uint)

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
