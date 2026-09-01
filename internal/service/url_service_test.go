package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vlad-sidius/go-url-shortener/internal/config"
	"github.com/vlad-sidius/go-url-shortener/internal/repository"
)

// mock for hash generator
type mockGenerator struct {
	results   []string
	callCount int
}

func (m *mockGenerator) Generate() string {
	if m.callCount < len(m.results) {
		res := m.results[m.callCount]
		m.callCount++
		return res
	}

	return ""
}

// test cases

func TestURLServiceLive_CreateShortCode(t *testing.T) {
	testCases := []struct {
		name        string
		repo        urlRepo
		generator   hashGenerator
		originalURL string
		expectedURL string
		expectedErr error
	}{
		{
			name:        "success on first try",
			repo:        repository.NewMemURLRepo(),
			generator:   &mockGenerator{results: []string{"abc123"}},
			originalURL: "https://example.com/long-url",
			expectedURL: "https://short.io/abc123",
			expectedErr: nil,
		},
		{
			name: "success after one collision",
			repo: func() urlRepo {
				r := repository.NewMemURLRepo()
				r.Put("abc123", "https://example.com/other-url")
				return r
			}(),
			generator:   &mockGenerator{results: []string{"abc123", "def456"}},
			originalURL: "https://example.com/long-url",
			expectedURL: "https://short.io/def456",
			expectedErr: nil,
		},
		{
			name: "retries exhausted due to continuous collisions",
			repo: func() urlRepo {
				r := repository.NewMemURLRepo()
				r.Put("collision", "https://example.com/other-url")
				return r
			}(),
			generator: func() hashGenerator {
				// Provide 10 collisions to exhaust retriesLimit (which is 10)
				results := make([]string, 10)

				for i := range results {
					results[i] = "collision"
				}

				return &mockGenerator{results: results}
			}(),
			originalURL: "https://example.com/long-url",
			expectedURL: "",
			expectedErr: errors.New("failed to generate unique hash, too many collisions"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf := &config.URLServiceConfig{BaseURL: "https://short.io"}
			urlService := NewURLServiceLive(conf, tc.repo, tc.generator)

			shortURL, err := urlService.CreateShortCode(tc.originalURL)

			if tc.expectedErr != nil {
				assert.Error(t, err, "expect error")
			}

			assert.Equal(t, tc.expectedURL, shortURL, "short URL is invalid")
		})
	}
}

func TestURLServiceLive_ResolveOriginalURL(t *testing.T) {
	testCases := []struct {
		name        string
		repo        urlRepo
		hash        string
		expectedURL string
		expectedOk  bool
	}{
		{
			name: "found in repo",
			repo: func() urlRepo {
				r := repository.NewMemURLRepo()
				r.Put("abc123", "https://example.com/long-url")
				return r
			}(),
			hash:        "abc123",
			expectedURL: "https://example.com/long-url",
			expectedOk:  true,
		},
		{
			name:        "not found in repo",
			repo:        repository.NewMemURLRepo(),
			hash:        "xyz789",
			expectedURL: "",
			expectedOk:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf := &config.URLServiceConfig{BaseURL: "https://short.io"}
			urlService := NewURLServiceLive(conf, tc.repo, &mockGenerator{})

			originalURL, ok := urlService.ResolveOriginalURL(tc.hash)

			assert.Equal(t, tc.expectedOk, ok, fmt.Sprintf("expected ok to be %v, got %v", tc.expectedOk, ok))
			assert.Equal(t, tc.expectedURL, originalURL, "invalid original URL")
		})
	}
}
