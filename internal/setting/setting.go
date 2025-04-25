package setting

import (
	"fmt"
	"net/http"
	"socialAPI/internal/api"
	"socialAPI/internal/api/auth"
	"socialAPI/internal/api/chat"
	"socialAPI/internal/api/chat/ws"
	"socialAPI/internal/api/friendship"
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
	service api.Service
	cache   cache.CacheStore
	logger  *zap.SugaredLogger
	hub     *ws.Hub
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

func (a *App) SetupWS() {
	a.hub = ws.NewHub(a.logger)
	go a.hub.Run()
}

func (a *App) InitStorages(doMigrations bool) {
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
	if doMigrations {
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
	repo := repository.NewPostgresRepo(a.db)

	tokenService := shared.NewTokenService(a.cfg.Auth.AccessSecret, a.cfg.Auth.AccessTTL)
	authService := auth.NewAuthService(repo.Users(), repo.RefreshTokens(), a.cfg.Auth, a.cache, *tokenService, a.logger)
	userService := user.NewUserService(repo.Users(), a.logger)
	friendshipService := friendship.NewFriendshipService(repo.Friendship(), a.logger)
	chatService := chat.NewChatService(repo.Chats(), repo.Users(), a.hub, a.logger)

	a.service = api.NewService(authService, *tokenService, userService, friendshipService, chatService)
}

func (a App) MountRouter() *chi.Mux {
	authController := auth.NewAuthController(a.service.Auth(), a.service.Token(), a.logger)
	userController := user.NewAuthController(a.service.User(), a.service.Token(), a.logger)
	friendshipController := friendship.NewFriendshipController(a.service.Friendship(), a.service.Token(), a.logger)
	chatController := chat.NewChatController(a.service.Chat(), a.service.Token(), a.logger)

	r := chi.NewRouter()

	authController.RegisterRoutes(r)
	userController.RegisterRoutes(r)
	friendshipController.RegisterRoutes(r)
	chatController.RegisterRoutes(r)

	return r
}

func (a App) RunServer(r *chi.Mux) {
	http.ListenAndServe(a.cfg.Server.Addr, r)
}
