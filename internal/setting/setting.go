package setting

import (
	"fmt"
	au "socialAPI/internal/api/auth"
	"socialAPI/internal/api/friendship"
	srv "socialAPI/internal/api/service"
	"socialAPI/internal/api/service/auth"
	frSrv "socialAPI/internal/api/service/friendship"
	uSrv "socialAPI/internal/api/service/user"
	"socialAPI/internal/api/user"
	"socialAPI/internal/lib"
	"socialAPI/internal/setting/cfg"
	"socialAPI/internal/shared"
	"socialAPI/internal/storage"
	"socialAPI/internal/storage/cache"
	"socialAPI/internal/storage/repository"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	cfg     cfg.Config
	db      *gorm.DB
	service srv.Service
	cache   cache.CacheStore
	logger  *zap.SugaredLogger
}

func (a *App) LoadConfig() {
	a.cfg = cfg.Config{
		AppEnv: lib.GetStringFromEnv("APP_ENV", "development"),
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
			Port:     lib.GetStringFromEnv("DB_PORT", "5432"),
			User:     lib.GetStringFromEnv("DB_USER", "postgres"),
			Password: lib.GetStringFromEnv("DB_PASSWORD", "postgres"),
			Name:     lib.GetStringFromEnv("DB_NAME", "socialdb"),
			SSLMode:  lib.GetStringFromEnv("DB_SSLMODE", "disable"),
		},
		Redis: cfg.RedisConfig{
			Host:     lib.GetStringFromEnv("REDIS_HOST", "localhost"),
			Port:     lib.GetStringFromEnv("REDIS_PORT", "6379"),
			Password: lib.GetStringFromEnv("REDIS_PASSWORD", ""),
			DB:       lib.GetIntFromEnv("REDIS_DB", 0),
		},
	}
}

func (a *App) SetupLogger() {
	var err error
	a.logger, err = shared.InitLogger(a.cfg.AppEnv)

	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
}

func (a *App) InitStorages(madeMigrations bool) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		a.cfg.DB.Host,
		a.cfg.DB.Port,
		a.cfg.DB.User,
		a.cfg.DB.Password,
		a.cfg.DB.Name,
		a.cfg.DB.SSLMode,
	)

	var err error
	a.db, err = storage.BootstrapDatabase(dsn)
	if madeMigrations {
		storage.MadeMigrations(a.db)
	}

	if err != nil {
		a.logger.Panicw("Failed to initialize database", "error", err)
		return
	}

	redis, err := cache.NewRedis(a.cfg.Redis)
	if err != nil {
		a.logger.Panicw("Failed to initialize Redis", "error", err)
		return
	}

	a.cache = redis
}

func (a *App) MountServices() {
	postgresRepo := repository.NewPostgresRepo(a.db)

	tokenService := shared.NewTokenService(a.cfg.Auth.AccessSecret, a.cfg.Auth.AccessTTL)
	authService := auth.NewAuthService(postgresRepo.Users(), postgresRepo.RefreshTokens(), a.cfg.Auth, a.cache, *tokenService, a.logger)
	userService := uSrv.NewUserService(postgresRepo.Users(), a.logger)
	friendshipService := frSrv.NewFriendshipService(postgresRepo.Friendship(), a.logger)

	a.service = srv.NewService(authService, *tokenService, userService, friendshipService)
}

func (a App) MountRouter() *chi.Mux {
	authController := au.NewAuthController(a.service.Auth(), a.service.Token(), a.logger)
	userController := user.NewAuthController(a.service.User(), a.service.Token(), a.logger)
	friendshipController := friendship.NewFriendshipController(a.service.Friendship(), a.service.Token(), a.logger)

	r := chi.NewRouter()

	authController.RegisterRoutes(r)
	userController.RegisterRoutes(r)
	friendshipController.RegisterRoutes(r)

	return r
}

// TODO
func (a App) RunServer() {
}
