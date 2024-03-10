package main

import (
	"log/slog"

	"github.com/andyfilya/customslog"
)

func main() {
	customhandler := customslog.NewHandler(nil)
	logger := slog.New(customhandler)

	logger.Info("hello world!")
}
