package config

import (
	"flag"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestNewConfig(t *testing.T) {
	// Backup the original command line arguments and defer their restoration
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Backup the original environment and defer its restoration
	originalEnv := os.Environ()
	defer restoreEnv(originalEnv)

	tests := []struct {
		name      string
		setupArgs func()
		setupEnv  func()
		want      *Config
		wantErr   bool
	}{
		{
			name: "DefaultValues",
			setupArgs: func() {
				os.Args = []string{"cmd"}
			},
			setupEnv: func() {
				// Clear environment variables for this test
				clearEnv()
			},
			want: &Config{
				ServerAddress:   "localhost:8080",
				BaseAddress:     "http://localhost:8080",
				FilePath:        "/tmp/short-url-db.json",
				DatabaseAddress: "",
			},
			wantErr: false,
		},
		{
			name: "CustomCommandLineArguments",
			setupArgs: func() {
				os.Args = []string{"cmd", "-a", "127.0.0.1:9090", "-b", "http://127.0.0.1:9090", "-f", "/var/tmp/urls.json"}
			},
			setupEnv: func() {
				clearEnv()
			},
			want: &Config{
				ServerAddress:   "127.0.0.1:9090",
				BaseAddress:     "http://127.0.0.1:9090",
				FilePath:        "/var/tmp/urls.json",
				DatabaseAddress: "",
			},
			wantErr: false,
		},
		{
			name: "EnvironmentVariablesSet",
			setupArgs: func() {
				os.Args = []string{"cmd"}
			},
			setupEnv: func() {
				clearEnv()
				os.Setenv("SERVER_ADDRESS", "192.168.1.1:8080")
				os.Setenv("BASE_URL", "http://192.168.1.1:8080")
				os.Setenv("FILE_STORAGE_PATH", "/data/urls.json")
				os.Setenv("DATABASE_DSN", "user:password@/dbname")
			},
			want: &Config{
				ServerAddress:   "192.168.1.1:8080",
				BaseAddress:     "http://192.168.1.1:8080",
				FilePath:        "/data/urls.json",
				DatabaseAddress: "user:password@/dbname",
			},
			wantErr: false,
		},
		{
			name: "CombinationEnvAndArgs",
			setupArgs: func() {
				os.Args = []string{"cmd", "-a", "127.0.0.1:9090"}
			},
			setupEnv: func() {
				clearEnv()
				os.Setenv("BASE_URL", "http://192.168.1.1:8080")
				os.Setenv("FILE_STORAGE_PATH", "/data/urls.json")
			},
			want: &Config{
				ServerAddress:   "127.0.0.1:9090",
				BaseAddress:     "http://192.168.1.1:8080",
				FilePath:        "/data/urls.json",
				DatabaseAddress: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupArgs()
			tt.setupEnv()
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ExitOnError) // Reset the flag set

			got, err := NewConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func clearEnv() {
	os.Clearenv()
}

func restoreEnv(env []string) {
	clearEnv()
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		os.Setenv(parts[0], parts[1])
	}
}
