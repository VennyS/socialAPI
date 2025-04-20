package repository

import "gorm.io/gorm"

type Repository interface {
	Users() UserRepository
	RefreshTokens() RefreshTokenRepository
	// Chats() ChatRepository
	// Messages() MessageRepository
	// Notifications() NotificationRepository
}

type postgresRepo struct {
	users         UserRepository
	refreshTokens RefreshTokenRepository
	// chats         ChatRepository
	// messages      MessageRepository
	// notifications NotificationRepository
}

func NewPostgresRepo(db *gorm.DB) Repository {
	return &postgresRepo{
		users:         NewPostgresUserRepo(db),
		refreshTokens: NewPostgresRefreshTokenRepo(db),
		// chats:         NewPostgresChatRepo(db),
		// messages:      NewPostgresMessageRepo(db),
		// notifications: NewPostgresNotificationRepo(db),
	}
}

func (r *postgresRepo) Users() UserRepository {
	return r.users
}

func (r *postgresRepo) RefreshTokens() RefreshTokenRepository {
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
