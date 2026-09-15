package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/vlad-sidius/go-url-shortener/internal/model"
	"go.uber.org/zap"
)

const storageFilePerm = 0644

var ErrInsertCollision = errors.New("insert collision")

type MemURLRepo struct {
	m    sync.Mutex
	log  *zap.Logger
	data map[string]*model.ShortURLModel
}

func NewMemURLRepo(log *zap.Logger) *MemURLRepo {
	return &MemURLRepo{
		log:  log,
		data: make(map[string]*model.ShortURLModel),
	}
}

func (r *MemURLRepo) Put(value *model.ShortURLModel) {
	r.m.Lock()
	defer r.m.Unlock()

	r.data[value.ShortURL] = value
}

func (r *MemURLRepo) TryPut(value *model.ShortURLModel) error {
	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.data[value.ShortURL]; ok {
		return fmt.Errorf("%w: key %q already exists", ErrInsertCollision, value.ShortURL)
	}

	r.data[value.ShortURL] = value

	return nil
}

func (r *MemURLRepo) Get(key string) (*model.ShortURLModel, bool) {
	r.m.Lock()
	defer r.m.Unlock()

	value, ok := r.data[key]
	return value, ok
}

func (r *MemURLRepo) Save(filePath string) error {
	r.m.Lock()
	defer r.m.Unlock()

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, storageFilePerm)
	if err != nil {
		r.log.Error("Failed to open file for writing", zap.String("filePath", filePath))
		return err
	}

	defer file.Close()

	records := make([]model.ShortURLModel, 0, len(r.data))
	for _, rec := range r.data {
		records = append(records, *rec)
	}

	encoder := json.NewEncoder(file)
	return encoder.Encode(records)
}

func (r *MemURLRepo) Load(filePath string) error {
	r.m.Lock()
	defer r.m.Unlock()

	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, storageFilePerm)
	if err != nil {
		r.log.Error("Failed to open file for reading", zap.String("filePath", filePath))
		return err
	}

	defer file.Close()

	var records []model.ShortURLModel
	if err := json.NewDecoder(file).Decode(&records); err != nil {
		r.log.Error("Failed to parse JSON")
		return err
	}

	for _, rec := range records {
		r.data[rec.ShortURL] = &rec
	}

	return nil
}
