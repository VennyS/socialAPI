package lib

import (
	"log"
	"os"
	"strconv"
	"time"
)

func GetStringFromEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func GetIntFromEnv(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Invalid integer for %s: %v, using default value %d", key, err, fallback)
			return fallback
		}
		return intValue
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
