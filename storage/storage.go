package storage

import "errors"

var (
	ErrHashCollision = errors.New("hash collides with different url")
	ErrUrlNotFound   = errors.New("url not found")
)
