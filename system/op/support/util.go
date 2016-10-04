package support

import "fmt"

func OneSubject(op string, subjects []interface{}) (interface{}, error) {
	if len(subjects) != 1 {
		return nil, fmt.Errorf("'%s' allows only one subject, found %d", op, len(subjects))
	}
	return subjects[0], nil
}

func OneArg(op string, args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("'%s' accepts only one argument, found %d", op, len(args))
	}
	return args[0], nil
}

func ArgCount(op string, args []interface{}, count int) error {
	if len(args) != count {
		return fmt.Errorf("'%s' accepts %d arguments, found %d", op, count, len(args))
	}
	return nil
}
