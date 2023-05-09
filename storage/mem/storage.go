package mem

import (
	"sync"

	"github.com/heppu/miniurl/storage"
)

type Storage struct {
	urlsMu sync.RWMutex
	urls   map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		urls: make(map[string]string),
	}
}

func (s *Storage) AddUrl(url, hash string) error {
	s.urlsMu.Lock()
	defer s.urlsMu.Unlock()

	oldUrl, ok := s.urls[hash]
	if ok && oldUrl != url {
		return storage.ErrHashCollision
	}

	s.urls[hash] = url
	return nil
}

func (s *Storage) GetUrl(hash string) (string, error) {
	s.urlsMu.RLock()
	defer s.urlsMu.RUnlock()

	url, ok := s.urls[hash]
	if !ok {
		return "", storage.ErrUrlNotFound
	}

	return url, nil
}
