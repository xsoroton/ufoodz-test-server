package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ufoodz-test-server/configs"
	"github.com/ufoodz-test-server/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/ufoodz-test-server/mocks"
)

func init() {
	gin.SetMode(gin.TestMode)
	logrus.SetLevel(logrus.PanicLevel)
}

// TestDataRoute ...
func TestDataRoute(t *testing.T) {

	endPoint := "/data"
	router := gin.New()
	router.GET(endPoint, Data)

	Convey("Test Data endpoint get data from Company A and B", t, func() {
		// Company A
		countA := 5
		mockDataA := mocks.GenerateFakeQuotes(countA, "Company A")
		jsonDataA, err := json.Marshal(mockDataA)
		So(err, ShouldBeNil)
		a := mocks.MockRemoteServer(string(jsonDataA), http.StatusOK)
		defer a.Close()
		configs.CompanyAHost = a.URL

		// Company B
		countB := 3
		mockDataB := mocks.GenerateFakeQuotes(countB, "Company B")
		jsonDataB, err := json.Marshal(mockDataB)
		So(err, ShouldBeNil)
		b := mocks.MockRemoteServer(string(jsonDataB), http.StatusOK)
		defer b.Close()
		configs.CompanyBHost = b.URL

		// Test API
		req, _ := http.NewRequest(http.MethodGet, endPoint, nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		body, readErr := ioutil.ReadAll(resp.Body)
		So(readErr, ShouldBeNil)
		var respBody models.Quotes
		jsonError := json.Unmarshal(body, &respBody)
		So(jsonError, ShouldBeNil)

		So(resp.Code, ShouldEqual, http.StatusOK)
		So(len(respBody.Quotes), ShouldEqual, countA+countB)
	})

	Convey("Test Data endpoint get data from Company A and B fail", t, func() {
		// Company A
		countA := 5
		mockDataA := mocks.GenerateFakeQuotes(countA, "Company A")
		jsonDataA, err := json.Marshal(mockDataA)
		So(err, ShouldBeNil)
		a := mocks.MockRemoteServer(string(jsonDataA), http.StatusOK)
		defer a.Close()
		configs.CompanyAHost = a.URL

		// Company B
		b := mocks.MockRemoteServer(``, http.StatusUnauthorized)
		defer b.Close()
		configs.CompanyBHost = b.URL

		// Test API
		req, _ := http.NewRequest(http.MethodGet, endPoint, nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		body, readErr := ioutil.ReadAll(resp.Body)
		So(readErr, ShouldBeNil)
		var respBody models.Quotes
		jsonError := json.Unmarshal(body, &respBody)
		So(jsonError, ShouldBeNil)

		So(resp.Code, ShouldEqual, http.StatusOK)
		So(len(respBody.Quotes), ShouldEqual, countA)
	})
}
