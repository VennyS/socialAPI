package repository

import (
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	EmailExists(email string) (bool, error)
	GetAll(excludeID *uint) ([]User, error)
	IDsExists(IDs []uint) (bool, error)
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

func (repo userPostgresRepo) EmailExists(email string) (bool, error) {
	var count int64
	err := repo.db.Model(&User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo userPostgresRepo) GetAll(excludeID *uint) ([]User, error) {
	var users []User

	// Если excludeID не nil, исключаем пользователя с этим id из выборки
	if excludeID != nil {
		err := repo.db.Where("id != ?", *excludeID).Find(&users).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := repo.db.Find(&users).Error
		if err != nil {
			return nil, err
		}
	}

	return users, nil
}

func (repo userPostgresRepo) IDsExists(IDs []uint) (bool, error) {
	var count int64
	err := repo.db.Model(&User{}).Where("id IN ?", IDs).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == int64(len(IDs)), nil
}
