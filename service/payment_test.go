package service

import (
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckIfWorkflowPaymentRequired(t *testing.T) {
	t.Run("ShouldRequirePaymentIfWorkflowNotForFree", func(t *testing.T) {
		permissions := &model.Permissions{Owner: "1"}
		workflow := &model.WorkflowItem{Price: 2, Permissions: *permissions}
		assert.True(t, isPaymentRequired(false, workflow, "2"))
	})
	t.Run("ShouldNotRequirePaymentIfWorkflowNotForFreeButAlreadyStarted", func(t *testing.T) {
		permissions := &model.Permissions{Owner: "1"}
		workflow := &model.WorkflowItem{Price: 2, Permissions: *permissions}
		assert.False(t, isPaymentRequired(true, workflow, "2"))
	})
	t.Run("ShouldNotRequirePaymentIfWorkflowIsFree", func(t *testing.T) {
		permissions := &model.Permissions{Owner: "1"}
		workflow := &model.WorkflowItem{Price: 0, Permissions: *permissions}
		assert.False(t, isPaymentRequired(false, workflow, "2"))
	})
	t.Run("ShouldNotRequirePaymentForWorkflowOwner", func(t *testing.T) {
		permissions := &model.Permissions{Owner: "1"}
		workflow := &model.WorkflowItem{Price: 2, Permissions: *permissions}
		assert.False(t, isPaymentRequired(false, workflow, "1"))
	})
}
