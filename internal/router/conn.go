package router

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/dj95/deception-proxy/pkg/config"
	"github.com/dj95/deception-proxy/pkg/proxy"
)

// ConnHandler Return all configured proxy connections.
func ConnHandler(proxies []proxy.Proxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		// initialize the config list
		var configList []*config.Target

		// iterate through all configured proxies
		for _, proxy := range proxies {
			// get and save the config
			configList = append(configList, proxy.Config())
		}

		// return the config list as json response
		c.JSON(
			http.StatusOK,
			configList,
		)
	}
}

// ConnIDHandler Return the configuration for one proxy connection.
func ConnIDHandler(proxies []proxy.Proxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get the requested id from the path
		id := c.Param("id")

		// iterate through the given proxy configs
		for _, proxy := range proxies {
			// if the id of the config is not the id from the path...
			if proxy.Config().ID != id {
				// ...skip this proxy
				continue
			}

			// return the correct config as json response
			c.JSON(
				http.StatusOK,
				proxy.Config(),
			)

			// early return as the config was already sent
			// back to the user
			return
		}

		// if no config was found, respond with a status code 404
		// and an empty json object
		c.JSON(
			http.StatusNotFound,
			gin.H{},
		)
	}
}

// ConnUpdateHandler Update a single proxy configuration.
func ConnUpdateHandler(proxies []proxy.Proxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		// initialize the new config
		var config config.Target

		// try to bind the request body to the config
		if err := c.ShouldBindJSON(&config); err != nil {
			// if an error occured, a status json with status code
			// 404 (Bad request) is returned
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"status": "error",
					"msg":    "invalid request body",
				},
			)

			// skip further processing of this request
			return
		}

		// get the id from the request path
		id := c.Param("id")

		// declare a variable for saving the proxy that is updated
		var savedProxy proxy.Proxy

		// iterate through all given proxies
		for i, proxy := range proxies {
			// if the id does not exist...
			if proxy.Config().ID != id {
				// ...check the next proxy
				continue
			}

			// save the proxy
			savedProxy = proxy

			// shutdown an existing proxy with the id
			err := proxy.Shutdown()

			// return a status code for success
			if err != nil {
				c.JSON(
					http.StatusInternalServerError,
					gin.H{
						"status": "error",
						"msg":    "cannot shutdown proxy",
					},
				)

				// skip further processing
				return
			}

			// remove the stopped proxy
			proxies = append(proxies[:i], proxies[i+1:]...)
		}

		// create a new proxy
		updatedProxy, err := proxy.New(&config)

		// error handling
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"status": "error",
					"msg":    "cannot create updated proxy: " + err.Error(),
				},
			)

			// skip further processing
			return
		}

		// start the listener
		err = updatedProxy.StartListener()

		// error handling
		if err != nil {

			// if no proxy is saved...
			if savedProxy == nil {
				c.JSON(
					http.StatusInternalServerError,
					gin.H{
						"status": "error",
						"msg":    "cannot start updated proxy",
					},
				)
				// ...just return
				return
			}

			// start the listener again
			err = savedProxy.StartListener()

			// error handling
			if err != nil {
				c.JSON(
					http.StatusInternalServerError,
					gin.H{
						"status": "error",
						"msg":    "udated proxy and restored saved proxy cannot start",
					},
				)
			}

			// skip further processing
			return
		}

		// save the updated proxy
		proxies = append(proxies, updatedProxy)

		// return a status code for success
		c.JSON(
			http.StatusOK,
			gin.H{
				"status": "success",
				"msg":    "",
			},
		)
	}
}
