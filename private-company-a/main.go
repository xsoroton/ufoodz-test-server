package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/ufoodz-test-server/configs"
	"github.com/ufoodz-test-server/middleware"
	"github.com/ufoodz-test-server/models"
)

const (
	prodEnv            = "production"
	serviceDescription = "Private Rental Company A"
)

var myJSONComEndpoint = "https://api.myjson.com/bins/7c0qw"
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
	err := getMainEngine().Run(":" + strconv.Itoa(config.CompanyAPort))
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
	quotes, err := getQuoteFromMyJSON()
	if err != nil {
		logrus.Error("Fail to get data from ", myJSONComEndpoint)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// Keep Full Logs
	go func() {
		logrus.WithFields(
			logrus.Fields{
				"requestHeaders": c.Request.Header,
				"quotes":         quotes,
			},
		).Info("Data ", serviceDescription)
	}()

	c.JSON(http.StatusOK, quotes)
}

// getQuoteFromMyJSON get quote from https://api.myjson.com/bins/7c0qw as in example
func getQuoteFromMyJSON() (quotes []models.Quote, err error) {
	client := http.Client{Timeout: time.Second * 30}
	request, err := http.NewRequest(http.MethodGet, myJSONComEndpoint, nil)
	if err != nil {
		return quotes, err
	}
	response, err := client.Do(request)
	if err != nil {
		return quotes, err
	}
	defer response.Body.Close()

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return quotes, readErr
	}
	if response.StatusCode != http.StatusOK {
		err = errors.New("invalid HTTP response code")
		return
	}
	if len(body) == 0 {
		err = errors.New("empty response from " + myJSONComEndpoint)
		return
	}
	var quote models.Quote
	err = json.Unmarshal(body, &quote)
	if err != nil {
		return
	}
	// Because https://api.myjson.com/bins/7c0qw example have return only single quote object :(
	quotes = []models.Quote{quote}
	return
}
