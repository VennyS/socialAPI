package repository

import "gorm.io/gorm"

type Repository interface {
	Users() UserRepository
	RefreshTokens() RefreshTokenService
	// Chats() ChatRepository
	// Messages() MessageRepository
	// Notifications() NotificationRepository
}

type postgresRepo struct {
	users         UserRepository
	refreshTokens RefreshTokenService
	// chats         ChatRepository
	// messages      MessageRepository
	// notifications NotificationRepository
}

func NewPostgresRepo(db *gorm.DB) Repository {
	return &postgresRepo{
		users:         NewPostgresUserRepo(db),
		refreshTokens: NewPostgresRefreshtokenService(db),
		// chats:         NewPostgresChatRepo(db),
		// messages:      NewPostgresMessageRepo(db),
		// notifications: NewPostgresNotificationRepo(db),
	}
}

func (r *postgresRepo) Users() UserRepository {
	return r.users
}

func (r *postgresRepo) RefreshTokens() RefreshTokenService {
	return r.refreshTokens
}

// func (r *PostgresRepo) Chats() ChatRepository {
// 	return r.chats
// }

// func (r *PostgresRepo) Messages() MessageRepository {
// 	return r.messages
// }

// func (r *PostgresRepo) Notifications() NotificationRepository {
// 	return r.notifications
// }
