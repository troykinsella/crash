package std

import "os/exec"

type ShOp struct {}

func (o *ShOp) Exec(subjects []interface{}, args []interface{}) (bool, interface{}, error) {

	cmd := exec.Command(args[0].(string))

	out, err := cmd.Output()
	if err != nil {
		return false, nil, err
	}

	return true, string(out), nil
}
