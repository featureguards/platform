package app

import (
	"net/url"
	"platform/go/core/ids"
	"platform/go/core/mail"
	"platform/go/core/models"
	"platform/go/core/ory"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Config struct {
	PhysicalBits    int
	DSN             string
	KratosPublicURL string
	SmtpURL         *url.URL
}

type App struct {
	IDs  *ids.IDs
	DB   *gorm.DB
	Ory  *ory.Ory
	Mail *mail.Courier
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

	addedModels := make([]interface{}, 0, len(models.AllModels))
	for _, m := range models.AllModels {
		addedModels = append(addedModels, m)
	}
	if err := db.AutoMigrate(addedModels...); err != nil {
		return nil, errors.WithStack(err)
	}

	ory, err := ory.New(ory.Opts{KratosPublicURL: config.KratosPublicURL})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	courier, err := mail.New(mail.Opts{URL: config.SmtpURL})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	log.Info("Successfully applied database migrations.")

	return &App{IDs: id, DB: db, Ory: ory, Mail: courier}, nil
}
