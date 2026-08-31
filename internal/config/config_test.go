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

func TestParseCliArgs(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected Config
	}{
		{
			name: "default values",
			args: []string{"cmd"},
			expected: Config{
				address: "localhost:8080",
				baseURL: "http://localhost:8080",
			},
		},
		{
			name: "custom port and base URL",
			args: []string{"cmd", "-a", "localhost:9000", "-b", "https://short.example.com"},
			expected: Config{
				address: "localhost:9000",
				baseURL: "https://short.example.com",
			},
		},
		{
			name: "only custom port",
			args: []string{"cmd", "-a", "localhost:3000"},
			expected: Config{
				address: "localhost:3000",
				baseURL: "http://localhost:8080",
			},
		},
		{
			name: "only custom base URL",
			args: []string{"cmd", "-b", "http://mydomain.com"},
			expected: Config{
				address: "localhost:8080",
				baseURL: "http://mydomain.com",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resetForTesting()

			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			os.Args = tc.args

			config := ParseCliArgs()
			assert.Equal(t, tc.expected, *config)
		})
	}
}
