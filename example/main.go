package main

import (
	"log/slog"

	"github.com/andyfilya/customslog"
)

func main() {
	customhandler := customslog.NewHandler(nil)
	// info(customhandler)
	debug(customhandler)
	error(customhandler)
}

func info(c *customslog.CustomSlogHandler) {
	logger := slog.New(c)
	logger.Info("in function info.")
}

func debug(c *customslog.CustomSlogHandler) {
	logger := slog.New(c)
	logger.Debug("in function debug.")
}

func error(c *customslog.CustomSlogHandler) {
	logger := slog.New(c)
	logger.Error("in function error.")
}
