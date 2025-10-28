package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ListenAddr         string
	MaxConnections     int
	MaxConcurrentCalls int
	CheckInterval      int
	TargetAPIURL       string // 被监控API的URL
	LogLevel           string
	GinMode            string
}

func LoadConfig() (*Config, error) {
	// 尝试加载.env文件，如果不存在也不报错
	_ = godotenv.Load()

	config := &Config{
		ListenAddr:         getEnv("LISTEN_ADDR", ":8080"),
		MaxConnections:     getEnvAsInt("MAX_CONNECTIONS", 100),
		MaxConcurrentCalls: getEnvAsInt("MAX_CONCURRENT_CALLS", 10),
		CheckInterval:      getEnvAsInt("CHECK_INTERVAL", 5),                             // 秒
		TargetAPIURL:       getEnv("TARGET_API_URL", "http://localhost:8081/api/health"), // 默认目标API地址
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
