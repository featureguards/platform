package app

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initializeDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
