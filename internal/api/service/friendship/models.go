package friendship

import "socialAPI/internal/storage/repository"

type ChangeStatusRequest struct {
	Status repository.FriendshipStatus `json:"status" validate:"required"`
}

type FriendshipPostRequest struct {
	ReceiverID uint `json:"receiver_id" validate:"required"`
}
