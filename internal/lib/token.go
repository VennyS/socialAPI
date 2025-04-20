package lib

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateTokenPair(userID uint, accessTokenTTL time.Duration, accessSecret string) (string, string, error) {
	access, err := GenerateToken(userID, accessTokenTTL, accessSecret)
	if err != nil {
		return "", "", nil
	}

	refresh, err := GenerateRefreshToken()
	if err != nil {
		return "", "", nil
	}

	return access, refresh, err
}

func GenerateToken(userID uint, accessTokenTTL time.Duration, accessSecret string) (string, error) {
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(accessTokenTTL).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccess, err := accessToken.SignedString(accessSecret)
	if err != nil {
		return "", err
	}

	return signedAccess, err
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
