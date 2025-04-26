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
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthService interface {
	Authenticate(UserRequest) (*shared.TokenPair, *shared.HttpError)
	Register(r UserRequest) *shared.HttpError
	Refresh(r RefreshRequest) (*shared.TokenPair, *shared.HttpError)
	Revoke(r RefreshRequest) *shared.HttpError
}

type authService struct {
	userRepo     repository.UserRepository
	refreshRepo  repository.RefreshTokenService
	cfg          cfg.AuthConfig
	cache        cache.CacheStore
	tokenService shared.TokenService
	logger       *zap.SugaredLogger
}

func NewAuthService(userRepo repository.UserRepository, refreshRepo repository.RefreshTokenService, cfg cfg.AuthConfig, cache cache.CacheStore, tokenService shared.TokenService, logger *zap.SugaredLogger) AuthService {
	return &authService{userRepo: userRepo, refreshRepo: refreshRepo, cfg: cfg, cache: cache, tokenService: tokenService, logger: logger}
}

func (a authService) generateAndStoreTokens(id uint) (*shared.TokenPair, *shared.HttpError) {
	tokenPair, err := a.tokenService.GenerateTokenPair(id)
	if err != nil {
		a.logger.Errorw("Error generating token pair", "error", err)
		return nil, shared.InternalError
	}

	err = a.cache.Set(fmt.Sprintf("access_token:%d", id), tokenPair.AccessToken, a.cfg.AccessTTL)
	if err != nil {
		a.logger.Errorw("Error storing access token in cache", "error", err)
		return nil, shared.InternalError
	}

	err = a.refreshRepo.SetRefreshToken(id, tokenPair.RefreshToken, time.Now().Add(a.cfg.RefreshTTL))
	if err != nil {
		a.logger.Errorw("Error storing refresh token", "error", err)
		return nil, shared.InternalError
	}

	a.logger.Infow("Tokens generated and stored", "userID", id)
	return tokenPair, nil
}

func (a authService) Register(r UserRequest) *shared.HttpError {
	exists, err := a.userRepo.EmailExists(r.Email)
	if err != nil {
		a.logger.Errorw("Error checking if user exists", "email", r.Email, "error", err)
		return shared.InternalError
	}

	if exists {
		a.logger.Warnw("User already exists", "email", r.Email)
		return shared.NewHttpError("user already exists", http.StatusNotFound)
	}

	hashedPassword, err := lib.HashPassword(r.Password)
	if err != nil {
		a.logger.Errorw("Error hashing password", "error", err)
		return shared.InternalError
	}

	err = a.userRepo.Create(&repository.User{Email: r.Email, Password: hashedPassword})
	if err != nil {
		a.logger.Errorw("Error creating new user", "email", r.Email, "error", err)
		return shared.InternalError
	}

	a.logger.Infow("User successfully registered", "email", r.Email)
	return nil
}

func (a authService) Authenticate(r UserRequest) (*shared.TokenPair, *shared.HttpError) {
	user, err := a.userRepo.FindByEmail(r.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			a.logger.Warnw("User not found", "email", r.Email)
			return nil, shared.NewHttpError("user doesn't exist", http.StatusNotFound)
		}
		a.logger.Errorw("Error finding user by email", "email", r.Email, "error", err)
		return nil, shared.InternalError
	}

	err = lib.ComparePasswords(user.Password, r.Password)
	if err != nil {
		a.logger.Warnw("Invalid credentials", "email", r.Email)
		return nil, shared.InvalidCredentials
	}

	tokenPair, hErr := a.generateAndStoreTokens(user.ID)
	if hErr != nil {
		a.logger.Errorw("Error generating and storing tokens", "userID", user.ID)
		return tokenPair, shared.InternalError
	}

	a.logger.Infow("User authenticated successfully", "email", r.Email)
	return tokenPair, nil
}

func (a authService) Refresh(r RefreshRequest) (*shared.TokenPair, *shared.HttpError) {
	userID, err := a.refreshRepo.GetUserIDIfValid(r.Refresh)
	if err != nil {
		a.logger.Warnw("Invalid refresh token", "refreshToken", r.Refresh)
		return nil, shared.NewHttpError(err.Error(), http.StatusUnauthorized)
	}

	tokenPair, hErr := a.generateAndStoreTokens(userID)
	if hErr != nil {
		a.logger.Errorw("Error generating and storing tokens during refresh", "userID", userID)
		return nil, shared.InternalError
	}

	a.logger.Infow("Refresh token used successfully", "userID", userID)
	return tokenPair, nil
}

func (a authService) Revoke(r RefreshRequest) *shared.HttpError {
	userID, err := a.refreshRepo.GetUserIDIfValid(r.Refresh)
	if err != nil {
		a.logger.Warnw("Invalid refresh token during revoke", "refreshToken", r.Refresh)
		return shared.NewHttpError(err.Error(), http.StatusUnauthorized)
	}

	err = a.refreshRepo.RevokeRefreshToken(r.Refresh)
	if err != nil {
		a.logger.Errorw("Error revoking refresh token", "refreshToken", r.Refresh, "error", err)
		return shared.InternalError
	}

	err = a.cache.Delete(fmt.Sprintf("access_token:%d", userID))
	if err != nil {
		a.logger.Errorw("Error deleting access token from cache", "userID", userID, "error", err)
		return shared.InternalError
	}

	a.logger.Infow("Refresh token revoked and access token removed", "userID", userID)
	return nil
}
