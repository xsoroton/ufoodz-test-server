package main

import (
	"net/http"
	"strconv"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/ufoodz-test-server/configs"
	"github.com/ufoodz-test-server/middleware"
	"github.com/ufoodz-test-server/mocks"
)

const (
	prodEnv            = "production"
	serviceDescription = "Private Rental Company B"
)

var config = configs.EnvConfig{}

func init() {
	logrus.Infof("Starting %s µService", serviceDescription)
	err := env.Parse(&config)
	if err != nil {
		logrus.Panic("Fail to read env variables", err)
	}
	logLevel, logPlain := configs.ParseInitVars(config)
	logrus.SetLevel(logrus.Level(logLevel))
	if !logPlain {
		// Log to STDOUT as json
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}

func main() {
	// Start Gin
	err := getMainEngine().Run(":" + strconv.Itoa(config.CompanyBPort))
	if err != nil {
		logrus.Error("Fail to start Service: ", err)
	}
}

func getMainEngine() *gin.Engine {
	// Set ReleaseMode for Gin
	if config.Environment == prodEnv {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.Default()

	// Log Gin with Logrus
	engine.Use(middleware.GinLogrus())

	// Simple auth middleware example, expect token in header
	engine.Use(middleware.Auth())

	engine.GET("/data", Data)
	engine.GET("/healthz", Health)

	return engine
}

// Health ...
func Health(c *gin.Context) {
	c.String(http.StatusOK, http.StatusText(http.StatusOK))
}

// Data ...
func Data(c *gin.Context) {
	quotes := mocks.GenerateFakeQuotes(10, "Rent Company B")
	// Keep Full Logs
	go func() {
		logrus.WithFields(
			logrus.Fields{
				"requestHeaders": c.Request.Header,
				"quote":          quotes,
			},
		).Info("Data ", serviceDescription)
	}()

	c.JSON(http.StatusOK, quotes)
}
