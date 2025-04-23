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
	// Логируем начало запроса
	s.logger.Info("Получение всех пользователей", "excludeID", excludeID)

	// Выполняем запрос на получение пользователей
	users, err := s.userRepo.GetAll(excludeID)
	if err != nil {
		// Логируем ошибку, если что-то пошло не так
		s.logger.Errorw("Ошибка при получении пользователей", "excludeID", excludeID, "error", err)
		return nil, shared.InternalError
	}

	// Логируем успешное завершение
	s.logger.Info("Пользователи успешно получены", "excludeID", excludeID, "userCount", len(users))

	return users, nil
}
