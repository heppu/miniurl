package mem

import "errors"

var (
	ErrHashCollision = errors.New("hash collides with different url")
	ErrUrlNotFound   = errors.New("url not found")
)

type Storage struct {
	urls map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		urls: make(map[string]string),
	}
}

func (s *Storage) AddUrl(url, hash string) error {
	oldUrl, ok := s.urls[hash]
	if ok && oldUrl != url {
		return ErrHashCollision
	}

	s.urls[hash] = url
	return nil
}

func (s *Storage) GetUrl(hash string) (string, error) {
	url, ok := s.urls[hash]
	if !ok {
		return "", ErrUrlNotFound
	}

	return url, nil
}
