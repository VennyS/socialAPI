package auth

import (
	"errors"
	"fmt"
	"net/http"
	"socialAPI/internal/lib"
	"socialAPI/internal/setting/cfg"
	"socialAPI/internal/shared"
	"socialAPI/internal/storage/cache"
	"socialAPI/internal/storage/repository"
	r "socialAPI/internal/storage/repository"
	"time"

	"gorm.io/gorm"
)

type AuthService interface {
	Authenticate(UserRequest) (*shared.TokenPair, *shared.HttpError)
	Register(r UserRequest) *shared.HttpError
	Refresh(r RefreshRequest) (*shared.TokenPair, *shared.HttpError)
	Revoke(r RefreshRequest) *shared.HttpError
}

type authService struct {
	userRepo     r.UserRepository
	refreshRepo  r.RefreshTokenService
	cfg          cfg.AuthConfig
	cache        cache.CacheStore
	tokenService shared.TokenService
}

func NewAuthService(userRepo r.UserRepository, refreshRepo r.RefreshTokenService, cfg cfg.AuthConfig, cache cache.CacheStore, tokenService shared.TokenService) AuthService {
	return &authService{userRepo: userRepo, refreshRepo: refreshRepo, cfg: cfg, cache: cache, tokenService: tokenService}
}

func (a authService) generateAndStoreTokens(id uint) (*shared.TokenPair, *shared.HttpError) {
	tokenPair, err := a.tokenService.GenerateTokenPair(id)
	if err != nil {
		return nil, shared.InternalError
	}

	// Проверяем установку access токена в Redis
	err = a.cache.Set(fmt.Sprintf("access_token:%d", id), tokenPair.AccessToken, a.cfg.AccessTTL)
	if err != nil {
		return nil, shared.InternalError
	}

	err = a.refreshRepo.SetRefreshToken(id, tokenPair.RefreshToken, time.Now().Add(a.cfg.RefreshTTL))
	if err != nil {
		return nil, shared.InternalError
	}

	return tokenPair, nil
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

func (a authService) Authenticate(r UserRequest) (*shared.TokenPair, *shared.HttpError) {
	user, err := a.userRepo.FindByEmail(r.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.NewHttpError("user doesnt exits", http.StatusNotFound)
		}
		return nil, shared.InternalError
	}

	err = lib.ComparePasswords(user.Password, r.Password)
	if err != nil {
		return nil, shared.InvalidCredentials
	}

	tokenPair, hErr := a.generateAndStoreTokens(user.ID)
	if hErr != nil {
		return tokenPair, shared.InternalError
	}

	return tokenPair, nil
}

func (a authService) Refresh(r RefreshRequest) (*shared.TokenPair, *shared.HttpError) {
	userID, err := a.refreshRepo.GetUserIDIfValid(r.Refresh)
	if err != nil {
		return nil, shared.NewHttpError(err.Error(), http.StatusUnauthorized)
	}

	tokenPair, hErr := a.generateAndStoreTokens(userID)
	if hErr != nil {
		return nil, shared.InternalError
	}

	return tokenPair, nil
}

func (a authService) Revoke(r RefreshRequest) *shared.HttpError {
	userID, err := a.refreshRepo.GetUserIDIfValid(r.Refresh)
	if err != nil {
		return shared.NewHttpError(err.Error(), http.StatusUnauthorized)
	}

	err = a.refreshRepo.RevokeRefreshToken(r.Refresh)
	if err != nil {
		return shared.InternalError
	}

	err = a.cache.Delete(fmt.Sprintf("access_token:%d", userID))
	if err != nil {
		return shared.InternalError
	}

	return nil
}
