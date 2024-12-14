package handlers

import (
	"ecommerce-platform/internal/metrics"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func MetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		status := c.Writer.Status()
		method := c.Request.Method

		metrics.RequestCount.WithLabelValues(method, strconv.Itoa(status)).Inc()
		metrics.RequestDuration.WithLabelValues(method, strconv.Itoa(status)).Observe(time.Since(start).Seconds())
	}
}
