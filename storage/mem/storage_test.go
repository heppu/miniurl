package mem_test

import (
	"testing"

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
	assert.ErrorIs(t, err, mem.ErrUrlNotFound)
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
	require.ErrorIs(t, err, mem.ErrHashCollision)
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
