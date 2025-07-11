package middleware

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

type contextKey string

const (
	RequestIDKey contextKey = "requestID"
	LoggerKey    contextKey = "logger"
)

func LoggerMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestID := uuid.New().String()

		requestLogger := logger.With("requestID", requestID)

		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		ctx = context.WithValue(ctx, LoggerKey, requestLogger)
		r = r.WithContext(ctx)

		requestLogger.Info("Request started",
			"method", r.Method,
			"path", r.URL.Path,
		)

		rw := &responseWriter{ResponseWriter: w}

		start := time.Now()
		next.ServeHTTP(rw, r)

		requestLogger.Info("Request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.statusCode,
			"duration", time.Since(start),
		)
	})
}

// responseWriter для перехвата статус-кода ответа
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// GetRequestID извлекает requestID из контекста
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// GetLogger извлекает логгер из контекста
func GetLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(LoggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}
