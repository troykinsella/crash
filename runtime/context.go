package runtime

import (
	"github.com/troykinsella/crash"
	"github.com/troykinsella/crash/logging"
)

type Context struct {
	log     *logging.Logger
	vars    Variables
}

func newTestContext(options *crash.TestOptions, config *crash.Config) (*Context, error) {
	vars := newVariables()

	for k, v := range config.Variables {
		vars.Set(k, v)
	}

	for k, v := range options.Variables {
		vars.Set(k, v)
	}

	log := logging.NewLogger(options.LogLevel,
		options.Colorize,
		options.LogJson)

	return &Context{
		log: log,
		vars: vars,
	}, nil
}

func (ctx *Context) NewChild() *Context {
	return &Context{
		log: ctx.log,
		vars: ctx.vars.NewChild(),
	}
}

func (ctx *Context) Commit() {
	ctx.vars = ctx.vars.Commit()
}
