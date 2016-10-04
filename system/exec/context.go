package exec

import (
	"github.com/troykinsella/crash/util"
	"github.com/troykinsella/crash/system/op"
)

type Context struct {
	Vars util.Values
	Ops  *op.Registry
}

func NewContext(vars util.Values, ops *op.Registry) *Context {
	return &Context{
		Vars: vars,
		Ops: ops,
	}
}
