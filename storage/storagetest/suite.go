package storagetest

import (
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/heppu/miniurl/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Storage interface {
	AddUrl(url, hash string) error
	GetUrl(hash string) (url string, err error)
}

type SuiteFn func(t *testing.T, s Storage)

func RunSuite(t *testing.T, s Storage) {
	tests := []SuiteFn{
		AddAndGet,
		NotFound,
		HashCollisionWithDifferentUrl,
		HashCollisionWithSameUrl,
		Parallel_AddAndGet,
	}

	for _, fn := range tests {
		parts := strings.Split(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), ".")
		name := parts[len(parts)-1]
		t.Run(name, func(t *testing.T) {
			fn(t, s)
		})
	}
}

func AddAndGet(t *testing.T, s Storage) {
	const (
		url  = "https://google.com"
		hash = "hash-1"
	)

	err := s.AddUrl(url, hash)
	require.NoError(t, err)

	gotUrl, err := s.GetUrl(hash)
	assert.NoError(t, err)
	assert.Equal(t, url, gotUrl)
}

func NotFound(t *testing.T, s Storage) {
	url, err := s.GetUrl("hash")
	assert.ErrorIs(t, err, storage.ErrUrlNotFound)
	assert.Empty(t, url)
}

func HashCollisionWithDifferentUrl(t *testing.T, s Storage) {
	const (
		url1 = "https://google.com"
		url2 = "https://bing.com"
		hash = "hash-2"
	)

	err := s.AddUrl(url1, hash)
	require.NoError(t, err)
	err = s.AddUrl(url2, hash)
	require.ErrorIs(t, err, storage.ErrHashCollision)
}

func HashCollisionWithSameUrl(t *testing.T, s Storage) {
	const (
		url  = "https://google.com"
		hash = "hash-3"
	)

	err := s.AddUrl(url, hash)
	require.NoError(t, err)
	err = s.AddUrl(url, hash)
	require.NoError(t, err)
}

func Parallel_AddAndGet(t *testing.T, s Storage) {
	data := map[string]string{
		"hash-a": "url-a",
		"hash-b": "url-b",
		"hash-c": "url-c",
	}

	wg := sync.WaitGroup{}
	for hash, url := range data {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(url, hash string) {
				defer wg.Done()

				err := s.AddUrl(url, hash)
				require.NoError(t, err)
				gotUrl, err := s.GetUrl(hash)
				assert.NoError(t, err)
				assert.Equal(t, url, gotUrl)
			}(url, hash)
		}
	}

	wg.Wait()
}
