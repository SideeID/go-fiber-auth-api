package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port         string
	MongoURI     string
	DBName       string
	JWTSecret    string
	AppEnv       string
	RedisURL     string
	APIRateLimit int
	APITimeout   int
	LogLevel     string
}

func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		MongoURI:     getEnv("MONGODB_URI", ""),
		DBName:       getEnv("DB_NAME", "ujikom"),
		JWTSecret:    getEnv("JWT_SECRET", "ujikom-secret-key"),
		AppEnv:       getEnv("APP_ENV", "development"),
		RedisURL:     getEnv("REDIS_URL", "redis://localhost:6379"),
		APIRateLimit: getEnvAsInt("API_RATE_LIMIT", 100),
		APITimeout:   getEnvAsInt("API_TIMEOUT", 30),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
	}
}

func (c *Config) GetCurrentTime() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

func (c *Config) ValidateAtlasConnection() error {
	if c.MongoURI == "" {
		return fmt.Errorf("MONGODB_URI is required for Atlas connection")
	}
	
	if !strings.Contains(c.MongoURI, "mongodb+srv://") {
		return fmt.Errorf("invalid Atlas URI format, should start with mongodb+srv://")
	}
	
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}