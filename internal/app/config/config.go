package config

import (
	"flag"
	"fmt"
	"net/url"
	"strings"
)

type Config struct {
	ServerURL  string
	ServerPort string
	BaseURL    string
	BasePort   string
}

func NewConfig() (*Config, error) {
	c := &Config{
		ServerURL:  "localhost",
		ServerPort: "8080",
		BaseURL:    "localhost",
		BasePort:   "8080",
	}

	flag.Func("a", "address to run the HTTP server on", func(serverAddress string) error {
		serverURL, serverPort, err := splitAddress(serverAddress)
		if err != nil {
			return fmt.Errorf("invalid server address: %w", err)
		}
		c.ServerURL = serverURL
		c.ServerPort = serverPort
		return nil
	})

	flag.Func("b", "base address for resulting shortened URL", func(baseAddress string) error {
		baseURL, basePort, err := splitAddress(baseAddress)
		if err != nil {
			return fmt.Errorf("invalid base URL: %w", err)
		}
		c.BaseURL = baseURL
		c.BasePort = basePort
		return nil
	})

	flag.Parse()

	return c, nil
}

func splitAddress(address string) (string, string, error) {
	if !strings.Contains(address, "://") {
		address = "//" + address
	}
	u, err := url.Parse(address)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse URL: %w", err)
	}

	hostParts := strings.Split(u.Host, ":")
	if len(hostParts) != 2 {
		return "", "", fmt.Errorf("expected host:port format, got %s", u.Host)
	}

	return hostParts[0], hostParts[1], nil
}
