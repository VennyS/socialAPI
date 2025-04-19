package lib

import (
	"log"
	"os"
	"time"
)

func GetStringFromEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func GetDurationFromEnv(key string, fallback time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		duration, err := time.ParseDuration(value)
		if err != nil {
			log.Printf("Invalid duration for %s: %v, using default value %v", key, err, fallback)
			return fallback
		}
		return duration
	}
	return fallback
}
