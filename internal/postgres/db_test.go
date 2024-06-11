package postgres_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nachoconques0/schwarz-challenge/internal/errors"
	"github.com/nachoconques0/schwarz-challenge/internal/postgres"
)

func TestNewDB(t *testing.T) {
	opts := &postgres.DBOptions{
		Host:     "127.0.0.1",
		Port:     "5435",
		User:     "schwarz_svc",
		Password: "schwarz_svc",
		Database: "schwarz_svc",
		SSLMode:  "disable",
	}

	t.Run("when no config is provided", func(t *testing.T) {
		db, err := postgres.NewDB(&postgres.DBOptions{})
		assert.Nil(t, db, "db should be nil")
		assert.NotNil(t, err, "error should not be nil")
	})

	t.Run("when a db config is provided", func(t *testing.T) {
		db, err := postgres.NewDB(opts)
		assert.NotNil(t, db, "db should not be nil")
		assert.Nil(t, err, "error should be nil")

		sqlDB, err := db.DB()
		assert.Nil(t, err, "error should be nil")
		assert.Nil(t, sqlDB.Ping(), "db ping should be successful")
	})

	t.Run("when a db config is provided", func(t *testing.T) {
		db, err := postgres.NewDB(opts)
		assert.NotNil(t, db, "db should not be nil")
		assert.Nil(t, err, "error should be nil")

		sqlDB, err := db.DB()
		assert.Nil(t, err, "error should be nil")
		assert.Nil(t, sqlDB.Ping(), "db ping should be successful")
	})

	testCases := []struct {
		name string
		opts func() *postgres.DBOptions
		err  error
	}{
		{
			name: "when host is not provided",
			opts: func() *postgres.DBOptions {
				o := *opts
				o.Host = ""
				return &o
			},
			err: errors.NewWrongInput("host is required"),
		},
		{
			name: "when port is not provided",
			opts: func() *postgres.DBOptions {
				o := *opts
				o.Port = ""
				return &o
			},
			err: errors.NewWrongInput("port is required"),
		},
		{
			name: "when user is not provided",
			opts: func() *postgres.DBOptions {
				o := *opts
				o.User = ""
				return &o
			},
			err: errors.NewWrongInput("user is required"),
		},
		{
			name: "when password is not provided",
			opts: func() *postgres.DBOptions {
				o := *opts
				o.Password = ""
				return &o
			},
			err: errors.NewWrongInput("password is required"),
		},
		{
			name: "when database is not provided",
			opts: func() *postgres.DBOptions {
				o := *opts
				o.Database = ""
				return &o
			},
			err: errors.NewWrongInput("database is required"),
		},
		{
			name: "when SSL Mode is not provided",
			opts: func() *postgres.DBOptions {
				o := *opts
				o.SSLMode = ""
				return &o
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			options := tc.opts()
			db, err := postgres.NewDB(options)

			if tc.err != nil {
				assert.Nil(t, db, "db should be nil")
				assert.NotNil(t, err, "error should not be nil")
				assert.Equal(t, err, tc.err, "error should be equal")
			} else {
				assert.NotNil(t, db, "db should not be nil")
				assert.Nil(t, err, "error should be nil")
			}
		})
	}
}
