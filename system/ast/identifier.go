package ast

import (
	"fmt"
	"github.com/troykinsella/crash/system/exec"
	"strings"
)

type Identifier struct {
	Name string
}

func (ast *Identifier) Exec(ctx *exec.Context) (bool, interface{}, error) {
	result := ctx.Vars.Get(ast.Name)
	if result == nil {
		return false, nil, fmt.Errorf("Not found: %s", ast.Name)
	}
	return true, result, nil
}

func (ast *Identifier) Dump(indent int) {
	fmt.Printf("%sidentifier: %s\n", strings.Repeat(" ", indent + 2), ast.Name)
}
