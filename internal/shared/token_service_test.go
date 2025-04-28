package shared_test

import (
	"socialAPI/internal/shared"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

type tokenServiceCfg struct {
	secret string
	ttl    time.Duration
}

var (
	correctCFG = tokenServiceCfg{secret: "correct_secret", ttl: 15 * time.Minute}
	expiredCfg = tokenServiceCfg{
		secret: "correct_secret",
		ttl:    -15 * time.Minute,
	}
	noSecretCFG  = tokenServiceCfg{secret: "", ttl: 15 * time.Minute}
	invalidToken = "invalid_token"
)

func TestTokenService_GenerateAccessToken(t *testing.T) {
	tests := []struct {
		name              string
		cfg               tokenServiceCfg
		userID            uint
		wantGenErr        bool
		genErrMessage     string
		wantValidationErr bool
		validationErrMsg  string
	}{
		{
			name:              "successful token generation",
			cfg:               correctCFG,
			userID:            1,
			wantGenErr:        false,
			wantValidationErr: false,
		},
		{
			name:              "empty secret - should fail",
			cfg:               noSecretCFG,
			userID:            1,
			wantGenErr:        true,
			genErrMessage:     "signing key is empty",
			wantValidationErr: false,
		},
		{
			name:              "zero user ID - should work",
			cfg:               correctCFG,
			userID:            0,
			wantGenErr:        false,
			wantValidationErr: false,
		},
		{
			name:              "max uint user ID - should work",
			cfg:               correctCFG,
			userID:            ^uint(0),
			wantGenErr:        false,
			wantValidationErr: false,
		},
		{
			name:              "expired token config - should still generate token",
			cfg:               expiredCfg,
			userID:            1,
			wantGenErr:        false,
			wantValidationErr: true,
			validationErrMsg:  "token is expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := shared.NewTokenService(tt.cfg.secret, tt.cfg.ttl)

			token, err := ts.GenerateAccessToken(tt.userID)

			if tt.wantGenErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.genErrMessage)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			claims, valErr := ts.ValidateToken(token)

			if tt.wantValidationErr {
				assert.Error(t, valErr)
				assert.Contains(t, valErr.Error(), tt.validationErrMsg)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, valErr)
				assert.NotNil(t, claims)
				assert.Equal(t, tt.userID, claims.UserID)
			}
		})
	}
}

func TestTokenService_GenerateRefreshToken(t *testing.T) {
	tests := []struct {
		name       string
		wantGenErr bool
	}{
		{
			name:       "successful refresh token generation",
			wantGenErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := shared.NewTokenService(correctCFG.secret, correctCFG.ttl)

			token, err := ts.GenerateRefreshToken()

			if tt.wantGenErr {
				assert.Error(t, err)
				assert.Empty(t, token)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, token)
			assert.Len(t, token, 64)
		})
	}
}

func TestTokenService_ValidateToken(t *testing.T) {
	tests := []struct {
		name              string
		tokenString       string
		wantValidationErr bool
		validationErrMsg  string
	}{
		{
			name:              "valid token",
			wantValidationErr: false,
		},
		{
			name:              "token contains an invalid number of segments",
			tokenString:       invalidToken,
			wantValidationErr: true,
			validationErrMsg:  "token contains an invalid number of segments",
		},
		{
			name:              "token with unexpected signing method",
			tokenString:       generateTokenWithUnexpectedMethod(),
			wantValidationErr: true,
			validationErrMsg:  "signature is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := shared.NewTokenService(correctCFG.secret, correctCFG.ttl)

			var (
				claims *shared.Claims
				err    error
			)

			if tt.name == "valid token" {
				token, _ := ts.GenerateAccessToken(1)
				claims, err = ts.ValidateToken(token)
			} else {
				claims, err = ts.ValidateToken(tt.tokenString)
			}

			if tt.wantValidationErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.validationErrMsg)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
			}
		})
	}
}

func TestTokenService_GenerateTokenPair(t *testing.T) {
	cfg := struct {
		secret string
		ttl    time.Duration
	}{
		secret: "test_secret",
		ttl:    time.Hour,
	}

	tests := []struct {
		name           string
		userID         uint
		wantErr        bool
		wantErrMessage string
	}{
		{
			name:    "успешная генерация пары токенов",
			userID:  1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := shared.NewTokenService(cfg.secret, cfg.ttl)

			tokenPair, err := ts.GenerateTokenPair(tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrMessage)
				assert.Nil(t, tokenPair)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tokenPair)
				assert.NotEmpty(t, tokenPair.AccessToken)
				assert.NotEmpty(t, tokenPair.RefreshToken)

				_, err := ts.ValidateToken(tokenPair.AccessToken)
				assert.NoError(t, err)
			}
		})
	}
}

// Функция для генерации токена с неожиданным методом подписи
func generateTokenWithUnexpectedMethod() string {
	// Генерация токена с неправильным методом подписи
	claims := &shared.Claims{
		UserID: 1,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims) // Используется неправильный метод подписи
	tokenString, _ := token.SignedString([]byte("testSecret"))
	return tokenString
}
