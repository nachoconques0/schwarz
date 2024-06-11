package helpers

import (
	"fmt"

	"github.com/nachoconques0/schwarz-challenge/internal/errors"
	"github.com/nachoconques0/schwarz-challenge/internal/postgres"
	"gorm.io/gorm"
)

// Teardown is used to close db Connection and cleanup
type Teardown func()

// NewTestDB will be used to create a test TX
// and teardown logic for cleanup
func NewTestDB() (*gorm.DB, Teardown, error) {
	opts := &postgres.DBOptions{
		Host:     "127.0.0.1",
		Port:     "5435",
		User:     "schwarz_svc",
		Password: "schwarz_svc",
		Database: "schwarz_svc",
		SSLMode:  "disable",
	}

	db, err := postgres.NewDB(opts)
	if err != nil {
		return nil, nil, errors.NewInternalError(fmt.Sprintf("error generating test db connection: %s", err))
	}

	teardown := func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}

	tx := db.Begin()
	if tx.Error != nil {
		return nil, teardown, errors.NewInternalError(fmt.Sprintf("error starting test db transaction: %s", err))
	}

	teardown = func() {
		_ = tx.Rollback()
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}

	return tx, teardown, nil
}
