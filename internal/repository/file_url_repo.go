package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/vlad-sidius/go-url-shortener/internal/model"
	"go.uber.org/zap"
)

const storageFilePerm = 0644

type FileURLRepo struct {
	MemURLRepo
	filePath string
}

func NewFileURLRepo(log *zap.Logger, filePath string) *FileURLRepo {
	repo := FileURLRepo{
		MemURLRepo: MemURLRepo{
			log:  log,
			data: make(map[string]*model.ShortURLModel),
		},
		filePath: filePath,
	}

	if err := repo.load(); err != nil {
		log.Error("Failed to load data from file. Continue with empty storage.")
	}

	return &repo
}

func (r *FileURLRepo) TryPut(value *model.ShortURLModel) error {
	if err := r.MemURLRepo.TryPut(value); err != nil {
		return err
	}

	if err := r.appendRecord(value); err != nil {
		return fmt.Errorf("failed to persist new record %w", err)
	}

	return nil
}

func (r *FileURLRepo) appendRecord(value *model.ShortURLModel) error {
	r.m.Lock()
	defer r.m.Unlock()

	file, err := os.OpenFile(r.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, storageFilePerm)
	if err != nil {
		r.log.Error("Failed to open file for writing", zap.String("filePath", r.filePath))
		return err
	}

	defer file.Close()

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to serialize record: %w", err)
	}

	data = append(data, '\n')

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}

	return file.Sync()
}

func (r *FileURLRepo) load() error {
	r.m.Lock()
	defer r.m.Unlock()

	file, err := os.OpenFile(r.filePath, os.O_RDONLY|os.O_CREATE, storageFilePerm)

	if err == io.EOF {
		r.log.Error("File is empty", zap.String("filePath", r.filePath))
	} else if err != nil {
		r.log.Error("Failed to open file for reading", zap.String("filePath", r.filePath))
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var rec model.ShortURLModel
		if err := json.Unmarshal(line, &rec); err != nil {
			r.log.Warn("skipping malformed JSONL line",
				zap.Error(err),
				zap.String("line", string(line)),
			)

			continue
		}

		r.data[rec.ShortURL] = &rec
	}

	return scanner.Err()
}
