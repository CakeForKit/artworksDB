package tracing

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type ConfigTracer struct {
	ServiceName  string
	CollectorURL string // Изменяем JaegerURL на CollectorURL
	Enabled      bool
	Timeout      time.Duration
}

func DefaultConfigTracer() *ConfigTracer {
	return &ConfigTracer{
		ServiceName:  "artworks",
		CollectorURL: "jaeger:4318", // OTLP HTTP endpoint
		Enabled:      false,
		Timeout:      10 * time.Second,
	}
}

type Tracer struct {
	provider *sdktrace.TracerProvider
	enabled  bool
	tracer   trace.Tracer
}

func NewTracer(config *ConfigTracer) (*Tracer, error) {
	if !config.Enabled {
		return &Tracer{enabled: false}, nil
	}

	// Создаем OTLP HTTP экспортер
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(config.CollectorURL),
		otlptracehttp.WithInsecure(), // Для HTTP без TLS
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(config.ServiceName),
			attribute.String("environment", "testing"),
		)),
	)

	tracer := tp.Tracer("app-tracer")

	return &Tracer{
		provider: tp,
		enabled:  true,
		tracer:   tracer,
	}, nil
}

func (t *Tracer) Shutdown(ctx context.Context) {
	if t.provider != nil {
		// Даем время на отправку оставшихся данных
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = t.provider.Shutdown(shutdownCtx)
	}
}

func (t *Tracer) IsEnabled() bool {
	return t.enabled
}

func (t *Tracer) StartSpan(
	ctx context.Context,
	name string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	if !t.enabled {
		return ctx, nil
	}
	return t.tracer.Start(ctx, name, opts...)
}
