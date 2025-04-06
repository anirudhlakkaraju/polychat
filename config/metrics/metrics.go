package metrics

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// MetricsTracker implements MetricsTracker using slog.
type MetricsTracker struct{}

func GetMetricsTracker() *MetricsTracker {
	return &MetricsTracker{}
}

// Log applies the middleware to log request metrics.
func (t *MetricsTracker) Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		slog.InfoContext(c, "HTTP Request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", latency.Milliseconds(),
			"user_agent", c.Request.UserAgent(),
		)
	}
}
