package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLMemRepo_PutAndGet(t *testing.T) {
	repo := NewMemURLRepo()

	repo.Put("abc123", "https://example.com")
	url, ok := repo.Get("abc123")
	assert.True(t, ok)
	assert.Equal(t, "https://example.com", url)

	_, ok = repo.Get("nonexistent")
	assert.False(t, ok)

	repo.Put("", "https://empty.com")
	url, ok = repo.Get("")
	assert.True(t, ok)
	assert.Equal(t, "https://empty.com", url)
}

func TestURLMemRepo_Overwrite(t *testing.T) {
	repo := NewMemURLRepo()

	repo.Put("test", "https://first.com")
	url, ok := repo.Get("test")
	assert.True(t, ok)
	assert.Equal(t, "https://first.com", url)

	repo.Put("test", "https://second.com")
	url, ok = repo.Get("test")
	assert.True(t, ok)
	assert.Equal(t, "https://second.com", url)
}
