package httpie

import (
	"bytes"
	"net/http"
)

// Wraps an http.ResponseWriter and watches for changes to the response
type WatchedResponseWriter struct {
	statusCode   int
	bytesWritten int
	buffer       *bytes.Buffer
	response     http.ResponseWriter
}

// Capture the written status code
func (w *WatchedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

// Delegate the Header method to the wrapped response
func (w *WatchedResponseWriter) Header() http.Header {
	return w.response.Header()
}

// Capture the written bytes to a buffer
func (w *WatchedResponseWriter) Write(b []byte) (int, error) {
	w.bytesWritten += len(b)
	return w.buffer.Write(b)
}

// Return the captured status code
func (w *WatchedResponseWriter) StatusCode() int {
	return w.statusCode
}

// Return the number of bytes written
func (w *WatchedResponseWriter) BytesWritten() int {
	return w.bytesWritten
}

// Apply the captured status code and bytes to the wrapped response
func (w *WatchedResponseWriter) Apply() {
	w.response.WriteHeader(w.statusCode)
	w.response.Write(w.buffer.Bytes())
}

// Reset the status code, bytes written, and buffer
func (w *WatchedResponseWriter) Reset() {
	w.statusCode = 0
	w.bytesWritten = 0
	w.buffer.Reset()
}

// Create a new WatchedResponseWriter
func NewWatchedResponseWriter(response http.ResponseWriter) *WatchedResponseWriter {
	return &WatchedResponseWriter{response: response, buffer: bytes.NewBuffer([]byte{})}
}
