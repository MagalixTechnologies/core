package tracing

import (
	"fmt"
	"math"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func getOperationName(r *http.Request) string {
	reg := regexp.MustCompile("[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}")
	url := r.URL.Path
	url = string(reg.ReplaceAll([]byte(url), []byte(":id")))
	return fmt.Sprintf("%s::%s", r.Method, url)
}
func setHeadersTags(span opentracing.Span, r *http.Request) {
	for k, v := range r.Header {
		span.SetTag(fmt.Sprintf("header.%s", k), strings.Join(v, ", "))
	}
}
func GinTracerMiddleware(tr opentracing.Tracer, cfg Config) gin.HandlerFunc {
	if cfg.OperationName == "" {
		cfg.OperationName = "http.request"
	}
	if cfg.SampleRate == 0 || cfg.SampleRate > 1 {
		cfg.SampleRate = 1.0
	}

	count := int64(0)
	return func(c *gin.Context) {
		r := c.Request
		// Skip Tracer
		cfg.OperationName = getOperationName(r)
		if cfg.SkipFunc != nil && cfg.SkipFunc(r) {
			c.Next()
			return
		}
		// Sample portion of requests, or 1.0 will sample all
		atomic.AddInt64(&count, 1)
		if math.Mod(float64(count)*cfg.SampleRate, 1.0) != 0 {
			c.Next()
			return
		}
		atomic.StoreInt64(&count, 0)

		// Pass request through Tracer
		serverSpanCtx, _ := tr.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span, traceCtx := opentracing.StartSpanFromContextWithTracer(r.Context(), tr, cfg.OperationName, ext.RPCServerOption(serverSpanCtx))
		setHeadersTags(span, r)
		defer span.Finish()

		defer func() {
			if err := recover(); err != nil {
				ext.HTTPStatusCode.Set(span, uint16(500))
				ext.Error.Set(span, true)
				span.SetTag("error.type", "panic")
				span.LogKV(
					"event", "error",
					"error.kind", "panic",
					"message", err,
					"stack", string(debug.Stack()),
				)
				span.Finish()

				panic(err)
			}
		}()

		for k, v := range cfg.Tags {
			span.SetTag(k, v)
		}
		span.SetTag("service.name", cfg.ServiceName)
		span.SetTag("version", cfg.ServiceVersion)
		/*			ext.SpanKind.Set(span, ext.SpanKindRPCServerEnum)
		 */ext.HTTPMethod.Set(span, r.Method)
		ext.HTTPUrl.Set(span, r.URL.Path)

		resourceName := r.URL.Path
		resourceName = r.Method + " " + resourceName
		span.SetTag("resource.name", resourceName)

		// pass the span through the request context and serve the request to the next middleware
		r = r.WithContext(opentracing.ContextWithSpan(traceCtx, span))
		Inject(span, r)
		c.Next()

		// set the status code
		status := c.Writer.Status()
		ext.HTTPStatusCode.Set(span, uint16(status))

		if status >= 500 && status < 600 {
			// mark 5xx server error
			ext.Error.Set(span, true)
			span.SetTag("error.type", fmt.Sprintf("%d: %s", status, http.StatusText(status)))
			span.LogKV(
				"event", "error",
				"message", fmt.Sprintf("%d: %s", status, http.StatusText(status)),
			)
		}
	}
}
