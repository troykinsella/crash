package action

import (
	"github.com/troykinsella/crash/util"
	"github.com/troykinsella/crash/logging"
)

type Action interface {
	Run() (*Result, error)
}

type ActionConfig struct {
	Name   string
	Params util.Values
	Log    *logging.Logger
}

type Result struct {
	Data    map[string]interface{}
}

func NewAction(config *ActionConfig) Action {
	switch config.Name {
	case "http":
		return NewHttp(config)
	}
	return nil
}
