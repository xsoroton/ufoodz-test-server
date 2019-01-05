package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	logrus.SetLevel(logrus.Level(logrus.ErrorLevel))
	gin.SetMode(gin.TestMode)
}

func TestGinLogrusMiddleware(t *testing.T) {
	Convey("Test Gin Logrus Middleware", t, func() {
		r := gin.New()
		r.Use(GinLogrus())
		r.POST("/healthz")
		req, err := http.NewRequest(http.MethodPost, "/healthz", bytes.NewBuffer([]byte(`OK`)))
		So(err, ShouldBeNil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)
	})
}
