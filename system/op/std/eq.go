package std

import (
	"github.com/troykinsella/crash/system/op/support"
	"github.com/troykinsella/crash/system/data"
)

type EqOp struct {}

func (o *EqOp) Exec(subjects []interface{}, args []interface{}) (bool, interface{}, error) {
	subject, err := support.OneSubject("eq", subjects)
	if err != nil {
		return false, nil, err
	}

	match := false
	for _, arg := range args {
		m, err := data.LooseEqual(subject, arg)
		if err != nil {
			return false, nil, err
		}

		if m {
			match = true
			break
		}
	}

	return match, match, nil
}
