// Package router Implement a small webserver offering an api
// to control and setup different proxy connections.
package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dj95/chaos-proxy/internal/middleware"
	"github.com/dj95/chaos-proxy/pkg/proxy"
)

// Setup Initialize a new gin engine with routes and middlewares
func Setup(proxies []proxy.Proxy) (*gin.Engine, error) {
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
	apiGroup := engine.Group("/api")
	{
		// get all configured proxy connections
		apiGroup.GET("/conn", ConnHandler(proxies))

		// manage single proxy connections
		apiGroup.GET("/conn/:id", ConnIDHandler(proxies))
		apiGroup.POST("/conn/:id", ConnUpdateHandler(proxies))
	}

	// serve some swagger documentation
	if gin.Mode() == gin.DebugMode {
		engine.StaticFS("/doc", http.Dir("./web/swaggerui"))
		engine.StaticFile("/swagger.yaml", "./api/swagger.yaml")
	}

	// return the engine
	return engine, nil
}
