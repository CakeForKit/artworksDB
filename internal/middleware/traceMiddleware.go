package middleware

import (
	"github.com/CakeForKit/artworksDB.git/internal/tracing"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

func TraceMiddleware(tracer *tracing.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		spanName := "HTTP " + c.Request.Method + " " + c.Request.URL.Path
		ctx, span := tracer.StartSpan(c.Request.Context(), spanName)
		c.Request = c.Request.WithContext(ctx)
		c.Next()

		if span != nil {
			span.SetAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.path", c.Request.URL.Path),
				attribute.Int("http.status_code", c.Writer.Status()),
				attribute.String("http.route", c.FullPath()),
			)
			if c.Writer.Status() >= 400 {
				span.SetAttributes(attribute.Bool("error", true))
			}
			span.End()
		}
	}
}
