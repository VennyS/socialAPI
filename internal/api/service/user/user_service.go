package user

import (
	"socialAPI/internal/shared"
	r "socialAPI/internal/storage/repository"
)

type UserService interface {
	GetAllUsers(excludeID *uint) ([]r.User, *shared.HttpError)
}

type userService struct {
	userRepo r.UserRepository
}

func NewUserService(userRepo r.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s userService) GetAllUsers(excludeID *uint) ([]r.User, *shared.HttpError) {
	users, err := s.userRepo.GetAll(excludeID)
	if err != nil {
		return nil, shared.InternalError
	}

	return users, nil
}
