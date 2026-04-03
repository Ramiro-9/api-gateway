package router

import (
	"net/http"
	"time"

	"github.com/Ramiro-9/api-gateway/internal/config"
	"github.com/Ramiro-9/api-gateway/internal/middleware"
	"github.com/Ramiro-9/api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
)

func Setup() *gin.Engine {
	r := gin.New()

	// Servicios con timeout y soporte para múltiples URLs (load balancing)
	authService := proxy.NewService(
		"auth-api",
		10*time.Second,
		config.Cfg.AuthAPIURL,
	)
	cryptoService := proxy.NewService(
		"crypto-etl",
		30*time.Second,
		config.Cfg.CryptoETLURL,
	)

	// Middlewares globales
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.InjectAPIKey())
	r.Use(middleware.RateLimiter(limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}))

	// Health del gateway
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"gateway": "running",
		})
	})

	// Estado de los servicios
	r.GET("/services", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"services": []interface{}{
				authService.Status(),
				cryptoService.Status(),
			},
		})
	})

	// Métricas
	r.GET("/metrics", middleware.MetricsHandler())

	// Rutas públicas — auth API
	r.Any("/auth/*path", middleware.TrackMetrics("auth-api"), authService.ProxyHandler())

	// Rutas protegidas — crypto ETL
	protected := r.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		protected.Any("/pipeline/*path", middleware.TrackMetrics("crypto-etl"), cryptoService.ProxyHandler())
		protected.Any("/scheduler/*path",
			middleware.CacheResponse(30*time.Second),
			middleware.TrackMetrics("crypto-etl"),
			cryptoService.ProxyHandler(),
		)
	}

	// Rutas solo admin
	admin := r.Group("/")
	admin.Use(middleware.AuthRequired())
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.Any("/admin/*path", middleware.TrackMetrics("crypto-etl"), cryptoService.ProxyHandler())
	}

	return r
}
