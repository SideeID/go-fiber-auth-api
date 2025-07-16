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
	
	// School location configuration
	SchoolLatitude  float64
	SchoolLongitude float64
	SchoolRadius    float64 // in kilometers
	
	// Attendance configuration
	SchoolStartHour   int
	SchoolStartMinute int
	SchoolEndHour     int
	SchoolEndMinute   int
	LateThreshold     int // in minutes
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
		
		SchoolLatitude:  getEnvAsFloat("SCHOOL_LATITUDE", -8.1575),
		SchoolLongitude: getEnvAsFloat("SCHOOL_LONGITUDE", 113.722778),
		SchoolRadius:    getEnvAsFloat("SCHOOL_RADIUS", 0.1), // 100 meters
		
		SchoolStartHour:   getEnvAsInt("SCHOOL_START_HOUR", 7),
		SchoolStartMinute: getEnvAsInt("SCHOOL_START_MINUTE", 0),
		SchoolEndHour:     getEnvAsInt("SCHOOL_END_HOUR", 15),
		SchoolEndMinute:   getEnvAsInt("SCHOOL_END_MINUTE", 30),
		LateThreshold:     getEnvAsInt("LATE_THRESHOLD", 30), // 30 minutes
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

func (c *Config) GetSchoolLocation() (float64, float64, float64) {
	return c.SchoolLatitude, c.SchoolLongitude, c.SchoolRadius
}

func (c *Config) GetSchoolHours() (int, int, int, int) {
	return c.SchoolStartHour, c.SchoolStartMinute, c.SchoolEndHour, c.SchoolEndMinute
}

func (c *Config) GetLateThreshold() int {
	return c.LateThreshold
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

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}