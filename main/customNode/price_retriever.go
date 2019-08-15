package customNode

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"git.proxeus.com/core/central/sys/workflow"
)

type priceRetrieverNode struct{}

func NewPriceRetriever(n *workflow.Node) (workflow.NodeIF, error) {
	return &priceRetrieverNode{}, nil
}

func (me priceRetrieverNode) Execute(n *workflow.Node, dat interface{}) (proceed bool, err error) {
	log.Println("Retrieving price...")
	j, _ := json.Marshal(dat)
	log.Println("Before: " + string(j))
	cryptoComparePriceService := NewCryptoComparePriceService("121842ca2a76b199422d3351a9958d1fa70a2a6ad478b3681848d64a1792e621", "https://min-api.cryptocompare.com")
	chfXes, err := cryptoComparePriceService.GetPriceInFor("CHF", "XES")
	if err != nil {
		return false, err
	}
	data := dat.(map[string]interface{})
	formDataKey := "input"

	if reflect.ValueOf(data[formDataKey]).IsNil() {
		data[formDataKey] = make(map[string]interface{})
	}
	inputData := data[formDataKey].(map[string]interface{})
	chfXesString := fmt.Sprintf("%f", chfXes)

	inputData["CHFXES"] = chfXesString

	j, _ = json.Marshal(dat)
	log.Println("After: " + string(j))

	return true, nil
}

func (me priceRetrieverNode) Remove(n *workflow.Node) {}
func (me priceRetrieverNode) Close()                  {}
