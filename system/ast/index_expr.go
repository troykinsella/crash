package ast

import (
	"fmt"
	"github.com/troykinsella/crash/system/exec"
	"strings"
)

type IndexExpr struct {
	Operand *PrimaryExpr
	Index *PrimaryExpr
}

func (ast *IndexExpr) Exec(ctx *exec.Context) (bool, interface{}, error) {
	ok, operand, err := ast.Operand.Exec(ctx)
	if err != nil {
		return false, nil, err
	}
	rok := ok

	ok, index, err := ast.Index.Exec(ctx)
	if err != nil {
		return false, nil, err
	}
	rok = rok && ok

	val, err := extractValue(operand, index)
	if err != nil {
		return false, nil, err
	}

	return rok, val, nil
}

func (ast *IndexExpr) Dump(indent int) {
	fmt.Printf("%s- index_expr:\n", strings.Repeat(" ", indent))
	ast.Operand.Dump(indent + 4)
	ast.Index.Dump(indent + 4)
}
