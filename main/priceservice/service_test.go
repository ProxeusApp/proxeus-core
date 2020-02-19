package priceservice

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var priceService PriceService

const apiKey = "anyTestApiKey"

const (
	currencyCHF string = "CHF"
	currencyXES string = "XES"
)

func TestMain(m *testing.M) {
	code := m.Run()

	os.Exit(code)
}

func TestGetPriceInForShouldReturnValue(t *testing.T) {
	server := setupResponseFor(`{ "CHF": 0.0025 }`, http.StatusOK)
	defer server.Close()

	priceService = NewCryptoComparePriceService(apiKey, server.URL)

	price, err := priceService.GetPriceInFor(currencyCHF, currencyXES)
	if err != nil {
		t.Error(err)
	}

	if price != 0.0025 {
		t.Error("price is " + fmt.Sprintf("%f", price))
	}
}

func TestGetPriceInForShouldReturnErrorOnStatusCodeOtherThanOk(t *testing.T) {
	server := setupResponseFor(`{ "validJson": true, "GBP": 1.5 }`,
		http.StatusBadRequest)
	defer server.Close()

	priceService = NewCryptoComparePriceService(apiKey, server.URL)

	_, err := priceService.GetPriceInFor(currencyCHF, currencyXES)
	if err == nil {
		t.Error("Should return an error")
	}
}

func setupResponseFor(jsonResponse string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(statusCode)
		rw.Write([]byte(jsonResponse))
	}))
}
