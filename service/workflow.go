package service

import (
	extNode "github.com/ProxeusApp/proxeus-core/externalnode"
	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/workflow"
	"log"
)

type (
	WorkflowService interface {
		List(auth model.Auth, contains string, options storage.Options) ([]*model.WorkflowItem, error)
		ListPublished(auth model.Auth, contains string, options storage.Options) ([]*model.WorkflowItem, error)
		ListIds(auth model.Auth, contains string, options storage.Options) ([]string, error)
		GetAndPopulateOwner(auth model.Auth, id string) (*model.WorkflowItem, error)
		Get(auth model.Auth, id string) (*model.WorkflowItem, error)
		Publish(auth model.Auth, wfItem *model.WorkflowItem) map[string]interface{}
		Put(auth model.Auth, wfItem *model.WorkflowItem) error
		Delete(auth model.Auth, id string) error
		InstantiateExternalNode(auth model.Auth, nodeId, nodeName string) (*extNode.ExternalQuery, error)
	}

	DefaultWorkflowService struct {
		userService UserService
		*baseService
	}
)

func NewWorkflowService(system *sys.System, userService UserService) *DefaultWorkflowService {
	return &DefaultWorkflowService{userService: userService, baseService: &baseService{system: system}}
}

func (me *DefaultWorkflowService) List(auth model.Auth, contains string, options storage.Options) ([]*model.WorkflowItem, error) {
	return me.workflowDB().List(auth, contains, options)
}

func (me *DefaultWorkflowService) ListPublished(auth model.Auth, contains string, options storage.Options) ([]*model.WorkflowItem, error) {
	return me.workflowDB().ListPublished(auth, contains, options)
}

func (me *DefaultWorkflowService) Get(auth model.Auth, id string) (*model.WorkflowItem, error) {
	return me.workflowDB().Get(auth, id)
}

func (me *DefaultWorkflowService) GetAndPopulateOwner(auth model.Auth, id string) (*model.WorkflowItem, error) {
	workflow, err := me.workflowDB().Get(auth, id)
	if err != nil {
		return workflow, err
	}
	workflowOwner, err := me.userService.GetById(auth, workflow.Owner)
	if err != nil {
		return workflow, err
	}
	workflow.OwnerEthAddress = workflowOwner.EthereumAddr

	return workflow, err
}

func (me *DefaultWorkflowService) ListIds(auth model.Auth, contains string, options storage.Options) ([]string, error) {
	items, err := me.List(auth, contains, options)
	if err != nil {
		return nil, err
	}

	id := make([]string, len(items))
	if len(items) > 0 {
		for i, item := range items {
			id[i] = item.ID
		}

	}
	return id, err
}

func (me *DefaultWorkflowService) Publish(auth model.Auth, wfItem *model.WorkflowItem) map[string]interface{} {
	errs := map[string]interface{}{}
	collectError := func(err error, node *workflow.Node) {
		errs[node.ID] = struct {
			Error string
			Item  interface{}
		}{Error: err.Error(), Item: node}
	}
	//loop recursively and change permissions on all children
	wfItem.LoopNodes(nil, func(l *workflow.Looper, node *workflow.Node) bool {
		if node.Type == "form" {
			it, er := me.formDB().Get(auth, node.ID)
			if er != nil {
				collectError(er, node)
				return true //continue
			}
			if !it.Published {
				it.Published = true
				er = me.formDB().Put(auth, it)
				if er != nil {
					collectError(er, node)
				}
			}
		} else if node.Type == "template" {
			it, er := me.templateDB().Get(auth, node.ID)
			if er != nil {
				collectError(er, node)
				return true //continue
			}
			if !it.Published {
				it.Published = true
				er = me.templateDB().Put(auth, it)
				if er != nil {
					collectError(er, node)
				}
			}
		} else if node.Type == "workflow" { // deep dive...
			it, er := me.workflowDB().Get(auth, node.ID)
			if er != nil {
				collectError(er, node)
				return true //continue
			}
			if !it.Published {
				it.Published = true
				er = me.workflowDB().Put(auth, it)
				if er != nil {
					collectError(er, node)
				}
			}
			it.LoopNodes(l, nil)
		}
		return true //continue
	})
	return errs
}

func (me *DefaultWorkflowService) Put(auth model.Auth, wfItem *model.WorkflowItem) error {
	return me.workflowDB().Put(auth, wfItem)
}

func (me *DefaultWorkflowService) Delete(auth model.Auth, id string) error {
	return me.workflowDB().Delete(auth, id)
}

func (me *DefaultWorkflowService) InstantiateExternalNode(auth model.Auth, nodeId, nodeName string) (*extNode.ExternalQuery, error) {
	externalQuery, err := me.workflowDB().QueryFromInstanceID(auth, nodeId)
	if err == nil {
		log.Println("UpdateHandler externalNode instance exists, will not create new instance")
		return &externalQuery, nil
	}

	newExternalNode, err := me.workflowDB().NodeByName(auth, nodeName)
	if err != nil {
		log.Println("UpdateHandler instantiateExternalNode NodeByName err: ", err.Error())
		return nil, err
	}

	newExternalNodeInstance := &extNode.ExternalNodeInstance{
		ID:       nodeId,
		NodeName: newExternalNode.Name,
	}

	//PutExternalNodeInstance
	err = me.workflowDB().PutExternalNodeInstance(auth, newExternalNodeInstance)
	if err != nil {
		log.Println("UpdateHandler externalNode PutExternalNodeInstance error: ", err.Error())
		return nil, err
	}

	newExternalQuery := extNode.ExternalQuery{
		ExternalNode:         newExternalNode,
		ExternalNodeInstance: newExternalNodeInstance,
	}

	return &newExternalQuery, nil
}
