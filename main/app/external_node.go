package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ProxeusApp/proxeus-core/sys/workflow"

	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type externalNode struct {
	ctx  *DocumentFlowInstance
	ctx2 *ExecuteAtOnceContext
	n    *workflow.Node
}

func (ex externalNode) Execute(n *workflow.Node) (proceed bool, err error) {
	client := http.Client{Timeout: 5 * time.Second}
	q, err := ex.system().DB.Workflow.QueryFromInstanceID(ex.auth(), n.ID)
	if err != nil {
		return false, err
	}
	d, err := ex.data()
	if err != nil {
		return false, err
	}
	buf, err := json.Marshal(d)
	if err != nil {
		return false, err
	}
	r, err := client.Post(q.NextUrl(), "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return false, err
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return false, fmt.Errorf("bad node status code %d", r.StatusCode)
	}
	err = ex.putData(r.Body)
	return err == nil, err
}

func (ex externalNode) Remove(n *workflow.Node) {
	client := http.Client{Timeout: 5 * time.Second}
	q, err := ex.system().DB.Workflow.QueryFromInstanceID(ex.auth(), n.ID)
	if err != nil {
		log.Print("remove err ", err.Error())
		return
	}
	client.Post(q.RemoveUrl(), "application/json", bytes.NewBuffer([]byte("{}")))
}

func (ex externalNode) Close() {
	client := http.Client{Timeout: 5 * time.Second}
	q, err := ex.system().DB.Workflow.QueryFromInstanceID(ex.auth(), ex.n.ID)
	if err != nil {
		log.Print("close err ", err.Error())
		return
	}
	client.Post(q.CloseUrl(), "application/json", bytes.NewBuffer([]byte("{}")))
}

func (ex externalNode) system() *sys.System {
	if ex.ctx != nil {
		return ex.ctx.system
	}
	return ex.ctx2.c.System()
}

func (ex externalNode) auth() model.Auth {
	if ex.ctx != nil {
		return ex.ctx.auth
	}
	return ex.ctx2.c.Session(false)
}

func (ex externalNode) data() (dat map[string]interface{}, err error) {
	if ex.ctx != nil {
		return ex.ctx.DataCluster.GetAllData()
	}
	return ex.ctx2.data, nil
}

func (ex externalNode) putData(r io.Reader) error {
	var d map[string]interface{}
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, &d)
	if err != nil {
		return err
	}
	if ex.ctx != nil {
		ex.ctx.statusResult.UserData = d
		return ex.ctx.writeData(ex.n, d)
	}
	ex.ctx2.data = d
	return nil
}

func ProbeExternalNodes(s *sys.System) {
	for _, node := range s.DB.Workflow.ListExternalNodes() {
		log.Printf("[nodeservice] checking external node %s \n", node.Name)
		err := healthCheck(node.HealthUrl())
		if err != nil {
			log.Printf("[nodeservice] removing external node err %s \n", err)
			s.DB.Workflow.DeleteExternalNode(new(model.User), node.ID)
		}
	}
}

func healthCheck(url string) error {
	client := http.Client{Timeout: 5 * time.Second}
	var err error
	var r *http.Response
	for i := 0; i < 3; i++ {
		r, err = client.Get(url)
		if err == nil && r.StatusCode == http.StatusOK {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	var code int
	if r != nil {
		code = r.StatusCode
	}
	return fmt.Errorf("%s [code %d]", err.Error(), code)
}
