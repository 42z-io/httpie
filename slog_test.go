package httpie

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type requestLog struct {
	Time       time.Time
	Level      string
	Msg        string
	Method     string
	Path       string
	Query      map[string][]string
	RemoteAddr string `json:"remote_addr"`
	UserAgent  string `json:"user_agent"`
	Referer    string
}

type responseLog struct {
	requestLog
	Size     int
	Duration int
	Status   int
}

func TestLoggingRequestResponse(t *testing.T) {
	t.Parallel()
	writer := bytes.NewBufferString("")
	slogger := slog.New(slog.NewJSONHandler(writer, nil))

	middleware := LoggingMiddleware(slogger)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte("hello"))
	})

	r := httptest.NewRequest("PUT", "http://domain.com/path?query=hello", nil)
	r.RemoteAddr = "127.0.0.1"
	r.Header.Set("User-Agent", "test-agent")
	r.Header.Set("Referer", "http://google.com")
	w := httptest.NewRecorder()
	middleware(handler).ServeHTTP(w, r)

	assert.Equal(t, 201, w.Code)
	parts := strings.Split(strings.Trim(writer.String(), "\n"), "\n")
	assert.Len(t, parts, 2)
	var reqLog requestLog
	var resLog responseLog
	err := json.Unmarshal([]byte(parts[0]), &reqLog)
	assert.NoError(t, err)
	err = json.Unmarshal([]byte(parts[1]), &resLog)
	assert.NoError(t, err)

	assert.Equal(t, "http.request", reqLog.Msg)
	assert.NotEmpty(t, reqLog.Time)
	assert.Equal(t, "PUT", reqLog.Method)
	assert.Equal(t, "/path", reqLog.Path)
	assert.Equal(t, map[string][]string{"query": {"hello"}}, reqLog.Query)
	assert.Equal(t, "test-agent", reqLog.UserAgent)
	assert.Equal(t, "http://google.com", reqLog.Referer)
	assert.Equal(t, "INFO", reqLog.Level)

	assert.Equal(t, "http.response", resLog.Msg)
	assert.NotEmpty(t, resLog.Time)
	assert.Equal(t, "PUT", resLog.Method)
	assert.Equal(t, "/path", resLog.Path)
	assert.Equal(t, map[string][]string{"query": {"hello"}}, resLog.Query)
	assert.Equal(t, "test-agent", resLog.UserAgent)
	assert.Equal(t, "http://google.com", resLog.Referer)
	assert.Equal(t, "INFO", resLog.Level)
	assert.Equal(t, 201, resLog.Status)
	assert.Equal(t, 5, resLog.Size)
	assert.Greater(t, resLog.Duration, 0)
}

func TestLoggingCustomOpts(t *testing.T) {
	t.Parallel()
	writer := bytes.NewBufferString("")
	slogger := slog.New(slog.NewJSONHandler(writer, nil))

	middleware := LoggingMiddleware(slogger, LoggingOpts{
		LogRequest:  false,
		LogResponse: false,
		OnRequest:   nil,
		OnResponse:  nil,
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte("hello"))
	})

	r := httptest.NewRequest("PUT", "http://domain.com/path?query=hello", nil)
	w := httptest.NewRecorder()
	middleware(handler).ServeHTTP(w, r)

	assert.Equal(t, 201, w.Code)
	assert.Len(t, writer.String(), 0)
}
