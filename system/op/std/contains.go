package std

import (
	"strings"
	"github.com/troykinsella/crash/system/op/support"
	"github.com/troykinsella/crash/system/data"
)

type ContainsOp struct {}

func (o *ContainsOp) Exec(subjects []interface{}, args []interface{}) (bool, interface{}, error) {
	subject, err := support.OneSubject("contains", subjects)
	if err != nil {
		return false, nil, err
	}

	err = support.ArgCount("contains", args, 1)
	if err != nil {
		return false, nil, err
	}

	return o.contains(subject, args[0])
}

func (o *ContainsOp) contains(subject interface{}, target interface{}) (bool, interface{}, error) {
	switch t := target.(type) {
	case string:
		ss, err := data.ToString(subject)
		if err != nil {
			return false, nil, err
		}
		return strings.Contains(ss, t), nil, nil
	}
	return false, false, nil
}
