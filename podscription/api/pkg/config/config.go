package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Server Server `json:"server"`
	OpenAI OpenAI `json:"openai"`
	Store  Store  `json:"store"`
}

// Server holds server configuration
type Server struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// OpenAI holds OpenAI API configuration
type OpenAI struct {
	APIKey      string `json:"apiKey"`
	Model       string `json:"model"`
	Temperature float32 `json:"temperature"`
	MaxTokens   int    `json:"maxTokens"`
}

// Store holds data store configuration
type Store struct {
	Type string `json:"type"`
	Path string `json:"path,omitempty"`
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: Server{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnvAsInt("SERVER_PORT", 8080),
		},
		OpenAI: OpenAI{
			APIKey:      getEnv("OPENAI_API_KEY", ""),
			Model:       getEnv("OPENAI_MODEL", "gpt-3.5-turbo"),
			Temperature: getEnvAsFloat32("OPENAI_TEMPERATURE", 0.7),
			MaxTokens:   getEnvAsInt("OPENAI_MAX_TOKENS", 1000),
		},
		Store: Store{
			Type: getEnv("STORE_TYPE", "memory"),
			Path: getEnv("STORE_PATH", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
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

func getEnvAsFloat32(key string, defaultValue float32) float32 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseFloat(valueStr, 32); err == nil {
		return float32(value)
	}
	return defaultValue
}