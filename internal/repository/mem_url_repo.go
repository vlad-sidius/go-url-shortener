package repository

import (
	"errors"
	"fmt"
	"sync"

	"github.com/vlad-sidius/go-url-shortener/internal/model"
	"go.uber.org/zap"
)

const retriesLimit = 10

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
