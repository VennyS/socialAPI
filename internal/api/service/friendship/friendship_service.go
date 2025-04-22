package friendship

import (
	"errors"
	"net/http"
	"socialAPI/internal/shared"
	r "socialAPI/internal/storage/repository"

	"gorm.io/gorm"
)

type FriendshipService interface {
	SendFriendRequest(senderID, receiverID uint) *shared.HttpError
	GetAllFriends(senderID uint) ([]*r.FriendWithID, *shared.HttpError)
	PatchFriendship(friendshipID uint, request ChangeStatusRequest) *shared.HttpError
	// GetAllPendingRequest(receiverID uint)
}

type friendshipService struct {
	friendshipRepo r.FriendshipRepository
}

func NewFriendshipService(friendshipRepo r.FriendshipRepository) FriendshipService {
	return &friendshipService{friendshipRepo: friendshipRepo}
}

func (f friendshipService) SendFriendRequest(senderID, receiverID uint) *shared.HttpError {
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

func (f friendshipService) GetAllFriends(senderID uint) ([]*r.FriendWithID, *shared.HttpError) {
	users, err := f.friendshipRepo.GetAllFriends(senderID)

	if err != nil {
		return nil, shared.InternalError
	}

	return users, nil
}

func (f friendshipService) PatchFriendship(friendshipID uint, request ChangeStatusRequest) *shared.HttpError {
	err := f.friendshipRepo.SetStatus(friendshipID, request.Status)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return shared.NewHttpError("Friend request not found or not yours", http.StatusNotFound)
		}
		return shared.InternalError
	}

	return nil
}
