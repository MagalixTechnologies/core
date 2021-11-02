package tracing

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

type TracerWithConfig struct {
	Config *config.Configuration
	Tracer opentracing.Tracer
}

// Init returns an instance of Jaeger Tracer.
func Init(service string) (TracerWithConfig, io.Closer) {
	cfg, err := config.FromEnv()
	if err != nil {
		panic(fmt.Sprintf("ERROR: failed to read config from env vars: %v\n", err))
	}
	cfg.ServiceName = service
	tracer, closer, err := cfg.NewTracer()
	//tracer, closer, err := cfg.NewTracer()

	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return TracerWithConfig{
		Config: cfg,
		Tracer: tracer,
	}, closer
}
