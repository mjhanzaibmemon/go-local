// Package tracing contains helpers to initialize and configure OpenTracing (Jaeger) tracers.
package tracing

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// ServiceName is a typed alias for the service name used by the tracer.
type ServiceName string

// InitNonSamplingTracer initializes a Jaeger tracer that does not sample spans.
func InitNonSamplingTracer(serviceName ServiceName) (opentracing.Tracer, io.Closer, error) {
	return initTracer(serviceName, &config.SamplerConfig{
		Type:  jaeger.SamplerTypeConst,
		Param: 0,
	})
}

func initTracer(serviceName ServiceName, samplerConfig *config.SamplerConfig) (opentracing.Tracer, io.Closer, error) {
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read tracer config from env vars: %w", err)
	}

	cfg.ServiceName = string(serviceName)
	cfg.Sampler = samplerConfig

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create Jaeger tracer: %w", err)
	}

	return tracer, closer, nil
}
