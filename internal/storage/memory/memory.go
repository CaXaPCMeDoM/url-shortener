package memory

import (
	"sync"
	"url-shortener/internal/storage"
)

type Storage struct {
	mu   sync.RWMutex
	urls map[string]string
}

func New() *Storage {
	return &Storage{
		urls: make(map[string]string),
	}
}

func (s *Storage) SaveUrl(urlAbsolute string, alias string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.urls[alias]; exists {
		return storage.ErrURLExists
	}

	s.urls[alias] = urlAbsolute
	return nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.urls[alias]
	if !exists {
		return "", storage.ErrURLNotFound
	}

	return url, nil
}
