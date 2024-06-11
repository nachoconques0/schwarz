// Package main defines the entrypoint for the service.
// It defines all the required options to run the service and
// the main function that will be executed when the service is started.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nachoconques0/schwarz-challenge/internal/app"
)

func main() {
	opts := []app.Option{
		app.WithHTTPPort(os.Getenv("HTTP_PORT")),
		app.WithDBHost(os.Getenv("DB_HOST")),
		app.WithDBPort(os.Getenv("DB_PORT")),
		app.WithDBName(os.Getenv("DB_NAME")),
		app.WithDBUser(os.Getenv("DB_USER")),
		app.WithDBPassword(os.Getenv("DB_PASSWORD")),
	}

	application, err := app.New(opts...)
	if err != nil {
		log.Fatal(nil, fmt.Sprintf("error schwarz-challenge application: %s", err.Error()))
	}

	application.Start()
}

func LoadOrDefault(env, def string) string {
	val, ok := os.LookupEnv(env)
	if !ok {
		return def
	}

	return val
}

func LoadOrPanic(env string) string {
	val, ok := os.LookupEnv(env)
	if !ok {
		panic(fmt.Sprintf("missing '%s' environment variable\n", env))
	}

	return val
}
