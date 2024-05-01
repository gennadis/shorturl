package memstore

import "github.com/gennadis/shorturl/internal/app/storage"

type MemStore struct {
	data map[string]string
}

func New() *MemStore {
	return &MemStore{
		data: make(map[string]string),
	}
}

func (m *MemStore) Read(key string) (string, error) {
	originalURL, ok := m.data[key]
	if !ok {
		return "", storage.ErrorUnknownSlugProvided
	}
	return originalURL, nil
}

func (m *MemStore) Write(key string, value string) error {
	if key == "" {
		return storage.ErrorEmptySlugProvided
	}
	m.data[key] = value
	return nil
}
