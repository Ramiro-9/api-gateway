package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ServiceMetrics struct {
	TotalRequests  int64
	SuccessCount   int64
	ErrorCount     int64
	TotalLatencyMs int64
}

var (
	metricsMu sync.Mutex
	metrics   = map[string]*ServiceMetrics{}
	startTime = time.Now()
)

func getOrCreate(service string) *ServiceMetrics {
	if _, ok := metrics[service]; !ok {
		metrics[service] = &ServiceMetrics{}
	}
	return metrics[service]
}

func TrackMetrics(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		metricsMu.Lock()
		defer metricsMu.Unlock()

		m := getOrCreate(service)
		m.TotalRequests++
		m.TotalLatencyMs += latency.Milliseconds()

		if c.Writer.Status() >= 400 {
			m.ErrorCount++
		} else {
			m.SuccessCount++
		}
	}
}

func MetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		metricsMu.Lock()
		defer metricsMu.Unlock()

		result := map[string]interface{}{}
		for service, m := range metrics {
			avgLatency := int64(0)
			if m.TotalRequests > 0 {
				avgLatency = m.TotalLatencyMs / m.TotalRequests
			}
			result[service] = map[string]interface{}{
				"total_requests": m.TotalRequests,
				"success_count":  m.SuccessCount,
				"error_count":    m.ErrorCount,
				"avg_latency_ms": avgLatency,
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"uptime_seconds": time.Since(startTime).Seconds(),
			"services":       result,
		})
	}
}
