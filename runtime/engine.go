package runtime

import (
	"github.com/troykinsella/crash"
	"github.com/troykinsella/crash/util"
	"github.com/troykinsella/crash/action"
	"github.com/troykinsella/crash/system"
	"github.com/troykinsella/crash/logging"
	"fmt"
	"errors"
)

type engine struct {
	rootCtx *Context
	checkInterp *system.Interpreter
}

func newEngine(config *crash.Config, ctx *Context) (*engine, error) {

	i, err := system.NewInterpreter(ctx.log, config.Checks)
	if err != nil {
		return nil, err
	}

	return &engine{
		rootCtx: ctx,
		checkInterp: i,
	}, nil
}

func newRootStep(s *Steps) *StepExec {
	return NewStepExec(&Step{
		Serial: s,
	})
}

func (e *engine) Run(tp *TestPlan) (bool, error) {

	rok := true

	for _, plan := range(*tp.Plans) {
		plan := NewPlanExec(&plan)
		ok, err := e.runPlan(plan, e.rootCtx)
		if err != nil {
			return false, err
		}
		rok = rok && ok
	}

	return rok, nil
}

func (e *engine) runPlan(plan *PlanExec, ctx *Context) (bool, error) {

	ctx.log.Start(logging.PLAN, plan.plan.Name)

	root := newRootStep(plan.plan.Steps)
	root.Start()

	result := <- e.runStep(root, ctx)

	root.Finish()

	ctx.log.Finish(logging.PLAN, result.Ok, root.RunDuration(), plan.plan.Name)

	if result.Error != nil {
		return false, result.Error
	}

	return result.Ok, nil
}

func (e *engine) selectStep(step *StepExec, ctx *Context) (logging.MessageType, string) {
	s := step.Step
	switch {
	case s.Run != nil:
		return logging.ACTION, s.Run.Name
	case s.Serial != nil:
		return logging.SERIAL, "serial"
	case s.Parallel != nil:
		return logging.PARALLEL, "parallel"
	default:
		panic(errors.New("internal error"))
	}
}

func (e *engine) runStep(step *StepExec, ctx *Context) (chan *StepResult) {
	ch := make(chan *StepResult)
	s := step.Step

	mt, stepName := e.selectStep(step, ctx)
	var ch2 chan *StepResult

	ctx = ctx.NewChild()
	ctx.log.Start(mt, stepName)
	step.Start()

	switch {
	case s.Run != nil:
		ch2 = e.runAction(step, ctx)
	case s.Serial != nil:
		ch2 = e.runSerial(step, ctx)
	case s.Parallel != nil:
		ch2 = e.runParallel(step, ctx)
	}

	go func() {
		select {
		case result := <-ch2:
			step.Finish()
			ctx.Commit()
			ctx.log.Finish(mt, result.Ok, step.RunDuration(), stepName, result.Data)
			e.afterStep(step, ctx, result, ch)

		case <-util.Timeout(step.Step.Timeout):
		    ctx.log.Error(mt, step.Step.Timeout, fmt.Errorf("timed out"), stepName)
			ch <- &StepResult{
				Ok: false,
			}
		}
	}()

	return ch
}

func (e *engine) afterStep(step *StepExec,
                           ctx *Context,
                           result *StepResult,
                           ch chan *StepResult) {
	if step.Step.Checks != nil {
		e.doChecks(step, ctx, result, ch)
	} else {
		ch <- result
	}
}

func (e *engine) doChecks(step *StepExec,
                          ctx *Context,
                          result *StepResult,
                          ch chan *StepResult) {
	rok := true
	for _, check := range *step.Step.Checks {
		ok, _, msg, err := e.checkInterp.Statement(check, ctx.vars)
		if err != nil {
			ch <- &StepResult{
				Ok: false,
				Error: err,
			}
			return
		}

		e.rootCtx.log.Check(ok, msg, ctx.vars.AsMap())

		rok = rok && ok
	}

	result.Ok = rok
	ch <- result
}

func (e *engine) runAction(step *StepExec, ctx *Context) (chan *StepResult) {
	ch := make(chan *StepResult)
	s := step.Step.Run

	params := util.AsStringValues(s.Params)
	params = system.NewInterpolatedValues(params, ctx.vars)

	a := action.NewAction(&action.ActionConfig{
		Name:   s.Type,
		Params: params,
		Log: e.rootCtx.log,
	})

	go func() {
		result, err := a.Run()

		success := err == nil
		var data map[string]interface{}

		if success {
			data = result.Data
			ctx.vars.SetAll(data)
		}

		ch <- &StepResult{
			Ok: success,
			Data: data,
			Error: err,
		}
	}()

	return ch
}

func (e *engine) runSerial(step *StepExec, ctx *Context) (chan *StepResult) {
	ch := make(chan *StepResult)
	s := step.Step.Serial

	stepLen := len(*s)
	results := make([]*StepResult, stepLen)

	go func() {
		for i, config := range *s {
			subStep := NewStepExec(&config)
			subCh := e.runStep(subStep, ctx)
			results[i] = <- subCh
		}

		ch <- aggregateResults(results)
	}()

	return ch
}

func (e *engine) runParallel(step *StepExec, ctx *Context) (chan *StepResult) {
	ch := make(chan *StepResult)
	p := step.Step.Parallel

	channels := make([]chan *StepResult, len(*p))

	for i, config := range *p {
		subStep := NewStepExec(&config)
		channels[i] = e.runStep(subStep, ctx)
	}

	go func() {
		var results = make([]*StepResult, len(channels))
		for i, ch := range channels {
			results[i] = <- ch
		}

		ch <- aggregateResults(results)
	}()

	return ch
}

func aggregateResults(results []*StepResult) *StepResult {
	result := true
	for _, r := range results {
		if !r.Ok {
			result = false
			break
		}
	}
	return &StepResult{
		Ok: result,
	}
}
