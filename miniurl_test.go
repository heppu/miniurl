package miniurl_test

import (
	"fmt"
	"testing"

	"github.com/heppu/miniurl"
	"github.com/stretchr/testify/assert"
)

func TestHashLength(t *testing.T) {
	const (
		originalUrl    = "https://github.com/heppu"
		expectedLength = 32
	)

	hash := miniurl.Hash(originalUrl)
	assert.Len(t, hash, expectedLength)
}

func TestHashIsDeterministic(t *testing.T) {
	const originalUrl = "https://github.com/heppu"

	h1 := miniurl.Hash(originalUrl)
	h2 := miniurl.Hash(originalUrl)
	assert.Equal(t, h1, h2)
}

func ExampleHash() {
	originalUrl := "https://github.com/heppu"
	hash := miniurl.Hash(originalUrl)
	fmt.Println(hash)
	// output:
	// c04b9f2c60bbb4150abaf1f317d07fc1
}

func BenchmarkHash(b *testing.B) {
	input := "https://github.com/heppu"
	for n := 0; n < b.N; n++ {
		miniurl.Hash(input)
	}
}
