package runtime

import (
	"github.com/troykinsella/crash"
)

type TestRunner struct {
	config  *crash.Config
	context *Context
}

func NewTestRunner(options *crash.TestOptions, config *crash.Config) (*TestRunner, error) {
	context, err := newTestContext(options, config)
	if err != nil {
		return nil, err
	}

	return &TestRunner{
		config: config,
		context: context,
	}, nil
}

func (tr *TestRunner) createEngine() (*engine, error) {
	engine, err := newEngine(tr.config, tr.context)

	return engine, err
}

func (tr *TestRunner) Run() (bool, error) {
	tp, err := CompileTestPlan(tr.config)
	if err != nil {
		return false, err
	}

	engine, err := tr.createEngine()
	if err != nil {
		return false, err
	}

	return engine.Run(tp)
}

