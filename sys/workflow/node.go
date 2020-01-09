package workflow

import (
	"fmt"

	"github.com/ProxeusApp/proxeus-core/sys/model/compatability"
)

type (
	NodeIF interface {
		//Execute is being called when the node becomes the current target, doesn't matter whether it goes forward or backward.
		//If the node was executed before and returned proceed = false, the same instance is being used again.
		Execute(node *Node) (proceed bool, err error)
		//Remove is being called when this node is not part of the path anymore. In other words, when it doesn't exist in the State
		Remove(node *Node)
		//Close will be called on the instance before the next node is being executed or when the end is reached
		//Close is always the last called method either Execute()* -> Close() OR Execute()* -> Remove() -> Close().
		//Execute() can be called n times before Close() or Remove() as it is controlled by the node impl.
		//When it is called, the instance is released from the engine.
		Close()
	}

	InitImplFunc func(n *Node) (NodeIF, error)

	Node struct {
		ID           string                         `json:"id"`
		Name         string                         `json:"name"`
		Type         string                         `json:"type"`
		Detail       string                         `json:"detail,omitempty"`
		Data         compatability.CarriedStringMap `json:"data,omitempty"`
		Cases        []*Case                        `json:"cases,omitempty"`
		Connections  []*Connection                  `json:"conns,omitempty"`
		Position     Position                       `json:"p"`
		context      *context                       `json:"-"`
		internalNode bool                           `json:"-"`
		err          error                          `json:"-"`
		impl         NodeIF                         `json:"-"`
		background   bool                           `json:"-"`
		new          InitImplFunc
	}
)

func (n *Node) WFUniqueID() string {
	return n.HierarchyPath() + "." + n.ID
}

func (n *Node) HierarchyPath() string {
	path := n.context.id
	p := n.context.parent
	for p != nil {
		path = p.id + "." + path
		p = p.parent
	}
	return path
}

func (n *Node) String() string {
	return fmt.Sprintf("HierarchyPath:%s, ID: %s, Type: %s, Name: %s, Detail: %s", n.HierarchyPath(), n.ID, n.Type, n.Name, n.Detail)
}

func (n *Node) focus() {
	n.context.engine.target = n.context
	n.context.current = n.ID
}

func (n *Node) isCondition() bool {
	return n.Type == "condition"
}

func (n *Node) isWorkflow() bool {
	return n.Type == "workflow"
}

func (me *Node) getImpl() (NodeIF, error) {
	//instance is still valid for reuse
	if me.impl != nil {
		return me.impl, nil
	}
	if me.new == nil {
		return nil, ErrNodeImplementationNotProvided
	}
	//create new instance
	var err error
	me.impl, err = me.new(me)
	if err != nil {
		return nil, err
	}
	return me.impl, nil
}

func (me *Node) remove() {
	if !me.internalNode && me.impl != nil {
		me.impl.Remove(me)
		me.impl.Close()
		me.impl = nil
	}
}

func (me *Node) resetState() {
	if !me.internalNode && me.impl != nil {
		me.impl.Close()
		me.impl = nil
	}
}

func (me *Node) toStep() Step {
	return Step{NodeID: me.ID, Type: me.Type, HierarchyPath: me.HierarchyPath(), Name: me.Name, Detail: me.Detail}
}

func (me *Node) close() {
	me.context = nil
	if me.impl != nil {
		me.impl.Close()
		me.impl = nil
	}
	me.err = nil
	me.new = nil
}
