package mocks

import (
	"bytes"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/icrowley/fake"
	"github.com/ufoodz-test-server/models"
)

// GenerateFakeQuotes ...
func GenerateFakeQuotes(n int, company string) (quotes []models.Quote) {
	for i := 1; i <= n; i++ {
		quotes = append(quotes, models.Quote{
			Quote: models.QuoteDetails{
				Name:     company + " " + fake.Word(),
				Amount:   strconv.FormatFloat(20+rand.Float64()*(300-25), 'f', 2, 64),
				Currency: fake.CurrencyCode(),
			},
		})
	}
	return
}

// MockRemoteServer ...
func MockRemoteServer(response string, responseCode int) *httptest.Server {
	handler := handler(response, responseCode)
	ts := httptest.NewServer(http.HandlerFunc(handler))
	return ts
}

func handler(response string, responseCode int) (handler func(w http.ResponseWriter, req *http.Request)) {
	handler = func(w http.ResponseWriter, req *http.Request) {
		// mock response code
		w.WriteHeader(responseCode)
		// mock response body
		w.Write([]byte(response))
	}
	return
}

// BuildBody get Request body from string
func BuildBody(st string) *bytes.Buffer {
	return bytes.NewBuffer([]byte(st))
}
