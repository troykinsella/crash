package runtime

import (
	"github.com/troykinsella/crash/util"
)

type Variables interface {
	Names() []string
	Get(string) interface{}
	Set(string, interface{})
	GetString(name string) string
	SetAll(map[string]interface{})
	AsMap() map[string]interface{}
	NewChild() Variables
	Commit() Variables
}

type varNode struct {
	parent *varNode
	vars   map[string]interface{}
}

func newVariables() *varNode {
	return &varNode{
		vars: make(map[string]interface{}),
	}
}

func (v *varNode) Names() []string {
	return util.KeysForMapStringInterface(v.AsMap())
}

func (v *varNode) Get(name string) interface{} {
	result := v.vars[name]
	if result == nil && v.parent != nil {
		result = v.parent.Get(name)
	}
	return result
}

func (v *varNode) GetString(name string) string {
	return util.ToString(v.Get(name))
}

func (v *varNode) Set(name string, value interface{}) {
	v.vars[name] = value
}

func (v *varNode) SetAll(values map[string]interface{}) {
	for key, val := range values {
		v.vars[key] = val
	}
}

func (v *varNode) AsMap() (map[string]interface{}) {
	vars := make(map[string]interface{})

	stackSize := 0
	varCount := 0
	n := v
	for n != nil {
		stackSize++
		varCount += len(n.vars)
		n = n.parent
	}

	stack := make([]*varNode, stackSize)

	i := stackSize - 1
	n = v
	for n != nil {
		stack[i] = n
		n = n.parent
		i--
	}

	i = 0
	for i < stackSize {
		util.PutAllForMapStringInterface(vars, stack[i].vars)
		i++
	}

	return vars
}

func (v *varNode) NewChild() Variables {
	return &varNode{
		parent: v,
		vars: make(map[string]interface{}),
	}
}

func (v *varNode) Commit() Variables {
	v.parent.SetAll(v.vars)

	return v.parent
}
