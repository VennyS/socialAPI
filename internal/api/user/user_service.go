package user

import (
	"socialAPI/internal/shared"
	r "socialAPI/internal/storage/repository"

	"go.uber.org/zap"
)

type UserService interface {
	GetAllUsers(excludeID *uint) ([]r.User, *shared.HttpError)
}

type userService struct {
	userRepo r.UserRepository
	logger   *zap.SugaredLogger
}

func NewUserService(userRepo r.UserRepository, logger *zap.SugaredLogger) UserService {
	return &userService{userRepo: userRepo, logger: logger}
}

func (s userService) GetAllUsers(excludeID *uint) ([]r.User, *shared.HttpError) {
	s.logger.Infow("Fetching all users", "excludeID", excludeID)

	users, err := s.userRepo.GetAll(excludeID)
	if err != nil {
		s.logger.Errorw("Failed to fetch users", "excludeID", excludeID, "error", err)
		return nil, shared.InternalError
	}

	s.logger.Infow("Users successfully fetched", "excludeID", excludeID, "userCount", len(users))

	return users, nil
}
