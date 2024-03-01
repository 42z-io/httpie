package httpie

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

// Default attributes to log for an http request
func DefaultLogRequestAttr(ctx context.Context, r *http.Request, start time.Time) []any {
	return []any{
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Any("query", r.URL.Query()),
		slog.Time("time", start),
		slog.String("remote_addr", r.RemoteAddr),
		slog.String("user_agent", r.UserAgent()),
		slog.String("referer", r.Referer()),
	}
}

// Default attributes to log for a http response
func DefaultLogResponseAttr(ctx context.Context, r *http.Request, ww *WatchedResponseWriter, start time.Time) []any {
	// Get the current time and calculate the microseconds since the start time
	now := time.Now().UTC()
	diff := now.Sub(start).Microseconds()
	return []any{
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
	}
}

// DefualtLogRequest logs the http request to a slog.Logger
func DefaultLogRequest(ctx context.Context, slogger *slog.Logger, r *http.Request, start time.Time) {
	// Log the http request
	slogger.Info("http.request", DefaultLogRequestAttr(ctx, r, start)...)
}

// DefaultLogResponse logs the http response to a slog.Logger
func DefaultLogResponse(ctx context.Context, slogger *slog.Logger, r *http.Request, ww *WatchedResponseWriter, start time.Time) {
	// Log the http response
	slogger.Info("http.response", DefaultLogResponseAttr(ctx, r, ww, start)...)
}

// LoggingOpts are the options for the LoggingMiddleware
type LoggingOpts struct {
	// Should the request be logged?
	LogRequest bool
	// Should the response be logged?
	LogResponse bool
	// Handler to log the request
	OnRequest func(ctx context.Context, slogger *slog.Logger, r *http.Request, start time.Time)
	// Handler to log the response
	OnResponse func(ctx context.Context, slogger *slog.Logger, r *http.Request, ww *WatchedResponseWriter, start time.Time)
	// SetupContext is a function to setup the context before the request is logged, useful for things like user that might be set later
	SetupContext func(ctx context.Context) context.Context
}

// Default logging options
var DefaultLoggingOpts = LoggingOpts{
	LogRequest:   true,
	LogResponse:  true,
	OnRequest:    DefaultLogRequest,
	OnResponse:   DefaultLogResponse,
	SetupContext: nil,
}

// LoggingMiddleware logs the request and response of an http handler to a slog.Logger
func LoggingMiddleware(slogger *slog.Logger, opts ...LoggingOpts) func(http.Handler) http.Handler {
	var opt LoggingOpts
	if len(opts) > 0 {
		opt = opts[0]
	} else {
		opt = DefaultLoggingOpts
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if opt.SetupContext != nil {
				ctx = opt.SetupContext(ctx)
			}
			var start time.Time
			if opt.LogRequest || opt.LogResponse {
				start = time.Now().UTC()
			}
			if opt.LogRequest {
				opt.OnRequest(ctx, slogger, r, start)
			}
			ww := NewWatchedResponseWriter(w)
			next.ServeHTTP(ww, r.WithContext(ctx))
			ww.Apply()
			if opt.LogResponse {
				opt.OnResponse(ctx, slogger, r, ww, start)
			}
		})
	}
}
