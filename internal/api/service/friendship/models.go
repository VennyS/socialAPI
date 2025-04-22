package friendship

import "socialAPI/internal/storage/repository"

type ChangeStatusRequest struct {
	Status repository.FriendshipStatus `json:"status" validate:"required"`
}
