package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var server *httptest.Server
var priceService PriceService

const apiKey = "anyTestApiKey"

func TestMain(m *testing.M) {
	code := m.Run()

	server.Close()
	os.Exit(code)
}

func TestGetPriceInForShouldReturnValue(t *testing.T) {
	setupResponseFor(`{ "CHF": 0.0025 }`, http.StatusOK)
	priceService = NewCryptoComparePriceService(apiKey, server.URL)

	price, err := priceService.GetPriceInFor(CurrencyCHF, CurrencyXES)

	if err != nil {
		t.Error(err)
	}
	if price != 0.0025 {
		t.Error("price is " + fmt.Sprintf("%f", price))
	}
}

func TestGetPriceInForShouldReturnErrorOnStatusCodeOtherThanOk(t *testing.T) {
	setupResponseFor(`{ "validJson": true, "GBP": 1.5 }`,
		http.StatusBadRequest)
	priceService = NewCryptoComparePriceService(apiKey, server.URL)

	_, err := priceService.GetPriceInFor(CurrencyCHF, CurrencyXES)

	if err == nil {
		t.Error("Should return an error")
	}
}

func setupResponseFor(jsonResponse string, statusCode int) {
	server = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(statusCode)
		rw.Write([]byte(jsonResponse))
	}))
}
