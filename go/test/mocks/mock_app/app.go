package mock_app

import (
	"platform/go/core/app"
	"platform/go/core/ids"
	"platform/go/core/kv"
	"platform/go/test/mocks/mock_db"
	"platform/go/test/mocks/mock_ory"
	"platform/go/test/mocks/mock_redis"
	"testing"

	"github.com/alicebob/miniredis/v2"
)

type MockApp struct {
	*app.AppFacade

	mockRedis *miniredis.Miniredis
}

func New(t *testing.T) (*MockApp, error) {
	c, s, err := mock_redis.New(t)
	if err != nil {
		return nil, err
	}
	gormDB, err := mock_db.New(t)
	// gormDB, err := gorm.Open(postgres.Open(os.Getenv("APP_DSN")))
	if err != nil {
		return nil, err
	}
	id, err := ids.New(ids.IDsOpts{PhysicalBits: 2})
	if err != nil {
		return nil, err
	}

	kvStore, err := kv.New(kv.Opts{Redis: c})
	if err != nil {
		return nil, err
	}

	mockOry, err := mock_ory.New(t)
	if err != nil {
		return nil, err
	}

	a := app.NewWithOptions(app.WithDB(gormDB), app.WithIDs(id), app.WithKV(kvStore), app.WithRedis(c), app.WithOry(mockOry))
	mockApp := &MockApp{
		AppFacade: a,
		mockRedis: s,
	}

	if err := mockApp.Initialize(); err != nil {
		return nil, err
	}
	return mockApp, nil
}

func (ma *MockApp) Cleanup() {
}
