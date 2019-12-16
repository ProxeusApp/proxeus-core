package app

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"log"

	"github.com/ProxeusApp/proxeus-core/sys/workflow"
	"github.com/pkg/errors"
)

var forwardURL = os.Getenv("FF_IBM_SENDER_URL")

type IBMSenderNodeImpl struct {
	ctx *DocumentFlowInstance
}

func newIBMSenderNodeImpl(n *workflow.Node) (workflow.NodeIF, error) {
	return &IBMSenderNodeImpl{}, nil
}

func addConfigHeaders(req *http.Request) {
	req.Header.Set("clientid", os.Getenv("FF_IBM_SENDER_CLIENT_ID"))
	req.Header.Set("tenantid", os.Getenv("FF_IBM_SENDER_TENANT_ID"))
	req.Header.Set("secret", os.Getenv("FF_IBM_SENDER_SECRET"))
	req.Header.Set("oauthserverurl", os.Getenv("FF_IBM_SENDER_OAUTH_URL"))
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

	if os.Getenv("FF_IBM_SENDER_ENABLED") == "" {
		return false, errors.New("IBM Sender not enabled")
	}

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
