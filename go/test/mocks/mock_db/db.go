package mock_db

import (
	"database/sql"
	logx "log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/benbjohnson/clock"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(t *testing.T, cl clock.Clock) (*gorm.DB, error) {
	gormDB, err := gorm.Open(sqlite.Open(""), &gorm.Config{
		Logger: logger.New(logx.New(os.Stdout, "\r\n", logx.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,   // Slow SQL threshold
				LogLevel:                  logger.Silent, // Log level
				IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			}),
		NowFunc: cl.Now,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return gormDB, nil
}

func NewMock(t *testing.T) (*gorm.DB, *sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	t.Cleanup(func() { db.Close() })

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, nil, errors.WithStack(err)
	}
	return gormDB, db, mock, nil
}
