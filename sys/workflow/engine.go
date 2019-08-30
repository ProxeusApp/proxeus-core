package workflow

import (
	"sync"

	"github.com/robertkrimen/otto"
)

/**
In this package there are a few words about internal, background and foreground nodes. The meaning of them is here explained.
background node:
	A background node is not blocking the process.
	When you call Next, you get a background node with Current if it is a node at the end of the workflow otherwise you will get a foreground node if no error occurred.
	In the background nodes implementation there should be always returned < proceed > true if no error occurred otherwise it is configured wrong.
	Which means Execute is called once if no error occurred.
	Background nodes can be skipped in the Previous call to be able to get back to the last foreground node if no error occurs.
	Please note, there is no detection for wrong configuration of background nodes during runtime.
	If your node implementation changes and has an effect on the execution type, you must keep the background flag consistent.
foreground node:
	A foreground node is blocking the process.
	When you call Next usually you should get the foreground node by calling Current. Foreground nodes have states in the node implementation that causes Execute to be called more than once.
	For example a form node. It has at least two states. State one is presenting the form, state two or higher is validating the form.
internal node:
	Currently there are just two internal nodes, which is workflow and condition.
	These nodes are handled by the engine itself, they are visible in the Stack but they can also be visible in Current.
	For example if the root workflow starts with another workflow or a condition.

*/
type (
	Config struct {
		//GetWorkflow is your function that provides workflows for embedded use.
		GetWorkflow func(id string) (*Workflow, error)
		//State is provided by the workflow engine. If you pass it for a new instance you should get the same state as before.
		State []Step
		//GetData is used in the execution of the condition and it provides the data in the Execute function of the NodeIF.
		GetData func() interface{}
		//NodeImpl contains your implementations of the node types. The key of the map is the node type.
		NodeImpl map[string]*NodeDef
	}

	NodeDef struct {
		//Impl is an sample instance of the node impl. The purpose of this sample is only to get reflect.Type from your implementation.
		InitImplFunc func(n *Node) (impl NodeIF, err error)
		//Background true means it is a none blocking node.
		//Background false means the node has more than one inner state that will cause multiple executions.
		Background bool
	}

	Engine struct {
		jsParser       *JS
		getWorkflow    func(id string) (*Workflow, error)
		getData        func() interface{}
		onNodeImplInit func(n *Node, implInstance interface{}) error
		//holds only user relevant nodes, internally it is used to get back to the previous node
		steps     []Stack
		stepsSize int
		//holds every node execution
		stack     []Stack
		nodeImpls map[string]*NodeDef
		target    *context
		root      *context
		m         sync.Mutex
		hasPrev   bool
	}
)

func New(wf *Workflow, conf Config) (*Engine, error) {
	me := &Engine{}
	me.steps = make([]Stack, 0)
	me.stack = make([]Stack, 0)
	me.getWorkflow = conf.GetWorkflow
	me.getData = conf.GetData
	if me.getData == nil {
		me.getData = func() interface{} {
			return nil
		}
	}
	me.nodeImpls = conf.NodeImpl
	me.jsParser = NewJSParser()

	ctxt, err := newWorkflow(wf, me, nil)
	if err != nil {
		return nil, err
	}
	me.root = ctxt
	err = me.recoverState(&conf)
	if err != nil {
		return nil, err
	}
	return me, nil
}

//LoopNext is the same as Next but it combines hasNext and err into one bool so it can be used easier in a loop.
//To ensure if everything went well, you can just call Current() which provides the recent error
//This call is thread safe.
func (me *Engine) LoopNext() bool {
	hasNext, err := me.Next()
	return hasNext && err == nil
}

//Next is moving one node further in the workflow
//The recent error can be retrieved by Current() as well.
//This call is thread safe.
func (me *Engine) Next() (bool, error) {
	me.m.Lock()
	defer me.m.Unlock()
	return me.next()
}

func (me *Engine) next() (bool, error) {
	if me.root.hasNext {
		err := me.target.resolve(me.target.getCurrent())
		return me.root.hasNext, err
	}
	return me.root.hasNext, nil
}

//HasPrevious is providing the last state of being able to move one step back.
//The result is guaranteed.
//This call is thread safe.
func (me *Engine) HasPrevious() bool {
	me.m.Lock()
	defer me.m.Unlock()
	return me.hasPrev
}

//HasNext is providing the last state of being able to move forward.
//HasNext result is not guaranteed as the path can change during the execution with the provided data.
//This call is thread safe.
func (me *Engine) HasNext() bool {
	me.m.Lock()
	defer me.m.Unlock()
	return me.root.hasNext
}

//Stack is providing a slice containing all the executed nodes no matter if it was during forward or backward.
//It contains internal nodes like workflow or condition.
//This call is thread safe.
func (me *Engine) Stack() []Stack {
	me.m.Lock()
	defer me.m.Unlock()
	return me.stack
}

//State is providing a slice containing the path with background and foreground nodes.
//It doesn't contain any internal nodes like workflow or condition.
//This state can be provided during initialization in the config to recover the previous state.
//This call is thread safe.
func (me *Engine) State() []Step {
	me.m.Lock()
	defer me.m.Unlock()
	state := make([]Step, len(me.steps))
	for i, s := range me.steps {
		state[i] = s.Node.toStep()
	}
	return state
}

//Current is providing the target node and the recent error.
//The target node can be an internal node or background node as well. It depends on the structure of the workflow.
//This call is thread safe.
func (me *Engine) Current() (*Node, error) {
	me.m.Lock()
	defer me.m.Unlock()
	return me.target.curr()
}

//LoopPrevious is the same as Previous but it combines hasPrev and err into one bool, so it can be used easier in a loop.
//To ensure if everything went well, you can just call Current() which provides the recent error
//This call is thread safe.
func (me *Engine) LoopPrevious(skipBackgroundNodes bool) bool {
	hasPrev, err := me.Previous(skipBackgroundNodes)
	return hasPrev && err == nil
}

//Previous is moving one or more steps back. A step consists of background nodes and foreground nodes.
//By skipping background nodes it moves back to a foreground node if there is one otherwise it stops at the beginning.
//This call is thread safe.
func (me *Engine) Previous(skipBackgroundNodes bool) (bool, error) {
	me.m.Lock()
	defer me.m.Unlock()
	if !me.hasPrev {
		return me.hasPrev, nil
	}
	var n *Node
	var lastNode *Node
	foregroundCount := 0
	i := me.stepsSize - 1
	count := 0
	for i >= 0 {
		if lastNode != nil {
			n.remove()
		}
		n = me.steps[i].Node
		if !n.background {
			foregroundCount++
		}
		if count == 0 {
			n.remove()
		} else {
			n.resetState()
			_, err := me.execute(n, false)
			if err != nil {
				me.setCurrent(i)
				return me.hasPrev, err
			}
		}
		if skipBackgroundNodes && foregroundCount < 2 {
			lastNode = n
			count++
			i--
			continue
		}
		if !skipBackgroundNodes && count == 0 {
			i--
			count++
			continue
		}
		break
	}
	if i <= 0 {
		//reset from scratch
		me.stepsSize = 0
		me.steps = make([]Stack, 0)
		me.root.start()
		me.hasPrev = false
		me.root.hasNext = true
	} else {
		me.setCurrent(i)
	}

	if me.stepsSize > 0 {
		return me.hasPrev, me.steps[me.stepsSize-1].Node.err
	}
	return me.hasPrev, nil
}

func (me *Engine) setCurrent(i int) {
	if i < 0 {
		i = 0
	}
	me.stepsSize = i + 1
	me.steps = me.steps[:me.stepsSize]
	me.hasPrev = me.stepsSize > 0
	me.root.hasNext = true
}

func (me *Engine) recoverState(conf *Config) error {
	lsteps := len(conf.State)
	if lsteps > 0 {
		callNext := true
		acceptableState := false
		for {
			if lsteps > len(me.steps) && callNext {
				if hasNext, err := me.next(); hasNext {
					if err != nil {
						return err
					}
				} else {
					break
				}
			}
			cl := len(me.steps)
			correctCount := 0
			for a, b := range conf.State {
				if a+1 >= cl {
					break
				}
				if me.steps[a].is(&b) {
					correctCount++
				}
			}
			if cl-1 == correctCount {
				if !callNext || lsteps-1 == correctCount {
					acceptableState = true
					break
				}
			}
			if cl-1 > correctCount {
				callNext = false
				_, err := me.Previous(false)
				if err != nil {
					return err
				}
			}
		}
		if !acceptableState {
			//not acceptable state! just reset the workflow but don't return an error
			//as it can happen when the workflow was completely changed after this state was recorded
			for me.LoopPrevious(false) {
			}
			me.stack = make([]Stack, 0)
		}
	}
	return nil
}

func (me *Engine) setupNodes(ctx *context) error {
	if len(me.nodeImpls) == 0 {
		return ErrNodeImplementationNotProvided
	}
	for _, item := range ctx.flow.Nodes {
		item.internalNode = item.isCondition() || item.isWorkflow()
		if !item.internalNode {
			if s, ok := me.nodeImpls[item.Type]; ok && s.InitImplFunc != nil {
				item.new = s.InitImplFunc
				item.background = s.Background
			} else {
				return ErrNodeImplementationNotProvided
			}
		}
		item.context = ctx
	}
	return nil
}

func (me *Engine) execute(nn *Node, considerSteps bool) (proceed bool, err error) {
	nn.focus()
	if nn.internalNode {
		proceed = true
	} else {
		var nImpl NodeIF
		nImpl, err = nn.getImpl()
		if err == nil {
			shouldAddStep := me.stepsSize == 0 || me.steps[me.stepsSize-1].Node != nn
			if considerSteps && shouldAddStep {
				if me.stepsSize > 0 {
					//ensure last node is closed
					n := me.steps[me.stepsSize-1].Node
					n.resetState()
				}
			}
			proceed, err = nImpl.Execute(nn)
			nn.err = err
			if considerSteps && shouldAddStep {
				me.steps = append(me.steps, Stack{Node: nn})
				me.stepsSize++
				me.hasPrev = me.stepsSize > 0
			}
		}
	}
	me.stack = append(me.stack, Stack{Node: nn})
	return
}

func (me *Engine) checkNodeImpl(n *Node) bool {
	if len(me.nodeImpls) > 0 {
		if s, ok := me.nodeImpls[n.Type]; ok && s.InitImplFunc != nil {
			return true
		}
	}
	return false
}

func (me *Engine) setJSGlobalData() error {
	return me.jsParser.SetGlobal(me.getData())
}

func (me *Engine) jsEval(script *string) (otto.Value, error) {
	return me.jsParser.Run(*script)
}

func (me Stack) is(step *Step) bool {
	return me.Node.ID == step.NodeID && me.Node.HierarchyPath() == step.HierarchyPath
}

func (me *Engine) Close() error {
	if me.jsParser != nil {
		me.jsParser.Close()
		me.jsParser = nil
	}
	me.stack = nil
	me.getWorkflow = nil
	me.getData = nil
	me.steps = nil
	me.nodeImpls = nil
	me.target = nil
	if me.root != nil {
		me.root.close()
	}
	return nil
}
