package app

import (
	"errors"
	"github.com/troykinsella/crash/runtime"
	"github.com/troykinsella/crash"
	"os"
	"fmt"
)

type Crash struct {
}

func New() *Crash {
	return &Crash{}
}

func exists(f string) bool {
	_, err := os.Stat(f)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

func findCrashFile(f string) (string, error) {
	if f != "" {
		if !exists(f) {
			return "", fmt.Errorf("test file not found: %s", f)
		}
		return f, nil
	}

	so := [...]string{
		"Crashfile",
		"Crashfile.yml",
		"Crashfile.yaml",
	}

	for _, f := range so {
		if exists(f) {
			return f, nil
		}
	}

	return "", errors.New("cannot locate Crashfile")
}

func (*Crash) Test(options *crash.TestOptions) (bool, error) {
	cf, err:= findCrashFile(options.Crashfile)
	if err != nil {
		return false, err
	}

	config := crash.NewConfig()
	err = config.UnmarshalYAMLFile(cf)
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

func (*Crash) Validate(f string) error {
	cf, err:= findCrashFile(f)
	if err != nil {
		return err
	}

	config := crash.NewConfig()
	err = config.UnmarshalYAMLFile(cf)
	if err != nil {
		return err
	}

	return runtime.ValidateConfig(config)
}
