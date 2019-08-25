package router

import (
	"github.com/gin-gonic/gin"

	"github.com/dj95/deception-proxy/internal/middleware"
)

// Setup Initialize a new gin engine with routes and middlewares
func Setup() *gin.Engine {
	// create a new empty engine
	engine := gin.New()

	// register middlewares
	engine.Use(
		gin.Recovery(),
		middleware.Logging(),
		middleware.Deception(),
	)

	// return the engine
	return engine
}
