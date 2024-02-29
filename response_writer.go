package httpie

import (
	"bytes"
	"net/http"
)

type WatchedResponseWriter struct {
	statusCode   int
	bytesWritten int
	buffer       *bytes.Buffer
	response     http.ResponseWriter
}

func (w *WatchedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *WatchedResponseWriter) Header() http.Header {
	return w.response.Header()
}

func (w *WatchedResponseWriter) Write(b []byte) (int, error) {
	w.bytesWritten += len(b)
	return w.buffer.Write(b)
}

func (w *WatchedResponseWriter) StatusCode() int {
	return w.statusCode
}

func (w *WatchedResponseWriter) BytesWritten() int {
	return w.bytesWritten
}

func (w *WatchedResponseWriter) Apply() {
	w.response.WriteHeader(w.statusCode)
	w.response.Write(w.buffer.Bytes())
}

func (w *WatchedResponseWriter) Reset() {
	w.statusCode = 0
	w.bytesWritten = 0
	w.buffer.Reset()
}

func NewWatchedResponseWriter(response http.ResponseWriter) *WatchedResponseWriter {
	return &WatchedResponseWriter{response: response, buffer: bytes.NewBuffer([]byte{})}
}
