package runtime

import (
	"github.com/troykinsella/crash"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Timed interface {

}

type TestRunner struct {
	config  *crash.Config
	context *Context
}

func loadVariablesFiles(files []string) (map[string]interface{}, error) {
	vars := make(map[string]interface{})
	for _, file := range files {
		fileBytes, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(fileBytes, vars)
		if err != nil {
			return nil, err
		}
	}
	return vars, nil
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
	engine, err := tr.createEngine()
	if err != nil {
		return false, err
	}

	return engine.Run(tr.config)
}

