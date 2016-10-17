package runtime

import (
	"time"
	"github.com/troykinsella/crash/util"
)

type StepExec struct {
	stopWatch *util.StopWatch
	Step *Step
}

type StepResult struct {
	Ok    bool
	Data  map[string]interface{}
	Error error
}

func NewStepExec(step *Step) *StepExec {
	return &StepExec{
		stopWatch: util.NewStopWatch(),
		Step: step,
	}
}

func (s *StepExec) Start() {
	s.stopWatch.Start()
}

func (s *StepExec) Finish() {
	s.stopWatch.Stop()
}

func (s *StepExec) RunDuration() time.Duration {
	return s.stopWatch.Since()
}

