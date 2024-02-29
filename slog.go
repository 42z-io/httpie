package httpie

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

func DefaultLogRequest(ctx context.Context, slogger *slog.Logger, r *http.Request, start time.Time) {
	// Log the http request
	slogger.Info("http.request",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Any("query", r.URL.Query()),
		slog.Time("time", start),
		slog.String("remote_addr", r.RemoteAddr),
		slog.String("user_agent", r.UserAgent()),
		slog.String("referer", r.Referer()),
	)
}

func DefaultLogResponse(ctx context.Context, slogger *slog.Logger, r *http.Request, ww *WatchedResponseWriter, start time.Time) {
	// Get the current time and calculate the microseconds since the start time
	now := time.Now().UTC()
	diff := now.Sub(start).Microseconds()

	// Log the http response
	slogger.Info("http.response",
		slog.Int("status", ww.StatusCode()),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Any("query", r.URL.Query()),
		slog.Time("time", now),
		slog.String("remote_addr", r.RemoteAddr),
		slog.String("user_agent", r.UserAgent()),
		slog.String("referer", r.Referer()),
		slog.Int("size", ww.BytesWritten()),
		slog.Int64("duration", diff),
	)
}

type LoggingOpts struct {
	LogRequest  bool
	LogResponse bool
	OnRequest   func(ctx context.Context, slogger *slog.Logger, r *http.Request, start time.Time)
	OnResponse  func(ctx context.Context, slogger *slog.Logger, r *http.Request, ww *WatchedResponseWriter, start time.Time)
}

var DefaultLoggingOpts = LoggingOpts{
	LogRequest:  true,
	LogResponse: true,
	OnRequest:   DefaultLogRequest,
	OnResponse:  DefaultLogResponse,
}

func LoggingMiddleware(slogger *slog.Logger, opts ...LoggingOpts) func(http.Handler) http.Handler {
	var opt LoggingOpts
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = DefaultLoggingOpts
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var start time.Time
			if opt.LogRequest || opt.LogResponse {
				start = time.Now().UTC()
			}
			if opt.LogRequest {
				opt.OnRequest(r.Context(), slogger, r, start)
			}
			ww := NewWatchedResponseWriter(w)
			next.ServeHTTP(ww, r)
			ww.Apply()
			if opt.LogResponse {
				opt.OnResponse(r.Context(), slogger, r, ww, start)
			}
		})
	}
}
