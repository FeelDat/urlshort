package config

import (
	"flag"
	"github.com/caarlos0/env"
	"log"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseAddress   string `env:"BASE_URL"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
}

func NewConfig() (*Config, error) {
	c := &Config{}

	flag.StringVar(&c.ServerAddress, "a", ":8080", "server address")
	flag.StringVar(&c.BaseAddress, "b", "http://localhost:8080", "base url for short links reply")
	flag.StringVar(&c.FilePath, "f", "/tmp/short-url-db.json", "path to store file with shorten url")
	flag.Parse()

	err := env.Parse(c)
	if err != nil {
		log.Fatal(err)
	}

	return c, nil
}
