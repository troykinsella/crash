package runtime

import (
	"github.com/troykinsella/crash"
	"time"
	"github.com/troykinsella/crash/system"
	"fmt"
	"errors"
)

type TestPlan struct {
	Variables map[string]string
	Plans     *Plans
}

func ValidateConfig(config *crash.Config) error {
	_, err := CompileTestPlan(config)
	if err != nil {
		return err
	}
	return nil
}

func CompileTestPlan(config *crash.Config) (*TestPlan, error) {
	if config.Plans == nil {
		return nil, errors.New("Missing required property: plans")
	}

	plans := make(Plans, len(*config.Plans))
	for i, pc := range *config.Plans {
		p, err := newPlan(&pc)
		if err != nil {
			return nil, err
		}
		plans[i] = *p
	}

	return &TestPlan{
		Variables: config.Variables,
		Plans: &plans,
	}, nil
}

func newPlan(config *crash.PlanConfig) (*Plan, error) {
	if config.Steps == nil {
		return nil, fmt.Errorf("Plan has no steps: %s", config.Name)
	}

	steps, err := newSteps(config.Steps)
	if err != nil {
		return nil, err
	}

	return &Plan{
		Name: config.Name,
		Steps: steps,
	}, nil
}

func newSteps(configs *crash.StepConfigs) (*Steps, error) {
	steps := make(Steps, len(*configs))
	for i, sc := range *configs {
		s, err := newStep(&sc)
		if err != nil {
			return nil, err
		}
		steps[i] = *s
	}
	return &steps, nil
}

func parseDuration(s string) (time.Duration, error) {
	if s == "" {
		return -1, nil
	}
	return time.ParseDuration(s)
}

func newStep(config *crash.StepConfig) (*Step, error) {
	ss, err := system.NewScriptSet(config.Checks)
	if err != nil {
		return nil, err
	}

	to, err := parseDuration(config.Timeout)
	if err != nil {
		return nil, err
	}

	if config.Run != nil {
		as, err := newActionStep(config.Run)
		if err != nil {
			return nil, err
		}
		return &Step{
			Run: as,
			Checks: ss,
			Timeout: to,
		}, nil
	}

	if config.Serial != nil {
		s, err := newSteps(config.Serial)
		if err != nil {
			return nil, err
		}
		return &Step{
			Serial: s,
			Checks: ss,
			Timeout: to,
		}, nil
	}

	if config.Parallel != nil {
		p, err := newSteps(config.Parallel)
		if err != nil {
			return nil, err
		}
		return &Step{
			Parallel: p,
			Checks: ss,
			Timeout: to,
		}, nil
	}

	return nil, fmt.Errorf("wtf")
}

func newActionStep(config *crash.ActionConfig) (*ActionStep, error) {
	// TODO: some validation
	return &ActionStep{
		Name: config.Name,
		Type: config.Type,
		Params: config.Params,
	}, nil
}

type Plan struct {
	Name  string
	Steps *Steps
}

type Plans []Plan

type Step struct {
	Run      *ActionStep
	Serial   *Steps
	Parallel *Steps

	Checks  *system.ScriptSet
	Timeout time.Duration
}

type Steps []Step

type ActionStep struct {
	Name   string
	Type   string
	Params map[string]string
}
