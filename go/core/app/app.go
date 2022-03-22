package app

import (
	"stackv2/go/core/ids"
	"stackv2/go/core/models"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Config struct {
	PhysicalBits int
	DSN          string
}

type App struct {
	IDs *ids.IDs
	DB  *gorm.DB
}

func Initialize(config Config) (*App, error) {
	db, err := initializeDB(config.DSN)
	if err != nil {
		return nil, err
	}

	id, err := ids.New(ids.IDsOpts{PhysicalBits: config.PhysicalBits})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.User{}, &models.Project{}, &models.Environment{}); err != nil {
		return nil, err
	}

	log.Info("Successfully applied database migrations.")

	return &App{IDs: id, DB: db}, nil
}
