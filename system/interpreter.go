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
		ch, err := NewScript(str)
		if err != nil {
			return nil, err
		}
		op := NewOpAdapter(result, ch)
		reg.Register(name, op)
	}

	return result, nil
}

func (interp *Interpreter) Statement(s *Script, vars util.Values) (bool, interface{}, string, error) {

	ctx := exec.NewContext(vars, interp.ops)

	ok, result, msg, err := s.n.Exec(ctx)
	if err != nil {
		return false, nil, "", err
	}
	if msg == "" {
		msg = s.str
	}
	return ok, result, msg, nil
}
