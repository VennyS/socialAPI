package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ChatDTO - это структура для отправки данных о чате без лишней информации
type ChatDTO struct {
	ID        uint         `json:"id"`
	Messages  []MessageDTO `json:"messages,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// MessageDTO - это структура для отправки данных о сообщении без лишней информации
type MessageDTO struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	Sender    SenderDTO `json:"sender"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SenderDTO - структура для отправки данных о пользователе (отправителе)
type SenderDTO struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

// ConvertToDTO преобразует Chat в ChatDTO
func (chat *Chat) ConvertToDTO() *ChatDTO {
	var messageDTOs []MessageDTO
	for _, message := range chat.Messages {
		messageDTOs = append(messageDTOs, MessageDTO{
			ID:        message.ID,
			Content:   message.Content,
			Sender:    message.Sender.ConvertToDTO(), // Преобразуем отправителя
			CreatedAt: message.CreatedAt,
			UpdatedAt: message.UpdatedAt,
		})
	}

	return &ChatDTO{
		ID:        chat.ID,
		Messages:  messageDTOs,
		CreatedAt: chat.CreatedAt,
		UpdatedAt: chat.UpdatedAt,
	}
}

func (user *User) ConvertToDTO() SenderDTO {
	return SenderDTO{
		ID:    user.ID,
		Email: user.Email,
	}
}

type ChatRepository interface {
	GetOne(chatID uint) (*Chat, error)
	GetAll() ([]*Chat, error)
	Create(name *string, userIDs []uint) error
	ExistsSetUserIDs(userIDs []uint) (bool, error)
	ExistsID(chatID uint) (bool, error)
	Update(id uint, name *string, userIDs []uint) error
}

type chatPostgresRepo struct {
	db *gorm.DB
}

func NewPostgresChatRepo(db *gorm.DB) ChatRepository {
	return chatPostgresRepo{db: db}
}

func (repo chatPostgresRepo) GetOne(chatID uint) (*Chat, error) {
	var chat *Chat
	err := repo.db.Preload("Messages.Sender").First(&chat, chatID).Error
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (repo chatPostgresRepo) GetAll() ([]*Chat, error) {
	var chats []*Chat
	err := repo.db.Preload("Messages.Sender").Find(&chats).Error
	if err != nil {
		return nil, err
	}
	return chats, nil
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

func (repo chatPostgresRepo) ExistsSetUserIDs(userIDs []uint) (bool, error) {
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

	if chatCount == 0 {
		return false, nil
	}

	return true, fmt.Errorf("chat with same userIDs already exists")
}

func (repo chatPostgresRepo) ExistsID(chatID uint) (bool, error) {
	var count int64
	err := repo.db.
		Model(&Chat{}).
		Where("id = ?", chatID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo chatPostgresRepo) Update(id uint, name *string, userIDs []uint) error {
	var chat Chat
	if err := repo.db.Preload("Users").First(&chat, id).Error; err != nil {
		return err
	}

	if name != nil {
		chat.Name = *name
	}

	if err := repo.db.Model(&chat).Update("name", chat.Name).Error; err != nil {
		return err
	}

	users := make([]User, len(userIDs))
	for i, id := range userIDs {
		users[i] = User{ID: id}
	}

	if err := repo.db.Model(&chat).Association("Users").Replace(users); err != nil {
		return err
	}

	return nil
}
