package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RefreshTokenService interface {
	SetRefreshToken(userID uint, token string, expiresAt time.Time) error
	FindByUserID(userID uint) (*RefreshToken, error)
	GetUserIDIfValid(token string) (uint, error)
	RevokeRefreshToken(token string) error
}

type refreshTokenPostgresRepo struct {
	db *gorm.DB
}

func NewPostgresRefreshtokenService(db *gorm.DB) RefreshTokenService {
	return refreshTokenPostgresRepo{db: db}
}

func (repo refreshTokenPostgresRepo) SetRefreshToken(userID uint, token string, expiresAt time.Time) error {
	refreshToken := RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		Revoked:   false,
	}

	// Выполняем upsert: если запись существует, обновляем её, если нет — создаём
	if err := repo.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},                                   // Указываем, по какому столбцу проверяется уникальность
		DoUpdates: clause.AssignmentColumns([]string{"token", "expires_at", "revoked"}), // Обновляем указанные столбцы
	}).Create(&refreshToken).Error; err != nil {
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

func (repo refreshTokenPostgresRepo) GetUserIDIfValid(token string) (uint, error) {
	var refreshToken RefreshToken
	err := repo.db.Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		return 0, err
	}

	if refreshToken.Revoked {
		return 0, errors.New("refresh token revoked")
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return 0, errors.New("refresh token expired")
	}

	return refreshToken.UserID, nil
}

func (repo refreshTokenPostgresRepo) RevokeRefreshToken(token string) error {
	result := repo.db.Model(&RefreshToken{}).
		Where("token = ?", token).
		Update("revoked", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("token not found")
	}

	return nil
}

func (repo refreshTokenPostgresRepo) UpdateRefreshToken(userID uint, token string, expiresAt time.Time) error {
	result := repo.db.Model(&RefreshToken{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"token":      token,
			"expires_at": expiresAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("token not found")
	}

	return nil
}
