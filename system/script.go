package system

import (
	"strings"
	"bufio"
	"errors"
	"github.com/troykinsella/crash/system/ast"
	"github.com/troykinsella/crash/system/parser"
)

type Script struct {
	str string
	n   *ast.Statement
}

type ScriptSet []*Script

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

func NewScriptSet(expressions []string) (*ScriptSet, error) {
	if expressions == nil {
		return nil, nil
	}
	scripts := make(ScriptSet, len(expressions))
	for i, e := range expressions {
		ch, err := NewScript(e)
		if err != nil {
			return nil, err
		}
		scripts[i] = ch
	}
	return &scripts, nil
}

func NewScript(script string) (*Script, error) {
	if script == "" {
		return nil, errors.New("nil script")
	}

	buf := bufio.NewReader(strings.NewReader(script))
	parser := parser.New(buf)

	n, err := parser.Statement()
	if err != nil {
		return nil, err
	}

	return &Script{
		str: script,
		n: n,
	}, nil
}

