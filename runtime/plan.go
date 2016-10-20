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
	if configs == nil {
		return nil, nil
	}

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

func parseTimeout(s string) (time.Duration, error) {
	if s == "" {
		return -1, nil
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return -1, err
	}
	if d <= 0 {
		return -1, fmt.Errorf("timeout must be greater than zero: %s", s)
	}

	return d, nil
}

func newStep(config *crash.StepConfig) (*Step, error) {
	as, err := newActionStep(config.Run)
	if err != nil {
		return nil, err
	}

	ss, err := newSteps(config.Serial)
	if err != nil {
		return nil, err
	}

	ps, err := newSteps(config.Parallel)
	if err != nil {
		return nil, err
	}

	ch, err := system.NewScriptList(config.Checks, system.STMT)
	if err != nil {
		return nil, err
	}

	to, err := parseTimeout(config.Timeout)
	if err != nil {
		return nil, err
	}

	wd, err := newWithDirective(config.With)
	if err != nil {
		return nil, err
	}

	if as == nil && ss == nil && ps == nil {
		return nil, errors.New("require action, serial, or parallel step")
	}

	return &Step{
		Run: as,
		Serial: ss,
		Parallel: ps,
		Checks: ch,
		Timeout: to,
		With: wd,
	}, nil
}

func newWithDirective(config *crash.WithConfig) (*WithDirective, error) {
	if config == nil {
		return nil, nil
	}

	it, err := system.NewScript(config.Item, system.EXPR)
	if err != nil {
		return nil, err
	}

	ls, err := system.NewScriptList(config.List, system.EXPR)
	if err != nil {
		return nil, err
	}

	return &WithDirective{
		Item: it,
		List: ls,
		Map: config.Map,

		As: config.As,
	}, nil
}

func newActionStep(config *crash.ActionConfig) (*ActionStep, error) {
	if config == nil {
		return nil, nil
	}

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

	Checks  *system.ScriptList
	Timeout time.Duration
	With    *WithDirective
}

type Steps []Step

type ActionStep struct {
	Name   string
	Type   string
	Params map[string]string
}

type WithDirective struct {
	Item    *system.Script
	List    *system.ScriptList
	Map     map[string]string

	As      string
}
