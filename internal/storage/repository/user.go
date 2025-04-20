package repository

import "gorm.io/gorm"

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	Exists(email string) (bool, error)
}

type userPostgresRepo struct {
	db *gorm.DB
}

func NewPostgresUserRepo(db *gorm.DB) UserRepository {
	return userPostgresRepo{db: db}
}

func (repo userPostgresRepo) FindByEmail(email string) (*User, error) {
	var user User
	if err := repo.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo userPostgresRepo) Create(user *User) error {
	if err := repo.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
func (repo userPostgresRepo) Exists(email string) (bool, error) {
	var count int64
	err := repo.db.Model(&User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
