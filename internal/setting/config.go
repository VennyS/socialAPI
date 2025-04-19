package setting

import "time"

type Config struct {
	Server ServerConfig
	Auth   AuthConfig
	DB     DBConfig
}

type ServerConfig struct {
	Addr string
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
