package ast

import (
	"fmt"
	"github.com/troykinsella/crash/system/exec"
	"strings"
)

type Statement struct {
	Operation *Operation
	Message   *IString
}

func (ast *Statement) Exec(ctx *exec.Context) (bool, interface{}, string, error) {

	//fmt.Printf("AST:\n")
	//s.n.dump(0)

	ok, result, err := ast.Operation.Exec(ctx)
	if err != nil {
		return false, nil, "", err
	}

	var msg string
	if ast.Message != nil {
		_, msg, err = ast.Message.Exec(ctx)
		if err != nil {
			return false, nil, "", err
		}
	}

	return ok, result, msg, nil
}

func (ast *Statement) Dump(indent int) {
	fmt.Printf("%s- statement:\n", strings.Repeat(" ", indent))
	ast.Operation.Dump(indent + 4)
	if ast.Message != nil {
		ast.Message.Dump(indent + 4)
	}
}
