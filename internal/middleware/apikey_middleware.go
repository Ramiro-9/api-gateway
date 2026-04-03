package middleware

import (
	"github.com/Ramiro-9/api-gateway/internal/config"
	"github.com/gin-gonic/gin"
)

// InjectAPIKey inyecta la API key interna en cada request que sale hacia los servicios
func InjectAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Header.Set("X-Internal-API-Key", config.Cfg.InternalAPIKey)
		c.Next()
	}
}
