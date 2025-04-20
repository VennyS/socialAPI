package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	SetRefreshToken(userID uint, token string, expiresAt time.Time) error
	FindByUserID(userID uint) (*RefreshToken, error)
	IsTokenValid(token string) error
	RevokeRefreshToken(userID uint) error
	UpdateRefreshToken(userID uint, token string) error
}

type refreshTokenPostgresRepo struct {
	db *gorm.DB
}

func NewPostgresRefreshTokenRepo(db *gorm.DB) RefreshTokenRepository {
	return refreshTokenPostgresRepo{db: db}
}

func (repo refreshTokenPostgresRepo) SetRefreshToken(userID uint, token string, expiresAt time.Time) error {
	refreshToken := RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	if err := repo.db.Create(&refreshToken).Error; err != nil {
		return err
	}

	return nil
}

func (repo refreshTokenPostgresRepo) FindByUserID(userID uint) (*RefreshToken, error) {
	var refreshToken RefreshToken
	err := repo.db.Where("UserID = ?", userID).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (repo refreshTokenPostgresRepo) IsTokenValid(token string) error {
	var refreshToken RefreshToken
	err := repo.db.Where("Token = ?", token).First(&refreshToken).Error
	if err != nil {
		return err
	}

	if refreshToken.Revoked {
		return errors.New("token revoked")
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return errors.New("token expired")
	}

	return nil
}

func (repo refreshTokenPostgresRepo) RevokeRefreshToken(userID uint) error {
	result := repo.db.Model(&RefreshToken{}).
		Where("UserID = ?", userID).
		Update("Revoked", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("token not found")
	}

	return nil
}

func (repo refreshTokenPostgresRepo) UpdateRefreshToken(userID uint, token string) error {
	result := repo.db.Model(&RefreshToken{}).
		Where("UserID = ?", userID).
		Update("Token", token)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("token not found")
	}

	return nil
}
