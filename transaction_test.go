package httpie

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TxMock struct {
	sql.Tx
	mock.Mock
}

func (t *TxMock) Rollback() error {
	args := t.Called()
	return args.Error(0)
}

func (t *TxMock) Commit() error {
	args := t.Called()
	return args.Error(0)
}

func TestMiddlewareCommit(t *testing.T) {
	t.Parallel()
	m := new(TxMock)
	m.On("Rollback").Return(nil)
	m.On("Commit").Return(nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx := r.Context().Value(TransactionCtxKey).(driver.Tx)
		tx.Commit()
		w.WriteHeader(201)
	})

	r := httptest.NewRequest("PUT", "http://example.com", nil)
	w := httptest.NewRecorder()
	middleware := TransactionalMiddleware(func(ctx context.Context) (driver.Tx, error) {
		return m, nil
	})

	middleware(handler).ServeHTTP(w, r)
	assert.Equal(t, 201, w.Code)

	m.AssertExpectations(t)
}

func TestMiddlewareErrGetTx(t *testing.T) {
	t.Parallel()
	m := new(TxMock)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx := r.Context().Value(TransactionCtxKey).(driver.Tx)
		tx.Commit()
		w.WriteHeader(201)
	})

	r := httptest.NewRequest("PUT", "http://example.com", nil)
	w := httptest.NewRecorder()
	middleware := TransactionalMiddleware(func(ctx context.Context) (driver.Tx, error) {
		return nil, errors.New("get tx error")
	})

	middleware(handler).ServeHTTP(w, r)
	assert.Equal(t, 500, w.Code)

	m.AssertExpectations(t)
}

func TestMiddlewareSkip(t *testing.T) {
	t.Parallel()
	m := new(TxMock)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx := r.Context().Value(TransactionCtxKey)
		assert.Nil(t, tx)
	})

	r := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()
	middleware := TransactionalMiddleware(func(ctx context.Context) (driver.Tx, error) {
		return m, nil
	})

	middleware(handler).ServeHTTP(w, r)

	m.AssertExpectations(t)
}

func TestMiddlewareHttpErrRollback(t *testing.T) {
	t.Parallel()
	m := new(TxMock)
	m.On("Rollback").Return(nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	})

	r := httptest.NewRequest("PUT", "http://example.com", nil)
	w := httptest.NewRecorder()
	middleware := TransactionalMiddleware(func(ctx context.Context) (driver.Tx, error) {
		return m, nil
	})

	middleware(handler).ServeHTTP(w, r)

	m.AssertExpectations(t)
}

func TestMiddlewareHttpErrCommit(t *testing.T) {
	t.Parallel()
	m := new(TxMock)
	m.On("Commit").Return(errors.New("commit error"))
	m.On("Rollback").Return(nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
	})

	r := httptest.NewRequest("PUT", "http://example.com", nil)
	w := httptest.NewRecorder()
	middleware := TransactionalMiddleware(func(ctx context.Context) (driver.Tx, error) {
		return m, nil
	})

	middleware(handler).ServeHTTP(w, r)
	m.AssertExpectations(t)

	assert.Equal(t, 500, w.Code)
}

func TestMiddlewareRollbackErr(t *testing.T) {
	t.Parallel()
	m := new(TxMock)
	m.On("Commit").Return(nil)
	m.On("Rollback").Return(errors.New("rollback error"))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
	})

	r := httptest.NewRequest("PUT", "http://example.com", nil)
	w := httptest.NewRecorder()
	middleware := TransactionalMiddleware(func(ctx context.Context) (driver.Tx, error) {
		return m, nil
	})

	middleware(handler).ServeHTTP(w, r)
	assert.Equal(t, 500, w.Code)

	m.AssertExpectations(t)
}

func TestMiddlewareRollbackErrNormal(t *testing.T) {
	t.Parallel()
	m := new(TxMock)
	m.On("Commit").Return(nil)
	m.On("Rollback").Return(errors.New("already been committed"))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
	})

	r := httptest.NewRequest("PUT", "http://example.com", nil)
	w := httptest.NewRecorder()
	middleware := TransactionalMiddleware(func(ctx context.Context) (driver.Tx, error) {
		return m, nil
	})

	middleware(handler).ServeHTTP(w, r)
	assert.Equal(t, 201, w.Code)

	m.AssertExpectations(t)
}
