package config

import (
	"os"
	"strconv"
)

type Config struct {
	ListenAddr         string
	MaxConnections     int
	MaxConcurrentCalls int
	CheckInterval      int
	TargetAPIURL       string // 被监控API的URL，只包含协议头和主域名，如 http://localhost:8081
	LogLevel           string
	GinMode            string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		ListenAddr:         getEnv("LISTEN_ADDR", ":8080"),
		MaxConnections:     getEnvAsInt("MAX_CONNECTIONS", 100),
		MaxConcurrentCalls: getEnvAsInt("MAX_CONCURRENT_CALLS", 10),
		CheckInterval:      getEnvAsInt("CHECK_INTERVAL", 5),                  // 秒
		TargetAPIURL:       getEnv("TARGET_API_URL", "http://localhost:8081"), // 只识别到协议头和主域名
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		GinMode:            getEnv("GIN_MODE", "release"),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
