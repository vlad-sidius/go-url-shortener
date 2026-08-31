package service

import (
	"errors"
	"fmt"
	"net/url"
)

const retriesLimit = 10

type urlRepo interface {
	Put(key, value string)
	Get(key string) (string, bool)
}

type handlerConfig interface {
	BaseURL() string
}

type hashGenerator interface {
	Generate() (string, error)
}

type URLServiceLive struct {
	conf      handlerConfig
	repo      urlRepo
	generator hashGenerator
}

func NewURLServiceLive(conf handlerConfig, repo urlRepo, generator hashGenerator) *URLServiceLive {
	return &URLServiceLive{conf, repo, generator}
}

func (s *URLServiceLive) CreateShortCode(originalURL string) (string, error) {
	hashCode, err := s.generateShortCode(retriesLimit)
	if err != nil {
		return "", err
	}

	s.repo.Put(hashCode, originalURL)

	base, err := url.Parse(s.conf.BaseURL())
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
func (s *URLServiceLive) generateShortCode(retriesLeft int) (string, error) {
	if retriesLeft == 0 {
		return "", errors.New("failed to generate unique hash, too many collisions")
	}

	hash, err := s.generator.Generate()
	if err != nil {
		return "", err
	}

	if _, ok := s.repo.Get(hash); ok {
		return s.generateShortCode(retriesLeft - 1)
	}

	return hash, nil
}
