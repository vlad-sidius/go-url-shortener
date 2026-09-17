package repository

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vlad-sidius/go-url-shortener/internal/model"
	"go.uber.org/zap"
)

func TestFileURLRepo_TryPutAndGet(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_url_repo_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	repo := NewFileURLRepo(zap.NewNop(), tmpFile.Name())

	record := model.NewShortURLModel("abc123", "https://example.com")
	err = repo.TryPut(record)
	require.NoError(t, err)

	got, ok := repo.Get("abc123")
	assert.True(t, ok)
	assert.Equal(t, "https://example.com", got.OriginalURL)
}

func TestFileURLRepo_TryPutCollision(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_url_repo_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	repo := NewFileURLRepo(zap.NewNop(), tmpFile.Name())

	record := model.NewShortURLModel("abc123", "https://example.com")
	err = repo.TryPut(record)
	require.NoError(t, err)

	duplicate := model.NewShortURLModel("abc123", "https://other.com")
	err = repo.TryPut(duplicate)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInsertCollision)
}

func TestFileURLRepo_LoadFromExistingFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_url_repo_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	record := model.ShortURLModel{
		UUID:        "test-uuid",
		ShortURL:    "xyz789",
		OriginalURL: "https://loaded.com",
	}
	data, err := json.Marshal(record)
	require.NoError(t, err)
	data = append(data, '\n')
	_, err = tmpFile.Write(data)
	require.NoError(t, err)
	tmpFile.Close()

	repo := NewFileURLRepo(zap.NewNop(), tmpFile.Name())

	got, ok := repo.Get("xyz789")
	assert.True(t, ok)
	assert.Equal(t, "https://loaded.com", got.OriginalURL)
}

func TestFileURLRepo_PersistenceAfterRestart(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_url_repo_*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	repo1 := NewFileURLRepo(zap.NewNop(), tmpFile.Name())
	record := model.NewShortURLModel("persist123", "https://persist.com")
	err = repo1.TryPut(record)
	require.NoError(t, err)

	repo2 := NewFileURLRepo(zap.NewNop(), tmpFile.Name())
	got, ok := repo2.Get("persist123")
	assert.True(t, ok)
	assert.Equal(t, "https://persist.com", got.OriginalURL)
}
