package app

import (
	"net/url"
	"platform/go/core/ids"
	"platform/go/core/kv"
	"platform/go/core/mail"
	"platform/go/core/models"
	"platform/go/core/ory"

	"github.com/benbjohnson/clock"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Config struct {
	PhysicalBits    int
	DSN             string
	KratosPublicURL string
	SmtpURL         *url.URL
	RedisURL        *url.URL
}

type AppFacade struct {
	id    *ids.IDs
	db    *gorm.DB
	ory   *ory.Ory
	redis *redis.Client
	mail  *mail.Courier
	kv    *kv.KV
	clock clock.Clock
}

type App interface {
	IDs() *ids.IDs
	DB() *gorm.DB
	Ory() *ory.Ory
	Redis() *redis.Client
	Mail() *mail.Courier
	KV() *kv.KV
	Clock() clock.Clock

	Initialize() error
}

func NewWithConfig(config Config) (*AppFacade, error) {
	db, err := initializeDB(config.DSN)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	id, err := ids.New(ids.IDsOpts{PhysicalBits: config.PhysicalBits})
	if err != nil {
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

	redisPassword, _ := config.RedisURL.User.Password()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisURL.Host,
		Username: config.RedisURL.User.Username(),
		Password: redisPassword,
	})

	kvStore, err := kv.New(kv.Opts{Redis: redisClient})
	if err != nil {
		return nil, err
	}

	app := &AppFacade{id: id, db: db, ory: ory, mail: courier, redis: redisClient, kv: kvStore, clock: clock.New()}

	if err := app.Initialize(); err != nil {
		return nil, err
	}

	return app, nil
}

func (a *AppFacade) IDs() *ids.IDs {
	return a.id
}
func (a *AppFacade) DB() *gorm.DB {
	return a.db
}
func (a *AppFacade) Ory() *ory.Ory {
	return a.ory
}
func (a *AppFacade) Redis() *redis.Client {
	return a.redis
}
func (a *AppFacade) Mail() *mail.Courier {
	return a.mail
}
func (a *AppFacade) KV() *kv.KV {
	return a.kv
}
func (a *AppFacade) Clock() clock.Clock {
	return a.clock
}

type Options func(a *AppFacade)

func NewWithOptions(options ...Options) *AppFacade {
	a := &AppFacade{}
	for _, opt := range options {
		opt(a)
	}
	return a
}

func (a *AppFacade) Initialize() error {
	addedModels := make([]interface{}, 0, len(models.AllModels))
	for _, m := range models.AllModels {
		addedModels = append(addedModels, m)
	}
	if err := a.DB().AutoMigrate(addedModels...); err != nil {
		return errors.WithStack(err)
	}
	log.Info("Successfully applied database migrations.")
	return nil
}

func WithDB(db *gorm.DB) func(a *AppFacade) {
	return func(a *AppFacade) {
		a.db = db
	}
}

func WithRedis(r *redis.Client) Options {
	return func(a *AppFacade) {
		a.redis = r
	}
}

func WithIDs(id *ids.IDs) Options {
	return func(a *AppFacade) {
		a.id = id
	}
}

func WithKV(kvStore *kv.KV) Options {
	return func(a *AppFacade) {
		a.kv = kvStore
	}
}

func WithOry(oryClient *ory.Ory) Options {
	return func(a *AppFacade) {
		a.ory = oryClient
	}
}

func WithClock(c clock.Clock) Options {
	return func(a *AppFacade) {
		a.clock = c
	}
}

func init() {
	log.SetReportCaller(true)
}
