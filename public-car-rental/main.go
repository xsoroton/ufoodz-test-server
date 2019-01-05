package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
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
	serviceDescription = "OneDayCarRental.biz Company"
)

// TODO: set it from config
var (
	profitMargin = 15
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
	// Start Gin.
	err := getMainEngine().Run(":" + strconv.Itoa(config.Port))
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
	var quotes models.Quotes
	var quoteA []models.Quote
	var quoteB []models.Quote
	var errA error
	var errB error

	// PARALLEL CALL
	var wg sync.WaitGroup
	// number of parallels
	wg.Add(2)

	go func() {
		quoteA, errA = getDataFromRemoteHost(buildURL(configs.CompanyAHost, configs.CompanyAPort) + "/data")
		if errA != nil {
			logrus.Error("Fail to get data from getDataCompanyA")
		}
		wg.Done()
	}()

	go func() {
		quoteB, errB = getDataFromRemoteHost(buildURL(configs.CompanyBHost, configs.CompanyBPort) + "/data")
		if errB != nil {
			logrus.Error("Fail to get data from getDataCompanyB")
		}
		wg.Done()
	}()
	// wait for APIs responses
	wg.Wait()

	if errA != nil && errB != nil {
		logrus.Error("Could not get any data from remote services")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if errA == nil && errB == nil {
		var sum []models.Quote
		sum = append(sum, quoteA...)
		sum = append(sum, quoteB...)
		quotes = models.Quotes{Quotes: sum}
	} else {
		if errA == nil {
			quotes = models.Quotes{Quotes: quoteA}
		}
		if errB == nil {
			quotes = models.Quotes{Quotes: quoteB}
		}
	}

	quotes = addProfitMarign(quotes)

	// Keep Full Logs
	go func() {
		logrus.WithFields(
			logrus.Fields{
				"requestHeaders": c.Request.Header,
				"quoteA":         quoteA,
				"quoteB":         quoteB,
				"errA":           errA,
				"errB":           errB,
				"quotes":         quotes,
			},
		).Info("Data ", serviceDescription)
	}()

	c.JSON(http.StatusOK, quotes)
}

func addProfitMarign(quotes models.Quotes) models.Quotes {
	for k, quote := range quotes.Quotes {
		price, err := strconv.ParseFloat(quote.Quote.Amount, 64)
		if err != nil {
			logrus.Error("Parse string to float error ", err)
			continue
		}
		priceStr := strconv.FormatFloat((price*float64(profitMargin+100))/float64(100), 'f', 2, 64)
		quotes.Quotes[k].Quote.Amount = priceStr
		//logrus.Debug(quote.Quote.Amount + " converted to " + priceStr)
	}

	return quotes
}

func getDataFromRemoteHost(host string) (quote []models.Quote, err error) {
	logrus.Debug("http request to ", host)
	client := http.Client{Timeout: time.Second * 30}
	request, err := http.NewRequest(http.MethodGet, host, nil)
	if err != nil {
		return quote, err
	}
	// External services expected token in header
	request.Header.Set("token", configs.AuthToken)
	response, err := client.Do(request)
	if err != nil {
		return quote, err
	}
	defer response.Body.Close()

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		return quote, readErr
	}
	if response.StatusCode != http.StatusOK {
		err = errors.New("invalid HTTP response code")
		return
	}
	if len(body) == 0 {
		err = errors.New("empty response from " + host)
		return
	}
	err = json.Unmarshal(body, &quote)
	if err != nil {
		err = errors.New("fail Unmarshal")
		return
	}
	return
}

func buildURL(host string, port int) string {
	// if port number set and not standard (80|443), append port number in to to Request URL
	if port != 0 && port != 80 && port != 443 {
		host += ":" + strconv.Itoa(port)
	}
	return host
}
