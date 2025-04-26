package repository

import "gorm.io/gorm"

type MessageRepository interface {
	Create(chatID, senderID uint, content string) error
}

type messagePostgresRepo struct {
	db *gorm.DB
}

func NewPostgresMessageRepo(db *gorm.DB) MessageRepository {
	return messagePostgresRepo{db: db}
}

func (repo messagePostgresRepo) Create(chatID, senderID uint, content string) error {
	message := Message{ChatID: chatID, SenderID: senderID, Content: content}

	return repo.db.Create(message).Error
}
