package mock

import (
	"fmt"
	"testing"

	"github.com/cox96de/runner/db"
	"github.com/cox96de/runner/util"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gotest.tools/v3/assert"
)

// NewMockDB creates a new database with tables created from given models.
// It uses sqlite in memory as the database.
func NewMockDB(t *testing.T) *db.Client {
	t.Helper()
	file := util.RandomID("sql-mocker")
	conn, err := gorm.Open(
		sqlite.Open(fmt.Sprintf("file:%s?mode=memory", file)),
		&gorm.Config{
			SkipDefaultTransaction: true,
			Logger:                 nil, // TODO: add custom logger.
		},
	)
	assert.NilError(t, err)
	err = migrateModels(conn, &db.Pipeline{}, &db.Job{}, &db.Step{})
	assert.NilError(t, err)
	return db.NewClient(db.SQLite, conn)
}

func migrateModels(conn *gorm.DB, models ...interface{}) error {
	err := conn.AutoMigrate(models...)
	if err != nil {
		if sqlDB, err := conn.DB(); err == nil {
			_ = sqlDB.Close()
		}
		return errors.WithMessage(err, "failed to migrate models")
	}
	return nil
}
