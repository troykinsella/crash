package runtime

import (
	"github.com/troykinsella/crash"
	"github.com/troykinsella/crash/util"
	"github.com/troykinsella/crash/action"
	"github.com/troykinsella/crash/system"
	"github.com/troykinsella/crash/logging"
	"fmt"
	"errors"
	"github.com/troykinsella/crash/system/data"
)

type engine struct {
	rootCtx *Context
	interp  *system.Interpreter
}

func newEngine(config *crash.Config, ctx *Context) (*engine, error) {

	i, err := system.NewInterpreter(ctx.log, config.Checks)
	if err != nil {
		return nil, err
	}

	return &engine{
		rootCtx: ctx,
		interp: i,
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

	ctx.log.Finish(logging.PLAN, result.Ok, root.RunDuration(), plan.plan.Name, nil)

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

func (e *engine) iterateWithList(ss *system.ScriptList, ctx *Context) (util.Iterator, error) {
	_, list, _, err := e.interp.RunAll(ss, ctx.vars)
	if err != nil {
		return nil, err
	}

	return util.NewSliceIterator(list)
}

func (e *engine) iterateWithItem(s *system.Script, ctx *Context) (itr util.Iterator, err error) {
	_, result, _, err := e.interp.Run(s, ctx.vars)
	if err != nil {
		return nil, err
	}

	switch val := result.(type) {
	case []interface{}:
		itr, err = util.NewSliceIterator(val)
	case map[interface{}]interface{}:
		fmt.Printf("MAP\n")
	case int64:
		itr, err = util.NewRangeIterator(0, int(val), 1)
	default:
		n, err := data.ToInt(val)
		if err == nil {
			itr, err = util.NewRangeIterator(0, int(n), 1)
		}
	}
	return
}

func (e *engine) iterateWith(w *WithDirective, ctx *Context) (util.Iterator, error) {
	switch {
	case w.Item != nil:
		return e.iterateWithItem(w.Item, ctx)
	case w.List != nil:
		return e.iterateWithList(w.List, ctx)
	}
	return nil, errors.New("internal error")
}

func (e *engine) runStepLoop(se *StepExec, ctx *Context) (chan *StepResult) {
	ch := make(chan *StepResult)

	itr, err := e.iterateWith(se.Step.With, ctx)
	if err != nil {
		panic(err)
	}

	go func() {
		results := make([]*StepResult, 0)

		for itr.HasNext() {
			val := itr.Next()

			idx := se.Step.With.As
			if idx == "" {
				idx = "i"
			}

			vCtx := ctx.NewChild()
			vCtx.vars.Set(idx, val)

			vStep := NewStepExec(se.Step)
			vStep.looped = true

			subCh := e.runStep(vStep, vCtx)

			r := <- subCh
			results = append(results, r)
		}

		ch <- aggregateResults(results)
	}()

	return ch
}

func (e *engine) runStep(se *StepExec, ctx *Context) (chan *StepResult) {
	s := se.Step
	if s.With != nil && !se.looped {
		return e.runStepLoop(se, ctx)
	}

	ch := make(chan *StepResult)

	mt, stepName := e.selectStep(se, ctx)
	var ch2 chan *StepResult

	ctx = ctx.NewChild()
	ctx.log.Start(mt, stepName)
	se.Start()

	switch {
	case s.Run != nil:
		ch2 = e.runAction(se, ctx)
	case s.Serial != nil:
		ch2 = e.runSerial(se, ctx)
	case s.Parallel != nil:
		ch2 = e.runParallel(se, ctx)
	}

	go func() {
		select {
		case result := <-ch2:
			se.Finish()
			ctx.Commit()
			ctx.log.Finish(mt, result.Ok, se.RunDuration(), stepName, result.Data)
			e.afterStep(se, ctx, result, ch)

		case <-util.Timeout(se.Step.Timeout):
		    ctx.log.Error(mt, se.Step.Timeout, fmt.Errorf("timed out"), stepName)
			ch <- &StepResult{
				Ok: false,
			}
		}
	}()

	return ch
}

func (e *engine) afterStep(se *StepExec,
                           ctx *Context,
                           result *StepResult,
                           ch chan *StepResult) {
	if se.Step.Checks != nil {
		e.doChecks(se, ctx, result, ch)
	} else {
		ch <- result
	}
}

func (e *engine) doChecks(se *StepExec,
                          ctx *Context,
                          result *StepResult,
                          ch chan *StepResult) {
	rok := true
	for _, check := range *se.Step.Checks {
		ok, _, msg, err := e.interp.Run(check, ctx.vars)
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

func (e *engine) runAction(se *StepExec, ctx *Context) (chan *StepResult) {
	ch := make(chan *StepResult)

	go func() {
		s := se.Step.Run

		params := util.AsStringValues(s.Params)
		params = system.NewInterpolatedValues(params, ctx.vars)

		a := action.NewAction(&action.ActionConfig{
			Name:   s.Type,
			Params: params,
			Log: e.rootCtx.log,
		})

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

func (e *engine) runSerial(se *StepExec, ctx *Context) (chan *StepResult) {
	ch := make(chan *StepResult)

	go func() {
		s := se.Step.Serial

		stepLen := len(*s)
		results := make([]*StepResult, stepLen)

		for i, config := range *s {
			subStep := NewStepExec(&config)
			subCh := e.runStep(subStep, ctx)
			results[i] = <- subCh
		}

		ch <- aggregateResults(results)
	}()

	return ch
}

func (e *engine) runParallel(se *StepExec, ctx *Context) (chan *StepResult) {
	ch := make(chan *StepResult)

	go func() {
		p := se.Step.Parallel

		channels := make([]chan *StepResult, len(*p))

		for i, config := range *p {
			subStep := NewStepExec(&config)
			channels[i] = e.runStep(subStep, ctx)
		}

		var results = make([]*StepResult, len(channels))

		// fan in
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
