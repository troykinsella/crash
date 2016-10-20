package system

import (
	"github.com/troykinsella/crash/util"
	"github.com/troykinsella/crash/logging"
	"github.com/troykinsella/crash/system/op"
	"github.com/troykinsella/crash"
	"github.com/troykinsella/crash/system/exec"
)

type Interpreter struct {
	config *crash.Config
	log    *logging.Logger
	ops    *op.Registry
}

func NewInterpreter(log *logging.Logger,
         scripts map[string]string) (*Interpreter, error) {
	reg := op.NewRegistry()

	result := &Interpreter{
		log: log,
		ops: reg,
	}

	for name, str := range scripts {
		ch, err := NewScript(str, STMT)
		if err != nil {
			return nil, err
		}
		op := NewOpAdapter(result, ch)
		reg.Register(name, op)
	}

	return result, nil
}

func (interp *Interpreter) Run(s *Script, vars util.Values) (ok bool, result interface{}, msg string, err error) {

	ctx := exec.NewContext(vars, interp.ops)

	if s.expr != nil {
		ok, result, err = s.expr.Exec(ctx)
	} else if s.stmt != nil {
		ok, result, msg, err = s.stmt.Exec(ctx)
	}

	if msg == "" {
		msg = s.str
	}
	return
}

func (interp *Interpreter) RunAll(sl *ScriptList, vars util.Values) (bool, []interface{}, []string, error) {
	rok := true
	n := len(*sl)
	results := make([]interface{}, n)
	msgs := make([]string, n)

	for i, s := range *sl {
		ok, result, msg, err := interp.Run(s, vars)
		if err != nil {
			return false, nil, nil, err
		}

		rok = rok && ok
		results[i] = result
		msgs[i] = msg
	}

	return rok, results, msgs, nil
}
