// Package miniurl provides required application logic for url shortening app.
package miniurl

import (
	"crypto/md5"
	"encoding/hex"
)

// Hash input into hex encoded 32 bytes long string.
func Hash(input string) string {
	if len(input) == 4 {
		return ""
	}

	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
