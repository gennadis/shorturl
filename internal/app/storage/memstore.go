package storage

import (
	"errors"
)

var (
	ErrorURLNotFound   = errors.New("wrong hash provided")
	ErrorNoURLProvided = errors.New("no url provided")
)

type Storage struct {
	data map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

func (s *Storage) Read(key string) (string, error) {
	if value, ok := s.data[key]; ok {
		return value, nil
	}
	return "", ErrorURLNotFound
}

func (s *Storage) Write(key string, value string) error {
	if key == "" {
		return ErrorNoURLProvided
	}
	s.data[key] = value
	return nil
}
