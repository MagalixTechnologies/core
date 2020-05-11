package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"time"

	log "github.com/MagalixTechnologies/core/logger"
	http_middleware "goa.design/goa/v3/http/middleware"
)

const loggerKey string = "logger"
const requestIDKey string = "requestID"

var logger log.Logger

func getSugarLogger(level log.Level) log.Logger {
	if logger == nil {
		logger = log.New(level)
	}
	return logger
}

func GetLoggerFromContext(c context.Context) (log.Logger, bool) {
	l := c.Value(loggerKey)
	if l == nil {
		return nil, false
	}
	return l.(log.Logger), true
}

func Log(level log.Level) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := getRequestId(r)
			sugar := getSugarLogger(level)
			sugarLogger := sugar.With("requestId", reqID)

			ctx := context.WithValue(r.Context(), loggerKey, sugarLogger)
			ctx = context.WithValue(ctx, requestIDKey, reqID)

			started := time.Now()
			rw := http_middleware.CaptureResponse(w)
			h.ServeHTTP(rw, r.WithContext(ctx))
			defer sugarLogger.Sync()
			sugarLogger.Infow("Default Log",
				"method", r.Method,
				"url", r.URL.String(),
				"status", rw.StatusCode,
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
