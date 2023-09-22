package middleware

import (
	"companies-service/pkg/metric"
	"time"

	"github.com/gin-gonic/gin"
)

// Prometheus metrics middleware
func (mw *MiddlewareManager) MetricsMiddleware(metrics metric.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		status := c.Writer.Status()

		metrics.ObserveResponseTime(status, c.Request.Method, c.Request.URL.Path,
			time.Since(start).Seconds())
		metrics.IncHits(status, c.Request.Method, c.Request.URL.Path)
	}
}
