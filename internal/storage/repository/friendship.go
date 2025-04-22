package repository

import "gorm.io/gorm"

type FriendshipRepository interface {
	SendRequest(friendship *Friendship) error
	GetAllFriends(userID uint) ([]*User, error)
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

func (repo friendshipPostgresRepo) GetAllFriends(userID uint) ([]*User, error) {
	var friendships []Friendship

	// Получаем все дружбы, где userID участвует и статус accepted
	err := repo.db.Preload("Sender").Preload("Receiver").
		Where("(sender_id = ? OR receiver_id = ?) AND status = ?", userID, userID, StatusFriendship).
		Find(&friendships).Error
	if err != nil {
		return nil, err
	}

	friends := []*User{}
	for _, f := range friendships {
		if f.SenderID == userID {
			friends = append(friends, &f.Receiver)
		} else {
			friends = append(friends, &f.Sender)
		}
	}

	return friends, nil
}
