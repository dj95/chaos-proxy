package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthzHandler Indicates if the service is running healthy
func HealthzHandler(c *gin.Context) {
	c.String(
		http.StatusOK,
		"OK",
	)
}
