package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ufoodz-test-server/configs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GinLogrus map Gin data to Logrus logs
func GinLogrus() gin.HandlerFunc {
	logger := logrus.StandardLogger()
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		// log latency in milliseconds
		latencyMilliseconds := time.Since(start) / 1000000

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.String()
		userAgent := c.Request.UserAgent()

		msg := fmt.Sprintf(
			"%s \"%s %s\" %d %d %s",
			clientIP,
			method,
			path,
			statusCode,
			latencyMilliseconds,
			userAgent,
		)

		logWithFields := logger.WithFields(logrus.Fields{
			"method":  method,
			"path":    path,
			"latency": latencyMilliseconds,
			"ip":      clientIP,
			"comment": comment,
			"status":  statusCode,
			"header":  c.Request.Header,
		})

		if statusCode > 499 {
			logWithFields.Warn(msg)
		} else {
			logWithFields.Info(msg)
		}
	}
}

// Auth ...
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token != configs.AuthToken {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}
