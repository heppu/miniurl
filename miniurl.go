// Package miniurl provides building blocks for url shortening app.
package miniurl

import (
	"crypto/md5"
	"encoding/hex"
)

// Hash produces deterministic hex encoded 32 bytes long string from input.
func Hash(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
