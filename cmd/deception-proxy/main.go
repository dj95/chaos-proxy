package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/dj95/chaos-proxy/internal/router"
	"github.com/dj95/chaos-proxy/pkg/config"
	"github.com/dj95/chaos-proxy/pkg/proxy"
)

func init() {
	// set the config name
	viper.SetConfigName("config")

	// add config paths
	viper.AddConfigPath(".")
	viper.AddConfigPath("/")

	// add command line flags
	initializeCommandFlags()

	if viper.IsSet("config") {
		viper.SetConfigFile(viper.GetString("config"))
	}

	// read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Warnf("Cannot read config file: %s", err.Error())
	}

	// read overrides from the environment
	viper.SetEnvPrefix("DP")
	viper.AutomaticEnv()

	// create an env replacer in order to reach core variables
	replacer := strings.NewReplacer(
		"CORE.", "CORE_",
		"CONN.", "CONN_",
		"LATENCY.", "LATENCY_",
	)
	viper.SetEnvKeyReplacer(replacer)

	// set the default log level and mode
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{})
	gin.SetMode(gin.ReleaseMode)

	// activate the debug mode
	if viper.GetString("core.log_level") == "debug" {
		log.SetLevel(log.DebugLevel)
		gin.SetMode(gin.DebugMode)
	}

	// set the json formatter if configured
	if viper.GetString("core.log_format") == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	// open the io writer for the log file
	file, err := os.OpenFile(
		"chaos-proxy.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)

	// create the log output
	logOutput := io.MultiWriter(os.Stdout, file)

	// if no error occurred...
	if err != nil {
		logOutput = os.Stdout
		log.Info("failed to log to file, using default stderr")
	}

	// set the stdout + file logger
	log.SetOutput(logOutput)
}

func main() {
	// get the targets from the configuration
	targets := config.Parse(
		viper.GetStringMap("conn"),
	)

	var runningProxies []proxy.Proxy

	// start the proxy servers for every target
	for name, target := range targets {
		log.Infof("starting target: %s", name)

		// create the proxy
		proxyConn, err := proxy.New(target)

		// error handling
		if err != nil {
			log.Errorf("%s", err.Error())

			continue
		}

		// start the listener
		err = proxyConn.StartListener()

		// error handling
		if err != nil {
			log.WithFields(log.Fields{
				"target": target,
			}).Errorf("cannot start proxy: %s", err.Error())

			continue
		}

		// save the proxy for later usage
		runningProxies = append(runningProxies, proxyConn)
	}

	// initialize a new router for the api
	router, _ := router.Setup(
		runningProxies,
	)

	// start the router
	router.Run(fmt.Sprintf(
		"%s:%d",
		viper.GetString("core.address"),
		viper.GetInt("core.port"),
	))
}

func initializeCommandFlags() {
	// create a new flag for docker health checks
	pflag.String("config", "", "choose the config file")
	pflag.Bool("healthcheck", false, "run a healthcheck")

	// parse the pflags
	pflag.Parse()

	// bind the pflags
	viper.BindPFlags(pflag.CommandLine)
}
