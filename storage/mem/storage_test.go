package mem_test

import (
	"sync"
	"testing"

	"github.com/heppu/miniurl/storage"
	"github.com/heppu/miniurl/storage/mem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_AddAndGet(t *testing.T) {
	const (
		url  = "https://google.com"
		hash = "hash"
	)

	s := mem.NewStorage()
	err := s.AddUrl(url, hash)
	require.NoError(t, err)

	gotUrl, err := s.GetUrl(hash)
	assert.NoError(t, err)
	assert.Equal(t, url, gotUrl)
}

func TestStorage_NotFound(t *testing.T) {
	s := mem.NewStorage()
	url, err := s.GetUrl("hash")
	assert.ErrorIs(t, err, storage.ErrUrlNotFound)
	assert.Empty(t, url)
}

func TestStorage_HashCollisionWithDifferentUrl(t *testing.T) {
	const (
		url1 = "https://google.com"
		url2 = "https://bing.com"
		hash = "hash"
	)

	s := mem.NewStorage()
	err := s.AddUrl(url1, hash)
	require.NoError(t, err)
	err = s.AddUrl(url2, hash)
	require.ErrorIs(t, err, storage.ErrHashCollision)
}

func TestStorage_HashCollisionWithSameUrl(t *testing.T) {
	const (
		url  = "https://google.com"
		hash = "hash"
	)

	s := mem.NewStorage()
	err := s.AddUrl(url, hash)
	require.NoError(t, err)
	err = s.AddUrl(url, hash)
	require.NoError(t, err)
}

func TestStorage_Parallel_AddAndGet(t *testing.T) {
	data := map[string]string{
		"hash-a": "url-a",
		"hash-b": "url-b",
		"hash-c": "url-c",
	}

	s := mem.NewStorage()
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
