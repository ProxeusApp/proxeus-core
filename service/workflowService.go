package service

import (
	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type (
	WorkflowService interface {
		List(auth model.Auth, contains string, options storage.Options) ([]*model.WorkflowItem, error)
		ListIds(auth model.Auth, contains string, options storage.Options) ([]string, error)
		GetAndPopulateOwner(auth model.Auth, id string) (*model.WorkflowItem, error)
		Get(auth model.Auth, id string) (*model.WorkflowItem, error)
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
