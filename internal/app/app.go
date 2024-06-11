// Package app provides a server application that handles request and related entities.
// The Application struct defines the core of the server, including its http server, databases, services,
// and infrastructure configurations.
//
// The Application is initialized through the New function, which accepts various options through the
// Option functional parameter. Once initialized, the server can be started by calling the Start method.
// The server can be gracefully shut down through the Stop method, which waits for all active requests to
// finish before shutting down.
package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nachoconques0/schwarz-challenge/internal/coupon"
	"github.com/nachoconques0/schwarz-challenge/internal/http"
	shoppingcart "github.com/nachoconques0/schwarz-challenge/internal/shopping_cart"
)

const defaultTimeout = 5 * time.Second

// Application holds the basic structure that
// defines the application
type Application struct {
	httpPort string
	server   *http.Server

	// value used for determine the gap of time
	// required for shutdown the application
	timeout time.Duration

	// DB configuration
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string

	// HTTP Endpoints
	ShoppingCartHTTPEndpoint string

	// internal domain services
	shoppingCartService shoppingcart.Service
	couponService       coupon.Service

	// domain repositories
	shoppingCartRepo shoppingcart.Repository
	couponRepo       coupon.Repository

	// HTTP Controllers
	shoppingCartCtrl shoppingcart.Server
}

// New function builds a new application applying
// the given application options
func New(opts ...Option) (*Application, error) {
	a := &Application{
		timeout: defaultTimeout,
	}
	for _, o := range opts {
		o(a)
	}
	db, err := a.setupInfra()
	if err != nil {
		return nil, fmt.Errorf("error while setting up the infra - %w", err)
	}

	if err := a.setupDomain(db); err != nil {
		return nil, fmt.Errorf("error while setting up the domain - %w", err)
	}

	err = a.setupHTTPServer()
	if err != nil {
		return nil, fmt.Errorf("error while setting up the http server - %w", err)
	}

	return a, nil
}

// Start methods will start the application, this means
// that it will run all the servers that defines the application itself
// it will be running until a shutdown signal is received, it returns an error
// in case the application can't start
func (a *Application) Start() {
	go func() {
		slog.Info(fmt.Sprintf("HTTP server: starting at %s\n", a.server.Addr))
		if err := a.server.Run(); err != nil {
			slog.Error(fmt.Sprintf("Application: error running the http server: %s", err))
		}
	}()

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGTERM, os.Interrupt)
	defer signal.Stop(quitCh)

	slog.Info("Application: running...")
	// Block waiting for shutdown
	<-quitCh
	slog.Info("Application: stopping")

	// Start context with cancellation set for the defined timeout
	slog.Info(fmt.Sprintf("Application: stopping in %s...\n", a.timeout.String()))

	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	if err := a.server.Stop(ctx); err != nil {
		slog.Error(fmt.Sprintf("Application: error stopping application: %s", err))
	}

	<-ctx.Done()
	slog.Info("Application: stopped")
}
