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

type TimingConfig struct {
	ServiceName  string
	CollectorURL string
	Timeout      time.Duration
}

func DefaultTimingConfig() *TimingConfig {
	return &TimingConfig{
		ServiceName:  "artworks-timing", // Отдельное имя сервиса
		CollectorURL: "jaeger:4318",
		Timeout:      10 * time.Second,
	}
}

// TimingTracer - трейсер для замера времени выполнения запросов
// Всегда включен, если инициализирован
type TimingTracer struct {
	provider *sdktrace.TracerProvider
	tracer   trace.Tracer
}

func NewTimingTracer(config *TimingConfig) (*TimingTracer, error) {
	// Создаем OTLP HTTP экспортер
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(config.CollectorURL),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP timing exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(config.ServiceName),
			attribute.String("environment", "testing"),
			attribute.String("tracer.type", "request-timing"), // Помечаем тип трейсера
		)),
	)

	tracer := tp.Tracer("timing-tracer")

	return &TimingTracer{
		provider: tp,
		tracer:   tracer,
	}, nil
}

func (t *TimingTracer) Shutdown(ctx context.Context) {
	if t.provider != nil {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		_ = t.provider.Shutdown(shutdownCtx)
	}
}

func (t *TimingTracer) StartSpan(
	ctx context.Context,
	name string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}

// StartRequestSpan - специализированный метод для HTTP запросов
func (t *TimingTracer) StartRequestSpan(
	ctx context.Context,
	spanName string,
	method string,
	path string,
) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, spanName,
		trace.WithAttributes(
			attribute.String("http.method", method),
			attribute.String("http.path", path),
			attribute.String("span.type", "timing"), // Помечаем как timing span
		),
	)
}
