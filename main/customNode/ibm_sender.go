package customNode

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
	forwardURL = "https://48h-gruendung-dev.eu-de.mybluemix.net/proxeusAPI"
)

type IBMSenderNodeImpl struct{}

func NewIBMSenderNodeImpl(n *workflow.Node) (workflow.NodeIF, error) {
	return &IBMSenderNodeImpl{}, nil
}

func addConfigHeaders(req *http.Request) {
	req.Header.Set("clientid", "b26447db-68bc-4ff5-add4-4bdd2a901314")
	req.Header.Set("tenantid", "a6cc00be-a1b1-46ad-860a-17be60852af2")
	req.Header.Set("secret", "ZTk0OGE1MTEtZjEwZS00YWUyLTk3ZDQtYTY0NTczNDgwNGQx")
	req.Header.Set("oauthserverurl", "https://appid-oauth.eu-de.bluemix.net/oauth/v3/a6cc00be-a1b1-46ad-860a-17be60852af2")
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

func (me *IBMSenderNodeImpl) Execute(n *workflow.Node, dat interface{}) (proceed bool, err error) {
	b, err := json.Marshal(changeDataBeforeSend(dat))
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
