package app

import (
	"fmt"

	"github.com/nachoconques0/schwarz-challenge/internal/postgres"
	"gorm.io/gorm"
)

// setupInfra will start all the infra components for this service, and
// it will return the postgres db instance
func (a *Application) setupInfra() (*gorm.DB, error) {
	db, err := postgres.NewDB(
		&postgres.DBOptions{
			Host:     a.dbHost,
			Port:     a.dbPort,
			User:     a.dbUser,
			Password: a.dbPassword,
			Database: a.dbName,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error while setting up the database - %w", err)
	}
	return db, nil
}
