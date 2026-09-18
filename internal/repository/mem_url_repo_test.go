package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vlad-sidius/go-url-shortener/internal/model"
	"go.uber.org/zap"
)

func TestURLMemRepo_PutAndGet(t *testing.T) {
	repo := NewMemURLRepo(zap.NewNop())

	repo.Put(model.NewShortURLModel("abc123", "https://example.com"))
	record, ok := repo.Get("abc123")
	assert.True(t, ok)
	assert.Equal(t, "https://example.com", record.OriginalURL)

	_, ok = repo.Get("nonexistent")
	assert.False(t, ok)

	repo.Put(model.NewShortURLModel("", "https://empty.com"))
	record, ok = repo.Get("")
	assert.True(t, ok)
	assert.Equal(t, "https://empty.com", record.OriginalURL)
}

func TestURLMemRepo_Overwrite(t *testing.T) {
	repo := NewMemURLRepo(zap.NewNop())

	repo.Put(model.NewShortURLModel("test", "https://first.com"))
	record, ok := repo.Get("test")
	assert.True(t, ok)
	assert.Equal(t, "https://first.com", record.OriginalURL)

	repo.Put(model.NewShortURLModel("test", "https://second.com"))
	record, ok = repo.Get("test")
	assert.True(t, ok)
	assert.Equal(t, "https://second.com", record.OriginalURL)
}
