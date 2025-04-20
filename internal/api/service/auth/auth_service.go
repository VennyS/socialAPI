package auth

import (
	"errors"
	"net/http"
	"socialAPI/internal/lib"
	"socialAPI/internal/setting/cfg"
	"socialAPI/internal/shared"
	r "socialAPI/internal/storage/repository"
	"time"

	"gorm.io/gorm"
)

type AuthService interface {
	Authenticate(UserRequest) (string, string, *shared.HttpError)
}

type authService struct {
	userRepo    r.UserRepository
	refreshRepo r.RefreshTokenRepository
	cfg         cfg.AuthConfig
}

func NewAuthService(userRepo r.UserRepository, refreshRepo r.RefreshTokenRepository, cfg cfg.AuthConfig) AuthService {
	return &authService{userRepo: userRepo, refreshRepo: refreshRepo, cfg: cfg}
}

func (a authService) Authenticate(r UserRequest) (string, string, *shared.HttpError) {
	user, err := a.userRepo.FindByEmail(r.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", shared.NewHttpError("user doesnt exits", http.StatusNotFound)
		}
		return "", "", shared.InternalError
	}

	err = lib.ComparePasswords(user.Password, r.Password)
	if err != nil {
		return "", "", shared.NewHttpError("invalid credentials", http.StatusUnauthorized)
	}

	access, refresh, err := lib.GenerateTokenPair(user.ID, a.cfg.AccessTTL, a.cfg.AccessSecret)
	if err != nil {
		return "", "", shared.InternalError
	}

	err = a.refreshRepo.SetRefreshToken(user.ID, refresh, time.Now().Add(a.cfg.RefreshTTL))
	if err != nil {
		return "", "", shared.InternalError
	}

	return access, refresh, nil
}
