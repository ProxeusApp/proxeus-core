package customNode

import (
	"git.proxeus.com/core/central/sys/workflow"
)

func List(nodeType string) *workflow.Node {
	var repositories = make(map[string]*workflow.Node)
	repositories["ibmsender"] = &workflow.Node{
		ID:     "1234123-1234123",
		Name:   "IBM Sender",
		Detail: "sends all workflow data to an IBM service",
		Type:   "ibmsender",
	}

	repositories["mailsender"] = &workflow.Node{
		ID:     "1234123-1234124",
		Name:   "Mail Sender",
		Detail: "sends an email",
		Type:   "mailsender",
	}

	repositories["priceretriever"] = &workflow.Node{
		ID:     "3",
		Name:   "Price retriever",
		Detail: "Retrieves CHF/XES price",
		Type:   "priceretriever",
	}
	return repositories[nodeType]
}
