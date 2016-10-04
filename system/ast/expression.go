package ast

import (
	"fmt"
	"github.com/troykinsella/crash/system/exec"
	"strings"
)

type Expression struct {
	PExpr   *PrimaryExpr
}

func (ast *Expression) Exec(ctx *exec.Context) (bool, interface{}, error) {
	ok, result, err := ast.PExpr.Exec(ctx)
	return ok, result, err
}

func (ast *Expression) Dump(indent int) {
	fmt.Printf("%s- expression:\n", strings.Repeat(" ", indent))
	ast.PExpr.Dump(indent + 2)
}
