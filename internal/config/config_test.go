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
					BaseURL: "http://localhost:8080",
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
					BaseURL: "https://short.example.com",
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
					BaseURL: "http://localhost:8080",
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
					BaseURL: "http://mydomain.com",
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
					BaseURL: "http://localhost:8080",
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
					BaseURL: "https://short.example.com",
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
					BaseURL: "http://localhost:8080",
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
					BaseURL: "http://mydomain.com",
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
				_ = os.Setenv(k, v)
			}

			defer func() {
				for k := range tt.vars {
					_ = os.Unsetenv(k)
				}
			}()

			config := InitConfig()
			assert.Equal(t, tt.expected, *config)
		})
	}
}
