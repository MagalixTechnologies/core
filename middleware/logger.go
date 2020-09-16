package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"time"

	logger "github.com/MagalixTechnologies/core/logger"
)

const loggerKey string = "logger"
const requestIDKey string = "requestID"

// Level log level
type Level int8

// Supported log levels
const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
)

var log logger.Logger

func getLogger(level Level) logger.Logger {
	if log == nil {
		log = logger.New(logger.Level(level))
	}
	return log
}

func GetLoggerFromContext(c context.Context) (logger.Logger, bool) {
	l := c.Value(loggerKey)
	if l == nil {
		return nil, false
	}
	return l.(logger.Logger), true
}

func Log(level Level) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := getRequestId(r)
			sugar := getLogger(level)
			sugarLogger := sugar.With("requestId", reqID)

			ctx := context.WithValue(r.Context(), loggerKey, sugarLogger)
			ctx = context.WithValue(ctx, requestIDKey, reqID)

			started := time.Now()
			rw := captureResponse(w)
			h.ServeHTTP(rw, r.WithContext(ctx))
			defer sugarLogger.Sync()
			sugarLogger.Infow("Default Log",
				"method", r.Method,
				"endpoint", r.URL.String(),
				"StatusCode", r.Response.StatusCode,
				"duration", time.Since(started).String(),
			)
		})
	}
}

func getRequestId(r *http.Request) interface{} {
	reqID := r.Header.Get("X-Request-Id")
	if reqID == "" {
		reqID = shortID()
	}

	return reqID
}

func shortID() string {
	b := make([]byte, 6)
	io.ReadFull(rand.Reader, b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// ResponseCapture is a http.ResponseWriter which captures the response status
// code and content length.
type ResponseCapture struct {
	http.ResponseWriter
	StatusCode    int
	ContentLength int
}

// CaptureResponse creates a ResponseCapture that wraps the given ResponseWriter.
func captureResponse(w http.ResponseWriter) *ResponseCapture {
	return &ResponseCapture{ResponseWriter: w}
}
