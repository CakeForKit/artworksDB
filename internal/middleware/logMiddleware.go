package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type loggerCtx int

const (
	LoggerKey loggerCtx = iota
)

type MiddlewareLogger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
}

func getDebugInfo(c *gin.Context) []interface{} {
	debugInfo := []interface{}{
		"method", c.Request.Method,
		"url", c.Request.URL.String(),
		"headers", c.Request.Header,
		"query", c.Request.URL.Query(),
		"postForm", c.Request.PostForm,
		"clientIP", c.ClientIP(),
		"cookies", c.Request.Cookies(),
		"fullPath", c.FullPath(),
		"handlerName", c.HandlerName(),
		"status", c.Writer.Status(),
		"remoteAddr", c.Request.RemoteAddr,
	}
	return debugInfo
}

func LogMiddleware(projLogger MiddlewareLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		projLogger.Infow("Incoming request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			// "ip", c.ClientIP(),
			// "user-agent", c.Request.UserAgent(),
		)

		debugInfo := getDebugInfo(c)
		projLogger.Debugf("Incoming request", debugInfo...)

		// c.Set(LoggerKey, projLogger)
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, LoggerKey, projLogger)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
		// Логируем результат обработки
		latency := time.Since(start)
		status := c.Writer.Status()
		logFields := []interface{}{
			"status", status,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"latency", latency,
			"ip", c.ClientIP(),
		}
		if len(c.Errors) > 0 {
			projLogger.Errorw("Request failed",
				append(logFields, "errors", c.Errors.String())...,
			)
		} else {
			if status >= 500 {
				projLogger.Errorw("Server error", logFields...)
			} else if status >= 400 {
				projLogger.Warnw("Client error", logFields...)
			} else {
				projLogger.Infow("Request completed", logFields...)
			}
		}
	}
}
