package router

import (
	"github.com/gin-gonic/gin"

	"github.com/dj95/deception-proxy/internal/middleware"
)

func Setup() *gin.Engine {
	// create a new empty engine
	engine := gin.New()

	// register middlewares
	engine.Use(
		gin.Recovery(),
		middleware.LoggingMiddleware(),
		middleware.Deception(),
	)

	// return the engine
	return engine
}
