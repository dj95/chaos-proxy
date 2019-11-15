package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthzHandler Indicates if the service is running healthy
func HealthzHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(
			http.StatusOK,
			"OK",
		)
	}
}
