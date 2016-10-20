package ast

import (
	"fmt"
	"github.com/troykinsella/crash/system/exec"
	"strings"
	"bytes"
	"github.com/troykinsella/crash/system/data"
)

type IString struct {
	Str string
	Expr *Expression
	Ident *Identifier
	Next *IString
}

func (ast *IString) Exec(ctx *exec.Context) (bool, string, error) {
	if ast.Expr == nil && ast.Ident == nil {
		return true, ast.Str, nil
	}

	rok := true
	var evaled interface{}
	var buf bytes.Buffer
	buf.WriteString(ast.Str)

	if ast.Expr != nil {
		ok, v, err := ast.Expr.Exec(ctx)
		if err != nil {
			return false, "", err
		}
		evaled = v
		rok = rok && ok
	} else if ast.Ident != nil {
		ok, v, err := ast.Ident.Exec(ctx)
		if err != nil {
			return false, "", err
		}
		evaled = v
		rok = rok && ok
	}

	rstr, err := data.ToString(evaled)
	if err != nil {
		return false, "", err
	}
	buf.WriteString(rstr)

	if ast.Next != nil {
		ok, s, err := ast.Next.Exec(ctx)
		if err != nil {
			return false, "", err
		}
		rok = rok && ok
		buf.WriteString(s)
	}

	return rok, buf.String(), nil
}

func (ast *IString) Dump(indent int) {
	fmt.Printf("%sistring:\n", strings.Repeat(" ", indent))
	fmt.Printf("%sstring: '%s'\n", strings.Repeat(" ", indent + 2), ast.Str)

	if ast.Expr != nil {
		ast.Expr.Dump(indent + 2)
	}
	if ast.Next != nil {
		ast.Next.Dump(indent + 2)
	}
}
