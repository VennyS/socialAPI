package friendship

import (
	"net/http"
	"socialAPI/internal/api/middleware"
	"socialAPI/internal/lib"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (f FriendshipController) SendRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := r.Context().Value(middleware.DataKey).(FriendshipPostRequest)
		senderID := r.Context().Value(middleware.UserIDKey).(uint)

		f.logger.Infow("Send friend request", "senderID", senderID, "receiverID", req.ReceiverID)

		hErr := f.friendshipService.SendFriendRequest(senderID, req.ReceiverID)
		if hErr != nil {
			f.logger.Warnw("Failed to send friend request", "senderID", senderID, "receiverID", req.ReceiverID, "error", hErr.Error())
			lib.SendMessage(w, r, hErr.StatusCode, hErr.Error())
			return
		}

		f.logger.Infow("Friend request sent successfully", "senderID", senderID, "receiverID", req.ReceiverID)
		lib.SendMessage(w, r, http.StatusOK, "sent successfully")
	}
}

func (f FriendshipController) GetFriendsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		senderID := r.Context().Value(middleware.UserIDKey).(uint)
		statusParam := r.URL.Query().Get("status")

		f.logger.Infow("Get friends request", "senderID", senderID, "status", statusParam)

		users, err := f.friendshipService.GetAllFriends(senderID, statusParam)
		if err != nil {
			f.logger.Warnw("Failed to retrieve friends", "senderID", senderID, "status", statusParam, "error", err.Error())
			lib.SendMessage(w, r, err.StatusCode, err.Error())
			return
		}

		f.logger.Infow("Successfully retrieved friends", "senderID", senderID, "status", statusParam)
		render.Status(r, http.StatusOK)
		render.JSON(w, r, users)
	}
}

func (f FriendshipController) PutStatusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatIDParam := chi.URLParam(r, "id")

		chatIDUint64, err := strconv.ParseUint(chatIDParam, 10, 32)
		if err != nil {
			f.logger.Warnw("Invalid chat ID parameter", "chatID", chatIDParam, "error", err.Error())
			lib.SendMessage(w, r, http.StatusBadRequest, "Invalid id parameter")
			return
		}
		chatID := uint(chatIDUint64)

		req := r.Context().Value(middleware.DataKey).(ChangeStatusRequest)

		f.logger.Infow("Change friendship status", "chatID", chatID, "status", req.Status)

		hErr := f.friendshipService.PatchFriendship(chatID, req)
		if hErr != nil {
			f.logger.Warnw("Failed to change friendship status", "chatID", chatID, "error", hErr.Error())
			lib.SendMessage(w, r, hErr.StatusCode, hErr.Error())
			return
		}

		f.logger.Infow("Friendship status changed successfully", "chatID", chatID, "status", req.Status)
		lib.SendMessage(w, r, http.StatusOK, "changed successfully")
	}
}
