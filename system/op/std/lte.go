package std

import (
	"github.com/troykinsella/crash/system/op/support"
	"github.com/troykinsella/crash/system/data"
)

type LteOp struct{}

func (o *LteOp) Exec(subjects []interface{}, args []interface{}) (bool, interface{}, error) {
	subject, err := support.OneSubject("lte", subjects)
	if err != nil {
		return false, nil, err
	}

	arg, err := support.OneArg("lte", args)
	if err != nil {
		return false, nil, err
	}

	c, err := data.Compare(subject, arg)
	if err != nil {
		return false, nil, err
	}

	return c <= 0, nil, nil
}
