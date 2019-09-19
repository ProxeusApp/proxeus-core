package workflow

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"git.proxeus.com/core/central/sys/validate"
)

//Use this when you are extending the tests
const testVerbose = false

type (
	FormImpl struct {
		n         *Node
		presented bool
	}
	UserImpl struct {
		n *Node
	}
	TemplateImpl struct {
		n *Node
	}
)

type FormFieldError struct { //TODO this is supposed to show how custom errors can be received in Next, Previous or Current of the workflow but not in use yet
	ValidationError map[string]validate.Error `json:"validationError"`
}

func (me *FormFieldError) Error() string {
	bts, err := json.Marshal(me)
	if err != nil {
		return ""
	}
	return string(bts)
}

//key = node.ID, value = true executed | false closed
//TODO this is supposed to help checking if the NodeIF calls are correct but not in use yet
type MyNodeState map[string]bool

var nodeStateMap = MyNodeState{}

func (me MyNodeState) onlyTrue(id string) bool {
	trueCount := 0
	targetTrue := false
	for k, v := range me {
		if v {
			trueCount++
		}
		if k == id && v {
			targetTrue = true
		}
	}
	return trueCount == 1 && targetTrue
}

func NewFormImpl(n *Node) (NodeIF, error) {
	return &FormImpl{}, nil
}
func NewTmplImpl(n *Node) (NodeIF, error) {
	return &TemplateImpl{}, nil
}
func NewUserImpl(n *Node) (NodeIF, error) {
	return &UserImpl{}, nil
}

func (me *FormImpl) Execute(n *Node) (bool, error) {
	if me.n == nil {
		me.n = n
	}
	nodeStateMap[n.ID] = true
	if !me.presented {
		//present
		if testVerbose {
			log.Println("--->WF TEST IMPL [form] Execute present state", n)
		}
		me.presented = true
		return false, nil
	}
	//validate
	if testVerbose {
		log.Println("--->WF TEST IMPL [form] Execute validate state", n)
	}
	return true, nil
}

func (me *FormImpl) Remove(n *Node) {
	if testVerbose {
		log.Println("--->WF TEST IMPL [form] Remove", n)
	}
}

func (me *FormImpl) Close() {
	if testVerbose {
		log.Println("--->WF TEST IMPL [form] Close", me.n)
	}
	nodeStateMap[me.n.ID] = false
	me.n = nil
}

func (me *UserImpl) Execute(n *Node) (bool, error) {
	if me.n == nil {
		me.n = n
	}
	nodeStateMap[n.ID] = true
	if testVerbose {
		log.Println("--->WF TEST IMPL [user] Execute", n)
	}
	return true, nil
}

func (me *UserImpl) Remove(n *Node) {
	if testVerbose {
		log.Println("--->WF TEST IMPL [user] Remove", n)
	}
}

func (me *UserImpl) Close() {
	if testVerbose {
		log.Println("--->WF TEST IMPL [user] Close", me.n)
	}
	nodeStateMap[me.n.ID] = false
	me.n = nil
}

func (me *TemplateImpl) Execute(n *Node) (bool, error) {
	if me.n == nil {
		me.n = n
	}
	nodeStateMap[n.ID] = true
	if testVerbose {
		log.Println("--->WF TEST IMPL [template] Execute", n)
	}
	return true, nil
}

func (me *TemplateImpl) Remove(n *Node) {
	if testVerbose {
		log.Println("--->WF TEST IMPL [template] Remove", n)
	}
}

func (me *TemplateImpl) Close() {
	if testVerbose {
		log.Println("--->WF TEST IMPL [template] Close", me.n)
	}
	nodeStateMap[me.n.ID] = false
	me.n = nil
}

var getWF = func(id string) (*Workflow, error) {
	if id != "" {
		if id == "complex" {
			wd := &Workflow{}
			err := json.Unmarshal([]byte(complexWF), &wd)
			return wd, err
		} else if id == "simple" {
			wd := &Workflow{}
			err := json.Unmarshal([]byte(simpleWF), &wd)
			return wd, err
		} else if id == "deep" {
			wd := &Workflow{}
			err := json.Unmarshal([]byte(deepWF), &wd)
			return wd, err
		} else if id == "deep2" {
			wd := &Workflow{}
			err := json.Unmarshal([]byte(deep2WF), &wd)
			return wd, err
		} else if id == "deep3" {
			wd := &Workflow{}
			err := json.Unmarshal([]byte(deep3WF), &wd)
			return wd, err
		} else if id == "deep4" {
			wd := &Workflow{}
			err := json.Unmarshal([]byte(deep4WF), &wd)
			return wd, err
		}
	}
	if testVerbose {
		log.Println("GetWorkflow dows not exist", id)
	}
	return nil, os.ErrNotExist
}

func wfWithDefaultConfig(wf *Workflow) (*Engine, error) {
	conf := Config{
		GetWorkflow: getWF,
		GetData: func() interface{} {
			return map[string]interface{}{"input": map[string]interface{}{"someVar": "someValue"}}
		},
		NodeImpl: map[string]*NodeDef{"form": {InitImplFunc: NewFormImpl, Background: false}, "template": {InitImplFunc: NewTmplImpl, Background: true}, "user": {InitImplFunc: NewUserImpl, Background: true}},
	}
	return New(wf, conf)
}

func TestDeepForwardOnlyLoop(t *testing.T) {
	wd := &Workflow{}
	err := json.Unmarshal([]byte(deepWF), &wd)
	if err != nil {
		t.Error(err)
	}
	p, err := wfWithDefaultConfig(wd)
	if err != nil {
		t.Error(err)
	}
	forwards := 1
	for p.LoopNext() {
		forwards++
		if forwards >= 50 {
			break
		}
	}
	if forwards >= 50 {
		n, err := p.Current()
		t.Error("index to high", n, err)
	}
	currentState := p.State()
	if len(currentState) != 5 {
		t.Error("expected 5", len(currentState))
	}
	forwards--
	for p.LoopPrevious(false) {
		forwards--
		if forwards < 0 {
			break
		}
	}
	if forwards != 1 {
		t.Error("expected 1", forwards)
	}

	currentState = p.State()
	if len(currentState) != 0 {
		t.Error("expected 0", len(currentState))
	}
	p.Close()
}

func TestDeepForwardOnly(t *testing.T) {
	wd := &Workflow{}
	err := json.Unmarshal([]byte(deepWF), &wd)
	if err != nil {
		t.Error(err)
	}
	p, err := wfWithDefaultConfig(wd)
	if err != nil {
		t.Error(err)
	}
	n, err := p.Current()
	if err != nil {
		t.Error(err, n)
	}
	if err != nil || n == nil || n.ID != "frm.1" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n, err)
	}
	hasNext, err := p.Next()
	if hasNext == false {
		t.Error("next can't be true! ")
	}
	n, err = p.Current()
	if err != nil || n == nil || n.ID != "frm.1" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n, err)
	}

	hasNext, err = p.Next()
	if hasNext == false {
		t.Error("next can't be true! ")
	}
	n, err = p.Current()
	if err != nil || n == nil || n.ID != "frm.2" || n.HierarchyPath() != "root.deep2" {
		t.Error("Current has an unexpected ID! ", n, err)
	}

	hasNext, err = p.Next()
	if hasNext == false {
		t.Error("next can't be true! ")
	}
	n, err = p.Current()
	if err != nil || n == nil || n.ID != "frm.3" || n.HierarchyPath() != "root.deep2.deep3" {
		t.Error("Current has an unexpected ID! ", n, err)
	}

	hasNext, err = p.Next()
	if hasNext == false {
		t.Error("next can't be true! ")
	}
	n, err = p.Current()
	if err != nil || n == nil || n.ID != "frm.2-2" || n.HierarchyPath() != "root.deep2" {
		t.Error("Current has an unexpected ID! ", n, err)
	}

	hasNext, err = p.Next()
	if hasNext == true {
		t.Error("next can't be true! ")
	}
	n, err = p.Current()
	if err != nil || n == nil || n.ID != "frm.2-2" || n.HierarchyPath() != "root.deep2" {
		t.Error("Current has an unexpected ID! ", n, err)
	}

	state := p.State()
	deepEndTest(p, t)

	//try to get to the same state
	conf := Config{
		GetWorkflow: getWF,
		State:       state,
		GetData: func() interface{} {
			return map[string]interface{}{"input": map[string]interface{}{"someVar": "someValue"}}
		},
		NodeImpl: map[string]*NodeDef{"form": {InitImplFunc: NewFormImpl, Background: false}, "template": {InitImplFunc: NewTmplImpl, Background: true}, "user": {InitImplFunc: NewUserImpl, Background: true}},
	}

	//wd = &Workflow{}
	//err = json.Unmarshal([]byte(deepWF), &wd)
	//if err != nil {
	//	t.Error(err)
	//}
	p, err = New(wd, conf)
	if err != nil {
		t.Error(err)
	}
	deepEndTest(p, t)
}

func deepEndTest(p *Engine, t *testing.T) {
	//nothing should happen as it reached already the end
	hasNext, err := p.Next()
	if hasNext == true || err != nil {
		t.Error(hasNext, err)
	}
	n, err := p.Current()
	if err != nil || n == nil || n.ID != "frm.2-2" || n.HierarchyPath() != "root.deep2" {
		t.Error("Current has an unexpected ID! ", n, err)
	}

	stack := p.Stack()
	i := 0
	stackSize := len(stack)
	if stackSize != 11 {
		t.Error("wrong stack size!")
		if testVerbose {
			for i, b := range p.State() {
				log.Println("State", i, b)
			}
			for i, b := range stack {
				log.Println("Stack", i, b.Node.ID, b.Node.HierarchyPath())
			}
		}
	}
	if stack[i].Node.ID != "frm.1" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm.1" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "deep2" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm.2" || stack[i].Node.HierarchyPath() != "root.deep2" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm.2" || stack[i].Node.HierarchyPath() != "root.deep2" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "deep3" || stack[i].Node.HierarchyPath() != "root.deep2" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm.3" || stack[i].Node.HierarchyPath() != "root.deep2.deep3" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm.3" || stack[i].Node.HierarchyPath() != "root.deep2.deep3" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "tmpl.3" || stack[i].Node.HierarchyPath() != "root.deep2.deep3" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm.2-2" || stack[i].Node.HierarchyPath() != "root.deep2" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm.2-2" || stack[i].Node.HierarchyPath() != "root.deep2" {
		t.Error("stack wrong!")
	}
	i++
	p.Close()
}

func TestSimpleForwardOnly(t *testing.T) {
	wd := &Workflow{}
	err := json.Unmarshal([]byte(deep4WF), &wd)
	if err != nil {
		t.Error(err)
	}
	p, err := wfWithDefaultConfig(wd)
	if err != nil {
		t.Error(err)
	}
	n, err := p.Current()
	if err != nil {
		t.Error(err, n)
	}
	if err != nil || n == nil || n.ID != "frm.4" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n, err)
	}
	hasNext, err := p.Next()
	if hasNext == false {
		t.Error("next can't be false! ")
	}
	n, err = p.Current()
	if err != nil || n == nil || n.ID != "frm.4" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n, err)
	}

	hasNext, err = p.Next()
	if hasNext == true {
		t.Error("next can't be true! ")
	}
	n, err = p.Current()
	if err != nil || n == nil || n.ID != "frm.4" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n, err)
	}

	//nothing should happen as it reached already the end
	hasNext, err = p.Next()
	if hasNext == true {
		t.Error("next can't be true! ")
	}
	n, err = p.Current()
	if err != nil || n == nil || n.ID != "frm.4" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n, err)
	}

	stack := p.Stack()
	i := 0
	stackSize := len(stack)
	if stackSize != 2 {
		t.Error("wrong stack size!")
		if testVerbose {
			for i, b := range stack {
				log.Println("Stack", i, b.Node.ID, b.Node.HierarchyPath())
			}
		}
	}
	if stack[i].Node.ID != "frm.4" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm.4" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	p.Close()
}

func TestOneBackgroundNodeForwardOnly(t *testing.T) {
	wd := &Workflow{}
	err := json.Unmarshal([]byte(oneBackgroundNode), &wd)
	if err != nil {
		t.Error(err)
	}
	p, err := wfWithDefaultConfig(wd)
	if err != nil {
		t.Error(err)
	}
	n, err := p.Current()
	if err != nil {
		t.Error(err, n)
	}
	if n == nil {
		t.Error("Current can't be nil! ", n)
	}
	if n.ID != "tmpl1" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}
	if p.HasNext() == false {
		t.Error("next can't be false! ")
	}
	hasNext, err := p.Next()
	if hasNext == true {
		t.Error("next can't be true! ")
	}
	n, err = p.Current()
	if err != nil {
		t.Error(err, n)
	}
	if n.ID != "tmpl1" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}

	//nothing should happen as it reached already the end
	hasNext, err = p.Next()
	if hasNext == true {
		t.Error("next can't be true! ")
	}
	n, err = p.Current()
	if n.ID != "tmpl1" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}

	stack := p.Stack()
	i := 0
	stackSize := len(stack)
	if stackSize != 1 {
		t.Error("wrong stack size!")
		if testVerbose {
			for i, b := range stack {
				log.Println("Stack", i, b.Node.ID, b.Node.HierarchyPath())
			}
		}
	}
	if stack[i].Node.ID != "tmpl1" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	p.Close()
}

func TestConditionOnlyWithConnections(t *testing.T) {
	myData := map[string]interface{}{"input": map[string]interface{}{"someVar": "someValue"}}
	wd := &Workflow{}
	err := json.Unmarshal([]byte(conditionOnlyWithConnections), &wd)
	if err != nil {
		t.Error(err)
	}
	conf := Config{
		GetWorkflow: getWF,
		GetData: func() interface{} {
			return myData
		},
		NodeImpl: map[string]*NodeDef{"form": {InitImplFunc: NewFormImpl, Background: false}, "template": {InitImplFunc: NewTmplImpl, Background: true}, "user": {InitImplFunc: NewUserImpl, Background: true}},
	}
	p, err := New(wd, conf)
	if err != nil {
		t.Error(err)
	}
	n, err := p.Current()
	if n.ID != "cond1" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}
	_, err = p.Next()
	if err != nil {
		t.Error("next shouldn't return an err!", err)
	}
	n, err = p.Current()
	if n.ID != "form1" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}
	var hasNext bool
	hasNext, err = p.Next()
	if hasNext == true {
		t.Error("next shouldn't return true!", err)
	}
	n, err = p.Current()
	if n.ID != "form1" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}

	ii := 1
	for p.LoopPrevious(true) {
		if ii >= 50 {
			t.Error("error called to often", ii)
			break
		}
	}

	//change data to walk through the other path
	if m, ok := myData["input"].(map[string]interface{}); ok {
		m["someVar"] = "some other value"
	} else {
		t.Error("couldn't change myData")
	}

	hasNext, err = p.Next()
	if hasNext == false || err != nil {
		t.Error(hasNext, err)
	}
	n, err = p.Current()
	if n.ID != "form2" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}
	hasNext, err = p.Next()
	if hasNext == true || err != nil {
		t.Error("next shouldn't return true!", err)
	}
	n, err = p.Current()
	if n.ID != "form2" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}
	hasNext, err = p.Next()
	if hasNext == true {
		t.Error("next shouldn't return true!", err)
	}
	n, err = p.Current()
	if n.ID != "form2" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}

	currentState := p.State()
	currentStateSize := len(currentState)
	if currentStateSize != 1 {
		t.Error("wrong currentState size!")
		if testVerbose {
			for i, b := range currentState {
				log.Println("currentState", i, b)
			}
		}
	}
	stack := p.Stack()
	i := 0
	stackSize := len(stack)
	if stackSize != 6 {
		t.Error("wrong stack size!")
		if testVerbose {
			for i, b := range stack {
				log.Println("Stack", i, b.Node.ID, b.Node.HierarchyPath())
			}
		}
	}
	if stack[stackSize-1].Node.ID != "form2" || stack[stackSize-1].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	p.Close()
}

func TestConditionOnlyNoConnections(t *testing.T) {
	wd := &Workflow{}
	err := json.Unmarshal([]byte(conditionOnly), &wd)
	if err != nil {
		t.Error(err)
	}
	p, err := wfWithDefaultConfig(wd)
	if err != nil {
		t.Error(err)
	}
	n, err := p.Current()
	if err != nil {
		t.Error(err, n)
	}
	_, err = p.Next()

	if err != ErrConditionConnectionsMissing {
		t.Error("next should return error!", err)
	}
	n, err = p.Current()

	if err == nil {
		t.Error(err, n)
	}
	//nothing should happen as it reached already the end
	_, err = p.Next()
	if err != ErrConditionConnectionsMissing {
		t.Error("next should return error!", err)
	}
	n, err = p.Current()
	if err != ErrConditionConnectionsMissing {
		t.Error("Current should be nil! ", err, n)
	}

	stack := p.Stack()
	stackSize := len(stack)
	if stackSize != 0 {
		t.Error("wrong stack size!")
		if testVerbose {
			for i, b := range stack {
				log.Println("Stack", i, b.Node.ID, b.Node.HierarchyPath())
			}
		}
	}
	p.Close()
}

func TestOneEmbeddedWFForwardOnly(t *testing.T) {
	wd := &Workflow{}
	err := json.Unmarshal([]byte(oneEmbeddedWF), &wd)
	if err != nil {
		t.Error(err)
	}
	p, err := wfWithDefaultConfig(wd)
	n, err := p.Current()
	if err != nil {
		t.Error(err, n)
	}
	if n.ID != "simple" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}

	hasNext, err := p.Next()
	if hasNext == false || err != nil {
		t.Error("next call wrong result", hasNext, err)
	}
	n, err = p.Current()
	if err != nil {
		t.Error(err, n)
	}
	if n.ID != "frm33" || n.HierarchyPath() != "root.simple" {
		t.Error("Current has an unexpected ID! ", n)
	}
	if p.HasNext() == false {
		t.Error("next can't be false! ")
	}
	hasNext, err = p.Next()
	if hasNext == true || err != nil {
		t.Error(hasNext, err)
	}
	n, err = p.Current()
	if err != nil {
		t.Error(err, n)
	}
	if n.ID != "tmpl1" || n.HierarchyPath() != "root.simple" {
		t.Error("Current has an unexpected ID! ", n)
	}

	//nothing should happen as it reached already the end
	hasNext, err = p.Next()
	if hasNext == true {
		t.Error("next can't be true! ")
	}
	n, err = p.Current()
	if n.ID != "tmpl1" || n.HierarchyPath() != "root.simple" {
		t.Error("Current has an unexpected ID! ", n)
	}

	stack := p.Stack()
	i := 0
	stackSize := len(stack)
	if stackSize != 4 {
		t.Error("wrong stack size!")
		if testVerbose {
			for i, b := range stack {
				log.Println("Stack", i, b.Node.ID, b.Node.HierarchyPath())
			}
		}
	}
	if stack[i].Node.ID != "simple" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm33" || stack[i].Node.HierarchyPath() != "root.simple" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm33" || stack[i].Node.HierarchyPath() != "root.simple" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "tmpl1" || stack[i].Node.HierarchyPath() != "root.simple" {
		t.Error("stack wrong!")
	}
	i++
	p.Close()
}

func TestBackgroundNodesOnlyForwardOnly(t *testing.T) {
	wd := &Workflow{}
	err := json.Unmarshal([]byte(backgroundNodesOnly), &wd)
	if err != nil {
		t.Error(err)
	}
	p, err := wfWithDefaultConfig(wd)
	if err != nil {
		t.Error(err)
	}
	n, err := p.Current()
	if err != nil {
		t.Error(err, n)
	}
	if n == nil {
		t.Error("Current can't be nil! ", n)
	}
	if n.ID != "tmpl1" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}
	if p.HasNext() == false {
		t.Error("next can't be false! ")
	}
	hasNext, err := p.Next()
	if hasNext == true {
		t.Error("next can't be true! ")
	}
	n, err = p.Current()
	if err != nil {
		t.Error(err, n)
	}
	if n.ID != "tmpl3" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}

	//nothing should happen as it reached already the end
	hasNext, err = p.Next()
	if hasNext == true {
		t.Error("next can't be true! ")
	}
	n, err = p.Current()
	if n.ID != "tmpl3" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}

	currentStateSize := len(p.State())
	if currentStateSize != 3 {
		t.Error("wrong currentState size!")
		if testVerbose {
			for i, b := range p.State() {
				log.Println("State", i, b)
			}
		}
	}

	stack := p.Stack()
	i := 0
	stackSize := len(stack)
	if stackSize != 3 {
		t.Error("wrong stack size!")
		if testVerbose {
			for i, b := range stack {
				log.Println("Stack", i, b.Node.ID, b.Node.HierarchyPath())
			}
		}
	}
	if stack[i].Node.ID != "tmpl1" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "tmpl2" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "tmpl3" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	p.Close()
}

func TestComplexWithTwoEmbeddedWorkflowsForwardOnly(t *testing.T) {
	wd := &Workflow{}
	err := json.Unmarshal([]byte(complexWithTwoEmbeddedWorkflows), &wd)
	if err != nil {
		t.Error(err)
	}
	p, err := wfWithDefaultConfig(wd)
	if err != nil {
		t.Error(err)
	}
	n, err := p.Current()
	if err != nil {
		t.Error(err, n)
	}
	if n == nil {
		t.Error("Current can't be nil! ", n)
	}
	if n.ID != "tmpl5" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}
	hasNext, err := p.Next()
	if hasNext == false || err != nil {
		t.Error(hasNext, err)
	}
	n, err = p.Current()
	if err != nil {
		t.Error(err, n)
	}
	if n.ID != "frm5" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}
	hasNext, err = p.Next()
	if hasNext == false {
		t.Error("next can't be false! ")
	}
	n, err = p.Current()
	if err != nil {
		t.Error(err, n)
	}
	if n.ID != "frm33" || n.HierarchyPath() != "root.simple" {
		t.Error("Current has an unexpected ID! ", n)
	}
	hasNext, err = p.Next()
	n, err = p.Current()
	if n.ID != "frm1" || n.HierarchyPath() != "root.complex" {
		t.Error("Current has an unexpected ID! ", n)
	}
	hasNext, err = p.Next()
	n, err = p.Current()
	if n.ID != "tmpl1" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}
	hasNext, err = p.Next()
	if hasNext == true {
		t.Error("next can't be true! ")
	}
	n, err = p.Current()
	if n.ID != "tmpl1" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}

	//nothing should happen as it reached already the end
	hasNext, err = p.Next()
	if hasNext == true {
		t.Error("next can't be true! ")
	}
	n, err = p.Current()
	if n.ID != "tmpl1" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}

	stack := p.Stack()
	i := 0
	stackSize := len(stack)
	if stackSize != 15 {
		t.Error("wrong stack size!")
		if testVerbose {
			for i, b := range stack {
				log.Println("Stack", i, b.Node.ID, b.Node.HierarchyPath())
			}
		}
	}
	if stack[i].Node.ID != "tmpl5" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm5" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm5" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "keq8ej0rc9182zm5vvfd6" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "simple" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm33" || stack[i].Node.HierarchyPath() != "root.simple" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm33" || stack[i].Node.HierarchyPath() != "root.simple" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "tmpl1" || stack[i].Node.HierarchyPath() != "root.simple" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "tmpl3" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "complex" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "gkvl7jheopp97keeex3sm" || stack[i].Node.HierarchyPath() != "root.complex" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "usr1" || stack[i].Node.HierarchyPath() != "root.complex" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm1" || stack[i].Node.HierarchyPath() != "root.complex" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm1" || stack[i].Node.HierarchyPath() != "root.complex" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "tmpl1" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	p.Close()
}

func TestComplexWithTwoEmbeddedWorkflowsForwardBackward(t *testing.T) {
	wd := &Workflow{}
	err := json.Unmarshal([]byte(complexWithTwoEmbeddedWorkflows), &wd)
	if err != nil {
		t.Error(err)
	}
	p, err := wfWithDefaultConfig(wd)
	if err != nil {
		t.Error(err)
	}
	n, err := p.Current()
	if err != nil {
		t.Error(err, n)
	}
	if n == nil {
		t.Error("Current can't be nil! ", n)
	}
	if n.ID != "tmpl5" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}
	forward := 0
	for p.LoopNext() {
		forward++
		if forward >= 50 {
			break
		}
	}

	if forward != 3 {
		t.Error("forward wrong, expected 2", forward)
	}

	backwards := 0
	for p.LoopPrevious(true) {
		backwards++
		if backwards >= 50 {
			break
		}
	}

	// 2 because there are 3 foreground nodes and the last one counts as the target
	if backwards != 2 {
		t.Error("backwards wrong, expected 2", backwards)
	}

	forward = 0
	for p.LoopNext() {
		forward++
		if forward >= 50 {
			t.Error("error called to often", forward)
			break
		}
	}

	backwards = 0
	for p.LoopPrevious(false) {
		backwards++
		if backwards >= 50 {
			t.Error("error called to often", backwards)
			break
		}
	}

	if backwards != 6 {
		t.Error("backwards wrong, expected 6", backwards)
	}

	forward = 0
	for p.LoopNext() {
		forward++
		if forward >= 50 {
			t.Error("error called to often", forward)
			break
		}
	}

	stack := p.Stack()

	i := 0
	stackSize := len(stack)
	if stackSize != 59 {
		t.Error("wrong stack size!")
		if testVerbose {
			for i, b := range p.State() {
				log.Println("State", i, b)
			}
			for i, b := range stack {
				log.Println("Stack", i, b.Node.ID, b.Node.HierarchyPath())
			}
		}
	}
	if stack[i].Node.ID != "tmpl5" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm5" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm5" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "keq8ej0rc9182zm5vvfd6" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "simple" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm33" || stack[i].Node.HierarchyPath() != "root.simple" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm33" || stack[i].Node.HierarchyPath() != "root.simple" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "tmpl1" || stack[i].Node.HierarchyPath() != "root.simple" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "tmpl3" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "complex" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "gkvl7jheopp97keeex3sm" || stack[i].Node.HierarchyPath() != "root.complex" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "usr1" || stack[i].Node.HierarchyPath() != "root.complex" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm1" || stack[i].Node.HierarchyPath() != "root.complex" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "frm1" || stack[i].Node.HierarchyPath() != "root.complex" {
		t.Error("stack wrong!")
	}
	i++
	if stack[i].Node.ID != "tmpl1" || stack[i].Node.HierarchyPath() != "root" {
		t.Error("stack wrong!")
	}
	i++
	p.Close()
}

func TestComplexWithTwoEmbeddedWorkflowsForwardBackwardButNotFully(t *testing.T) {
	wd := &Workflow{}
	err := json.Unmarshal([]byte(complexWithTwoEmbeddedWorkflows), &wd)
	if err != nil {
		t.Error(err)
	}
	p, err := wfWithDefaultConfig(wd)
	if err != nil {
		t.Error(err)
	}
	n, err := p.Current()
	if err != nil {
		t.Error(err, n)
	}
	if n == nil {
		t.Error("Current can't be nil! ", n)
	}
	if n.ID != "tmpl5" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}
	forward := 0
	for p.LoopNext() {
		forward++
		if forward == 2 {
			break
		}
	}
	if forward != 2 {
		t.Error("expected 2", forward)
	}

	hasPrev, err := p.Previous(true)
	if err != nil || hasPrev == false {
		t.Error(err, hasPrev)
	}

	n, err = p.Current()
	if err != nil || n == nil {
		t.Error(err, n)
	}

	if n.ID != "frm5" || n.HierarchyPath() != "root" {
		t.Error("Current has an unexpected ID! ", n)
	}

	var hasNext bool
	hasNext, err = p.Next()
	if err != nil || hasNext == false {
		t.Error(err, hasNext)
	}
	n, err = p.Current()
	if err != nil || n == nil {
		t.Error(err, n)
	}
	if n.ID != "frm33" || n.HierarchyPath() != "root.simple" {
		t.Error("Current has an unexpected ID! ", n)
	}
	//presenting frm33

	forward = 1
	for p.LoopNext() {
		forward++
		if forward >= 50 {
			break
		}
	}
	// 2 because there are 3 foreground nodes and the last one counts as the target
	if forward != 2 {
		t.Error("forward wrong, expected 2", forward)
	}

	stack := p.State()

	i := 0
	stackSize := len(stack)
	if stackSize != 8 {
		t.Error("wrong State size!")
		if testVerbose {
			for i, b := range p.State() {
				log.Println("State", i, b)
			}
			for i, b := range stack {
				log.Println("Stack", i, b)
			}
		}
	}

	if stack[i].NodeID != "tmpl5" || stack[i].HierarchyPath != "root" {
		t.Error("stack wrong!", stack[i])
	}
	i++
	if stack[i].NodeID != "frm5" || stack[i].HierarchyPath != "root" {
		t.Error("stack wrong!", stack[i])
	}
	i++
	if stack[i].NodeID != "frm33" || stack[i].HierarchyPath != "root.simple" {
		t.Error("stack wrong!", stack[i])
	}
	i++
	if stack[i].NodeID != "tmpl1" || stack[i].HierarchyPath != "root.simple" {
		t.Error("stack wrong!", stack[i])
	}
	i++
	if stack[i].NodeID != "tmpl3" || stack[i].HierarchyPath != "root" {
		t.Error("stack wrong!", stack[i])
	}
	i++
	if stack[i].NodeID != "usr1" || stack[i].HierarchyPath != "root.complex" {
		t.Error("stack wrong!", stack[i])
	}
	i++
	if stack[i].NodeID != "frm1" || stack[i].HierarchyPath != "root.complex" {
		t.Error("stack wrong!", stack[i])
	}
	i++
	if stack[i].NodeID != "tmpl1" || stack[i].HierarchyPath != "root" {
		t.Error("stack wrong!", stack[i])
	}
	i++
	p.Close()
}

var complexWithTwoEmbeddedWorkflows = `
{
        "flow": {
            "start": {
                "node": "tmpl5"
            },
            "nodes": {
                "complex": {
                    "id": "complex",
                    "name": "Complex Name",
                    "type": "workflow",
                    "detail": "More about complex..",
					"conns": [{"id": "tmpl1"}]
                },
                "frm5": {
                    "id": "frm5",
                    "name": "frm5",
                    "type": "form",
                    "detail": "bla bla5",
                    "conns": [{"id": "keq8ej0rc9182zm5vvfd6"}]
                },
                "keq8ej0rc9182zm5vvfd6": {
                    "id": "keq8ej0rc9182zm5vvfd6",
                    "name": "condition",
                    "type": "condition",
                    "data": {
                        "js": "\nfunction condition(){\n  if( input[\"someVar\"] == \"someValue\" ){\n    /*must match the value in the cases*/\n    return \"someValue\";\n  }else{\n    /*must match the value in the cases*/\n    return \"something else\";\n  }\n}\n                                        "
                    },
                    "cases": [{"name": "someValue","value": "someValue"},{"name": "something else","value": "something else"}],
                    "conns": [{"id": "tmpl1","value": "something else"},{"id": "simple","value": "someValue"}]
                },
                "simple": {
                    "id": "simple",
                    "name": "Simple Name",
                    "type": "workflow",
                    "detail": "More about simple..",
                    "conns": [{"id": "tmpl3"}]
                },
                "tmpl1": {
                    "id": "tmpl1",
                    "name": "tmpl1",
                    "type": "template",
                    "detail": "bla bla1"
                },
                "tmpl3": {
                    "id": "tmpl3",
                    "name": "tmpl3",
                    "type": "template",
                    "detail": "bla bla3",
                    "conns": [{"id": "complex"}]
				},
                "tmpl5": {
                    "id": "tmpl5",
                    "name": "tmpl5",
                    "type": "template",
                    "detail": "bla bla5",
                    "conns": [{"id": "frm5"}]
                }
            }
        },
        "progressFlow": []
    }
`

var complexWF = `
{
        "flow": {
            "start": {
                "node": "gkvl7jheopp97keeex3sm"
            },
            "nodes": {
                "frm1": {
                    "id": "frm1",
                    "name": "My Form Name",
                    "type": "form",
                    "detail": "bla bla1"
                },
                "gkvl7jheopp97keeex3sm": {
                    "id": "gkvl7jheopp97keeex3sm",
                    "name": "condition",
                    "type": "condition",
                    "data": {
                        "js": "\nfunction condition(){\n  if( input[\"someVar\"] == \"someValue\" ){\n    /*must match the value in the cases*/\n    return true;\n  }else{\n    /*must match the value in the cases*/\n    return;\n  }\n}\n                                        "
                    },
                    "cases": [{"name": "Some value explanation","value": true},{"name": "Something else explanation","value": false}],
                    "conns": [{"id": "usr1","value": true},{"id": "tmpl11","value": false}]
                },
                "tmpl11": {
                    "id": "tmpl11",
                    "name": "tmpl11",
                    "type": "template",
                    "detail": "bla bla1"
                },
                "usr1": {
                    "id": "usr1",
                    "name": "usr1",
                    "type": "user",
                    "detail": "bla bla1",
					"conns": [{"id": "frm1"}]
                }
            }
        },
        "progressFlow": []
    }
`

var simpleWF = `

{
        "flow": {
            "start": {
                "node": "frm33"
            },
            "nodes": {
                "frm33": {
                    "id": "frm33",
                    "name": "My Form Name",
                    "type": "form",
                    "detail": "bla bla1",
                    "conns": [{"id": "tmpl1"}]
                },
                "tmpl1": {
                    "id": "tmpl1",
                    "name": "tmpl1",
                    "type": "template",
                    "detail": "bla bla1"
                }
            }
        },
        "progressFlow": []
    }
`
var oneBackgroundNode = `
{
        "flow": {
            "start": {
                "node": "tmpl1"
            },
            "nodes": {
                "tmpl1": {
                    "id": "tmpl1",
                    "name": "tmpl1",
                    "type": "template",
                    "detail": "bla bla1"
                }
            }
        },
        "progressFlow": []
    }
`
var backgroundNodesOnly = `
{
        "flow": {
            "start": {
                "node": "tmpl1"
            },
            "nodes": {
                "tmpl1": {
                    "id": "tmpl1",
                    "name": "tmpl1",
                    "type": "template",
                    "detail": "bla bla1",
					"conns": [{"id": "tmpl2"}]
                },
                "tmpl2": {
                    "id": "tmpl2",
                    "name": "tmpl2",
                    "type": "template",
                    "detail": "bla bla1",
					"conns": [{"id": "tmpl3"}]
                },
                "tmpl3": {
                    "id": "tmpl3",
                    "name": "tmpl3",
                    "type": "template",
                    "detail": "bla bla1"
                }
            }
        },
        "progressFlow": []
    }
`
var conditionOnly = `
{
	"flow": {
		"start": {
			"node": "cond1"
		},
		"nodes": {
			"cond1": {
				"id": "cond1",
				"name": "condition",
				"type": "condition",
				"data": {
					"js": "\nfunction condition(){\n  if( input[\"someVar\"] == \"someValue\" ){\n    /*must match the value in the cases*/\n    return true;\n  }else{\n    /*must match the value in the cases*/\n    return;\n  }\n}\n                                        "
				},
				"cases": [{"name": "Some value explanation","value": true},{"name": "Something else explanation","value": false}],
				"conns": []
			}
		}
	},
	"progressFlow": []
}
`
var conditionOnlyWithConnections = `
{
	"flow": {
		"start": {
			"node": "cond1"
		},
		"nodes": {
			"cond1": {
				"id": "cond1",
				"name": "condition",
				"type": "condition",
				"data": {
					"js": "\nfunction condition(){\n  if( input[\"someVar\"] == \"someValue\" ){\n    /*must match the value in the cases*/\n    return true;\n  }else{\n    /*must match the value in the cases*/\n    return false;\n  }\n}\n                                        "
				},
				"cases": [{"name": "Some value explanation","value": true},{"name": "Something else explanation","value": false}],
				"conns": [{"id": "form1","value": true},{"id": "form2","value": false}]
			},
			"form1": {
				"id": "form1",
				"name": "My form1 Name",
				"type": "form"
			},
			"form2": {
				"id": "form2",
				"name": "My form2 Name",
				"type": "form"
			}
		}
	},
	"progressFlow": []
}
`

var oneEmbeddedWF = `
{
	"flow": {
		"start": {
			"node": "simple"
		},
		"nodes": {
			"simple": {
				"id": "simple",
				"name": "Simple Name",
				"type": "workflow",
				"detail": "More about simple.."
			}
		}
	},
	"progressFlow": []
}
`
var deepWF = `
{
	"flow": {
		"start": {
			"node": "frm.1"
		},
		"nodes": {
			"frm.1": {
				"id": "frm.1",
				"name": "My Form Name",
				"type": "form",
				"detail": "bla bla1",
				"conns": [{"id": "deep2"}]
			},
			"deep2": {
				"id": "deep2",
				"name": "Deep2 Name",
				"type": "workflow",
				"detail": "More about deep2.."
			}
		}
	},
	"progressFlow": []
}
`
var deep2WF = `
{
	"flow": {
		"start": {
			"node": "frm.2"
		},
		"nodes": {
			"frm.2": {
				"id": "frm.2",
				"name": "My Form Name",
				"type": "form",
				"detail": "bla bla1",
				"conns": [{"id": "deep3"}]
			},
			"deep3": {
				"id": "deep3",
				"name": "Deep3 Name",
				"type": "workflow",
				"detail": "More about deep3..",
				"conns": [{"id": "frm.2-2"}]
			},
			"frm.2-2": {
				"id": "frm.2-2",
				"name": "My Form Name",
				"type": "form",
				"detail": "bla bla1"
			}
		}
	},
	"progressFlow": []
}
`
var deep3WF = `
{
	"flow": {
		"start": {
			"node": "frm.3"
		},
		"nodes": {
			"frm.3": {
				"id": "frm.3",
				"name": "My Form Name",
				"type": "form",
				"detail": "bla bla1",
				"conns": [{"id": "tmpl.3"}]
			},
			"tmpl.3": {
				"id": "tmpl.3",
				"name": "My Templ Name",
				"type": "template",
				"detail": "bla bla1"
			}
		}
	},
	"progressFlow": []
}
`
var deep4WF = `
{
	"flow": {
		"start": {
			"node": "frm.4"
		},
		"nodes": {
			"frm.4": {
				"id": "frm.4",
				"name": "My Form Name",
				"type": "form",
				"detail": "bla bla1"
			}
		}
	},
	"progressFlow": []
}
`
