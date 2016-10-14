package action

import (
	"os/exec"
	"fmt"
	"github.com/troykinsella/crash/logging"
	"bytes"
	"strings"
)

type Shell struct {
	config *ActionConfig
}

func (s *Shell) Run() (*Result, error) {

	str := s.config.Params.GetString("command")
	if str == "" {
		return nil, fmt.Errorf("command parameter required")
	}
	s.config.Log.Start(logging.INFO, str)

	cmd := exec.Command("sh", "-c", str)
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf

	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	r := s.genResult(outBuf, errBuf)
	return r, nil
}

func (s *Shell) genResult(outBuf bytes.Buffer, errBuf bytes.Buffer) *Result {
	data := make(map[string]interface{})
	data["out"] = strings.TrimSpace(outBuf.String())
	data["err"] = strings.TrimSpace(errBuf.String())

	return &Result{
		Data:    data,
	}
}

func NewShell(config *ActionConfig) *Shell {
	return &Shell{
		config: config,
	}
}
