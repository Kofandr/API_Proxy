package logger

import (
	"log/slog"
	"os"
)

// пока без конфига чтобы нормально его с такском сделать

func New(level string) *slog.Logger {
	otps := &slog.HandlerOptions{}

	switch level {
	case "DEBUG":
		otps.Level = slog.LevelDebug
	case "WARN":
		otps.Level = slog.LevelWarn
	case "ERROR":
		otps.Level = slog.LevelError
	default:
		otps.Level = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, otps))

}
