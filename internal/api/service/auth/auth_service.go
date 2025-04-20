package auth

import (
	"errors"
	"net/http"
	"socialAPI/internal/lib"
	"socialAPI/internal/setting/cfg"
	"socialAPI/internal/shared"
	"socialAPI/internal/storage/repository"
	r "socialAPI/internal/storage/repository"
	"time"

	"gorm.io/gorm"
)

type AuthService interface {
	Authenticate(UserRequest) (string, string, *shared.HttpError)
	Register(r UserRequest) *shared.HttpError
	Refresh(r RefreshRequest) (string, string, *shared.HttpError)
	Revoke(r RefreshRequest) *shared.HttpError
}

type authService struct {
	userRepo    r.UserRepository
	refreshRepo r.RefreshTokenRepository
	cfg         cfg.AuthConfig
}

func NewAuthService(userRepo r.UserRepository, refreshRepo r.RefreshTokenRepository, cfg cfg.AuthConfig) AuthService {
	return &authService{userRepo: userRepo, refreshRepo: refreshRepo, cfg: cfg}
}

func (a authService) Register(r UserRequest) *shared.HttpError {
	exists, err := a.userRepo.Exists(r.Email)
	if err != nil {
		return shared.InternalError
	}

	if exists {
		return shared.NewHttpError("user already exits", http.StatusNotFound)
	}

	hashedPassword, err := lib.HashPassword(r.Password)

	if err != nil {
		return shared.InternalError
	}

	err = a.userRepo.Create(&repository.User{Email: r.Email, Password: hashedPassword})

	if err != nil {
		return shared.InternalError
	}

	return nil
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

func (a authService) Refresh(r RefreshRequest) (string, string, *shared.HttpError) {
	userID, err := a.refreshRepo.GetUserIDIfValid(r.Refresh)
	if err != nil {
		return "", "", shared.NewHttpError(err.Error(), http.StatusUnauthorized)
	}

	access, newRefresh, err := lib.GenerateTokenPair(userID, a.cfg.AccessTTL, a.cfg.AccessSecret)
	if err != nil {
		return "", "", shared.InternalError
	}

	err = a.refreshRepo.UpdateRefreshToken(userID, newRefresh, time.Now().Add(a.cfg.RefreshTTL))
	if err != nil {
		return "", "", shared.InternalError
	}

	return access, newRefresh, nil
}

func (a authService) Revoke(r RefreshRequest) *shared.HttpError {
	err := a.refreshRepo.RevokeRefreshToken(r.Refresh)

	if err != nil {
		return shared.InternalError
	}

	return nil
}
