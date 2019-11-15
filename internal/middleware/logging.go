package middleware

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/dj95/chaos-proxy/pkg/metrics"
)

// Logging Logs every request
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		log.WithFields(log.Fields{
			"remote_addr": c.Request.RemoteAddr,
			"protocol":    c.Request.Proto,
			"method":      c.Request.Method,
			"user_agent":  c.Request.UserAgent(),
			"host":        c.Request.Host,
			"status":      c.Writer.Status(),
		}).Info("Received request")

		// skip further processing for the healthz route
		if c.Request.RequestURI == "/healthz" {
			return
		}

		// skip further processing for the metrics route
		if c.Request.RequestURI == "/metrics" {
			return
		}

		// increse the http requests total metric
		metrics.Inc("http.request.total")
	}
}
