package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	
	"github.com/flack/chaos-engineering-as-a-platform/pkg/monitoring"
)

// MetricsMiddleware adds Prometheus metrics to API requests
func MetricsMiddleware(metrics *monitoring.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Process the request
		c.Next()
		
		// Record metrics after the request is processed
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()
		
		// Record request count
		metrics.APIRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			string(rune(status)),
		).Inc()
		
		// Record request duration
		metrics.APIRequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		).Observe(duration)
	}
}