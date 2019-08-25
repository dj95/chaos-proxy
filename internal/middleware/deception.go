package middleware

import (
	"net/http/httputil"
	"net/url"
	_ "time"

	_ "github.com/cevatbarisyilmaz/lossy"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	decProxy "github.com/dj95/deception-proxy/internal/proxy"
)

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
	proxy.Transport = decProxy.NewLossyRoundTrip()

	// return the handler function
	return func(c *gin.Context) {
		// proxy the incoming request
		proxy.ServeHTTP(c.Writer, c.Request)

		// abort further processing
		c.Done()
	}
}
