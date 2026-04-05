package config

import (
	"strconv"
	"time"
)

func getEnvAsInt(key string, fallback int) int {
	value := getEnv(key, "")
	if value == "" {
		return fallback
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return result
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	value := getEnv(key, "")
	if value == "" {
		return fallback
	}

	result, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return result
}
