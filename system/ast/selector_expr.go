package ast

import (
	"fmt"
	"github.com/troykinsella/crash/system/exec"
	"strings"
)

type SelectorExpr struct {
	Operand *PrimaryExpr
	Ident   string
}

func (ast *SelectorExpr) Exec(ctx *exec.Context) (bool, interface{}, error) {
	ok, operand, err := ast.Operand.Exec(ctx)
	if err != nil {
		return false, operand, err
	}

	val, err := extractValue(operand, ast.Ident)
	if err != nil {
		return false, nil, err
	}

	return ok, val, nil
}

func (ast *SelectorExpr) Dump(indent int) {
	fmt.Printf("%s- selector:\n", strings.Repeat(" ", indent))
	fmt.Printf("%sidentifier: %s\n", strings.Repeat(" ", indent + 4), ast.Ident)
}
