package middleware

import (
	"time"

	"github.com/CakeForKit/artworksDB.git/internal/tracing"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

// OpenTelemetryMiddleware создает span для каждого HTTP запроса
func MetricsMiddleware(timingTracer *tracing.TimingTracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// span для замера времени
		spanName := "http.request.timing"
		ctx, span := timingTracer.StartRequestSpan(
			c.Request.Context(),
			spanName,
			c.Request.Method,
			c.Request.URL.Path,
		)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
		duration := time.Since(start)

		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
			attribute.Float64("http.duration_ms", float64(duration.Milliseconds())),
		)
		span.End()
	}
}
