package workflow

import (
	"errors"
)

type (
	Step struct {
		NodeID        string `json:"nID"`
		HierarchyPath string `json:"wfID"`
		Name          string `json:"name"`
		Type          string `json:"type"`
		Detail        string `json:"detail,omitempty"`
	}
	Stack struct {
		Node *Node
	}
	Workflow struct {
		Flow *Flow `json:"flow"`
	}
	Flow struct {
		Start *Start           `json:"start"`
		Nodes map[string]*Node `json:"nodes"`
	}
	Position struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}
	Start struct {
		NodeID   string   `json:"node"`
		Position Position `json:"p"`
	}
	Case struct {
		Name   string                 `json:"name"`
		Value  interface{}            `json:"value"`
		Detail string                 `json:"detail,omitempty"`
		Data   map[string]interface{} `json:"data,omitempty"`
	}
	Connection struct {
		NodeID    string      `json:"id"`
		CaseValue interface{} `json:"value,omitempty"`
	}
	//ensures entities are looped once to prevent from an endless loop in nested patterns
	Looper struct {
		accessedOnes map[string]bool
		cb           LooperCallback
	}
	LooperCallback func(l *Looper, node *Node) bool

	context struct {
		id       string
		flow     *Flow
		current  string
		hasNext  bool
		children map[string]*context
		parent   *context
		engine   *Engine
	}
)

var (
	ErrWorkflowMissing               = errors.New("cannot run without Workflow")
	ErrStartNodeMissing              = errors.New("cannot run without the start node")
	ErrConditionDataMissing          = errors.New("condition data missing")
	ErrConditionCaseConnection       = errors.New("condition cases are not leading in any case to a node")
	ErrConditionConnectionsMissing   = errors.New("condition connections missing")
	ErrNotExist                      = errors.New("does not exist")
	ErrConfigGetWorkflowMissing      = errors.New("GetWorkflow missing in your config")
	ErrNodeImplementationNotProvided = errors.New("node impl not provided")
)

func newWorkflow(wf *Workflow, engine *Engine, parent *context) (*context, error) {
	if wf == nil {
		return nil, ErrWorkflowMissing
	}

	if wf.Flow == nil || wf.Flow.Start == nil || wf.Flow.Start.NodeID == "" {
		return nil, ErrStartNodeMissing
	}

	me := context{
		id:      "root",
		flow:    wf.Flow,
		hasNext: true,
		parent:  parent,
		engine:  engine,
	}

	err := engine.removeUselessNodes(&me)
	if err != nil {
		return nil, err
	}

	err = engine.setupNodes(&me)
	if err != nil {
		return nil, err
	}

	me.start()

	return &me, nil
}

func (me *context) start() {
	me.engine.target = me
	me.current = me.flow.Start.NodeID
}

func (me *context) curr() (n *Node, err error) {
	n = me.getCurrent()
	if n != nil {
		if n.err != nil {
			err = n.err
		}
	} else {
		err = ErrNotExist
	}
	return
}

func (me *context) end() {
	me.hasNext = false
}

func (me *context) resolve(n *Node) error {
	if n != nil {
		if n.internalNode {
			if n.isCondition() {
				nn, err := me.condition(n)
				if err != nil {
					return err
				}
				return me.resolve(nn)
			} else if n.isPlaceholder() {
				return nil
			} else if n.isWorkflow() {
				return me.stepIntoWorkflow(n)
			}
		} else {
			var proceed bool
			var err error
			proceed, err = me.engine.execute(n, true)
			if err != nil {
				return err
			}
			if proceed {
				n.resetState()
				if me.hasNext = len(n.Connections) > 0; me.hasNext {
					con := n.Connections[0]
					if con != nil {
						return me.resolve(me.getNode(&con.NodeID))
					}
				}
			} else {
				return nil
			}
		}
	}
	if n == nil {
		n = me.getCurrent()
	}
	if me.isNotRoot() { //step out of workflow
		me.end()
		n.resetState()
		t := me.parent
		tn := t.getCurrent()
		for {
			if t.hasNext = len(tn.Connections) > 0; t.hasNext {
				con := tn.Connections[0]
				if con != nil {
					return t.resolve(t.getNode(&con.NodeID))
				}
			} else if t.isNotRoot() {
				t = t.parent
				tn = t.getCurrent()
			} else {
				break
			}
		}
		return nil
	}
	me.end()
	n.resetState()
	return nil
}

func (me *context) isNotRoot() bool {
	return me.parent != nil
}

func (me *context) getCurrent() *Node {
	return me.getNode(&me.current)
}

func (me *context) getNode(id *string) *Node {
	return me.flow.Nodes[*id]
}

func (me *context) condition(n *Node) (*Node, error) {
	if len(n.Connections) > 0 {
		n.err = nil
		if n.Data == nil {
			n.err = ErrConditionDataMissing
			return nil, n.err
		}
		jScript, ok := n.Data["js"].(string)
		if ok && jScript != "" {
			err := me.engine.setJSGlobalData()
			if err != nil {
				n.err = err
				return nil, err
			}
			jScript = jScript + "; if(!input){input={};}condition();"
			v, err := me.engine.jsEval(&jScript)
			var caseValue interface{}
			if err != nil {
				n.err = err
				return nil, err
			}
			if v.IsNumber() {
				caseValue, err = v.ToFloat()
			} else {
				caseValue, err = v.Export()
			}
			if err != nil {
				n.err = err
				return nil, err
			}
			me.engine.execute(n, false)
			for _, conn := range n.Connections {
				if conn != nil && conn.CaseValue == caseValue {
					return me.getNode(&conn.NodeID), nil
				}
			}
			n.err = ErrConditionCaseConnection
			return nil, n.err
		} else {
			n.err = ErrConditionDataMissing
			return nil, n.err
		}
	} else {
		n.err = ErrConditionConnectionsMissing
		return nil, n.err
	}
}

func (me *context) stepIntoWorkflow(n *Node) error {
	n.err = nil
	var newContext *context
	if me.children == nil {
		me.children = make(map[string]*context)
	} else {
		newContext = me.children[n.ID]
	}
	if newContext == nil {
		if me.engine.getWorkflow == nil {
			n.err = ErrConfigGetWorkflowMissing
			return n.err
		}
		var wf *Workflow
		wf, n.err = me.engine.getWorkflow(n.ID)
		if n.err != nil {
			return n.err
		}
		me.engine.execute(n, false)
		newContext, n.err = newWorkflow(wf, me.engine, me)
		if n.err != nil {
			return n.err
		}
		newContext.id = n.ID
		me.children[n.ID] = newContext
	} else {
		//reuse embedded context
		me.engine.execute(n, false)
		newContext.start()
	}
	return newContext.resolve(newContext.getCurrent())
}

func (me *Flow) close() {
	for _, v := range me.Nodes {
		if v != nil {
			v.close()
		}
	}
}

// Loop is making node loops easier
// looper keeps track of accessed entities, cb can be nil after the first call to simplify nested calls
func (me *Workflow) Loop(looper *Looper, cb LooperCallback) {
	if looper != nil && looper.accessedOnes == nil {
		looper.accessedOnes = map[string]bool{}
	}
	for k, v := range me.Flow.Nodes {
		if looper != nil {
			if looper.accessedOnes[k] {
				continue
			}
			looper.accessedOnes[k] = true
			if cb == nil {
				cb = looper.cb
			} else if looper.cb == nil {
				looper.cb = cb
			}
		}
		if cb == nil {
			break
		}
		if !cb(looper, v) {
			break
		}
	}
}

func (me *context) close() {
	if me.flow != nil {
		me.flow.close()
	}
	for _, v := range me.children {
		if v != nil {
			v.close()
		}
	}
	me.children = nil
	me.parent = nil
	me.engine = nil
}
