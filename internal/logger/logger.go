// Package logger provides utilities for initializing and managing logging.
//
// It uses the uber-go/zap logging library for efficient, structured logging.
package logger

import "go.uber.org/zap"

// InitLogger initializes a new zap SugaredLogger with the specified logging level.
//
// The function accepts a string representation of the logging level (e.g., "info", "debug", "error") and returns a SugaredLogger.
// The SugaredLogger wraps the base Logger functionality and provides a more developer-friendly API.
//
// If there's an issue parsing the level or building the logger, an error is returned.
func InitLogger(level string) (*zap.SugaredLogger, error) {

	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return zl.Sugar(), nil
}
