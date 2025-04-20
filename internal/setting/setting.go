package setting

import (
	"fmt"
	au "socialAPI/internal/api/auth"
	srv "socialAPI/internal/api/service"
	"socialAPI/internal/api/service/auth"
	"socialAPI/internal/lib"
	"socialAPI/internal/setting/cfg"
	"socialAPI/internal/storage"
	"socialAPI/internal/storage/repository"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type App struct {
	cfg     cfg.Config
	db      *gorm.DB
	service srv.Service
}

func (a *App) LoadConfig() {
	a.cfg = cfg.Config{
		Server: cfg.ServerConfig{
			Addr: lib.GetStringFromEnv("ADDR", ":8080"),
		},
		Auth: cfg.AuthConfig{
			AccessTTL:    lib.GetDurationFromEnv("ACCESS_TTL", 15*time.Minute),
			RefreshTTL:   lib.GetDurationFromEnv("REFRESH_TTL", 720*time.Hour),
			AccessSecret: lib.GetStringFromEnv("ACCESS_SECRET", "supersecretaccess"),
		},
		DB: cfg.DBConfig{
			Host:     lib.GetStringFromEnv("DB_HOST", "localhost"),
			Port:     lib.GetStringFromEnv("DB_PORT", "5433"),
			User:     lib.GetStringFromEnv("DB_USER", "postgres"),
			Password: lib.GetStringFromEnv("DB_PASSWORD", "postgres"),
			Name:     lib.GetStringFromEnv("DB_NAME", "socialdb"),
			SSLMode:  lib.GetStringFromEnv("DB_SSLMODE", "disable"),
		},
	}
}

func (a *App) ConnectDB() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		a.cfg.DB.Host,
		a.cfg.DB.Port,
		a.cfg.DB.User,
		a.cfg.DB.Password,
		a.cfg.DB.Name,
		a.cfg.DB.SSLMode,
	)

	a.db = storage.BootstrapDatabase(dsn)
}

func (a *App) MountServices() {
	repo := repository.NewPostgresRepo(a.db)

	authService := auth.NewAuthService(repo.Users(), repo.RefreshTokens(), a.cfg.Auth)

	a.service = srv.NewService(authService)
}

func (a App) MountRouter() *chi.Mux {
	authController := au.NewAuthController(a.service.Auth())
	r := chi.NewRouter()
	authController.RegisterRoutes(r)

	return r
}

// TODO
func (a App) RunServer() {
}
