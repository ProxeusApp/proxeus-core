package externalnode

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/model"

	"github.com/ProxeusApp/proxeus-core/main/www"

	"github.com/ProxeusApp/proxeus-core/sys/workflow"
)

func List(c *www.Context, nodeType string) []*workflow.Node {
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
		for _, node := range c.System().DB.Workflow.ListExternalNodes() {
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

func ProbeExternalNodes(s *sys.System) {
	for _, node := range s.DB.Workflow.ListExternalNodes() {
		log.Printf("[nodeservice] checking external node %s \n", node.Name)
		err := healthCheck(node.HealthUrl())
		if err != nil {
			log.Printf("[nodeservice] removing external node err %s \n", err)
			err = s.DB.Workflow.DeleteExternalNode(new(model.User), node.ID)
			if err != nil {
				log.Printf("[nodeservice] unable to remove external node err %s", err)
			}
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
