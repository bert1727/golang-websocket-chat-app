package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// APIKey     string
	// APISecret  string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	Port       int
	GoEnv      string
	JWTSecret  string
	LogLevel   string
}

func New() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("cannot load .env file")
	}

	// port, err := strconv.Atoi(getOrDefault("PORT", "8080"))
	// if err != nil {
	// 	log.Fatal("Invalid PORT value")
	// }
	port := 0

	return &Config{
		DBHost:     getOrDefault("DB_HOST", "localhost"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		Port:       port,
		GoEnv:      getOrDefault("GO_ENV", "development"),
		JWTSecret:  os.Getenv("JWTSecret"),
		LogLevel:   os.Getenv("LOG_LEVEL"),
	}
}

func getOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
