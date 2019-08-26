package router

import (
	"github.com/gin-gonic/gin"

	"github.com/dj95/deception-proxy/internal/middleware"
)

// Setup Initialize a new gin engine with routes and middlewares
func Setup() (*gin.Engine, error) {
	// create a new empty engine
	engine := gin.New()

	// initialize the deception middleware
	deceptionMiddleware, err := middleware.Deception()

	// error checking
	if err != nil {
		return nil, err
	}

	// register middlewares
	engine.Use(
		gin.Recovery(),
		middleware.Logging(),
		deceptionMiddleware,
	)

	// register routes
	engine.GET("/healthz", HealthzHandler)

	// return the engine
	return engine, nil
}
