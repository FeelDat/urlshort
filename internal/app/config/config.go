// Package config provides functions for working with configuration settings.
//
// This package is used for extracting configuration data from environment variables and command line arguments.
// The primary use of this package is to create an instance of Config with the current application settings.
package config

import (
	"flag"
	"github.com/caarlos0/env"
	"log"
)

// Config contains application settings.
//
// These settings can be set via environment variables or command line arguments.
type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`    // Server address.
	BaseAddress     string `env:"BASE_URL"`          // Base URL for responses with short links.
	FilePath        string `env:"FILE_STORAGE_PATH"` // Path for saving a file with short URLs.
	DatabaseAddress string `env:"DATABASE_DSN"`      // Database address.
}

// NewConfig creates and returns a new instance of Config.
//
// The function also reads command line arguments and sets configuration values based on environment variables.
// In case of an error parsing the environment variables, the function terminates the program.
func NewConfig() (*Config, error) {
	c := &Config{}

	// Define command line flags.
	flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "server address")
	flag.StringVar(&c.BaseAddress, "b", "http://localhost:8080", "base URL for responses with short links")
	flag.StringVar(&c.FilePath, "f", "/tmp/short-url-db.json", "path for saving a file with short URLs")
	flag.StringVar(&c.DatabaseAddress, "d", "", "database address")
	flag.Parse()

	// Parse environment variables.
	err := env.Parse(c)
	if err != nil {
		log.Fatal(err) // Terminate execution in case of error.
	}

	return c, nil
}
