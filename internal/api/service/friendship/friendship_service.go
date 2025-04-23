package friendship

import (
	"errors"
	"net/http"
	"socialAPI/internal/lib"
	"socialAPI/internal/shared"
	r "socialAPI/internal/storage/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FriendshipService interface {
	SendFriendRequest(senderID, receiverID uint) *shared.HttpError
	GetAllFriends(senderID uint, statusParam string) ([]*r.FriendWithID, *shared.HttpError)
	PatchFriendship(friendshipID uint, request ChangeStatusRequest) *shared.HttpError
	// GetAllPendingRequest(receiverID uint)
}

type friendshipService struct {
	friendshipRepo r.FriendshipRepository
	logger         *zap.SugaredLogger
}

func NewFriendshipService(friendshipRepo r.FriendshipRepository, logger *zap.SugaredLogger) FriendshipService {
	return &friendshipService{friendshipRepo: friendshipRepo, logger: logger}
}

func (f friendshipService) SendFriendRequest(senderID, receiverID uint) *shared.HttpError {
	f.logger.Infow("Attempting to send friend request", "senderID", senderID, "receiverID", receiverID)

	if senderID == receiverID {
		f.logger.Warnw("Attempt to send a friend request to oneself", "senderID", senderID)
		return shared.NewHttpError("Cannot send friend request to yourself", http.StatusBadRequest)
	}

	exists, err := f.friendshipRepo.Exists(senderID, receiverID)
	if err != nil {
		f.logger.Errorw("Error checking if friendship exists", "senderID", senderID, "receiverID", receiverID, "error", err)
		return shared.InternalError
	}

	if exists {
		f.logger.Warnw("Friend request already exists", "senderID", senderID, "receiverID", receiverID)
		return shared.NewHttpError("friendship exists", http.StatusConflict)
	}

	friendShip := r.Friendship{SenderID: senderID, ReceiverID: receiverID, Status: r.StatusPending}
	err = f.friendshipRepo.SendRequest(&friendShip)
	if err != nil {
		f.logger.Errorw("Error sending friend request", "senderID", senderID, "receiverID", receiverID, "error", err)
		return shared.InternalError
	}

	f.logger.Infow("Friend request successfully sent", "senderID", senderID, "receiverID", receiverID)

	return nil
}

func (f friendshipService) GetAllFriends(senderID uint, statusParam string) ([]*r.FriendWithID, *shared.HttpError) {
	f.logger.Infow("Getting list of friends", "senderID", senderID, "statusParam", statusParam)

	var statusPtr *r.FriendshipStatus
	if statusParam != "" {
		status := r.FriendshipStatus(statusParam)
		statusPtr = &status
		if !lib.IsValidStatus(*statusPtr) {
			f.logger.Warnw("Invalid status for request", "statusParam", statusParam)
			return nil, shared.NewHttpError("incorrect status", http.StatusBadRequest)
		}
	}

	users, err := f.friendshipRepo.GetAllFriends(senderID, statusPtr)
	if err != nil {
		f.logger.Errorw("Error getting list of friends", "senderID", senderID, "statusParam", statusParam, "error", err)
		return nil, shared.InternalError
	}

	f.logger.Infow("Successfully retrieved list of friends", "senderID", senderID, "userCount", len(users))

	return users, nil
}

func (f friendshipService) PatchFriendship(friendshipID uint, request ChangeStatusRequest) *shared.HttpError {
	f.logger.Infow("Updating friendship status", "friendshipID", friendshipID, "status", request.Status)

	err := f.friendshipRepo.SetStatus(friendshipID, request.Status)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			f.logger.Warnw("Friend request not found or not yours", "friendshipID", friendshipID)
			return shared.NewHttpError("Friend request not found or not yours", http.StatusNotFound)
		}
		f.logger.Errorw("Error updating friendship status", "friendshipID", friendshipID, "status", request.Status, "error", err)
		return shared.InternalError
	}

	f.logger.Infow("Friendship status successfully updated", "friendshipID", friendshipID, "status", request.Status)

	return nil
}
