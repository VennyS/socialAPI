package storage

import (
	"fmt"
	repo "socialAPI/internal/storage/repository"

	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"

	"gorm.io/gorm"
)

func BootstrapDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func MadeMigrations(db *gorm.DB) {
	err := db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'friendship_status') THEN
				CREATE TYPE friendship_status AS ENUM ('pending', 'rejected', 'friendship');
			END IF;
		END
		$$;
	`).Error
	if err != nil {
		panic(fmt.Sprintf("Error creating enum type: %v", err))
	}

	if err := db.AutoMigrate(&repo.User{}, &repo.Chat{}, &repo.Message{}, &repo.Friendship{}, &repo.RefreshToken{}); err != nil {
		panic(fmt.Sprintf("Migrations went wrong: %v", err))
	}

	fmt.Println("Migration success")
}
