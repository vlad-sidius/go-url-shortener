package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func resetForTesting() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestInitConfigFromCliArgs(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected Config
	}{
		{
			name: "default values",
			args: []string{"cmd"},
			expected: Config{
				ServerConf: ServerConfig{
					Address: "localhost:8080",
				},
				URLServiceConf: URLServiceConfig{
					BaseURL:     "http://localhost:8080",
					StoragePath: "fileDb.json",
				},
			},
		},
		{
			name: "custom port and base URL",
			args: []string{"cmd", "-a", "localhost:9000", "-b", "https://short.example.com"},
			expected: Config{
				ServerConf: ServerConfig{
					Address: "localhost:9000",
				},
				URLServiceConf: URLServiceConfig{
					BaseURL:     "https://short.example.com",
					StoragePath: "fileDb.json",
				},
			},
		},
		{
			name: "only custom port",
			args: []string{"cmd", "-a", "localhost:3000"},
			expected: Config{
				ServerConf: ServerConfig{
					Address: "localhost:3000",
				},
				URLServiceConf: URLServiceConfig{
					BaseURL:     "http://localhost:8080",
					StoragePath: "fileDb.json",
				},
			},
		},
		{
			name: "only custom base URL",
			args: []string{"cmd", "-b", "http://mydomain.com"},
			expected: Config{
				ServerConf: ServerConfig{
					Address: "localhost:8080",
				},
				URLServiceConf: URLServiceConfig{
					BaseURL:     "http://mydomain.com",
					StoragePath: "fileDb.json",
				},
			},
		},
		{
			name: "only custom storage path",
			args: []string{"cmd", "-f", "storage.data"},
			expected: Config{
				ServerConf: ServerConfig{
					Address: "localhost:8080",
				},
				URLServiceConf: URLServiceConfig{
					BaseURL:     "http://localhost:8080",
					StoragePath: "storage.data",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resetForTesting()

			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			os.Args = tc.args

			config := InitConfig()
			assert.Equal(t, tc.expected, *config)
		})
	}
}

func TestInitConfigFromEnvVars(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		vars     map[string]string
		expected Config
	}{
		{
			name: "default values",
			args: []string{"cmd"},
			vars: map[string]string{},
			expected: Config{
				ServerConf: ServerConfig{
					Address: "localhost:8080",
				},
				URLServiceConf: URLServiceConfig{
					BaseURL:     "http://localhost:8080",
					StoragePath: "fileDb.json",
				},
			},
		},
		{
			name: "custom port and base URL",
			args: []string{"cmd"},
			vars: map[string]string{addressEnv: "localhost:9000", baseURLEnv: "https://short.example.com"},
			expected: Config{
				ServerConf: ServerConfig{
					Address: "localhost:9000",
				},
				URLServiceConf: URLServiceConfig{
					BaseURL:     "https://short.example.com",
					StoragePath: "fileDb.json",
				},
			},
		},
		{
			name: "only custom port",
			args: []string{"cmd"},
			vars: map[string]string{addressEnv: "localhost:3000"},
			expected: Config{
				ServerConf: ServerConfig{
					Address: "localhost:3000",
				},
				URLServiceConf: URLServiceConfig{
					BaseURL:     "http://localhost:8080",
					StoragePath: "fileDb.json",
				},
			},
		},
		{
			name: "only custom base URL",
			args: []string{"cmd"},
			vars: map[string]string{baseURLEnv: "http://mydomain.com"},
			expected: Config{
				ServerConf: ServerConfig{
					Address: "localhost:8080",
				},
				URLServiceConf: URLServiceConfig{
					BaseURL:     "http://mydomain.com",
					StoragePath: "fileDb.json",
				},
			},
		},
		{
			name: "only custom storage path",
			args: []string{"cmd"},
			vars: map[string]string{storagePathEnv: "storage.data"},
			expected: Config{
				ServerConf: ServerConfig{
					Address: "localhost:8080",
				},
				URLServiceConf: URLServiceConfig{
					BaseURL:     "http://localhost:8080",
					StoragePath: "storage.data",
				},
			},
		},
		{
			name: "overridden variables from ENV",
			args: []string{"cmd", "-a", "localhost:9000", "-b", "https://short.example.com", "-f", "storage.data"},
			vars: map[string]string{addressEnv: "localhost:3000", baseURLEnv: "http://mydomain.com", storagePathEnv: "store.db"},
			expected: Config{
				ServerConf: ServerConfig{
					Address: "localhost:3000",
				},
				URLServiceConf: URLServiceConfig{
					BaseURL:     "http://mydomain.com",
					StoragePath: "store.db",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetForTesting()

			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			os.Args = tt.args

			for k, v := range tt.vars {
				t.Setenv(k, v)
			}

			config := InitConfig()
			assert.Equal(t, tt.expected, *config)
		})
	}
}
