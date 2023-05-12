package main

import (
	"os"

	"github.com/heppu/miniurl"
	"golang.org/x/exp/slog"
)

func main() {
	err := miniurl.Run()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
