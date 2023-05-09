package miniurl_test

import (
	"testing"

	"github.com/heppu/miniurl"
	"github.com/stretchr/testify/assert"
)

func TestHashLength(t *testing.T) {
	const (
		input          = "https://github.com/heppu"
		expectedLength = 32
	)

	output := miniurl.Hash(input)
	assert.Len(t, output, expectedLength)
}

func TestHashIsDeterministic(t *testing.T) {
	const input = "https://github.com/heppu"

	h1 := miniurl.Hash(input)
	h2 := miniurl.Hash(input)
	assert.Equal(t, h1, h2)
}
