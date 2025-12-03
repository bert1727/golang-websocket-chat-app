package db

import (
	"fmt"

	"github.com/bert1727/ChatApp/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {
	db, err := connectDB()
	if err != nil {
		fmt.Println("cannot connect to db, error:", err)
	}

	if err = db.AutoMigrate(&models.User{}, &models.Message{}); err != nil {
		fmt.Println("error occured", err)
	}

	return db
}

func connectDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
}
