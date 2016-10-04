package app

import (
	"errors"
	"github.com/troykinsella/crash/runtime"
	"github.com/troykinsella/crash"
)

type Crash struct {
}

func New() (*Crash, error) {
	return &Crash{}, nil
}

func (*Crash) Test(options *crash.TestOptions) (bool, error) {
	if options.InputFile == "" {
		return false, errors.New("Must supply input file")
	}

	config := crash.NewConfig()
	err := config.UnmarshalYAMLFile(options.InputFile)
	if err != nil {
		return false, err
	}

	tr, err := runtime.NewTestRunner(options, config)
	if err != nil {
		return false, err
	}

	ok, err := tr.Run()
	if err != nil {
		return false, err
	}

	return ok, nil
}
