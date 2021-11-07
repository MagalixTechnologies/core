package tracing

import (
	"fmt"
	"math"
	"net/http"
	"runtime/debug"
	"sync/atomic"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Config struct {
	// Service details to record
	ServiceName    string
	ServiceVersion string

	// Operation name
	//
	// The span operation name record for the http request trace.
	// Default value if empty is set to "http.request"
	OperationName string

	// SampleRate sample rate, value between 0 to 1.0
	//
	// Only span a percentage of the spans. Default value is
	// set to 1.0
	SampleRate float64

	// Skip particular requests from the Tracer
	SkipFunc func(r *http.Request) bool

	// SetTagFunc
	//
	// ...
	// SetTagFunc func(r *http.Request) map[string]interface{}

	// Tags
	//
	// Extra tags to include with a span
	Tags map[string]interface{}
}

func Tracer(tr opentracing.Tracer, cfg Config) func(next http.Handler) http.Handler {
	if cfg.OperationName == "" {
		cfg.OperationName = "http.request"
	}
	if cfg.SampleRate == 0 || cfg.SampleRate > 1 {
		cfg.SampleRate = 1.0
	}

	count := int64(0)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cfg.OperationName = getOperationName(r)
			if cfg.SkipFunc != nil && cfg.SkipFunc(r) {
				next.ServeHTTP(w, r)
				return
			}

			// Sample portion of requests, or 1.0 will sample all
			atomic.AddInt64(&count, 1)
			if math.Mod(float64(count)*cfg.SampleRate, 1.0) != 0 {
				next.ServeHTTP(w, r)
				return
			}
			atomic.StoreInt64(&count, 0)
			// Pass request through Tracer
			serverSpanCtx, _ := tr.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
			span, traceCtx := opentracing.StartSpanFromContextWithTracer(r.Context(), tr, cfg.OperationName, ext.RPCServerOption(serverSpanCtx))
			defer span.Finish()

			setHeadersTags(span, r)

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
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			r = r.WithContext(opentracing.ContextWithSpan(traceCtx, span))
			Inject(span, r)
			next.ServeHTTP(ww, r)

			// set the resource name as we get it only once the handler is executed
			// resourceName := chi.RouteContext(r.Context()).RoutePattern()
			// if resourceName == "" {
			// 	resourceName = r.URL.Path
			// }

			// set the status code
			status := ww.Status()
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
		})
	}
}
