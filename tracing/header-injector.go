package tracing

import (
	"context"
	"net/http"

	"github.com/opentracing/opentracing-go"
)

func Inject(span opentracing.Span, request *http.Request) error {
	return span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(request.Header))
}

func HeaderInjector(ctx context.Context, req *http.Request) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "client")
	defer span.Finish()
	if err := Inject(span, req); err != nil {
		return err
	}
	return nil
}
