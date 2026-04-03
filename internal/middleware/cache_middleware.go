package middleware

import (
	"bytes"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type cacheEntry struct {
	body       []byte
	statusCode int
	expiresAt  time.Time
}

type responseCapture struct {
	gin.ResponseWriter
	body       bytes.Buffer
	statusCode int
}

func (r *responseCapture) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *responseCapture) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

var (
	cacheMu sync.RWMutex
	cache   = map[string]*cacheEntry{}
)

func CacheResponse(ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Solo cachear GETs
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		key := c.Request.URL.String()

		// Verificar si hay cache válido
		cacheMu.RLock()
		entry, found := cache[key]
		cacheMu.RUnlock()

		if found && time.Now().Before(entry.expiresAt) {
			c.Data(entry.statusCode, "application/json", entry.body)
			c.Header("X-Cache", "HIT")
			c.Abort()
			return
		}

		// Capturar la respuesta
		capture := &responseCapture{
			ResponseWriter: c.Writer,
			statusCode:     200,
		}
		c.Writer = capture
		c.Header("X-Cache", "MISS")

		c.Next()

		// Guardar en cache solo respuestas exitosas
		if capture.statusCode < 400 {
			cacheMu.Lock()
			cache[key] = &cacheEntry{
				body:       capture.body.Bytes(),
				statusCode: capture.statusCode,
				expiresAt:  time.Now().Add(ttl),
			}
			cacheMu.Unlock()
		}
	}
}
