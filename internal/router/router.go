package router

import (
	"github.com/gin-gonic/gin"

	"github.com/dj95/deception-proxy/internal/middleware"
)

// Setup Initialize a new gin engine with routes and middlewares
func Setup() (*gin.Engine, error) {
	// create a new empty engine
	engine := gin.New()

	// register middlewares
	engine.Use(
		gin.Recovery(),
		middleware.Logging(),
	)

	// register routes
	engine.GET("/healthz", HealthzHandler)

	// define an api group for target based routes
	targetGroup := engine.Group("/api/{target}")
	{
		// TODO: get config
		// TODO: set different single parameters
		targetGroup.GET("/config")
	}

	// return the engine
	return engine, nil
}
