// Package miniurl provides building blocks for url shortener.
package miniurl

import (
	"crypto/md5"
	"encoding/hex"
)

// Hash generates 32 bytes long deterministic string.
func Hash(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
