package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type A struct{}

func ConnectDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
}

func AutoMigrate(db *gorm.DB, models ...any) error {
	return db.AutoMigrate(models...)
}
