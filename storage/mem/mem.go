package mem

import "errors"

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
		return errors.New("hash collides with different url")
	}

	s.urls[hash] = url
	return nil
}

func (s *Storage) GetUrl(hash string) (string, error) {
	url, ok := s.urls[hash]
	if !ok {
		return "", errors.New("url not found")
	}

	return url, nil
}
