//go:build dev
// +build dev

package main

import (
	"os"
)

func init() {
	os.Setenv("HTTP_PORT", "8080")
	os.Setenv("DB_HOST", "127.0.0.1")
	os.Setenv("DB_PORT", "5434")
	os.Setenv("DB_NAME", "schwarz_svc")
	os.Setenv("DB_USER", "schwarz_svc")
	os.Setenv("DB_PASSWORD", "schwarz_svc")
}
