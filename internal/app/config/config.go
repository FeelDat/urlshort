package config

import (
	"flag"
	"github.com/caarlos0/env"
	"log"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseAddress   string `env:"BASE_URL"`
}

func NewConfig() (*Config, error) {
	c := &Config{}

	err := env.Parse(c)
	if err != nil {
		log.Fatal(err)
	}

	flag.Func("a", "address to run the HTTP server on", func(input string) error {
		if input != "" {
			c.ServerAddress = input
		}
		return nil
	})

	flag.Func("b", "base address for resulting shortened URL", func(input string) error {
		if input != "" {
			c.BaseAddress = input
		}
		return nil
	})

	flag.Parse()

	// Only if no environment variable and no flag provided, use default
	if c.ServerAddress == "" {
		c.ServerAddress = "localhost:8080"
	}
	if c.BaseAddress == "" {
		c.BaseAddress = "localhost:8080"
	}

	return c, nil
}
