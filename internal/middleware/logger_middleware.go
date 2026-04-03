package middleware

import (
	"fmt"
	"time"

	"github.com/Ramiro-9/api-gateway/internal/logger"
	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		msg := fmt.Sprintf(
			"[GATEWAY] %s %s | IP: %s | Status: %d | Latency: %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Writer.Status(),
			latency,
		)
		logger.Info(msg)
	}
}
