package app

import (
	"stackv2/go/core/ids"
	"stackv2/go/core/models"
	"stackv2/go/core/ory"

	kratos "github.com/ory/kratos-client-go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Config struct {
	PhysicalBits    int
	DSN             string
	KratosPublicURL string
}

type App struct {
	IDs *ids.IDs
	DB  *gorm.DB
	Ory *kratos.APIClient
}

func Initialize(config Config) (*App, error) {
	db, err := initializeDB(config.DSN)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	id, err := ids.New(ids.IDsOpts{PhysicalBits: config.PhysicalBits})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	addedModels := make([]interface{}, len(models.AllModels))
	for i, m := range models.AllModels {
		addedModels[i] = m
	}
	if err := db.AutoMigrate(addedModels...); err != nil {
		return nil, errors.WithStack(err)
	}

	ory, err := ory.New(ory.Opts{KratosPublicURL: config.KratosPublicURL})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	log.Info("Successfully applied database migrations.")

	return &App{IDs: id, DB: db, Ory: ory}, nil
}
