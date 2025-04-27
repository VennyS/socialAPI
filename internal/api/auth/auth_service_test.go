package auth_test

import (
	"errors"
	"fmt"
	"socialAPI/internal/api/auth"
	"socialAPI/internal/mocks"
	"socialAPI/internal/setting/cfg"
	"socialAPI/internal/shared"
	"socialAPI/internal/storage/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type authServiceMocks struct {
	userRepo    *mocks.UserRepository
	refreshRepo *mocks.RefreshTokenService
	cacheStore  *mocks.CacheStore
	tokenSvc    *mocks.TokenService
	authSvc     auth.AuthService
	cfg         cfg.AuthConfig
}

func setupAuthService() authServiceMocks {
	userRepo := new(mocks.UserRepository)
	refreshRepo := new(mocks.RefreshTokenService)
	cacheStore := new(mocks.CacheStore)
	tokenSvc := new(mocks.TokenService)
	logger := zap.NewNop().Sugar()

	config := cfg.AuthConfig{
		AccessTTL:  time.Minute * 15,
		RefreshTTL: time.Hour * 24,
	}

	authSvc := auth.NewAuthService(userRepo, refreshRepo, config, cacheStore, tokenSvc, logger)

	return authServiceMocks{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		cacheStore:  cacheStore,
		tokenSvc:    tokenSvc,
		authSvc:     authSvc,
		cfg:         config,
	}
}

var (
	passwordExample  = "1234"
	user             = repository.User{ID: 1, Email: "new@example.com", Password: `$2a$10$HBNNE9kQTwYKgvD08SnePeHwhGInHdvplfVGkKVqv1uvEsKdNzVpO`}
	errExample       = errors.New("example error")
	tokenPairExample = shared.TokenPair{AccessToken: "a", RefreshToken: "b"}
)

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name       string
		email      string
		setupMock  func(m *mocks.UserRepository)
		wantErr    bool
		errMessage string
	}{
		{
			name:  "user already exists",
			email: "test@example.com",
			setupMock: func(m *mocks.UserRepository) {
				m.On("EmailExists", "test@example.com").Return(true, nil)
			},
			wantErr:    true,
			errMessage: "user already exists",
		},
		{
			name:  "error checking user existence",
			email: "error@example.com",
			setupMock: func(m *mocks.UserRepository) {
				m.On("EmailExists", "error@example.com").Return(false, errors.New("db error"))
			},
			wantErr:    true,
			errMessage: shared.InternalError.Error(),
		},
		{
			name:  "successfully register user",
			email: "new@example.com",
			setupMock: func(m *mocks.UserRepository) {
				m.On("EmailExists", "new@example.com").Return(false, nil)
				m.On("Create", mock.AnythingOfType("*repository.User")).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks := setupAuthService()
			tt.setupMock(mocks.userRepo)

			err := mocks.authSvc.Register(auth.UserRequest{
				Email:    tt.email,
				Password: "password",
			})

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errMessage, err.Error())
			} else {
				assert.Nil(t, err)
			}
			mocks.userRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Authenticate(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(m authServiceMocks)
		wantTokens bool
		wantErr    bool
		errMessage string
	}{
		{
			name: "user doesn't exist",
			setup: func(m authServiceMocks) {
				m.userRepo.On("FindByEmail", user.Email).Return(nil, gorm.ErrRecordNotFound)
			},
			wantTokens: false,
			wantErr:    true,

			errMessage: "user doesn't exist",
		},
		{
			name: "error finding user by email",
			setup: func(m authServiceMocks) {
				m.userRepo.On("FindByEmail", user.Email).Return(nil, errExample)
			},
			wantTokens: false,
			wantErr:    true,

			errMessage: shared.InternalError.Error(),
		},
		{
			name: "invalid credentials",
			setup: func(m authServiceMocks) {
				m.userRepo.On("FindByEmail", user.Email).Return(&user, nil)
			},
			wantTokens: false,
			wantErr:    true,
			errMessage: shared.InvalidCredentials.Error(),
		},
		{
			name: "error generating token pair",
			setup: func(m authServiceMocks) {
				m.userRepo.On("FindByEmail", user.Email).Return(&user, nil)
				m.tokenSvc.On("GenerateTokenPair", user.ID).Return(nil, shared.InternalError)
			},
			wantTokens: false,
			wantErr:    true,

			errMessage: shared.InternalError.Error(),
		},
		{
			name: "error storing access token in cache",
			setup: func(m authServiceMocks) {
				m.userRepo.On("FindByEmail", user.Email).Return(&user, nil)
				m.tokenSvc.On("GenerateTokenPair", user.ID).Return(&tokenPairExample, nil)
				m.cacheStore.On("Set", fmt.Sprintf("access_token:%d", user.ID), tokenPairExample.AccessToken, m.cfg.AccessTTL).Return(errExample)
			},
			wantTokens: false,
			wantErr:    true,

			errMessage: shared.InternalError.Error(),
		},
		{
			name: "error storing refresh token",
			setup: func(m authServiceMocks) {
				m.userRepo.On("FindByEmail", user.Email).Return(&user, nil)
				m.tokenSvc.On("GenerateTokenPair", user.ID).Return(&tokenPairExample, nil)
				m.cacheStore.On("Set", fmt.Sprintf("access_token:%d", user.ID), tokenPairExample.AccessToken, m.cfg.AccessTTL).Return(nil)
				m.refreshRepo.On("SetRefreshToken", user.ID, tokenPairExample.RefreshToken, mock.AnythingOfType("time.Time")).Return(errExample)
			},
			wantTokens: false,
			wantErr:    true,

			errMessage: shared.InternalError.Error(),
		},
		{
			name: "user authenticated successfully",
			setup: func(m authServiceMocks) {
				m.userRepo.On("FindByEmail", user.Email).Return(&user, nil)
				m.tokenSvc.On("GenerateTokenPair", user.ID).Return(&tokenPairExample, nil)
				m.cacheStore.On("Set", fmt.Sprintf("access_token:%d", user.ID), tokenPairExample.AccessToken, m.cfg.AccessTTL).Return(nil)
				m.refreshRepo.On("SetRefreshToken", user.ID, tokenPairExample.RefreshToken, mock.AnythingOfType("time.Time")).Return(nil)
			},
			wantTokens: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := setupAuthService()
			tt.setup(m)

			req := auth.UserRequest{
				Email:    user.Email,
				Password: passwordExample,
			}

			if tt.name == "invalid credentials" {
				req.Password = "wrong password"
			}

			tokenPair, err := m.authSvc.Authenticate(req)

			if tt.wantTokens {
				assert.NotNil(t, tokenPair)
				assert.NotEmpty(t, tokenPair.AccessToken)
				assert.NotEmpty(t, tokenPair.RefreshToken)
			} else {
				assert.Nil(t, tokenPair)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errMessage, err.Error())
			} else {
				assert.Nil(t, err)
			}

			m.userRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Refresh(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(m authServiceMocks)
		errMessage string
		wantErr    bool
		wantTokens bool
	}{
		{
			name: "invalid refresh token",
			setup: func(m authServiceMocks) {
				m.refreshRepo.On("GetUserIDIfValid", tokenPairExample.RefreshToken).Return(uint(0), errExample)
			},
			wantErr:    true,
			wantTokens: false,
		},
		{
			name: "error generating token pair",
			setup: func(m authServiceMocks) {
				m.refreshRepo.On("GetUserIDIfValid", tokenPairExample.RefreshToken).Return(uint(1), nil)
				m.tokenSvc.On("GenerateTokenPair", user.ID).Return(nil, shared.InternalError)
			},
			wantTokens: false,
			wantErr:    true,
			errMessage: shared.InternalError.Error(),
		},
		{
			name: "error storing access token in cache",
			setup: func(m authServiceMocks) {
				m.refreshRepo.On("GetUserIDIfValid", tokenPairExample.RefreshToken).Return(uint(1), nil)
				m.tokenSvc.On("GenerateTokenPair", user.ID).Return(&tokenPairExample, nil)
				m.cacheStore.On("Set", fmt.Sprintf("access_token:%d", user.ID), tokenPairExample.AccessToken, m.cfg.AccessTTL).Return(errExample)
			},
			wantTokens: false,
			wantErr:    true,
			errMessage: shared.InternalError.Error(),
		},
		{
			name: "error storing refresh token",
			setup: func(m authServiceMocks) {
				m.refreshRepo.On("GetUserIDIfValid", tokenPairExample.RefreshToken).Return(uint(1), nil)
				m.tokenSvc.On("GenerateTokenPair", user.ID).Return(&tokenPairExample, nil)
				m.cacheStore.On("Set", fmt.Sprintf("access_token:%d", user.ID), tokenPairExample.AccessToken, m.cfg.AccessTTL).Return(nil)
				m.refreshRepo.On("SetRefreshToken", user.ID, tokenPairExample.RefreshToken, mock.AnythingOfType("time.Time")).Return(errExample)
			},
			wantTokens: false,
			wantErr:    true,
			errMessage: shared.InternalError.Error(),
		},
		{
			name: "error storing refresh token",
			setup: func(m authServiceMocks) {
				m.refreshRepo.On("GetUserIDIfValid", tokenPairExample.RefreshToken).Return(uint(1), nil)
				m.tokenSvc.On("GenerateTokenPair", user.ID).Return(&tokenPairExample, nil)
				m.cacheStore.On("Set", fmt.Sprintf("access_token:%d", user.ID), tokenPairExample.AccessToken, m.cfg.AccessTTL).Return(nil)
				m.refreshRepo.On("SetRefreshToken", user.ID, tokenPairExample.RefreshToken, mock.AnythingOfType("time.Time")).Return(errExample)
			},
			wantTokens: false,
			wantErr:    true,
			errMessage: shared.InternalError.Error(),
		},
		{
			name: "refresh token used successfully",
			setup: func(m authServiceMocks) {
				m.refreshRepo.On("GetUserIDIfValid", tokenPairExample.RefreshToken).Return(uint(1), nil)
				m.tokenSvc.On("GenerateTokenPair", user.ID).Return(&tokenPairExample, nil)
				m.cacheStore.On("Set", fmt.Sprintf("access_token:%d", user.ID), tokenPairExample.AccessToken, m.cfg.AccessTTL).Return(nil)
				m.refreshRepo.On("SetRefreshToken", user.ID, tokenPairExample.RefreshToken, mock.AnythingOfType("time.Time")).Return(nil)

			},
			wantTokens: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := setupAuthService()
			tt.setup(m)

			req := auth.RefreshRequest{
				Refresh: tokenPairExample.RefreshToken,
			}

			tokenPair, err := m.authSvc.Refresh(req)

			if tt.wantTokens {
				assert.NotNil(t, tokenPair)
				assert.NotEmpty(t, tokenPair.AccessToken)
				assert.NotEmpty(t, tokenPair.RefreshToken)
			} else {
				assert.Nil(t, tokenPair)
			}

			if tt.wantErr && tt.name != "invalid refresh token" {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errMessage, err.Error())
			} else if tt.wantErr && tt.name == "invalid refresh token" {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			m.userRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Revoke(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(m authServiceMocks)
		errMessage string
		wantErr    bool
	}{
		{
			name: "invalid refresh token",
			setup: func(m authServiceMocks) {
				m.refreshRepo.On("GetUserIDIfValid", tokenPairExample.RefreshToken).Return(uint(0), errExample)
			},
			wantErr: true,
		},
		{
			name: "error revoking refresh token",
			setup: func(m authServiceMocks) {
				m.refreshRepo.On("GetUserIDIfValid", tokenPairExample.RefreshToken).Return(uint(1), nil)
				m.refreshRepo.On("RevokeRefreshToken", tokenPairExample.RefreshToken).Return(errExample)
			},
			wantErr:    true,
			errMessage: shared.InternalError.Error(),
		},
		{
			name: "refresh token revoked and access token removed",
			setup: func(m authServiceMocks) {
				m.refreshRepo.On("GetUserIDIfValid", tokenPairExample.RefreshToken).Return(uint(1), nil)
				m.refreshRepo.On("RevokeRefreshToken", tokenPairExample.RefreshToken).Return(nil)
				m.cacheStore.On("Delete", fmt.Sprintf("access_token:%d", uint(1))).Return(nil)
			},
			wantErr:    false,
			errMessage: shared.InternalError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := setupAuthService()
			tt.setup(m)

			req := auth.RefreshRequest{
				Refresh: tokenPairExample.RefreshToken,
			}

			err := m.authSvc.Revoke(req)

			if tt.wantErr && tt.name != "invalid refresh token" {
				assert.NotNil(t, err)
				assert.Equal(t, tt.errMessage, err.Error())
			} else if tt.wantErr && tt.name == "invalid refresh token" {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			m.userRepo.AssertExpectations(t)
		})
	}
}
