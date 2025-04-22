package friendship

import (
	"errors"
	"net/http"
	"socialAPI/internal/shared"
	r "socialAPI/internal/storage/repository"

	"gorm.io/gorm"
)

type FriendshipService interface {
	SendFriendRequst(senderID, receiverID uint) *shared.HttpError
	// GetAllPendingRequest(receiverID uint)
}

type friendshipService struct {
	friendshipRepo r.FriendshipRepository
}

func NewFriendshipService(friendshipRepo r.FriendshipRepository) FriendshipService {
	return &friendshipService{friendshipRepo: friendshipRepo}
}

func (f friendshipService) SendFriendRequst(senderID, receiverID uint) *shared.HttpError {
	if senderID == receiverID {
		return shared.NewHttpError("Cannot send friend request to yourself", http.StatusBadRequest)
	}

	friendShip := r.Friendship{SenderID: senderID, ReceiverID: receiverID, Status: r.StatusPending}

	err := f.friendshipRepo.SendRequest(&friendShip)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return shared.NewHttpError("friendship already exitst", http.StatusConflict)
		}

		return shared.InternalError
	}

	return nil
}
