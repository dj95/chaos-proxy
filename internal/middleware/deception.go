package middleware

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/dj95/deception-proxy/internal/roundtrip"
)

// Deception Proxy requests with different latency, packet loss rate and bandwidth
func Deception() gin.HandlerFunc {
	// read requested values from the config
	target := viper.GetString("conn.target")

	// parse the target url
	targetURL, err := url.Parse(target)

	// error handling
	if err != nil {
		log.Fatalf("Wrong configured target url: %s", err.Error())
	}

	// create a new reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// set custom round tripper with lossy
	proxy.Transport = roundtrip.NewLossy()

	// return the handler function
	return func(c *gin.Context) {
		// do not proxy the healthz route
		if c.Request.RequestURI == "/healthz" {
			c.Next()

			return
		}

		// proxy the incoming request
		proxy.ServeHTTP(c.Writer, c.Request)

		// abort further processing
		c.Done()
	}
}
