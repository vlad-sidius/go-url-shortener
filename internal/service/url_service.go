package service

import (
	"fmt"
	"net/url"

	"github.com/vlad-sidius/go-url-shortener/internal/config"
)

const retriesLimit = 10

type urlRepo interface {
	Put(key, value string)
	TryPut(key, value string) error
	Get(key string) (string, bool)
}

type hashGenerator interface {
	Generate() (string, error)
}

type URLServiceLive struct {
	conf      *config.URLServiceConfig
	repo      urlRepo
	generator hashGenerator
}

func NewURLServiceLive(conf *config.URLServiceConfig, repo urlRepo, generator hashGenerator) *URLServiceLive {
	return &URLServiceLive{conf, repo, generator}
}

func (s *URLServiceLive) CreateShortCode(originalURL string) (string, error) {
	hashCode, err := s.saveShortCode(originalURL, retriesLimit)
	if err != nil {
		return "", err
	}

	base, err := url.Parse(s.conf.BaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	return base.JoinPath(hashCode).String(), nil
}

func (s *URLServiceLive) ResolveOriginalURL(hash string) (string, bool) {
	originalURL, ok := s.repo.Get(hash)
	if ok {
		return originalURL, true
	}

	return "", false
}

// generateShortCode generates unique hash
func (s *URLServiceLive) saveShortCode(originalURL string, retriesLeft int) (string, error) {
	hash, err := s.generator.Generate()
	if err != nil {
		return "", err
	}

	if err := s.repo.TryPut(hash, originalURL); err != nil {
		return s.saveShortCode(originalURL, retriesLeft-1)
	}

	return hash, nil
}
