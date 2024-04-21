package logger

import (
	"io"
	"log/slog"
)

const (
	Local = iota
	Prod
)

func New(env int, w io.Writer) *slog.Logger {
	switch env {
	case Local:
		return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case Prod:
		return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{
			Level: slog.LevelError,
		}))
	default:
		return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}
}
