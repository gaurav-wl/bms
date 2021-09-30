package middlewareProvider

import (
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	HttpRequestMethod     = "http.request.method"
	HttpRequestUserAgent  = "http.request.user_agent"
	HttpRequestScheme     = "http.request.scheme"
	HttpRequestRemoteAddr = "http.request.remote_addr"
	HttpRequestURI        = "http.request.uri"
	HttpRequestDuration   = "http.request_duration"
	HttpResponseStatus    = "http.response.status"
	HttpResponseSize      = "http.response.size"
	Stack                 = "stack"
	Error                 = "error"
)

// StructuredLogger is a simple but powerful middleware logger backed by logrus.
type StructuredLogger struct{}

func NewStructuredLogger() *StructuredLogger {
	return &StructuredLogger{}
}

// NewLogEntry creates a new middleware.logEntry based on the contents of the
// http.Request.
func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	fieldMap := map[string]interface{}{
		HttpRequestScheme:     scheme,
		HttpRequestMethod:     r.Method,
		HttpRequestRemoteAddr: r.RemoteAddr,
		HttpRequestUserAgent:  r.UserAgent(),
		HttpRequestURI:        fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI),
	}

	return &StructuredLoggerEntry{
		Fields: fieldMap,
	}
}

// StructuredLoggerEntry is a single log entry in a StructuredLogger.
type StructuredLoggerEntry struct {
	Fields logrus.Fields
}

// Write is run at the end of a requests processing.
func (e *StructuredLoggerEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, _ interface{}) {
	fieldMap := e.Fields

	fieldMap[HttpResponseSize] = bytes
	fieldMap[HttpResponseStatus] = status
	fieldMap[HttpRequestDuration] = float64(elapsed.Nanoseconds() / 1000000)

	logrus.WithFields(fieldMap).Infof("complete URI [%s]", fieldMap[HttpRequestURI])
}

// Panic is run when a request panics during its processing.
func (e *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	fieldMap := e.Fields
	fieldMap[Stack] = string(stack)
	fieldMap[Error] = v
	logrus.WithFields(fieldMap).Error("request panicked!")
}
