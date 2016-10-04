package op

import (
	"github.com/troykinsella/crash/system/op/std"
	"fmt"
)

type Op interface {
	Exec(subjects []interface{}, args []interface{}) (bool, interface{}, error)
}

type Registry struct {
	registry map[string]Op
}

func NewRegistry() *Registry {
	return &Registry{
		registry: make(map[string]Op),
	}
}

func (or *Registry) Get(name string) (Op, error) {
	switch name {
	case "contains": return &std.ContainsOp{}, nil
	case "eq":       return &std.EqOp{}, nil
	case "gt":       return &std.GtOp{}, nil
	case "gte":      return &std.GteOp{}, nil
	case "lt":       return &std.LtOp{}, nil
	case "lte":      return &std.LteOp{}, nil
	case "in":       return &std.InOp{}, nil
	case "sh":       return &std.ShOp{}, nil
	}

	r := or.registry[name]
	if r == nil {
		return nil, fmt.Errorf("Operation not found: %s", name)
	}

	return r, nil
}

func (or *Registry) Register(name string, op Op) {
	or.registry[name] = op
}
