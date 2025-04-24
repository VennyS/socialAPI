package repository

import (
	"fmt"

	"gorm.io/gorm"
)

type ChatRepository interface {
	Create(name *string, userIDs []uint) error
	Exists(userIDS []uint) (bool, error)
}

type chatPostgresRepo struct {
	db *gorm.DB
}

func NewPostgresChatRepo(db *gorm.DB) ChatRepository {
	return chatPostgresRepo{db: db}
}

func (repo chatPostgresRepo) Create(name *string, userIDs []uint) error {
	chat := Chat{
		Name:  "",
		Users: make([]User, len(userIDs)),
	}

	if name != nil {
		chat.Name = *name
	}

	for i, id := range userIDs {
		chat.Users[i] = User{ID: id}
	}

	if err := repo.db.Create(&chat).Error; err != nil {
		return err
	}

	return nil
}

func (repo chatPostgresRepo) Exists(userIDs []uint) (bool, error) {
	var chatCount int64

	err := repo.db.
		Table("user_chats").
		Select("chat_id").
		Where("user_id IN ?", userIDs).
		Group("chat_id").
		Having(
			"COUNT(user_id) = ? AND COUNT(user_id) = (SELECT COUNT(*) FROM user_chats WHERE chat_id = user_chats.chat_id)",
			len(userIDs),
		).Count(&chatCount).Error

	if err != nil {
		return true, err
	}

	if chatCount == 64 {
		return false, nil
	}

	return true, fmt.Errorf("chat with same userIDs already exists")
}
