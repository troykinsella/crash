package runtime

import (
	"github.com/troykinsella/crash"
	"github.com/troykinsella/crash/util"
	"github.com/troykinsella/crash/action"
	"github.com/troykinsella/crash/system"
	"github.com/troykinsella/crash/logging"
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

	ctx.log.Finish(logging.PLAN, root.RunDuration(), plan.plan.Name)

	if result.Error != nil {
		return false, result.Error
	}

	return result.Success, nil
}

func (e *engine) runStep(step *StepExec, ctx *Context) (chan *StepResult) {
	ch := make(chan *StepResult)
	s := step.Step

	var ch2 chan *StepResult

	ctx = ctx.NewChild()
	step.Start()

	if s.Run != nil {
		ch2 = e.runAction(step, ctx)
	} else if s.Serial != nil {
		ch2 = e.runSerial(step, ctx)
	} else if s.Parallel != nil {
		ch2 = e.runParallel(step, ctx)
	}

	go func() {
		select {
		case result := <-ch2:
			step.Finish()
			ctx.Commit()
			e.afterStep(step, ctx, result, ch)
		case <-util.Timeout(step.Step.Timeout):
		    // TODO: better error
			ch <- &StepResult{
				Success: false,
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
				Success: false,
				Error: err,
			}
			return
		}

		e.rootCtx.log.Check(ok, msg, ctx.vars.AsMap())

		rok = rok && ok
	}

	result.Success = rok
	ch <- result
}

func (e *engine) runAction(step *StepExec, ctx *Context) (chan *StepResult) {

	ch := make(chan *StepResult)
	s := step.Step.Run

	ctx.log.Start(logging.ACTION, s.Name)

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
			Success: success,
			Error: err,
		}

		ctx.log.Finish(logging.ACTION, step.RunDuration(), s.Name, data)
	}()

	return ch
}

func (e *engine) runSerial(step *StepExec, ctx *Context) (chan *StepResult) {

	ctx.log.Start(logging.SERIAL, "serial...")
	sw := util.NewStopWatch().Start()

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

		ctx.log.Finish(logging.SERIAL, sw.Stop().Time(), "")
	}()

	return ch
}

func (e *engine) runParallel(step *StepExec, ctx *Context) (chan *StepResult) {
	ctx.log.Start(logging.PARALLEL, "parallel...")
	sw := util.NewStopWatch().Start()

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

		ctx.log.Finish(logging.PARALLEL, sw.Stop().Time(), "")
	}()

	return ch
}

func aggregateResults(results []*StepResult) *StepResult {
	result := true
	for _, r := range results {
		if !r.Success {
			result = false
			break
		}
	}
	return &StepResult{
		Success: result,
	}
}
