package crash

import "github.com/troykinsella/crash/logging"

type TestOptions struct {
	Crashfile      string
	LogLevel       logging.Level
	Colorize       bool
	LogJson        bool
	Variables      map[string]string
	VariablesFiles []string
}
