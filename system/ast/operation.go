package ast

import (
	"fmt"
	"github.com/troykinsella/crash/system/exec"
	"strings"
)

type Operation struct {
	Subjects *ExpressionList
	Negate bool
	Operation string
	Arguments *ExpressionList
}

var zeroArgs = make([]interface{}, 0)

func (ast *Operation) Exec(ctx *exec.Context) (bool, interface{}, error) {

	rok := true

	ok, subjects, err := ast.Subjects.Exec(ctx)
	if err != nil {
		return false, nil, err
	}
	rok = rok && ok

	op, err := ctx.Ops.Get(ast.Operation)
	if err != nil {
		return false, nil, err
	}

	var args interface{}
	if ast.Arguments != nil {
		ok, args, err = ast.Arguments.Exec(ctx)
		if err != nil {
			return false, nil, err
		}
		rok = rok && ok
	} else {
		args = zeroArgs
	}

	ok, result, err := op.Exec(subjects.([]interface{}), args.([]interface{}))
	if err != nil {
		return false, nil, err
	}
	rok = rok && ok

	if ast.Negate {
		rok = !rok
	}

	return rok, result, nil
}

func (ast *Operation) Dump(indent int) {
	fmt.Printf("%soperation:\n", strings.Repeat(" ", indent))
	ast.Subjects.Dump(indent + 2)
	fmt.Printf("%sfunction: %s\n", strings.Repeat(" ", indent + 2), ast.Operation)
	if ast.Arguments != nil {
		ast.Arguments.Dump(indent + 2)
	}
}
