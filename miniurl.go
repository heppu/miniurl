// Package miniurl provides building blocks for url shortening app.
package miniurl

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/heppu/miniurl/api"
	"github.com/heppu/miniurl/storage"
	"github.com/heppu/miniurl/storage/mem"
	"golang.org/x/exp/slog"
)

type Storage interface {
	AddUrl(url, hash string) error
	GetUrl(hash string) (url string, err error)
}

type App struct {
	storage Storage
}

func Run() error {
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	app := &App{storage: mem.NewStorage()}
	srv := api.NewServer(addr, app)

	closeCh := make(chan os.Signal, 1)
	signal.Notify(closeCh, os.Interrupt)
	go func() {
		<-closeCh
		slog.Info("Closing server")
		if err := srv.Stop(); err != nil {
			slog.Error(err.Error())
		}
	}()

	slog.Info("Starting app", slog.String("LISTEN_ADDR", addr))
	return srv.Start()
}

func (a *App) AddUrl(url string) (string, error) {
	seed := url
	for i := 0; i < 100; i++ {
		hash := Hash(seed)[0:5]
		err := a.storage.AddUrl(url, hash)
		switch {
		case err == nil:
			return hash, nil
		case !errors.Is(err, storage.ErrHashCollision):
			return "", fmt.Errorf("failed to store url: %w", err)
		default:
			seed = hash
		}
	}

	err := fmt.Errorf("unable to get new hash for url: %s", url)
	slog.Error(err.Error())
	return "", err
}

// Hash produces deterministic hex encoded 32 bytes long string from input.
func (a *App) GetUrl(hash string) (string, error) {
	return a.storage.GetUrl(hash)
}

// Hash input into hex encoded 32 bytes long string.
func Hash(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
