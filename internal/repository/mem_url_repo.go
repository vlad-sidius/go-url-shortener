package repository

import (
	"errors"
	"fmt"
	"sync"
)

var ErrInsertCollision = errors.New("insert collision")

type MemURLRepo struct {
	m    sync.Mutex
	data map[string]string
}

func NewMemURLRepo() *MemURLRepo {
	return &MemURLRepo{
		data: make(map[string]string),
	}
}

func (r *MemURLRepo) Put(key, value string) {
	r.m.Lock()
	defer r.m.Unlock()

	r.data[key] = value
}

func (r *MemURLRepo) TryPut(key, value string) error {
	r.m.Lock()
	defer r.m.Unlock()

	if _, ok := r.data[key]; ok {
		return fmt.Errorf("%w: key %q already exists", ErrInsertCollision, key)
	}

	r.data[key] = value

	return nil
}

func (r *MemURLRepo) Get(key string) (string, bool) {
	r.m.Lock()
	defer r.m.Unlock()

	value, ok := r.data[key]
	return value, ok
}
