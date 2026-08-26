package repository

import "sync"

type MemURLRepo struct {
	m    sync.RWMutex
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

func (r *MemURLRepo) Get(key string) (string, bool) {
	r.m.RLock()
	defer r.m.RUnlock()

	value, ok := r.data[key]
	return value, ok
}
