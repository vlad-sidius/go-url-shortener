package service

import (
	"fmt"
	"net/url"

	"github.com/vlad-sidius/go-url-shortener/internal/config"
	"github.com/vlad-sidius/go-url-shortener/internal/model"
	"go.uber.org/zap"
)

const retriesLimit = 10

type urlRepo interface {
	Put(value *model.ShortURLModel)
	TryPut(value *model.ShortURLModel) error
	Get(key string) (*model.ShortURLModel, bool)
	Save(fileName string) error
	Load(fileName string) error
}

type hashGenerator interface {
	Generate() string
}

type URLServiceLive struct {
	log       *zap.Logger
	conf      *config.URLServiceConfig
	repo      urlRepo
	generator hashGenerator
}

func NewURLServiceLive(log *zap.Logger, conf *config.URLServiceConfig, repo urlRepo, generator hashGenerator) *URLServiceLive {
	return &URLServiceLive{log, conf, repo, generator}
}

func (s *URLServiceLive) CreateShortCode(originalURL string) (string, error) {
	hashCode, err := s.saveShortCode(originalURL, retriesLimit)
	if err != nil {
		return "", err
	}

	err = s.repo.Save(s.conf.StoragePath)
	if err != nil {
		s.log.Error("Failed to save data in file storage", zap.Error(err))
	}

	base, err := url.Parse(s.conf.BaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	return base.JoinPath(hashCode).String(), nil
}

func (s *URLServiceLive) ResolveOriginalURL(hash string) (string, bool) {
	record, ok := s.repo.Get(hash)
	if ok {
		return record.OriginalURL, true
	}

	return "", false
}

// saveShortCode tries to generate unique hash and save via repository,
// it performs retriesLimit attempts
func (s *URLServiceLive) saveShortCode(originalURL string, retriesLeft int) (string, error) {
	if retriesLeft == 0 {
		return "", fmt.Errorf("failed to generate hash in %d tries", retriesLimit)
	}

	hash := s.generator.Generate()

	if err := s.repo.TryPut(model.NewShortURLModel(hash, originalURL)); err != nil {
		return s.saveShortCode(originalURL, retriesLeft-1)
	}

	return hash, nil
}
