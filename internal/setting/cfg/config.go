package cfg

import "time"

type Config struct {
	AppEnv string
	Server ServerConfig
	Auth   AuthConfig
	DB     DBConfig
	Redis  RedisConfig
}

type ServerConfig struct {
	Addr             string
	OriginsSeparator string
	AllowedOrigins   []string
}

type AuthConfig struct {
	AccessTTL    time.Duration
	RefreshTTL   time.Duration
	AccessSecret string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}
