package httpie

import (
	"context"
	"database/sql/driver"
	"log/slog"
	"net/http"
	"strings"
)

type ctxKey int

var TransactionCtxKey ctxKey = 0

func TransactionalMiddleware(getTx func(ctx context.Context) (driver.Tx, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Debug("middleware.Transactional", slog.String("state", "start"))

			// Skip HTTP methods that don't need a transaction
			if r.Method == "GET" || r.Method == "OPTIONS" || r.Method == "HEAD" || r.Method == "TRACE" {
				slog.Debug("middleware.Transactional", slog.String("state", "skip"), slog.Any("method", r.Method))
				next.ServeHTTP(w, r)
				return
			}

			// Wrap the response writer to capture the status code
			ww := NewWatchedResponseWriter(w)

			defer func() {
				ww.Apply()
			}()

			slog.Debug("middleware.Transactional", slog.String("state", "begin"))
			// Begin the transaction
			tx, err := getTx(r.Context())
			if err != nil {
				slog.Error("middleware.Transactional", slog.String("state", "begin"), slog.Any("err", err))
				WriteErr(ww, err)
				return
			}

			// Rollback the transaction at the end of the request - if we comitted this is fine
			defer func() {
				err := tx.Rollback()
				if err != nil {
					if !strings.Contains(err.Error(), "already been committed") {
						slog.Error("middleware.Transactional", slog.String("state", "rollback"), slog.Any("err", err))
						ww.Reset()
						WriteErr(ww, err)
					}
				}
			}()

			// Attach the transaction to the context so it can be used in downstream handlers
			ctx := context.WithValue(r.Context(), TransactionCtxKey, tx)

			// Call the next http handler
			next.ServeHTTP(ww, r.WithContext(ctx))

			// If we hit an error in the http handler then we don't want to commit the transaction
			statusCode := ww.StatusCode()
			if statusCode >= 400 {
				slog.Error("middleware.Transactional", slog.String("state", "request"), slog.Int("status", statusCode))
				return
			}

			// Commit the transaction
			err = tx.Commit()
			if err != nil {
				slog.Error("middleware.Transactional", slog.String("state", "commit"), slog.Any("err", err))
				ww.Reset()
				WriteErr(ww, err)
				return
			}

			slog.Debug("middleware.Transactional", slog.String("state", "end"))
		})
	}
}
