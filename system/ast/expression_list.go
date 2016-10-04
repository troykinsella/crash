package ast

import (
	"fmt"
	"github.com/troykinsella/crash/system/exec"
	"strings"
)

type ExpressionList struct {
	Expressions []*Expression
}

func (ast *ExpressionList) Exec(ctx *exec.Context) (bool, interface{}, error) {
	expressions := make([]interface{}, len(ast.Expressions))

	rok := true

	for i, expr := range ast.Expressions {
		ok, r, err := expr.Exec(ctx)
		if err != nil {
			return false, nil, err
		}
		rok = rok && ok
		expressions[i] = r
	}

	return rok, expressions, nil
}

func (ast *ExpressionList) Dump(indent int) {
	fmt.Printf("%sexpression_list:\n", strings.Repeat(" ", indent))
	for _, subject := range ast.Expressions {
		subject.Dump(indent)
	}
}
