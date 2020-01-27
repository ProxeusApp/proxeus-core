package app

import (
	"fmt"
	"log"

	"github.com/ProxeusApp/proxeus-core/sys/workflow"
)

type priceRetrieverNode struct {
	ctx  *DocumentFlowInstance
	ctx2 *ExecuteAtOnceContext
}

func (me priceRetrieverNode) Execute(n *workflow.Node) (proceed bool, err error) {
	log.Println("Retrieving price...")
	cryptoComparePriceService := NewCryptoComparePriceService("API_KEY", "https://min-api.cryptocompare.com")
	chfXes, err := cryptoComparePriceService.GetPriceInFor("CHF", "XES")
	if err != nil {
		return false, err
	}

	chfXesString := fmt.Sprintf("%f", chfXes)
	fmt.Println(chfXesString)
	/**
	IMPORTANT: dat is read only and is provided by GetData you set up the workflow with.
	The workflow engine is not doing anything with your data but resolving the condition.
	Reading and writing it to a database must be handled outside of the workflow.
	For example -> outside = ctx
	*/

	/**
	example of how to read:
	Please note, readData provides simple access to your specific data structure.
	Accessing your complex data structure works just like you would access it in javascript like:
		input.myArray[3].myObjectInsideArray.firstName
	*/
	//reading only input.CHFXES
	//var val interface{}
	//val, err = me.ctx.readData("input.CHFXES")
	//if err != nil {
	//	return false, err
	//}

	/**
	example 1 of how to update:
	Please note, writeField writes it inside -> input:{} as it was used only for the form updates.
	Please change if you need to use another context than input! Input was initially meant for form input only.
	*/
	if me.ctx != nil {
		err = me.ctx.writeField(n, "CHFXES", chfXesString)
		if err != nil {
			return false, err
		}
	} else {
		me.ctx2.data["CHFXES"] = chfXesString
	}

	/**
	example 2 of how to update more generic to pre-fill all form fields that match approximately:
	// 1. CTX_CHFXES with chfxes
	// 2. Consumer_CHFXES with chfxes
	// 3. ProducerCHFXES with chfxes
	// 4. CHFXES_Price with chfxes
	// 5. PriceCHFXES with chfxes
	*/
	//err = me.ctx.updateFormFieldsContaining(n, chfXesString, "CHFXES", "price")
	// This example would match only with 4. and 5. from above.
	//if err != nil {
	//	return false, err
	//}

	return true, nil
}

func (me priceRetrieverNode) Remove(n *workflow.Node) {}
func (me priceRetrieverNode) Close()                  {}
