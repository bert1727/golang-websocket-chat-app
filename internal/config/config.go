package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/bert1727/ChatApp/internal/db"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var (
	once   sync.Once
	config *Config
)

type Config struct {
	// APIKey     string
	// APISecret  string
	DBHost                 string
	DBUser                 string
	DBPassword             string
	DBName                 string
	Port                   int
	GoEnv                  string
	JWTSecret              string
	LogLevel               string
	DB                     *gorm.DB
	HTTPAccessTokenExpire  int
	HTTPRefreshTokenExpire int
}

func New() *Config {
	once.Do(func() {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("cannot load .env file")
		}

		db := db.SetupDB()
		// port, err := strconv.Atoi(getOrDefault("PORT", "8080"))
		// if err != nil {
		// 	log.Fatal("Invalid PORT value")
		// }
		port := 0

		config = &Config{
			DBHost:                 getOrDefault("DB_HOST", "localhost"),
			DBUser:                 os.Getenv("DB_USER"),
			DBPassword:             os.Getenv("DB_PASSWORD"),
			DBName:                 os.Getenv("DB_NAME"),
			HTTPAccessTokenExpire:  toInt(os.Getenv("HTTP_ACCESS_TOKEN_EXPIRE")),
			HTTPRefreshTokenExpire: toInt(os.Getenv("HTTP_REFRESH_TOKEN_EXPIRE")),
			Port:                   port,
			GoEnv:                  getOrDefault("GO_ENV", "development"),
			JWTSecret:              os.Getenv("JWTSecret"),
			LogLevel:               os.Getenv("LOG_LEVEL"),
			DB:                     db,
		}
	})

	return config
}

func getOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse string to int from .env, field:%s err:%v", s, err))
	}
	return i
}
