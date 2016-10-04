package ast

import (
	"fmt"
	"github.com/troykinsella/crash/system/exec"
	"github.com/troykinsella/crash/system/token"
	"strings"
)

type Literal struct {
	Type token.Token
	Str string
	Int int64
}

func (ast *Literal) Exec(ctx *exec.Context) (bool, interface{}, error) {
	switch ast.Type {
	case token.STRING: return true, ast.Str, nil
	case token.NUMBER: return true, ast.Int, nil
	case token.TRUE, token.FALSE: return true, ast.Type == token.TRUE, nil
	}
	return false, nil, fmt.Errorf("Internal error: invalid token: %s", ast.Type.String())
}

func (ast *Literal) Dump(indent int) {
	fmt.Printf("%sliteral:\n", strings.Repeat(" ", indent))
	switch ast.Type {
	case token.STRING: fmt.Printf("%sstring: %s\n", strings.Repeat(" ", indent + 2), ast.Str)
	case token.NUMBER: fmt.Printf("%snumber: %d\n", strings.Repeat(" ", indent + 2), ast.Int)
	case token.TRUE, token.FALSE: fmt.Printf("%sboolean: %s\n", strings.Repeat(" ", indent + 2), ast.Type == token.TRUE)
	}
}
