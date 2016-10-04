package std

import (
	"github.com/troykinsella/crash/system/op/support"
	"github.com/troykinsella/crash/system/data"
)

type InOp struct {}

func (o *InOp) Exec(subjects []interface{}, args []interface{}) (bool, interface{}, error) {
	subject, err := support.OneSubject("in", subjects)
	if err != nil {
		return false, nil, err
	}

	err = support.ArgCount("in", args, 2)
	if err != nil {
		return false, nil, err
	}

	return o.inRange(subject, args[0], args[1])
}

func (o *InOp) inRange(subject interface{}, min interface{}, max interface{}) (bool, interface{}, error) {
	lc, err := data.Compare(subject, min)
	if err != nil {
		return false, nil, err
	}
	if lc < 0 {
		return false, nil, nil
	}

	rc, err := data.Compare(subject, max)
	if err != nil {
		return false, nil, err
	}
	if rc > 0 {
		return false, nil, nil
	}
	return true, true, nil
}
