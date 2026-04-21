package db

import (
	"fmt"
	"os"

	"github.com/bert1727/ChatApp/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {
	db, err := connectDB()
	if err != nil {
		fmt.Println("failed to connect db, error:", err)
		panic("fail")
	}

	db.Exec("PRAGMA foreign_keys = ON")

	if err = db.AutoMigrate(&domain.User{}, &domain.Message{}); err != nil {
		panic(fmt.Sprintf("can't AutoMigrate err: %v", err))
	}

	return db
}

func connectDB() (*gorm.DB, error) {
	connection := os.Getenv("SQLITE_CONNECTION")
	return gorm.Open(sqlite.Open(connection), &gorm.Config{
		TranslateError: true,
	})
}
