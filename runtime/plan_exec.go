package runtime

import (
	"github.com/troykinsella/crash"
)

type PlanExec struct {
	config *crash.PlanConfig
}

func NewPlanExec(plan *crash.PlanConfig) *PlanExec {
	return &PlanExec{
		config: plan,
	}
}






