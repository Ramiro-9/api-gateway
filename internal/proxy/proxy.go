package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/Ramiro-9/api-gateway/internal/circuitbreaker"
	"github.com/Ramiro-9/api-gateway/internal/logger"
	"github.com/gin-gonic/gin"
)

type Service struct {
	Name    string
	URLs    []string
	current int
	CB      *circuitbreaker.CircuitBreaker
	Timeout time.Duration
}

func NewService(name string, timeout time.Duration, urls ...string) *Service {
	return &Service{
		Name:    name,
		URLs:    urls,
		CB:      circuitbreaker.New(name, 5, 30*time.Second),
		Timeout: timeout,
	}
}

// Siguiente URL en round-robin
func (s *Service) nextURL() string {
	url := s.URLs[s.current%len(s.URLs)]
	s.current++
	return url
}

func (s *Service) ProxyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !s.CB.Allow() {
			logger.Error(fmt.Sprintf("Circuit breaker abierto para %s", s.Name))
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "servicio temporalmente no disponible",
				"service": s.Name,
			})
			c.Abort()
			return
		}

		rawURL := s.nextURL()
		target, err := url.Parse(rawURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "URL de servicio inválida"})
			c.Abort()
			return
		}

		// Cliente HTTP con timeout
		transport := &http.Transport{}
		client := &http.Client{
			Timeout:   s.Timeout,
			Transport: transport,
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.Transport = client.Transport

		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			s.CB.Failure()
			logger.Error(fmt.Sprintf("Error en proxy %s: %v", s.Name, err))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error":"servicio no disponible"}`))
		}

		c.Request.Host = target.Host

		start := time.Now()
		proxy.ServeHTTP(c.Writer, c.Request)
		latency := time.Since(start)

		status := c.Writer.Status()
		if status >= 500 {
			s.CB.Failure()
		} else {
			s.CB.Success()
		}

		logger.Request(c.Request.Method, c.Request.URL.Path, c.ClientIP(), status, latency)
	}
}

func (s *Service) Status() map[string]interface{} {
	return map[string]interface{}{
		"name":    s.Name,
		"urls":    s.URLs,
		"state":   string(s.CB.GetState()),
		"timeout": s.Timeout.String(),
	}
}
