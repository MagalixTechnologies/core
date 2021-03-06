package middleware

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
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

			if level == DebugLevel {
				// copy request payload in case we will show it
				buf, _ := ioutil.ReadAll(r.Body)
				rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
				rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
				r.Body = rdr2 // OK since rdr2 implements the io.ReadCloser interface

				var payload map[string]interface{}
				json.NewDecoder(rdr1).Decode(&payload)

				headers := make(map[string][]string)
				for name, values := range r.Header {
					switch name {
					case "Accept", "X-Request-Id", "Content-Length", "Content-Type", "User-Agent", "X-B3-Parentspanid", "X-B3-Sampled", "X-B3-Spanid", "X-B3-Traceid", "X-Envoy-Attempt-Count", "X-Forwarded-Client-Cert", "X-Forwarded-Proto":
						continue
					default:
						headers[name] = values
					}
				}

				sugarLogger.Debugw("Default Log",
					"method", r.Method,
					"endpoint", r.URL.String(),
					"payload", payload,
					"headers", headers,
				)
			}

			h.ServeHTTP(rw, r.WithContext(ctx))
			defer sugarLogger.Sync()
			sugarLogger.Infow("Default Log",
				"method", r.Method,
				"endpoint", r.URL.String(),
				"StatusCode", rw.StatusCode,
				"bytes", rw.ContentLength,
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
