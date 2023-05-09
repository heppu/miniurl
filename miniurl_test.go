package miniurl_test

import (
	"fmt"
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

func ExampleHash() {
	const input = "https://github.com/heppu"
	hash := miniurl.Hash(input)
	fmt.Println(hash)
	// output:
	// c04b9f2c60bbb4150abaf1f317d07fc1
}
