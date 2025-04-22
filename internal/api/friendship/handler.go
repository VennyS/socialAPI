package friendship

import (
	"net/http"
	"socialAPI/internal/api"
	"socialAPI/internal/lib"
	"strconv"

	"github.com/go-chi/chi/v5"
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

		hErr := f.friendshipService.SendFriendRequst(senderID, receiverID)
		if hErr != nil {
			lib.SendMessage(w, r, hErr.StatusCode, hErr.Error())
			return
		}

		lib.SendMessage(w, r, http.StatusOK, "sent succefully")
	}
}
