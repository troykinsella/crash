package util

import "time"

type StopWatch struct {
	start time.Time
	finish time.Time
}

func NewStopWatch() *StopWatch {
	return &StopWatch{}
}

func (t *StopWatch) Start() *StopWatch {
	t.start = time.Now()
	return t
}

func (t *StopWatch) Stop() *StopWatch {
	t.finish = time.Now()
	return t
}

func (t *StopWatch) Since() time.Duration {
	return time.Now().Sub(t.start)
}

func (t *StopWatch) Time() time.Duration {
	return t.finish.Sub(t.start)
}
