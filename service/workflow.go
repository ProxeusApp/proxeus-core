package service

import (
	"log"
	"strings"

	"github.com/ProxeusApp/proxeus-core/sys"

	extNode "github.com/ProxeusApp/proxeus-core/externalnode"
	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/workflow"
)

type (

	// WorkflowService is an interface that provides workflow functions
	WorkflowService interface {

		// List returns a list of all WorkflowItem that match "contains" and the provided storage.Options
		List(auth model.Auth, contains string, options storage.Options) ([]*model.WorkflowItem, error)

		// ListPublished returns a list of published WorkflowItem that match "contains" and the provided storage.Options
		ListPublished(auth model.Auth, contains string, options storage.Options) ([]*model.WorkflowItem, error)

		// ListIds returns a list of all workflow ids that match contains and the provided storage.Options
		ListIds(auth model.Auth, contains string, options storage.Options) ([]string, error)

		// GetAndPopulateOwner a workflow by the provided id and sets the OwnerEthAddress
		GetAndPopulateOwner(auth model.Auth, id string) (*model.WorkflowItem, error)

		// Get returns a workflow by the provided id
		Get(auth model.Auth, id string) (*model.WorkflowItem, error)

		// Publish publishes a workflowItem
		Publish(auth model.Auth, wfItem *model.WorkflowItem) map[string]interface{}

		// Put saves a WorkflowItem
		Put(auth model.Auth, wfItem *model.WorkflowItem) error

		// Delete removes a WorkflowItem
		Delete(auth model.Auth, id string) error

		// InstantiateExternalNode creates a new instance of an external node
		InstantiateExternalNode(auth model.Auth, nodeId, nodeName string) (*extNode.ExternalQuery, error)

		// CopyWorkflows copies the workflow and related forms and templates to the new user
		CopyWorkflows(rootUser, newUser *model.User)

		// GetPublished returns a workflow item matching the supplied filter options that if it is flagged as published
		GetPublished(auth model.Auth, id string) (*model.WorkflowItem, error)
	}

	DefaultWorkflowService struct {
		userService UserService
	}
)

func NewWorkflowService(userService UserService) *DefaultWorkflowService {
	return &DefaultWorkflowService{userService: userService}
}

// List returns a list of all WorkflowItem that match "contains" and the provided storage.Options
func (me *DefaultWorkflowService) List(auth model.Auth, contains string, options storage.Options) ([]*model.WorkflowItem, error) {
	if auth.(*sys.Session) == nil {
		return nil, model.ErrAuthorityMissing
	}
	return workflowDB().List(auth, contains, options)
}

// ListPublished returns a list of published WorkflowItem that match "contains" and the provided storage.Options
func (me *DefaultWorkflowService) ListPublished(auth model.Auth, contains string, options storage.Options) ([]*model.WorkflowItem, error) {
	if auth.(*sys.Session) == nil {
		return nil, model.ErrAuthorityMissing
	}
	return workflowDB().ListPublished(auth, contains, options)
}

// Get returns a workflow by the provided id
func (me *DefaultWorkflowService) Get(auth model.Auth, id string) (*model.WorkflowItem, error) {
	return workflowDB().Get(auth, id)
}

// GetAndPopulateOwner a workflow by the provided id and sets the OwnerEthAddress
func (me *DefaultWorkflowService) GetAndPopulateOwner(auth model.Auth, id string) (*model.WorkflowItem, error) {
	workflow, err := me.Get(auth, id)
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

// ListIds returns a list of all workflow ids that match contains and the provided storage.Options
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

// Publish publishes a workflowItem
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
			it, er := formDB().Get(auth, node.ID)
			if er != nil {
				collectError(er, node)
				return true //continue
			}
			if !it.Published {
				it.Published = true
				er = formDB().Put(auth, it)
				if er != nil {
					collectError(er, node)
				}
			}
		} else if node.Type == "template" {
			it, er := templateDB().Get(auth, node.ID)
			if er != nil {
				collectError(er, node)
				return true //continue
			}
			if !it.Published {
				it.Published = true
				er = templateDB().Put(auth, it)
				if er != nil {
					collectError(er, node)
				}
			}
		} else if node.Type == "workflow" { // deep dive...
			it, er := workflowDB().Get(auth, node.ID)
			if er != nil {
				collectError(er, node)
				return true //continue
			}
			if !it.Published {
				it.Published = true
				er = workflowDB().Put(auth, it)
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

// Put saves a WorkflowItem
func (me *DefaultWorkflowService) Put(auth model.Auth, wfItem *model.WorkflowItem) error {
	return workflowDB().Put(auth, wfItem)
}

// Delete removes a WorkflowItem
func (me *DefaultWorkflowService) Delete(auth model.Auth, id string) error {
	return workflowDB().Delete(auth, id)
}

// InstantiateExternalNode creates a new instance of an external node
func (me *DefaultWorkflowService) InstantiateExternalNode(auth model.Auth, nodeId, nodeName string) (*extNode.ExternalQuery, error) {
	externalNode, err := workflowDB().NodeByName(auth, nodeName)
	if err != nil {
		log.Println("UpdateHandler instantiateExternalNode NodeByName err: ", err.Error())
		return nil, err
	}

	externalNodeInstance, err := workflowDB().QueryFromInstanceID(auth, nodeId)
	if err == nil {
		log.Println("UpdateHandler externalNode instance exists, will not create new instance")
		return &extNode.ExternalQuery{externalNode, &externalNodeInstance}, nil
	}

	newExternalNodeInstance := &extNode.ExternalNodeInstance{
		ID:       nodeId,
		NodeName: nodeName,
	}

	//PutExternalNodeInstance
	err = workflowDB().PutExternalNodeInstance(auth, newExternalNodeInstance)
	if err != nil {
		log.Println("UpdateHandler externalNode PutExternalNodeInstance error: ", err.Error())
		return nil, err
	}

	newExternalQuery := extNode.ExternalQuery{
		ExternalNode:         externalNode,
		ExternalNodeInstance: newExternalNodeInstance,
	}

	return &newExternalQuery, nil
}

// CopyWorkflows copies the workflow and related forms and templates to the new user
func (me *DefaultWorkflowService) CopyWorkflows(rootUser, newUser *model.User) {
	log.Println("Copy workflows to new user, if any...")
	// If some default workflows have to be assigned to the user, then clone them
	settings, err := settingsDB().Get()
	if err != nil {
		log.Printf("Unable to get settingsDB, retrieve list of workflows. Please check the ids exist. Error: %s", err.Error())
		return
	}
	workflowIds := strings.Split(settings.DefaultWorkflowIds, ",")
	workflows, err := workflowDB().GetList(rootUser, workflowIds)
	if err != nil {
		log.Printf("Can't retrieve list of workflows (%v). Please check the ids exist. Error: %s", workflowIds, err.Error())
		return
	}
	for _, loopWorkflow := range workflows {
		w := loopWorkflow.Clone()
		w.OwnerEthAddress = newUser.EthereumAddr
		w.Owner = newUser.ID
		newNodes := make(map[string]*workflow.Node)
		oldToNewIdsMap := make(map[string]string)
		for oldId, node := range w.Data.Flow.Nodes {
			if node.Type == "form" {
				form, er := formDB().Get(rootUser, node.ID)
				if er != nil {
					log.Println(err.Error())
				}
				f := form.Clone()
				er = formDB().Put(newUser, &f)
				if er != nil {
					log.Println("can't put form" + err.Error())
				}

				oldToNewIdsMap[node.ID] = f.ID
				node.ID = f.ID
				newNodes[node.ID] = node
				delete(w.Data.Flow.Nodes, oldId)

			} else if node.Type == "template" {
				template, er := templateDB().Get(rootUser, node.ID)
				if er != nil {
					log.Println(err.Error())
				}
				t := template.Clone()
				er = templateDB().Put(newUser, &t)
				if er != nil {
					log.Println("can't put template" + err.Error())
				}
				oldToNewIdsMap[node.ID] = t.ID
				node.ID = t.ID
				newNodes[node.ID] = node
				delete(w.Data.Flow.Nodes, oldId)
			} else {
				newNodes[node.ID] = node
			}
		}
		oldStartNodeId := w.Data.Flow.Start.NodeID
		if _, ok := oldToNewIdsMap[oldStartNodeId]; ok {
			w.Data.Flow.Start.NodeID = oldToNewIdsMap[oldStartNodeId]
		}

		// Now go through all connections and map them with the new ids
		for _, node := range newNodes {
			for _, connection := range node.Connections {
				if _, ok := oldToNewIdsMap[connection.NodeID]; ok {
					connection.NodeID = oldToNewIdsMap[connection.NodeID]
				}
			}
		}
		w.Data.Flow.Nodes = newNodes
		workflowDB().Put(newUser, &w)
	}
}

// GetPublished returns a workflow item matching the supplied filter options that if it is flagged as published
func (me *DefaultWorkflowService) GetPublished(auth model.Auth, id string) (*model.WorkflowItem, error) {
	return workflowDB().GetPublished(auth, id)
}
