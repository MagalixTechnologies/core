package mw

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	log "github.com/MagalixTechnologies/core/logger"
	"go.uber.org/zap"
	http_middleware "goa.design/goa/v3/http/middleware"
	"goa.design/goa/v3/middleware"
	"io"
	"net/http"
	"time"
)

var logger *zap.Logger

const loggerKey string = "logger"

func getSugarLogger(level log.Level) *zap.SugaredLogger {
	if logger == nil {
		logger = log.New(level)
	}
	return logger.Sugar()
}

func GetLoggerFromContext(c context.Context) *zap.SugaredLogger {
	return c.Value(loggerKey).(*zap.SugaredLogger)
}

func Log(level log.Level) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := getRequestId(r)
			sugar := getSugarLogger(level)
			sugarLogger := sugar.With("requestId", reqID)
			ctx := context.WithValue(r.Context(), loggerKey, sugarLogger)
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
	reqID := r.Context().Value(middleware.RequestIDKey)
	if reqID == nil {
		reqID = shortID()
	}
	return reqID
}

func shortID() string {
	b := make([]byte, 6)
	io.ReadFull(rand.Reader, b)
	return base64.RawURLEncoding.EncodeToString(b)
}
