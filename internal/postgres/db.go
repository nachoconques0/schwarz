// Package postgres contains the implementation of the persistence layer
// for the service. It uses the GORM library as ORM and
// postgres as the database driver.
package postgres

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	internalErrors "github.com/nachoconques0/schwarz-challenge/internal/errors"
)

const (
	maxConnections = 1
)

// DBOptions defines the data needed to create a db connection
// SSLMode is a optional option that will take disable as default
type DBOptions struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (o *DBOptions) connection() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		o.Host,
		o.Port,
		o.Database,
		o.User,
		o.Password,
		o.SSLMode,
	)
}

// NewDB initialized a new GORM DB with the provided database
// configuration
func NewDB(opts *DBOptions) (*gorm.DB, error) {
	if opts.Host == "" {
		return nil, internalErrors.NewWrongInput("host is required")
	}
	if opts.Port == "" {
		return nil, internalErrors.NewWrongInput("port is required")
	}
	if opts.User == "" {
		return nil, internalErrors.NewWrongInput("user is required")
	}
	if opts.Password == "" {
		return nil, internalErrors.NewWrongInput("password is required")
	}
	if opts.Database == "" {
		return nil, internalErrors.NewWrongInput("database is required")
	}
	if opts.SSLMode == "" {
		// we should use prefer as the default SSL Mode
		// to enforce a secure connection
		opts.SSLMode = "prefer"
	}

	dbLogger := logger.New(
		log.Default(),
		logger.Config{
			IgnoreRecordNotFoundError: true,
			LogLevel:                  logger.Warn,
		},
	)
	db, err := gorm.Open(postgres.Open(opts.connection()), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		return nil, internalErrors.NewWrongInput(fmt.Sprintf("error openning db connection: %s", err))
	}

	con, err := db.DB()
	if err != nil {
		return nil, internalErrors.NewWrongInput(fmt.Sprintf("error getting db connection: %s", err))
	}
	con.SetMaxOpenConns(maxConnections)
	return db, nil
}
