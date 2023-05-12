// Package miniurl provides building blocks for url shortener.
package miniurl

import (
	"crypto/md5"
	"encoding/hex"
	"os"

	"github.com/heppu/miniurl/api"
	"golang.org/x/exp/slog"
)

func Run() error {
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	slog.Info("starting server", slog.String("LISTEN_ADDR", addr))
	srv := api.NewServer(addr, nil)
	return srv.Start()
}

// Hash generates 32 bytes long deterministic string.
func Hash(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
