package ast

import (
	"fmt"
	"github.com/troykinsella/crash/system/exec"
	"strings"
)

type PrimaryExpr struct {
	Literal  *Literal
	Ident    *Identifier

	Selector *SelectorExpr

	Index    *IndexExpr
}

func (ast *PrimaryExpr) Exec(ctx *exec.Context) (bool, interface{}, error) {
	if ast.Literal != nil {
		return ast.Literal.Exec(ctx)
	}
	if ast.Ident != nil {
		return ast.Ident.Exec(ctx)
	}
	if ast.Selector != nil {
		return ast.Selector.Exec(ctx)
	}
	if ast.Index != nil {
		return ast.Index.Exec(ctx)
	}

	return false, nil, fmt.Errorf("Internal error")
}

func (ast *PrimaryExpr) Dump(indent int) {
	fmt.Printf("%s- primary_expr:\n", strings.Repeat(" ", indent))
	if ast.Literal != nil {
		ast.Literal.Dump(indent + 4)
	}
	if ast.Ident != nil {
		ast.Ident.Dump(indent + 4)
	}
	if ast.Selector != nil {
		ast.Selector.Dump(indent + 4)
	}
	if ast.Index != nil {
		ast.Index.Dump(indent + 4)
	}
}
