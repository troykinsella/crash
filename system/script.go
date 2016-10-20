package system

import (
	"strings"
	"bufio"
	"github.com/troykinsella/crash/system/ast"
	"github.com/troykinsella/crash/system/parser"
	"errors"
)

type ScriptType uint8

const (
	STMT ScriptType = iota
	EXPR
)

type Script struct {
	str string
	stmt   *ast.Statement
	expr   *ast.Expression
}

type ScriptList []*Script

type OpAdapter struct {
	interp *Interpreter
	script *Script
}

func NewOpAdapter(interp *Interpreter, script *Script) *OpAdapter {
	return &OpAdapter{
		interp: interp,
		script: script,
	}
}

// implements Op
func (o *OpAdapter) Exec(subjects []interface{}, args []interface{}) (bool, interface{}, error) {
	return false, nil, nil
}

func NewScriptList(expressions []string, t ScriptType) (*ScriptList, error) {
	if expressions == nil {
		return nil, nil
	}
	scripts := make(ScriptList, len(expressions))
	for i, e := range expressions {
		ch, err := NewScript(e, t)
		if err != nil {
			return nil, err
		}
		if ch != nil {
			scripts[i] = ch
		}
	}
	return &scripts, nil
}

func NewScript(script string, t ScriptType) (*Script, error) {
	if script == "" {
		return nil, nil
	}

	buf := bufio.NewReader(strings.NewReader(script))
	parser := parser.New(buf)

	if t == STMT {
		n, err := parser.Statement()
		if err != nil {
			return nil, err
		}
		return &Script{
			str: script,
			stmt: n,
		}, nil

	} else if t == EXPR {
		n, err := parser.Expression()
		if err != nil {
			return nil, err
		}
		return &Script{
			str: script,
			expr: n,
		}, nil
	}

	return nil, errors.New("internal error")
}

