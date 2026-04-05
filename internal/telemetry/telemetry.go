package telemetry

import (
	"context"
	"io"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Telemetry struct {
	TracerProvider io.Closer
}

func Init(ctx context.Context, serviceName string) (*Telemetry, error) {
	tp, err := newTracerProvider(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tp)

	return &Telemetry{
		TracerProvider: tp,
	}, nil
}

func (t *Telemetry) Shutdown(ctx context.Context) error {
	if t == nil || t.TracerProvider == nil {
		return nil
	}

	return t.TracerProvider.Close()
}

func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}
