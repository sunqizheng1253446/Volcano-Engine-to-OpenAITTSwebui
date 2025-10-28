package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ListenAddr           string
	MaxConnections       int
	MaxConcurrentCalls   int
	CheckInterval        int
	ByteDanceAppID       string
	ByteDanceBearerToken string
	ByteDanceCluster     string
	ByteDanceVoiceType   string
	OpenAITTSAPIKey      string
	LogLevel             string
	GinMode              string
}

func LoadConfig() (*Config, error) {
	// 尝试加载.env文件，如果不存在也不报错
	_ = godotenv.Load()

	config := &Config{
		ListenAddr:           getEnv("LISTEN_ADDR", ":8080"),
		MaxConnections:       getEnvAsInt("MAX_CONNECTIONS", 100),
		MaxConcurrentCalls:   getEnvAsInt("MAX_CONCURRENT_CALLS", 10),
		CheckInterval:        getEnvAsInt("CHECK_INTERVAL", 5), // 秒
		ByteDanceAppID:       getEnv("BYTEDANCE_TTS_APP_ID", ""),
		ByteDanceBearerToken: getEnv("BYTEDANCE_TTS_BEARER_TOKEN", ""),
		ByteDanceCluster:     getEnv("BYTEDANCE_TTS_CLUSTER", ""),
		ByteDanceVoiceType:   getEnv("BYTEDANCE_TTS_VOICE_TYPE", ""),
		OpenAITTSAPIKey:      getEnv("OPENAI_TTS_API_KEY", ""),
		LogLevel:             getEnv("LOG_LEVEL", "info"),
		GinMode:              getEnv("GIN_MODE", "release"),
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