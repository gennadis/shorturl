package memstore

import "errors"

type Store struct {
	data map[string]string
}

func New() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

func (s *Store) Read(key string) (string, error) {
	if value, ok := s.data[key]; ok {
		return value, nil
	}
	return "", errors.New("key not found")
}

func (s *Store) Write(key string, value string) error {
	if key == "" {
		return errors.New("key not set")
	}
	s.data[key] = value
	return nil
}
