package service

import (
	"fmt"
	"github.com/ProxeusApp/proxeus-core/externalnode"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/workflow"
	"log"
	"net/http"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
)

type (
	NodeService interface {
		ProbeExternalNodes()
		List(nodeType string) []*workflow.Node
	}
	defaultNodeService struct {
		workflowService WorkflowService
	}
)

func NewNodeService(workflowService WorkflowService) *defaultNodeService {
	return &defaultNodeService{workflowService: workflowService}
}

//ProbeExternalNodes checks all registered external nodes health endpoint and deletes the ones that are offline
func (me *defaultNodeService) ProbeExternalNodes() {
	for _, node := range me.listExternalNodes() {
		me.probeExternalNode(node)
	}
}

//List returns a list of workflow nodes.
func (me *defaultNodeService) List(nodeType string) []*workflow.Node {
	var nodes []*workflow.Node
	switch nodeType {
	case "mailsender":
		nodes = append(nodes, &workflow.Node{
			ID:     "1234123-1234124",
			Name:   "Mail Sender",
			Detail: "sends an email",
			Type:   "mailsender",
		})
	case "priceretriever":
		nodes = append(nodes, &workflow.Node{
			ID:     "3",
			Name:   "Price retriever",
			Detail: "Retrieves CHF/XES price",
			Type:   "priceretriever",
		})
	case "ibmsender":
		if os.Getenv("FF_IBM_SENDER_ENABLED") == "true" {
			nodes = append(nodes, &workflow.Node{
				ID:     "1234123-1234123",
				Name:   "IBM Sender",
				Detail: "sends all workflow data to an IBM service",
				Type:   "ibmsender",
			})
		}
	case "externalNode":
		for _, node := range me.listExternalNodes() {
			id := uuid.NewV4().String()
			nodes = append(nodes, &workflow.Node{
				ID:     id,
				Name:   node.Name,
				Detail: node.Detail,
				Type:   "externalNode",
			})
		}
	}
	return nodes
}

func (me *defaultNodeService) probeExternalNode(node *externalnode.ExternalNode) {
	log.Printf("[nodeservice] checking external node %s \n", node.Name)
	err := me.healthCheck(node.HealthUrl())
	if err != nil {
		log.Printf("[nodeservice] removing external node err %s \n", err)
		me.deleteExternalNode(new(model.User), node.ID)
	}
}

func (me *defaultNodeService) healthCheck(url string) error {
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

func (me *defaultNodeService) deleteExternalNode(auth model.Auth, id string) error {
	return workflowDB().DeleteExternalNode(auth, id)
}

func (me *defaultNodeService) listExternalNodes() []*externalnode.ExternalNode {
	return workflowDB().ListExternalNodes()
}
