package repository

import (
	"gorm.io/gorm"
)

type FriendWithID struct {
	Friend       *User `json:"friend"`
	FriendshipID uint  `json:"friendship_id"`
}

type FriendshipRepository interface {
	SendRequest(friendship *Friendship) error
	GetAllFriends(userID uint) ([]*FriendWithID, error)
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

func (repo friendshipPostgresRepo) GetAllFriends(userID uint) ([]*FriendWithID, error) {
	var friendships []Friendship

	// Получаем все дружбы, где userID участвует и статус accepted
	err := repo.db.Preload("Sender").Preload("Receiver").
		Where("(sender_id = ? OR receiver_id = ?) AND status = ?", userID, userID, StatusFriendship).
		Find(&friendships).Error
	if err != nil {
		return nil, err
	}

	var friends []*FriendWithID
	for _, f := range friendships {
		var friend *User
		var friendshipID uint

		if f.SenderID == userID {
			friend = &f.Receiver
		} else {
			friend = &f.Sender
		}

		friendshipID = f.ID // ID дружбы

		friends = append(friends, &FriendWithID{
			Friend:       friend,
			FriendshipID: friendshipID,
		})
	}

	return friends, nil
}
