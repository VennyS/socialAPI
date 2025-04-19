package storage

import (
	"fmt"
	"log"
	repo "socialAPI/internal/storage/repository"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

func BootstrapDatabase(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error database connection: %v", err)
	}

	err = db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'friendship_status') THEN
				CREATE TYPE friendship_status AS ENUM ('pending', 'rejected', 'friendship');
			END IF;
		END
		$$;
	`).Error
	if err != nil {
		log.Fatalf("Error creating enum type: %v", err)
	}

	if err := db.AutoMigrate(&repo.User{}, &repo.Chat{}, &repo.Message{}, &repo.Friendship{}, &repo.RefreshToken{}); err != nil {
		log.Fatalf("Migrations went wrong: %v", err)
	}

	fmt.Println("Migration success")
	return db
}
