package shared

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenService interface {
	GenerateTokenPair(userID uint) (*TokenPair, error)
	GenerateAccessToken(userID uint) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

// TokenService handles JWT token generation and validation
type jwtTokenService struct {
	accessSecret   string
	accessTokenTTL time.Duration
}

// NewTokenService creates a new TokenService instance
func NewTokenService(accessSecret string, accessTokenTTL time.Duration) TokenService {
	return &jwtTokenService{
		accessSecret:   accessSecret,
		accessTokenTTL: accessTokenTTL,
	}
}

// TokenPair contains both access and refresh tokens
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// Claims represents the JWT claims structure
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// GenerateTokenPair generates both access and refresh tokens
func (ts *jwtTokenService) GenerateTokenPair(userID uint) (*TokenPair, error) {
	accessToken, err := ts.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := ts.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// generateAccessToken creates a new JWT access token
func (ts *jwtTokenService) GenerateAccessToken(userID uint) (string, error) {
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ts.accessTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(ts.accessSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// generateRefreshToken creates a secure random refresh token
func (ts *jwtTokenService) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// ValidateToken validates the JWT token and returns claims
func (ts *jwtTokenService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(ts.accessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
