package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const (
	CurrencyCHF string = "CHF"
	CurrencyXES string = "XES"
)

type PriceService interface {
	GetPriceInFor(toCurrency string, fromCurrency string) (float64, error)
}

func NewCryptoComparePriceService(apiKey, baseUrl string) PriceService {
	return &cryptoComparePriceService{
		apiKey:  apiKey,
		baseUrl: baseUrl,
	}
}

// CryptoCompare's HTTP Implementation
type cryptoComparePriceService struct {
	apiKey  string
	baseUrl string
}

func (me *cryptoComparePriceService) GetPriceInFor(to string, from string) (float64, error) {
	var value = 0.0
	resp, err := http.Get(me.baseUrl + "/data/price?fsym=" + from + "&tsyms=" + to + "&api_key=" + me.apiKey)
	if err != nil {
		return value, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return value, errors.New("server returned an unexpected answer: " + resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return value, err
	}
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return value, err
	}
	value, ok := response[to].(float64)
	if ok {
		return value, nil
	}
	return value, errors.New("can't convert response (" + fmt.Sprintf("%v", response) + ") " + err.Error())
}
