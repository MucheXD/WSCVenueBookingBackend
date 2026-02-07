package logger

import (
	"log/slog"
	"os"
)

func InitLogger() {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})
	slog.SetDefault(slog.New(h))
}
