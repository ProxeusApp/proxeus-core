package app

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"log"

	"git.proxeus.com/core/central/sys/workflow"
)

const (
	forwardURL = "<forwardurl>"
)

type IBMSenderNodeImpl struct {
	ctx *DocumentFlowInstance
}

func newIBMSenderNodeImpl(n *workflow.Node) (workflow.NodeIF, error) {
	return &IBMSenderNodeImpl{}, nil
}

func addConfigHeaders(req *http.Request) {
	req.Header.Set("clientid", "<clientid>")
	req.Header.Set("tenantid", "<tenantid>")
	req.Header.Set("secret", "<secret>")
	req.Header.Set("oauthserverurl", "<oauthserverurl>")
}

//data changes requested by customer
func changeDataBeforeSend(dat interface{}) interface{} {
	if m, ok := dat.(map[string]interface{}); ok {
		if d, ok := m["input"]; ok {
			bts, _ := json.Marshal(d)
			var dataCopy map[string]interface{}
			json.Unmarshal(bts, &dataCopy)
			if cs, ok := dataCopy["CapitalSource"]; ok {
				bts, _ := json.Marshal(cs)
				dataCopy["CapitalSource"] = string(bts)
			}
			return dataCopy
		}
	}
	return dat
}

func (me *IBMSenderNodeImpl) Execute(n *workflow.Node) (proceed bool, err error) {
	b, err := json.Marshal(changeDataBeforeSend(me.ctx.getData()))
	if err != nil {
		return false, err
	}
	req, err := http.NewRequest("POST", forwardURL, bytes.NewReader(b))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	addConfigHeaders(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		b2, _ := ioutil.ReadAll(io.LimitReader(resp.Body, 100*1024))
		log.Printf("SERVER NOT ACCEPTED '%s', RESPONSE '%s'\n", b, b2)
		return false, err
	}
	return true, nil
}

func (me *IBMSenderNodeImpl) Remove(n *workflow.Node) {}
func (me *IBMSenderNodeImpl) Close()                  {}
