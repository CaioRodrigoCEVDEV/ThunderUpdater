package logger

import (
	"io"
	"log/slog"
	"os"
)

// New cria um logger central usando slog
func New(level slog.Level, w io.Writer) *slog.Logger {
	if w == nil {
		w = os.Stdout
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(w, opts)
	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger
}
