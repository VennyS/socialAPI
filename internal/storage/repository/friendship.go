package repository

import "gorm.io/gorm"

type FriendshipRepository interface {
	SendRequest(friendship *Friendship) error
}

type friendshipPostgresRepo struct {
	db *gorm.DB
}

func NewPostgresFriendshipRepo(db *gorm.DB) FriendshipRepository {
	return friendshipPostgresRepo{db: db}
}

func (repo friendshipPostgresRepo) SendRequest(friendship *Friendship) error {
	if err := repo.db.Create(friendship).Error; err != nil {
		return err
	}

	return nil
}
