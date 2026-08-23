package repository

type MemURLRepo struct {
	data map[string]string
}

func NewMemURLRepo() *MemURLRepo {
	return &MemURLRepo{
		data: make(map[string]string),
	}
}

func (r *MemURLRepo) Put(key, value string) {
	r.data[key] = value
}

func (r *MemURLRepo) Get(key string) (string, bool) {
	value, ok := r.data[key]
	return value, ok
}
