package setting

import (
	"socialAPI/internal/lib"
	"time"
)

func LoadConfig() Config {
	return Config{
		Server: ServerConfig{
			Addr: lib.GetStringFromEnv("ADDR", ":8080"),
		},
		Auth: AuthConfig{
			AccessTTL:    lib.GetDurationFromEnv("ACCESS_TTL", 15*time.Minute),
			RefreshTTL:   lib.GetDurationFromEnv("REFRESH_TTL", 720*time.Hour),
			AccessSecret: lib.GetStringFromEnv("ACCESS_SECRET", "supersecretaccess"),
		},
		DB: DBConfig{
			Host:     lib.GetStringFromEnv("DB_HOST", "localhost"),
			Port:     lib.GetStringFromEnv("DB_PORT", "5433"),
			User:     lib.GetStringFromEnv("DB_USER", "postgres"),
			Password: lib.GetStringFromEnv("DB_PASSWORD", "postgres"),
			Name:     lib.GetStringFromEnv("DB_NAME", "socialdb"),
			SSLMode:  lib.GetStringFromEnv("DB_SSLMODE", "disable"),
		},
	}
}
