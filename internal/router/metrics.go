package router

import (
	"github.com/gin-gonic/gin"
	gometrics "github.com/rcrowley/go-metrics"

	"github.com/dj95/chaos-proxy/pkg/metrics"
)

// MetricsHandler Handle the metrics route for prometheus metrics.
func MetricsHandler() gin.HandlerFunc {
	return gin.WrapH(metrics.HTTPHandler(
		map[string]interface{}{
			"http.request.total": gometrics.NewCounter(),
		},
		"dj95",
		"chaos-proxy",
	))
}
