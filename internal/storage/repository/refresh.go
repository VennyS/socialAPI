package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RefreshTokenRepository interface {
	SetRefreshToken(userID uint, token string, expiresAt time.Time) error
	FindByUserID(userID uint) (*RefreshToken, error)
	GetUserIDIfValid(token string) (uint, error)
	RevokeRefreshToken(userID uint) error
	UpdateRefreshToken(userID uint, token string, expiresAt time.Time) error
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

	// Выполняем upsert: если запись существует, обновляем её, если нет — создаём
	if err := repo.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},                        // Указываем, по какому столбцу проверяется уникальность
		DoUpdates: clause.AssignmentColumns([]string{"token", "expires_at"}), // Обновляем указанные столбцы
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
		return 0, errors.New("token revoked")
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return 0, errors.New("token expired")
	}

	return refreshToken.UserID, nil
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
