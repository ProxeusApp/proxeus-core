package customNode

import (
	"os"

	uuid "github.com/satori/go.uuid"

	"github.com/ProxeusApp/proxeus-core/main/www"

	"github.com/ProxeusApp/proxeus-core/sys/workflow"
)

// Returns a list of all custom nodes
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
