package runtime

import (
	"time"
	"github.com/troykinsella/crash"
	"github.com/troykinsella/crash/util"
)

type StepExec struct {
	stopWatch *util.StopWatch
	Config *crash.StepConfig
}

type StepResult struct {
	Success bool
	Error error
}

func NewStepExec(config *crash.StepConfig) *StepExec {
	return &StepExec{
		stopWatch: util.NewStopWatch(),
		Config: config,
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

