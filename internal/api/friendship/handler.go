package friendship

import (
	"net/http"
	"socialAPI/internal/api"
	"socialAPI/internal/api/service/friendship"
	"socialAPI/internal/lib"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (f FriendshipController) SendRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		receiverIDParam := chi.URLParam(r, "id")

		receiverIDUint64, err := strconv.ParseUint(receiverIDParam, 10, 32)
		if err != nil {
			lib.SendMessage(w, r, http.StatusBadRequest, "Invalid id parameter")
			return
		}
		receiverID := uint(receiverIDUint64)

		senderID := r.Context().Value(api.UserIDKey).(uint)

		hErr := f.friendshipService.SendFriendRequest(senderID, receiverID)
		if hErr != nil {
			lib.SendMessage(w, r, hErr.StatusCode, hErr.Error())
			return
		}

		lib.SendMessage(w, r, http.StatusOK, "sent succefully")
	}
}

func (f FriendshipController) GetFriendsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		senderID := r.Context().Value(api.UserIDKey).(uint)

		users, err := f.friendshipService.GetAllFriends(senderID)
		if err != nil {
			lib.SendMessage(w, r, err.StatusCode, err.Error())
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, users)
	}
}

func (f FriendshipController) PutStatusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatIDParam := chi.URLParam(r, "id")

		chatIDUint64, err := strconv.ParseUint(chatIDParam, 10, 32)
		if err != nil {
			lib.SendMessage(w, r, http.StatusBadRequest, "Invalid id parameter")
			return
		}
		chatID := uint(chatIDUint64)

		req := r.Context().Value(api.DataKey).(friendship.ChangeStatusRequest)

		hErr := f.friendshipService.PatchFriendship(chatID, req)
		if hErr != nil {
			lib.SendMessage(w, r, hErr.StatusCode, hErr.Error())
		}

		lib.SendMessage(w, r, http.StatusOK, "changed succefully")
	}
}
