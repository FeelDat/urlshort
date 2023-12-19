package logger

import (
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestInitLogger(t *testing.T) {
	// Define a map for converting string levels to zapcore.Level
	stringToZapLevel := map[string]zapcore.Level{
		"debug":  zapcore.DebugLevel,
		"info":   zapcore.InfoLevel,
		"warn":   zapcore.WarnLevel,
		"error":  zapcore.ErrorLevel,
		"dpanic": zapcore.DPanicLevel,
		"panic":  zapcore.PanicLevel,
		"fatal":  zapcore.FatalLevel,
	}

	// Define the test cases
	testCases := []struct {
		name    string
		level   string
		wantErr bool
	}{
		{"ValidDebugLevel", "debug", false},
		{"ValidInfoLevel", "info", false},
		{"ValidWarnLevel", "warn", false},
		{"ValidErrorLevel", "error", false},
		{"ValidDPanicLevel", "dpanic", false},
		{"ValidPanicLevel", "panic", false},
		{"ValidFatalLevel", "fatal", false},
		{"InvalidLevel", "invalid", true}}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			gotLogger, err := InitLogger(tt.level)

			if (err != nil) != tt.wantErr {
				t.Errorf("InitLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Skip further checks if an error is expected
			if tt.wantErr {
				return
			}

			// Unwrap the SugaredLogger to access the base Logger
			unwrappedLogger := gotLogger.Desugar()

			// Check if the logger is set to the correct level
			if !unwrappedLogger.Core().Enabled(stringToZapLevel[tt.level]) {
				t.Errorf("Logger level is not set correctly for level %s", tt.level)
			}
		})
	}
}
