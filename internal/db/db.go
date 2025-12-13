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
		fmt.Println("cannot connect to db, error:", err)
		panic("")
	}

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
