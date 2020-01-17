package database

import (
	"testing"

	"github.com/ProxeusApp/proxeus-core/sys/workflow"

	"github.com/ProxeusApp/proxeus-core/storage"

	"github.com/ProxeusApp/proxeus-core/sys/model"

	. "github.com/onsi/gomega"
)

func TestWorkflow(t *testing.T) {
	RegisterTestingT(t)
	wo := testDBSet.Workflow

	options := storage.IndexOptions(0)

	item := &model.WorkflowItem{
		ID:   "1",
		Name: "some name",
		Data: &workflow.Workflow{},
	}
	item2 := &model.WorkflowItem{
		Name: "some name 2",
	}

	// add
	Expect(wo.Put(dummy, item)).To(Succeed())
	Expect(wo.Put(dummy, item2)).To(Succeed())

	// fetch item2.ID
	gotItems, _ := wo.List(dummy, item2.Name, options)
	item2 = gotItems[0]

	// modify
	item.Name = "some name 1"
	Expect(wo.Put(dummy, item)).To(Succeed())
	item2.Published = true
	item2.Owner = dummy.UserID()
	Expect(wo.Put(dummy, item2)).To(Succeed())

	// get
	Expect(wo.Get(dummy, item.ID)).To(equalJSON(item))
	Expect(wo.GetPublished(dummy, item2.ID)).To(equalJSON(item2))

	// list
	Expect(wo.List(dummy, "some", options)).
		To(equalJSON([]*model.WorkflowItem{item2, item}))
	Expect(wo.GetList(dummy, []string{item.ID, item2.ID})).
		To(equalJSON([]*model.WorkflowItem{item, item2}))
	Expect(wo.ListPublished(dummy, item2.Name, options)).
		To(equalJSON([]*model.WorkflowItem{item2}))

	// delete
	Expect(wo.Delete(dummy, item.ID)).To(Succeed())
	Expect(wo.List(dummy, "", options)).To(equalJSON([]*model.WorkflowItem{item2}))
}
