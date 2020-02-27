package app

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/storage/mock"
	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/ProxeusApp/proxeus-core/sys/workflow"
)

var auth = fakeAuth{}

func TestNewDocumentApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userDataItem := model.UserDataItem{}
	workflowDB := mock.NewMockWorkflowIF(ctrl)
	system := &sys.System{
		DB: &storage.DBSet{
			Workflow: workflowDB,
		},
	}

	t.Run("when given a nil system", func(t *testing.T) {
		document, err := NewDocumentApp(&userDataItem, auth, nil, "any", "any")

		assert.Equal(t, os.ErrInvalid, err, "should return os.ErrInvalid")
		assert.Nil(t, document, "should return an empty document")
	})

	t.Run("when given an empty wfid", func(t *testing.T) {
		document, err := NewDocumentApp(&userDataItem, auth, system, "", "any")

		assert.Equal(t, os.ErrInvalid, err, "should return os.ErrInvalid")
		assert.Nil(t, document, "should return an empty document")
	})

	t.Run("when given an existing wfid", func(t *testing.T) {
		workflowResponse := &model.WorkflowItem{
			Permissions: model.Permissions{},
			ID:          "",
			Name:        "",
			Detail:      "",
			Updated:     time.Time{},
			Created:     time.Time{},
			Price:       0,
			Data: &workflow.Workflow{
				Flow: &workflow.Flow{
					Start: &workflow.Start{
						NodeID:   "111",
						Position: workflow.Position{},
					},
					Nodes: nil,
				},
			},
			OwnerEthAddress: "",
			Deactivated:     false,
		}
		workflowDB.EXPECT().Get(auth, "any").Times(1).Return(workflowResponse, nil)

		document, err := NewDocumentApp(&userDataItem, auth, system, "any", "any")

		assert.Nil(t, err, "should return no error")
		assert.Equal(t, "any", document.WFID, "should have the correct WFID")
		assert.Equal(t, auth, document.auth, "should have given auth")
	})

}

func TestDocumentFlowInstance_Init(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workflowDB := mock.NewMockWorkflowIF(ctrl)
	system := &sys.System{
		DB: &storage.DBSet{
			Workflow: workflowDB,
		},
	}

	t.Run("when workflow is not found", func(t *testing.T) {
		workflowDB.EXPECT().Get(auth, "any").Return(nil, errors.New("not found"))

		documentFlowInstance := DocumentFlowInstance{WFID: "any"}
		err := documentFlowInstance.Init(auth, system)

		assert.NotNil(t, err, "should return error")
	})

	t.Run("when no error is returned & no workflow", func(t *testing.T) {
		workflowDB.EXPECT().Get(auth, "any").Return(nil, nil).Times(1)

		documentFlowInstance := DocumentFlowInstance{WFID: "any"}
		err := documentFlowInstance.Init(auth, system)
		assert.Equal(t, os.ErrNotExist, err)
	})

	t.Run("when workflow is returned", func(t *testing.T) {
		workflowDB.EXPECT().Get(auth, "any").Return(nil, nil).Times(1)

		documentFlowInstance := DocumentFlowInstance{WFID: "any"}
		err := documentFlowInstance.Init(auth, system)
		assert.Equal(t, os.ErrNotExist, err)
	})
}

func TestDocumentFlowInstance_isLangAvailable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	il8n := mock.NewMockI18nIF(ctrl)

	documentFlowInstance := DocumentFlowInstance{
		system: &sys.System{
			DB: &storage.DBSet{
				I18n: il8n,
			},
		},
	}

	t.Run("when retrieval of languages return error", func(t *testing.T) {
		il8n.EXPECT().GetLangs(true).Return(nil, errors.New("any")).Times(1)

		available, err := documentFlowInstance.isLangAvailable("en")

		assert.False(t, available, "should return available=false")
		assert.NotNil(t, err, "should return error")
	})

	t.Run("when retrieval of languages returns languages but not the one we're requesting", func(t *testing.T) {
		langs := []*model.Lang{
			{
				ID:      "",
				Code:    "it",
				Enabled: false,
			},
		}
		il8n.EXPECT().GetLangs(true).Return(langs, nil).Times(1)

		available, err := documentFlowInstance.isLangAvailable("en")

		assert.Nil(t, err, "should not return error")
		assert.False(t, available, "should return not available")
	})

	t.Run("when requested language is there", func(t *testing.T) {
		langs := []*model.Lang{
			{
				ID:      "",
				Code:    "en",
				Enabled: false,
			},
		}
		il8n.EXPECT().GetLangs(true).Return(langs, nil).Times(1)

		available, err := documentFlowInstance.isLangAvailable("en")
		assert.Nil(t, err, "should not return error")
		assert.True(t, available)
	})
}

// Fakes
type fakeAuth struct{}

func (me fakeAuth) UserID() string {
	return "5"
}
func (me fakeAuth) AccessRights() model.Role {
	return model.CREATOR
}
