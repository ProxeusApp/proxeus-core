package sys

import (
	"encoding/json"
	"log"
	"testing"

	"git.proxeus.com/core/central/sys/model"
	"git.proxeus.com/core/central/sys/workflow"
)

func TestNew(t *testing.T) {
	m := &model.WorkflowItem{}
	m.Owner = "lala"
	m.Data = &workflow.Workflow{Flow: &workflow.Flow{Start: &workflow.Start{NodeID: "1"}, Nodes: map[string]*workflow.Node{"1": {ID: "1"}}}}
	wfItem(m)
}

func wfItem(a model.PermissionItem) {
	bts, err := json.Marshal(a)
	log.Println(err, string(bts))
}
