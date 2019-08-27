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

	"github.com/dj95/deception-proxy/internal/router"
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

	// TODO: check configured values

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
		"deception-proxy.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)

	// create the log output
	logOutput := io.MultiWriter(os.Stdout, file)

	// if no error occured...
	if err != nil {
		logOutput = os.Stdout
		log.Info("Failed to log to file, using default stderr")
	}

	// set the stdout + file logger
	log.SetOutput(logOutput)
}

func main() {
	// create a new router
	router, err := router.Setup()

	// error handling
	if err != nil {
		log.Fatalf("Cannot initialize router: %s", err.Error())
	}

	// run the router
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
