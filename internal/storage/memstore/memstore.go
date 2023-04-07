package memstore

import (
	"errors"
)

type Storage struct {
	data map[string]string
}

func New() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

func (s *Storage) Read(key string) (string, error) {
	if value, ok := s.data[key]; ok {
		return value, nil
	}
	return "", errors.New("wrong hash")
}

func (s *Storage) Write(key string, value string) error {
	if key == "" {
		return errors.New("no key provided")
	}
	s.data[key] = value
	return nil
}
