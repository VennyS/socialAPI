package repository

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `json:"-"`

	Chats    []Chat    `json:"chats,omitempty" gorm:"many2many:user_chats;"`
	Messages []Message `json:"messages,omitempty" gorm:"foreignKey:SenderID"`
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Token     string    `gorm:"unique;not null" json:"token"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Revoked   bool      `gorm:"default:false" json:"revoked"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserID"`
}

type Chat struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Users     []User    `gorm:"many2many:user_chats;" json:"users,omitempty"`
	Messages  []Message `gorm:"foreignKey:ChatID" json:"messages,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ChatID    uint      `gorm:"not null" json:"chat_id"`
	SenderID  uint      `gorm:"not null" json:"sender_id"`
	Content   string    `gorm:"not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Chat   Chat `gorm:"foreignKey:ChatID" json:"chat,omitempty"`
	Sender User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Message   string    `gorm:"not null" json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
