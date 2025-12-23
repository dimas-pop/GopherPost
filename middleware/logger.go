package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		reqID, _ := r.Context().Value(RequestIDKey).(string)

		slog.InfoContext(r.Context(), "Incoming Request",
			"method", r.Method,
			"path", r.URL.Path,
			"request_id", reqID,
		)

		next.ServeHTTP(w, r)

		slog.InfoContext(r.Context(), "Request Completed",
			"method", r.Method,
			"path", r.URL.Path,
			"request_id", reqID,
			"duration", time.Since(start).String(),
		)
	})
}
